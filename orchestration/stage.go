package orchestration

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"agentdemo/tool"

	"github.com/mchenziyi/aisc/agents/prompts"
	"github.com/mchenziyi/aisc/state"
)

// ─── StageConfig ──────────────────────────────────────────────

// StageConfig 描述一个 Stage 的元信息。
// 所有 stage-specific 的值从这里取，StageRunner 不再硬编码。
type StageConfig struct {
	StageID        string   // "stage-requirement", "stage-api-design"
	StageName      string   // "Requirement", "API Design"（显示用）
	OwnerAgent     string   // "pm", "tech-lead"
	ArtifactName   string   // "prd", "api-spec"
	MeetingType    string   // "requirement_review", "api_review"
	MemoryTags     []string // 记忆标签
	MaxRounds      int      // 最大评审轮次

	// PromptDraft 返回该 Stage Owner 起草产物用的 system prompt
	PromptDraft func() (string, error)
	// PromptRevise 返回该 Stage Owner 修订产物用的 system prompt
	PromptRevise func() (string, error)

	// InputReader 读取本 Stage 的输入（如 requirement.md 或上一阶段冻结产物）
	InputReader func(root string) (string, error)

	// Tools 可用工具列表（Draft/Revise 时注入 Agent）。空=nil=纯文本模式
	Tools []tool.Tool

	// ReviewContentBuilder 构造评审时给 reviewer 看的内容。
	// 默认（nil）= 直接用 artifact 文本。Backend Stage 需要读取代码文件。
	ReviewContentBuilder func(root string, summary string) (string, error)

	// FreezeAction 冻结时的额外动作。
	// 默认（nil）= SaveFrozenArtifact(artifactName, artifact内容)。
	// Backend Stage 需要在这里快照 backend/ 目录。
	FreezeAction func(root string, summary string) error
}

// DefaultRequirementConfig 返回 Requirement Stage 的默认配置
func DefaultRequirementConfig() StageConfig {
	return StageConfig{
		StageID:     "stage-requirement",
		StageName:   "Requirement",
		OwnerAgent:  "pm",
		ArtifactName: "prd",
		MeetingType: "requirement_review",
		MemoryTags:  []string{"需求评审", "PRD"},
		MaxRounds:   5,
		PromptDraft: func() (string, error) { return prompts.Load("pm", "draft") },
		PromptRevise: func() (string, error) { return prompts.Load("pm", "revise") },
		InputReader: state.ReadRequirement,
	}
}

// DefaultAPIDesignConfig 返回 API Design Stage 的默认配置
func DefaultAPIDesignConfig() StageConfig {
	return StageConfig{
		StageID:      "stage-api-design",
		StageName:    "API Design",
		OwnerAgent:   "tech-lead",
		ArtifactName: "api-spec",
		MeetingType:  "api_review",
		MemoryTags:   []string{"API评审", "API Spec"},
		MaxRounds:    5,
		PromptDraft:  func() (string, error) { return prompts.Load("tech-lead", "api-design") },
		PromptRevise: func() (string, error) { return prompts.Load("tech-lead", "api-design-revise") },
		InputReader:  state.ReadFrozenPRD,
	}
}

// DefaultTechDesignConfig 返回 Tech Design Stage 的默认配置
func DefaultTechDesignConfig() StageConfig {
	return StageConfig{
		StageID:      "stage-tech-design",
		StageName:    "Tech Design",
		OwnerAgent:   "tech-lead",
		ArtifactName: "tech-design",
		MeetingType:  "tech_review",
		MemoryTags:   []string{"技术评审", "Tech Design"},
		MaxRounds:    5,
		PromptDraft:  func() (string, error) { return prompts.Load("tech-lead", "tech-design") },
		PromptRevise: func() (string, error) { return prompts.Load("tech-lead", "tech-design-revise") },
		InputReader:  state.ReadFrozenDesignDocs,
	}
}

// DefaultBackendConfig 返回 Backend Dev Stage 的默认配置
func DefaultBackendConfig() StageConfig {
	return StageConfig{
		StageID:      "stage-backend",
		StageName:    "Backend Dev",
		OwnerAgent:   "backend",
		ArtifactName: "backend",
		MeetingType:  "backend_review",
		MemoryTags:   []string{"后端开发", "Backend"},
		MaxRounds:    5,
		PromptDraft:  func() (string, error) { return prompts.Load("backend", "draft") },
		PromptRevise: func() (string, error) { return prompts.Load("backend", "revise") },
		InputReader:  state.ReadFrozenDesignDocs,
		Tools:        tool.AllBuiltInTools(),
		ReviewContentBuilder: func(root string, summary string) (string, error) {
			return state.ReadCodeDir(root, "backend")
		},
		FreezeAction: func(root string, summary string) error {
			if err := state.SaveFrozenArtifact(root, "backend", summary); err != nil {
				return err
			}
			return state.ArchiveDir(root, "backend", "backend")
		},
	}
}

// ─── StageRunner ──────────────────────────────────────────────

// StageRunner 驱动一个 Stage 完整执行。
type StageRunner struct {
	Root   string        // 项目根目录
	Orch   *Orchestrator // 评审编排器
	cfg    StageConfig
	stage  *state.Stage
}

// NewStageRunner 创建 Stage 执行器
func NewStageRunner(root string, orch *Orchestrator) *StageRunner {
	return &StageRunner{Root: root, Orch: orch}
}

// Run 执行指定 Stage 直至 freeze 或达到最大轮次。
func (sr *StageRunner) Run(ctx context.Context, cfg StageConfig) error {
	sr.cfg = cfg

	// 加载状态
	stage, err := state.LoadStage(sr.Root, cfg.StageID)
	if err != nil {
		return fmt.Errorf("load stage %s: %w", cfg.StageID, err)
	}
	sr.stage = stage

	if stage.Status == "frozen" {
		fmt.Printf("✅ %s Stage 已完成冻结\n", cfg.StageName)
		return nil
	}

	// 读取输入（Requirement: requirement.md; API Design: prd-frozen.md）
	input, err := cfg.InputReader(sr.Root)
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}

	// 断点续跑
	artifactExists := state.ArtifactExists(sr.Root, cfg.ArtifactName)
	roundNum := 1
	if artifactExists && stage.CurrentVersion > 1 {
		roundNum = stage.CurrentVersion
		fmt.Printf("📋 检测到已有 %s v%d，从第 %d 轮继续评审\n", cfg.ArtifactName, roundNum, roundNum)
	} else {
		stage.CurrentVersion = 1
		state.SaveStage(sr.Root, stage)
	}

	// 加载上一轮 Decision（Scope Lock）
	var prevDecision *Decision
	if dm, err := state.LoadDecisionMemory(sr.Root); err == nil {
		prevDecision = parseDecisionMap(dm)
	}

	for roundNum <= cfg.MaxRounds {
		reviewType := "全量深度评审"
		if roundNum > 1 {
			reviewType = "定向复核"
		}
		fmt.Printf("\n========== %s 第 %d 轮评审（%s）==========\n\n", cfg.StageName, roundNum, reviewType)

		// 获取或生成产物
		var artifact string
		if roundNum == 1 && !artifactExists {
			fmt.Printf("🚀 %s Agent 起草 %s v1...\n", cfg.OwnerAgent, cfg.ArtifactName)
			artifact, err = sr.generateArtifact(ctx, input)
			if err != nil {
				return fmt.Errorf("generate %s: %w", cfg.ArtifactName, err)
			}
			state.SaveArtifact(sr.Root, cfg.StageID, cfg.ArtifactName, artifact, 1)
			fmt.Printf("✅ %s v1 已保存\n", cfg.ArtifactName)
		} else {
			artifact, err = state.ReadArtifact(sr.Root, cfg.ArtifactName)
			if err != nil {
				return fmt.Errorf("read %s: %w", cfg.ArtifactName, err)
			}
		}

		// 创建 Meeting
		meeting := sr.createMeeting(roundNum)

		// 构造评审内容（Backend Stage 需要读取代码文件而非 summary）
		reviewContent := artifact
		if cfg.ReviewContentBuilder != nil {
			if built, err := cfg.ReviewContentBuilder(sr.Root, artifact); err == nil {
				reviewContent = built
			}
		}

		// 执行评审
		decision, reviews, err := sr.Orch.RunReviewRound(ctx, reviewContent, roundNum, prevDecision, stage.ReviewerAgents, cfg.ArtifactName)
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
			if err := sr.handleFreeze(ctx, artifact, decision, stage, meeting); err != nil {
				return fmt.Errorf("freeze: %w", err)
			}
			return nil

		case "revise":
			meeting.Meta.Status = "needs_revision"
			meeting.Meta.Decision = "revise"
			sr.saveMeetingWithDecision(meeting, decision)

			if roundNum >= cfg.MaxRounds {
				fmt.Printf("\n⚠️  已达最大评审轮次 (%d)，需要用户介入决策。\n", cfg.MaxRounds)
				return nil
			}

			fmt.Printf("🔧 修订 %s（%d 个行动项）...\n", cfg.ArtifactName, len(decision.ActionItems))
			artifact, err = sr.reviseArtifact(ctx, artifact, decision)
			if err != nil {
				return fmt.Errorf("revise %s: %w", cfg.ArtifactName, err)
			}
			stage.CurrentVersion++
			state.SaveArtifact(sr.Root, cfg.StageID, cfg.ArtifactName, artifact, stage.CurrentVersion)
			state.SaveStage(sr.Root, stage)
			fmt.Printf("✅ %s v%d 已保存\n", cfg.ArtifactName, stage.CurrentVersion)

			state.SaveDecisionMemory(sr.Root, decision)
			prevDecision = decision
			roundNum++

		case "reject":
			meeting.Meta.Status = "rejected"
			meeting.Body = fmt.Sprintf("## Decision (Reject)\n\n%s", decision.Summary)
			state.SaveMeeting(sr.Root, meeting)
			return fmt.Errorf("%s 被驳回: %s", cfg.ArtifactName, decision.Summary)

		default:
			return fmt.Errorf("未知决策类型: %s", decision.Type)
		}
	}

	return nil
}

// ─── freeze ───────────────────────────────────────────────────

func (sr *StageRunner) handleFreeze(
	ctx context.Context, artifact string,
	decision *Decision, stage *state.Stage, meeting *state.Meeting,
) error {
	cfg := sr.cfg

	// adopt 带 action_items → 先做静默修订再冻结，保存版本
	if decision.Type == "adopt" && len(decision.ActionItems) > 0 {
		fmt.Printf("🔧 adopt + %d 个微小修改 → 静默修订后冻结\n", len(decision.ActionItems))
		var err error
		artifact, err = sr.reviseArtifact(ctx, artifact, decision)
		if err != nil {
			return fmt.Errorf("静默修订失败: %w", err)
		}
		stage.CurrentVersion++
		state.SaveArtifact(sr.Root, cfg.StageID, cfg.ArtifactName, artifact, stage.CurrentVersion)
		state.SaveStage(sr.Root, stage)
	}
	state.SaveFrozenArtifact(sr.Root, cfg.ArtifactName, artifact)
	// 执行 Stage 特定的冻结动作（如代码快照）
	if cfg.FreezeAction != nil {
		if err := cfg.FreezeAction(sr.Root, artifact); err != nil {
			return fmt.Errorf("freeze action: %w", err)
		}
	}
	stage.Status = "frozen"
	state.SaveStage(sr.Root, stage)

	meeting.Meta.Status = "passed"
	meeting.Meta.Decision = "freeze"
	sr.saveMeetingWithDecision(meeting, decision)

	// 保存 reviewer memory
	for _, r := range meeting.Reviews {
		state.SaveMemory(sr.Root, r.AgentID, meeting.ID+"-review", &state.Memory{
			Type:    "decision",
			Title:   fmt.Sprintf("参与%s %s", cfg.StageName, meeting.ID),
			Content: truncate(r.Content, 2000),
			Relations: []state.Relation{{Type: "based_on", TargetType: "meeting", TargetID: meeting.ID}},
			Tags:     cfg.MemoryTags,
		})
	}

	state.DeleteDecisionMemory(sr.Root)
	fmt.Printf("✅ %s 已冻结！\n", cfg.ArtifactName)
	return nil
}

// ─── Helpers ──────────────────────────────────────────────────

func (sr *StageRunner) createMeeting(roundNum int) *state.Meeting {
	cfg := sr.cfg
	stage := sr.stage
	return &state.Meeting{
		ID: fmt.Sprintf("meeting-%03d", sr.nextMeetingCounter()),
		Meta: state.MeetingMeta{
			Round:          roundNum,
			ArtifactVersion: stage.CurrentVersion,
			Type:           cfg.MeetingType,
			Stage:          cfg.StageID,
			TargetArtifact: fmt.Sprintf(".aisc/stages/%02d-%s/artifact/%s-v%d%s", stage.Order, strings.ReplaceAll(strings.ToLower(cfg.StageName), " ", "-"), cfg.ArtifactName, stage.CurrentVersion, state.ArtifactExt(cfg.ArtifactName)),
			Moderator:      "project-manager",
			Participants:   strings.Join(stage.ReviewerAgents, ", "),
			Status:         "in_progress",
			CreatedAt:      time.Now().UTC().Format(time.RFC3339),
		},
	}
}

func (sr *StageRunner) saveMeetingWithDecision(meeting *state.Meeting, decision *Decision) {
	decisionJSON, _ := json.MarshalIndent(decision, "", "  ")
	meeting.Body = fmt.Sprintf("## Decision (%s)\n\n%s\n\n```json\n%s\n```", decision.Type, decision.Summary, string(decisionJSON))
	state.SaveMeeting(sr.Root, meeting)
}

func (sr *StageRunner) generateArtifact(ctx context.Context, input string) (string, error) {
	prompt, err := sr.cfg.PromptDraft()
	if err != nil {
		return "", err
	}
	if len(sr.cfg.Tools) > 0 {
		if tc, ok := sr.Orch.Client.(AgentClientWithTools); ok {
			return tc.RunWithTools(ctx, prompt, input, sr.cfg.Tools)
		}
	}
	return sr.Orch.Client.Run(ctx, prompt, input)
}

func (sr *StageRunner) reviseArtifact(ctx context.Context, artifact string, decision *Decision) (string, error) {
	prompt, err := sr.cfg.PromptRevise()
	if err != nil {
		return "", err
	}
	version := sr.stage.CurrentVersion + 1
	summary := decision.Summary
	actionText := ActionItemsText(decision.ActionItems)

	sysPrompt := fmt.Sprintf(prompt, version, summary, actionText, version)
	task := fmt.Sprintf("## 当前 %s (v%d)\n%s\n\n请根据上述 ActionItem 逐条修改，输出完整的 v%d。",
		sr.cfg.ArtifactName, sr.stage.CurrentVersion, artifact, version)
	if len(sr.cfg.Tools) > 0 {
		if tc, ok := sr.Orch.Client.(AgentClientWithTools); ok {
			return tc.RunWithTools(ctx, sysPrompt, task, sr.cfg.Tools)
		}
	}
	return sr.Orch.Client.Run(ctx, sysPrompt, task)
}

func (sr *StageRunner) nextMeetingCounter() int {
	counter := sr.stage.MeetingCounter
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

func parseDecisionMap(dm map[string]any) *Decision {
	d := &Decision{}
	if t, ok := dm["type"].(string); ok {
		d.Type = t
	}
	if s, ok := dm["summary"].(string); ok {
		d.Summary = s
	}
	if items, ok := dm["action_items"].([]any); ok {
		for _, item := range items {
			if m, ok := item.(map[string]any); ok {
				if desc, ok := m["description"].(string); ok {
					d.ActionItems = append(d.ActionItems, ActionItem{Description: desc})
				}
			}
		}
	}
	return d
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
