package orchestration

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"agentdemo/tool"

	"github.com/mchenziyi/aisc/agents/prompts"
	"github.com/mchenziyi/aisc/logger"
	"github.com/mchenziyi/aisc/state"
)

// ─── StageConfig ──────────────────────────────────────────────

// StageConfig 描述一个 Stage 的元信息。
// 所有 stage-specific 的值从这里取，StageRunner 不再硬编码。
type StageConfig struct {
	StageID      string   // "stage-requirement", "stage-api-design"
	StageName    string   // "Requirement", "API Design"（显示用）
	OwnerAgent   string   // "pm", "tech-lead"
	ArtifactName string   // "prd", "api-spec"
	MeetingType  string   // "requirement_review", "api_review"
	MemoryTags   []string // 记忆标签
	MaxRounds    int      // 最大评审轮次

	// PromptDraft 返回该 Stage Owner 起草产物用的 system prompt
	PromptDraft func() (string, error)
	// PromptRevise 返回该 Stage Owner 修订产物用的 system prompt
	PromptRevise func() (string, error)

	// PromptSmokeFix 返回冒烟测试失败后自动修复用的 prompt（可选，默认=PromptRevise）
	PromptSmokeFix func() (string, error)

	// ReviewPromptDir review prompt 文件所在的子目录（相对于 prompts/）。
	// 例如 "requirement", "api-design", "backend"。空字符串回退到 _shared。
	ReviewPromptDir string

	// InputReader 读取本 Stage 的输入（如 requirement.md 或上一阶段冻结产物）
	InputReader func(root string) (string, error)

	// Tools 可用工具列表（Draft/Revise 时注入 Agent）。空=nil=纯文本模式
	Tools []tool.Tool

	// ReviewContentBuilder 构造评审时给 reviewer 看的内容。
	// 默认（nil）= 直接用 artifact 文本。
	// roundNum=1 全量，roundNum>=2 只传涉及 action item 的文件。
	ReviewContentBuilder func(root string, summary string, roundNum int, prevDecision *Decision) (string, error)

	// FreezeAction 冻结时的额外动作。
	// 默认（nil）= SaveFrozenArtifact(artifactName, artifact内容)。
	// Backend Stage 需要在这里快照 backend/ 目录。
	FreezeAction func(root string, summary string) error

	// SmokeTester 冒烟测试（Draft/Revise 后自动执行）。
	// 返回 nil = 通过，error = 失败触发自动修复循环。
	SmokeTester func(root string) error

	// MaxSmokeRetries 冒烟失败最大自动修复次数（默认 3）
	MaxSmokeRetries int

	// TechStackExtractor 在 freeze 后从冻结文档提取技术栈（仅 Tech Design 需要）
	TechStackExtractor func(ctx context.Context, root string, runLLM func(string, string) (string, error)) error
}

// DefaultRequirementConfig 返回 Requirement Stage 的默认配置
func DefaultRequirementConfig() StageConfig {
	return StageConfig{
		StageID:         "stage-requirement",
		StageName:       "Requirement",
		OwnerAgent:      "pm",
		ArtifactName:    "prd",
		MeetingType:     "requirement_review",
		MemoryTags:      []string{"需求评审", "PRD"},
		MaxRounds:       5,
		PromptDraft:     func() (string, error) { return prompts.Load("pm", "draft") },
		PromptRevise:    func() (string, error) { return prompts.Load("pm", "revise") },
		InputReader:     state.ReadRequirement,
		ReviewPromptDir: "requirement",
	}
}

// DefaultAPIDesignConfig 返回 API Design Stage 的默认配置
func DefaultAPIDesignConfig() StageConfig {
	return StageConfig{
		StageID:         "stage-api-design",
		StageName:       "API Design",
		OwnerAgent:      "tech-lead",
		ArtifactName:    "api-spec",
		MeetingType:     "api_review",
		MemoryTags:      []string{"API评审", "API Spec"},
		MaxRounds:       5,
		PromptDraft:     func() (string, error) { return prompts.Load("tech-lead", "api-design") },
		PromptRevise:    func() (string, error) { return prompts.Load("tech-lead", "api-design-revise") },
		InputReader:     state.ReadFrozenPRD,
		ReviewPromptDir: "api-design",
	}
}

// DefaultTechDesignConfig 返回 Tech Design Stage 的默认配置
func DefaultTechDesignConfig() StageConfig {
	return StageConfig{
		StageID:         "stage-tech-design",
		StageName:       "Tech Design",
		OwnerAgent:      "tech-lead",
		ArtifactName:    "tech-design",
		MeetingType:     "tech_review",
		MemoryTags:      []string{"技术评审", "Tech Design"},
		MaxRounds:       5,
		PromptDraft:     func() (string, error) { return prompts.Load("tech-lead", "tech-design") },
		PromptRevise:    func() (string, error) { return prompts.Load("tech-lead", "tech-design-revise") },
		InputReader:        state.ReadFrozenPRDAndAPI,
		ReviewPromptDir:    "tech-design",
		TechStackExtractor: state.ExtractTechStack,
	}
}

// DefaultBackendConfig 返回 Backend Dev Stage 的默认配置
func DefaultBackendConfig() StageConfig {
	return StageConfig{
		StageID:         "stage-backend",
		StageName:       "Backend",
		OwnerAgent:      "backend",
		ArtifactName:    "backend",
		MeetingType:     "backend_review",
		MemoryTags:      []string{"后端开发", "Backend"},
		MaxRounds:       5,
		PromptDraft:     func() (string, error) { return prompts.Load("backend", "draft") },
		PromptRevise:    func() (string, error) { return prompts.Load("backend", "revise") },
		PromptSmokeFix:  func() (string, error) { return prompts.Load("backend", "fix") },
		InputReader:     state.ReadFrozenDesignDocs,
		ReviewPromptDir: "backend",
		Tools:           tool.AllBuiltInTools(),
		MaxSmokeRetries: 3,
		SmokeTester:     state.BackendSmokeTest,
		ReviewContentBuilder: func(root string, summary string, roundNum int, prevDecision *Decision) (string, error) {
			var actions []string
			if prevDecision != nil {
				for _, a := range prevDecision.ActionItems {
					actions = append(actions, a.Description)
				}
			}
			return state.SmartReadCodeDir(root, "backend", roundNum, actions)
		},
		FreezeAction: func(root string, summary string) error {
			code, err := state.ReadCodeDir(root, "backend")
			if err != nil {
				return err
			}
			if err := state.SaveFrozenArtifact(root, "backend", code); err != nil {
				return err
			}
			return state.ArchiveDir(root, "backend", "backend")
		},
	}
}

// ─── StageRunner ──────────────────────────────────────────────

// StageRunner 驱动一个 Stage 完整执行。
type StageRunner struct {
	Root  string        // 项目根目录
	Orch  *Orchestrator // 评审编排器
	cfg   StageConfig
	stage *state.Stage
	log   *logger.Logger
}

// NewStageRunner 创建 Stage 执行器
func NewStageRunner(root string, orch *Orchestrator) *StageRunner {
	return &StageRunner{Root: root, Orch: orch}
}

// Run 执行指定 Stage 直至 freeze 或达到最大轮次。
func (sr *StageRunner) Run(ctx context.Context, cfg StageConfig) error {
	sr.cfg = cfg

	// 获取文件锁 — 防止并发运行同一 Stage
	lock, err := state.LockStage(sr.Root, cfg.StageID)
	if err != nil {
		return fmt.Errorf("lock stage %s: %w", cfg.StageID, err)
	}
	defer state.UnlockStage(lock)

	// 初始化日志
	sr.log, err = logger.New(sr.Root, cfg.ArtifactName)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	defer sr.log.Close()
	sr.log.Info("stage_start")

	// 注入 logger 到 AgentClient
	if qc, ok := sr.Orch.Client.(*QiuQiuProClient); ok {
		qc.Log = sr.log
	}

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
		sr.log.Log(logger.INFO, "round_start", 0, logger.F{"round": roundNum, "type": reviewType})

		// 获取或生成产物
		var artifact string
		if roundNum == 1 && !artifactExists {
			fmt.Printf("🚀 %s Agent 起草 %s v1...\n", cfg.OwnerAgent, cfg.ArtifactName)
			t0 := time.Now()
			artifact, err = sr.generateArtifact(ctx, input)
			sr.log.Log(logger.INFO, "draft", time.Since(t0).Milliseconds(), nil)
			if err != nil {
				return fmt.Errorf("generate %s: %w", cfg.ArtifactName, err)
			}
			state.SaveArtifact(sr.Root, cfg.StageID, cfg.ArtifactName, artifact, 1)
			fmt.Printf("✅ %s v1 已保存\n", cfg.ArtifactName)

			// 冒烟测试
			artifact, err = sr.runSmokeLoop(ctx, artifact, roundNum)
			if err != nil {
				return fmt.Errorf("smoke: %w", err)
			}
			state.SaveArtifact(sr.Root, cfg.StageID, cfg.ArtifactName, artifact, 1)
		} else {
			artifact, err = state.ReadArtifact(sr.Root, cfg.ArtifactName)
			if err != nil {
				return fmt.Errorf("read %s: %w", cfg.ArtifactName, err)
			}
		}

		// 创建 Meeting
		meeting := sr.createMeeting(roundNum)

		// 构造评审内容（Backend Stage 需要读取代码文件）
		reviewContent := artifact
		if cfg.ReviewContentBuilder != nil {
			if built, err := cfg.ReviewContentBuilder(sr.Root, artifact, roundNum, prevDecision); err == nil {
				reviewContent = built
			}
		}

		// 执行评审
		t1 := time.Now()
		decision, reviews, err := sr.Orch.RunReviewRound(ctx, reviewContent, roundNum, prevDecision, stage.ReviewerAgents, cfg.ArtifactName, cfg.ReviewPromptDir)
		sr.log.Log(logger.INFO, "review_round", time.Since(t1).Milliseconds(), nil)
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
		sr.log.Log(logger.INFO, "decision", 0, logger.F{"type": decision.Type, "action_items": len(decision.ActionItems)})

		// 执行决策
		switch decision.Type {
		case "adopt", "freeze":
			if err := sr.handleFreeze(ctx, artifact, decision, stage, meeting); err != nil {
				sr.log.Error("freeze_failed", logger.F{"error": err.Error()})
				return fmt.Errorf("freeze: %w", err)
			}
			// 提取技术栈（仅 Tech Design Stage 有配置）
			if cfg.TechStackExtractor != nil {
				runLLM := func(sys, user string) (string, error) {
					return sr.Orch.Client.Run(ctx, sys, user)
				}
				if err := cfg.TechStackExtractor(ctx, sr.Root, runLLM); err != nil {
					sr.log.Error("techstack_extract", logger.F{"error": err.Error()})
				}
			}
			sr.log.Info("stage_frozen")
			return nil

		case "revise":
			meeting.Meta.Status = "needs_revision"
			meeting.Meta.Decision = "revise"
			if err := sr.saveMeetingWithDecision(meeting, decision); err != nil {
				sr.log.Error("save_meeting", logger.F{"error": err.Error()})
			}

			if roundNum >= cfg.MaxRounds {
				fmt.Printf("\n⚠️  已达最大评审轮次 (%d)，需要用户介入决策。\n", cfg.MaxRounds)
				return nil
			}

			fmt.Printf("🔧 修订 %s（%d 个行动项）...\n", cfg.ArtifactName, len(decision.ActionItems))
			sr.log.Log(logger.INFO, "revise_start", 0, logger.F{"action_items": len(decision.ActionItems)})
			t2 := time.Now()
			artifact, err = sr.reviseArtifact(ctx, artifact, decision)
			if err != nil {
				return fmt.Errorf("revise %s: %w", cfg.ArtifactName, err)
			}
			sr.log.Log(logger.INFO, "revise", time.Since(t2).Milliseconds(), nil)
			stage.CurrentVersion++
			state.SaveArtifact(sr.Root, cfg.StageID, cfg.ArtifactName, artifact, stage.CurrentVersion)
			state.SaveStage(sr.Root, stage)
			fmt.Printf("✅ %s v%d 已保存\n", cfg.ArtifactName, stage.CurrentVersion)

			// 冒烟测试
			artifact, err = sr.runSmokeLoop(ctx, artifact, roundNum)
			if err != nil {
				return fmt.Errorf("smoke: %w", err)
			}
			state.SaveArtifact(sr.Root, cfg.StageID, cfg.ArtifactName, artifact, stage.CurrentVersion)

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
		if _, err := state.SaveArtifact(sr.Root, cfg.StageID, cfg.ArtifactName, artifact, stage.CurrentVersion); err != nil {
			return fmt.Errorf("save artifact: %w", err)
		}
		if err := state.SaveStage(sr.Root, stage); err != nil {
			return fmt.Errorf("save stage: %w", err)
		}
	}
	if err := state.SaveFrozenArtifact(sr.Root, cfg.ArtifactName, artifact); err != nil {
		return fmt.Errorf("save frozen: %w", err)
	}
	// 执行 Stage 特定的冻结动作（如代码快照）
	if cfg.FreezeAction != nil {
		if err := cfg.FreezeAction(sr.Root, artifact); err != nil {
			return fmt.Errorf("freeze action: %w", err)
		}
	}
	stage.Status = "frozen"
	if err := state.SaveStage(sr.Root, stage); err != nil {
		return fmt.Errorf("save stage: %w", err)
	}

	meeting.Meta.Status = "passed"
	meeting.Meta.Decision = "freeze"
	if err := sr.saveMeetingWithDecision(meeting, decision); err != nil {
		return fmt.Errorf("save meeting: %w", err)
	}

	// 保存 reviewer memory
	for _, r := range meeting.Reviews {
		state.SaveMemory(sr.Root, r.AgentID, meeting.ID+"-review", &state.Memory{
			Type:      "decision",
			Title:     fmt.Sprintf("参与%s %s", cfg.StageName, meeting.ID),
			Content:   truncate(r.Content, 2000),
			Relations: []state.Relation{{Type: "based_on", TargetType: "meeting", TargetID: meeting.ID}},
			Tags:      cfg.MemoryTags,
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
			Round:           roundNum,
			ArtifactVersion: stage.CurrentVersion,
			Type:            cfg.MeetingType,
			Stage:           cfg.StageID,
			TargetArtifact:  fmt.Sprintf(".aisc/stages/%02d-%s/artifact/%s-v%d%s", stage.Order, strings.ReplaceAll(strings.ToLower(cfg.StageName), " ", "-"), cfg.ArtifactName, stage.CurrentVersion, state.ArtifactExt(cfg.ArtifactName)),
			Moderator:       "project-manager",
			Participants:    strings.Join(stage.ReviewerAgents, ", "),
			Status:          "in_progress",
			CreatedAt:       time.Now().UTC().Format(time.RFC3339),
		},
	}
}

func (sr *StageRunner) saveMeetingWithDecision(meeting *state.Meeting, decision *Decision) error {
	decisionJSON, _ := json.MarshalIndent(decision, "", "  ")
	meeting.Body = fmt.Sprintf("## Decision (%s)\n\n%s\n\n```json\n%s\n```", decision.Type, decision.Summary, string(decisionJSON))
	return state.SaveMeeting(sr.Root, meeting)
}

func (sr *StageRunner) generateArtifact(ctx context.Context, input string) (string, error) {
	prompt, err := sr.cfg.PromptDraft()
	if err != nil {
		return "", err
	}
	// 注入技术栈信息（从 Tech Design Stage 读取）。
	if ts := sr.loadTechStack(); ts != nil {
		prompt += fmt.Sprintf("\n\n# 当前项目技术栈\n- 语言: %s\n- 框架: %s\n- ORM: %s\n- 数据库: %s\n- 构建命令: %s\n- 测试命令: %s\n- 运行环境: %s",
			ts.Language, ts.Framework, ts.ORM, ts.Database,
			ts.BuildCommand, ts.TestCommand, ts.Runtime)
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
	// 注入技术栈
	if ts := sr.loadTechStack(); ts != nil {
		sysPrompt += fmt.Sprintf("\n\n# 当前项目技术栈\n- 语言: %s\n- 框架: %s\n- ORM: %s\n- 构建命令: %s\n- 测试命令: %s",
			ts.Language, ts.Framework, ts.ORM, ts.BuildCommand, ts.TestCommand)
	}
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

func (sr *StageRunner) loadTechStack() *state.TechStack {
	s, err := state.LoadStage(sr.Root, "stage-tech-design")
	if err != nil {
		return nil
	}
	return s.TechStack
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
