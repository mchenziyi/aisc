package orchestration

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/mchenziyi/aisc/agents/prompts"
)

// ─── Types ────────────────────────────────────────────────────

// Decision Moderator 的结构化决策
type Decision struct {
	Type        string       `json:"type"`
	Summary     string       `json:"summary"`
	ActionItems []ActionItem `json:"action_items"`
	Conflicts   []Conflict   `json:"conflicts"`
	FreezeCheck *FreezeCheck `json:"freeze_check,omitempty"`
	Raw         string       `json:"raw,omitempty"` // fallback 原始文本
}

// ActionItem 一个修改任务
type ActionItem struct {
	Description string `json:"description"`
}

// Conflict 评审冲突记录
type Conflict struct {
	Topic          string   `json:"topic"`
	Sides          []string `json:"sides"`
	Resolution     string   `json:"resolution"`
	EscalateToUser bool     `json:"escalate_to_user"`
}

// FreezeCheck 冻结条件检查
type FreezeCheck struct {
	AllBlockersResolved  bool `json:"all_blockers_resolved"`
	AllConflictsResolved bool `json:"all_conflicts_resolved"`
	ReadyForNextStage    bool `json:"ready_for_next_stage"`
}

// Review 一个评审人的意见
type Review struct {
	AgentID string `json:"agent_id"`
	Content string `json:"content"`
}

// ─── Orchestrator ─────────────────────────────────────────────

// Orchestrator 评审编排器
type Orchestrator struct {
	Client AgentClient
}

// New 使用默认 AgentClient（QiuQiuPro）创建 Orchestrator
func New(apiKey, model string) *Orchestrator {
	return &Orchestrator{Client: NewQiuQiuProClient(apiKey, model)}
}

// NewWithClient 使用指定 AgentClient 创建 Orchestrator
func NewWithClient(c AgentClient) *Orchestrator {
	return &Orchestrator{Client: c}
}

// RunReviewRound 执行一轮完整评审：并行审阅 → 汇总裁决
// roundNum=1 时用全量深度评审，roundNum>=2 时用定向复核。
func (o *Orchestrator) RunReviewRound(
	ctx context.Context,
	artifact string,
	roundNum int,
	prevDecision *Decision,
	reviewers []string,
	artifactLabel string,
	reviewPromptDir string,
) (*Decision, []Review, error) {

	// Step 1: 并行评审
	reviews, err := o.parallelReview(ctx, artifact, roundNum, prevDecision, reviewers, artifactLabel, reviewPromptDir)
	if err != nil {
		return nil, nil, fmt.Errorf("parallel review: %w", err)
	}

	// Step 2: 汇总裁决
	decision, err := o.consensus(ctx, artifact, reviews, prevDecision, artifactLabel)
	if err != nil {
		return nil, reviews, fmt.Errorf("consensus: %w", err)
	}

	return decision, reviews, nil
}

// ─── Parallel Review ──────────────────────────────────────────

func (o *Orchestrator) parallelReview(
	ctx context.Context,
	artifact string,
	roundNum int,
	prevDecision *Decision,
	reviewers []string,
	artifactLabel string,
	reviewPromptDir string,
) ([]Review, error) {

	reviews := make([]Review, len(reviewers))
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	for i, agentID := range reviewers {
		wg.Add(1)
		go func(idx int, id string) {
			defer wg.Done()

			review, err := o.reviewOnce(ctx, id, artifact, roundNum, prevDecision, artifactLabel, reviewPromptDir)
			if err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = fmt.Errorf("%s: %w", id, err)
				}
				mu.Unlock()
				return
			}
			reviews[idx] = Review{AgentID: id, Content: review}
		}(i, agentID)
	}

	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	return reviews, nil
}

func (o *Orchestrator) reviewOnce(
	ctx context.Context,
	agentID, artifact string,
	roundNum int,
	prevDecision *Decision,
	artifactLabel string,
	reviewPromptDir string,
) (string, error) {

	var sysPrompt string
	var err error

	if roundNum == 1 {
		sysPrompt, err = prompts.LoadReviewer(reviewPromptDir, agentID, "exhaustive")
	} else {
		sysPrompt, err = prompts.LoadReviewer(reviewPromptDir, agentID, "verification")
	}
	if err != nil {
		return "", err
	}

	// 构造 user prompt
	var task strings.Builder
	task.WriteString(fmt.Sprintf("请评审以下 %s：\n\n", artifactLabel))
	task.WriteString(artifact)

	if roundNum >= 2 && prevDecision != nil {
		prevJSON, _ := json.MarshalIndent(prevDecision, "", "  ")
		task.WriteString("\n\n## 上一轮 Decision（Scope Lock — 不能推翻的裁决）\n")
		task.WriteString(string(prevJSON))
		task.WriteString("\n\n请逐条验证 Action Items，只对遗漏的阻断级问题提新 blocker。")
	}

	return o.Client.Run(ctx, sysPrompt, task.String())
}

// ─── Consensus ────────────────────────────────────────────────

func (o *Orchestrator) consensus(
	ctx context.Context,
	artifact string,
	reviews []Review,
	prevDecision *Decision,
	artifactLabel string,
) (*Decision, error) {

	modPrompt, err := prompts.Load("moderator", "default")
	if err != nil {
		return nil, err
	}

	// 构造评审意见摘要
	var reviewText strings.Builder
	for _, r := range reviews {
		reviewText.WriteString(fmt.Sprintf("## %s 的评审意见\n%s\n\n---\n\n", r.AgentID, r.Content))
	}

	var task strings.Builder
	task.WriteString(fmt.Sprintf("请根据以下信息做出决策（只输出 JSON）：\n\n## %s 内容\n", artifactLabel))
	task.WriteString(artifact)
	task.WriteString("\n\n## 评审意见\n")
	task.WriteString(reviewText.String())

	// 注入上一轮 Decision 作为 Scope Lock
	if prevDecision != nil {
		prevJSON, _ := json.MarshalIndent(map[string]any{
			"type":         prevDecision.Type,
			"summary":      prevDecision.Summary,
			"action_items": prevDecision.ActionItems,
		}, "", "  ")
		task.WriteString("\n\n## 上一轮已裁决（Scope Lock — 不能推翻的结论）\n")
		task.WriteString(string(prevJSON))
		task.WriteString("\n\n注意：上述裁决结果在本轮不能被推翻。只能基于上一轮 action items 的完成情况做决策。")
	}

	raw, err := o.Client.Run(ctx, modPrompt, task.String())
	if err != nil {
		return nil, err
	}

	return ParseDecision(raw), nil
}

// ─── Decision Parser ──────────────────────────────────────────

// ParseDecision 从 LLM 输出中提取 JSON 决策。
// 兼容裸 JSON、markdown 代码块、嵌在文本中的 JSON。
func ParseDecision(raw string) *Decision {
	text := strings.TrimSpace(raw)
	candidates := []string{text}

	// 尝试提取 ```json ... ``` 代码块
	re := regexp.MustCompile("```(?:json)?\\s*\\n?(.*?)\\n?```")
	if m := re.FindStringSubmatch(text); m != nil {
		candidates = append([]string{strings.TrimSpace(m[1])}, candidates...)
	}

	// 尝试提取第一个 { 到最后一个 }
	re2 := regexp.MustCompile(`\{.*\}`)
	if m := re2.FindString(text); m != "" {
		candidates = append([]string{m}, candidates...)
	}

	for _, c := range candidates {
		var d Decision
		if err := json.Unmarshal([]byte(c), &d); err == nil && d.Type != "" {
			d.Type = strings.TrimSpace(strings.ToLower(d.Type))
			return &d
		}
	}

	// Fallback: 正则提取 type:
	fmt.Println("⚠️  JSON 解析失败，使用 fallback 正则提取 decision type")
	dtype := "unknown"
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "type:") || strings.HasPrefix(line, `"type"`) {
			dtype = strings.Trim(strings.SplitN(line, ":", 2)[1], ` ",`)
			dtype = strings.ToLower(dtype)
			break
		}
	}
	return &Decision{Type: dtype, Summary: text[:min(200, len(text))], Raw: text}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ─── Helpers ──────────────────────────────────────────────────

// ActionItemsText 将 action_items 格式化为序号列表
func ActionItemsText(items []ActionItem) string {
	var b strings.Builder
	for i, item := range items {
		b.WriteString(fmt.Sprintf("%d. %s\n", i+1, item.Description))
	}
	return b.String()
}
