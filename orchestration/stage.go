package orchestration

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"agentdemo/agent"

	"github.com/mchenziyi/aisc/agents/prompts"
	"github.com/mchenziyi/aisc/state"
)

// ─── StageRunner ──────────────────────────────────────────────

// StageRunner 驱动一个 Stage 完整执行。
type StageRunner struct {
	Root  string         // 项目根目录
	Orch  *Orchestrator  // 评审编排器
	stage *state.Stage
}

// NewStageRunner 创建 Stage 执行器
func NewStageRunner(root string, orch *Orchestrator) *StageRunner {
	return &StageRunner{Root: root, Orch: orch}
}

// Run 执行 Requirement Stage 直至 freeze 或达到最大轮次。
func (sr *StageRunner) Run(ctx context.Context) error {
	// 加载状态
	stage, err := state.LoadStage(sr.Root, "stage-requirement")
	if err != nil {
		return fmt.Errorf("load stage: %w", err)
	}
	sr.stage = stage

	if stage.Status == "frozen" {
		fmt.Printf("✅ 项目已完成冻结，PRD: %s/docs/prd-frozen.md\n", sr.Root)
		return nil
	}

	// 读取用户需求
	req, err := state.ReadRequirement(sr.Root)
	if err != nil {
		return fmt.Errorf("read requirement: %w", err)
	}

	// 断点续跑
	prdExists := state.PRDExists(sr.Root)
	roundNum := 1
	if prdExists && stage.CurrentVersion > 1 {
		roundNum = stage.CurrentVersion
		fmt.Printf("📋 检测到已有 PRD v%d，从第 %d 轮继续评审\n", roundNum, roundNum)
	} else {
		stage.CurrentVersion = 1
		state.SaveStage(sr.Root, stage)
	}

	// 加载上一轮 Decision（Scope Lock）
	var prevDecision *Decision
	if dm, err := state.LoadDecisionMemory(sr.Root); err == nil {
		if t, ok := dm["type"].(string); ok {
			prevDecision = &Decision{Type: t}
			if s, ok := dm["summary"].(string); ok {
				prevDecision.Summary = s
			}
			if items, ok := dm["action_items"].([]any); ok {
				for _, item := range items {
					if m, ok := item.(map[string]any); ok {
						if desc, ok := m["description"].(string); ok {
							prevDecision.ActionItems = append(prevDecision.ActionItems, ActionItem{Description: desc})
						}
					}
				}
			}
		}
	}

	maxRounds := 5
	for roundNum <= maxRounds {
		reviewType := "全量深度评审"
		if roundNum > 1 {
			reviewType = "定向复核"
		}
		fmt.Printf("\n========== 第 %d 轮评审（%s）==========\n\n", roundNum, reviewType)

		// 获取或生成 PRD
		var prd string
		if roundNum == 1 && !prdExists {
			fmt.Println("🚀 PM Agent 起草 PRD v1...")
			var err error
			prd, err = sr.generatePRD(ctx, req)
			if err != nil {
				return fmt.Errorf("generate PRD: %w", err)
			}
			state.SaveArtifact(sr.Root, "stage-requirement", "prd", prd, 1)
			fmt.Printf("✅ PRD v1 已保存\n")
		} else {
			var err error
			prd, err = state.ReadArtifact(sr.Root, "prd")
			if err != nil {
				return fmt.Errorf("read PRD: %w", err)
			}
		}

		// 创建 Meeting
		meeting := &state.Meeting{
			ID: fmt.Sprintf("meeting-%03d", sr.nextMeetingCounter()),
			Meta: state.MeetingMeta{
				Round:          roundNum,
				PRDVersion:     stage.CurrentVersion,
				Type:           "requirement_review",
				Stage:          "requirement",
				TargetArtifact: fmt.Sprintf(".aisc/stages/01-requirement/artifact/prd-v%d.md", stage.CurrentVersion),
				Moderator:      "project-manager",
				Participants:   strings.Join(stage.ReviewerAgents, ", "),
				Status:         "in_progress",
				CreatedAt:      time.Now().UTC().Format(time.RFC3339),
			},
		}

		// 执行评审
		decision, reviews, err := sr.Orch.RunReviewRound(ctx, prd, roundNum, prevDecision)
		if err != nil {
			return fmt.Errorf("review round %d: %w", roundNum, err)
		}

		// 填充 meeting
		meeting.Meta.ID = meeting.ID
		for _, r := range reviews {
			meeting.Reviews = append(meeting.Reviews, state.Review{AgentID: r.AgentID, Content: r.Content})
		}
		stage.MeetingIDs = append(stage.MeetingIDs, meeting.ID)
		stage.MeetingCounter = sr.meetingCounterValue()
		state.SaveStage(sr.Root, stage)

		fmt.Printf("✅ 决策: %s — %s\n", decision.Type, truncate(decision.Summary, 200))

		// 执行决策
		switch decision.Type {
		case "adopt", "freeze":
			// adopt 带 action_items → 先做静默修订再冻结，保存版本
			if decision.Type == "adopt" && len(decision.ActionItems) > 0 {
				fmt.Printf("🔧 adopt + %d 个微小修改 → 静默修订后冻结\n", len(decision.ActionItems))
				var err error
				prd, err = sr.revisePRD(ctx, prd, decision)
				if err == nil {
					stage.CurrentVersion++
					state.SaveArtifact(sr.Root, "stage-requirement", "prd", prd, stage.CurrentVersion)
					state.SaveStage(sr.Root, stage)
				}
			}
			state.SaveFrozenPRD(sr.Root, prd)
			stage.Status = "frozen"
			state.SaveStage(sr.Root, stage)

			meeting.Meta.Status = "passed"
			meeting.Meta.Decision = "freeze"
			decisionJSON, _ := json.MarshalIndent(decision, "", "  ")
			meeting.Body = fmt.Sprintf("## Decision (%s)\n\n%s\n\n```json\n%s\n```", decision.Type, decision.Summary, string(decisionJSON))
			state.SaveMeeting(sr.Root, meeting)

			// 保存 reviewer memory
			for _, r := range meeting.Reviews {
				state.SaveMemory(sr.Root, r.AgentID, meeting.ID+"-review", &state.Memory{
					Type:    "decision",
					Title:   fmt.Sprintf("参与需求评审 %s", meeting.ID),
					Content: truncate(r.Content, 2000),
					Relations: []state.Relation{{Type: "based_on", TargetType: "meeting", TargetID: meeting.ID}},
					Tags:     []string{"需求评审", "PRD"},
				})
			}

			state.DeleteDecisionMemory(sr.Root)
			fmt.Println("✅ PRD 已冻结！MVP 闭环完成！")
			return nil

		case "revise":
			meeting.Meta.Status = "needs_revision"
			meeting.Meta.Decision = "revise"
			decisionJSON, _ := json.MarshalIndent(decision, "", "  ")
			meeting.Body = fmt.Sprintf("## Decision (Revise)\n\n%s\n\n```json\n%s\n```", decision.Summary, string(decisionJSON))
			state.SaveMeeting(sr.Root, meeting)

			if roundNum >= maxRounds {
				fmt.Printf("\n⚠️  已达最大评审轮次 (%d)，需要用户介入决策。\n", maxRounds)
				return nil
			}

			fmt.Printf("🔧 修订 PRD（%d 个行动项）...\n", len(decision.ActionItems))
			prd, err = sr.revisePRD(ctx, prd, decision)
			if err != nil {
				return fmt.Errorf("revise PRD: %w", err)
			}
			stage.CurrentVersion++
			state.SaveArtifact(sr.Root, "stage-requirement", "prd", prd, stage.CurrentVersion)
			state.SaveStage(sr.Root, stage)
			fmt.Printf("✅ PRD v%d 已保存\n", stage.CurrentVersion)

			// 保存 Decision 作为下一轮的 Scope Lock
			state.SaveDecisionMemory(sr.Root, decision)
			prevDecision = decision
			roundNum++

		case "reject":
			meeting.Meta.Status = "rejected"
			meeting.Body = fmt.Sprintf("## Decision (Reject)\n\n%s", decision.Summary)
			state.SaveMeeting(sr.Root, meeting)
			return fmt.Errorf("PRD 被驳回: %s", decision.Summary)

		default:
			return fmt.Errorf("未知决策类型: %s", decision.Type)
		}
	}

	return nil
}

// ─── Helpers ──────────────────────────────────────────────────

func (sr *StageRunner) generatePRD(ctx context.Context, requirement string) (string, error) {
	pmPrompt, err := prompts.Load("pm", "draft")
	if err != nil {
		return "", err
	}
	a, err := agent.New(sr.Orch.APIKey, sr.Orch.Model, false)
	if err != nil {
		return "", err
	}
	a.SetSystemPrompt(pmPrompt)
	a.MaxMessages = 5000
	return a.Run(ctx, "请根据以下用户需求，输出完整 PRD：\n\n"+requirement)
}

func (sr *StageRunner) revisePRD(ctx context.Context, prd string, decision *Decision) (string, error) {
	revisePrompt, err := prompts.Load("pm", "revise")
	if err != nil {
		return "", err
	}
	version := sr.stage.CurrentVersion + 1
	summary := decision.Summary
	actionText := ActionItemsText(decision.ActionItems)

	sysPrompt := fmt.Sprintf(revisePrompt, version, summary, actionText, version)

	a, err := agent.New(sr.Orch.APIKey, sr.Orch.Model, false)
	if err != nil {
		return "", err
	}
	a.SetSystemPrompt(sysPrompt)
	a.MaxMessages = 5000

	task := fmt.Sprintf("## 当前 PRD (v%d)\n%s\n\n请根据上述 ActionItem 逐条修改 PRD，输出完整的 v%d。", sr.stage.CurrentVersion, prd, version)
	return a.Run(ctx, task)
}

func (sr *StageRunner) nextMeetingCounter() int {
	counter := sr.stage.MeetingCounter
	// 从 meeting_ids 推断种子
	for _, mid := range sr.stage.MeetingIDs {
		var n int
		fmt.Sscanf(mid, "meeting-%d", &n)
		if n > counter {
			counter = n
		}
	}
	return counter + 1
}

func (sr *StageRunner) meetingCounterValue() int {
	c := sr.stage.MeetingCounter
	for _, mid := range sr.stage.MeetingIDs {
		var n int
		fmt.Sscanf(mid, "meeting-%d", &n)
		if n > c {
			c = n
		}
	}
	return c
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
