package orchestration

import (
	"context"

	"agentdemo/agent"
)

// ─── AgentClient ──────────────────────────────────────────────

// AgentClient 封装对底层 Agent 运行时的调用。
// 当前实现基于 QiuQiuPro，将来可替换为其他运行时。
type AgentClient interface {
	Run(ctx context.Context, systemPrompt, userTask string) (string, error)
}

// QiuQiuProClient 基于 QiuQiuPro Agent 运行时的实现。
type QiuQiuProClient struct {
	APIKey string
	Model  string
}

func NewQiuQiuProClient(apiKey, model string) *QiuQiuProClient {
	return &QiuQiuProClient{APIKey: apiKey, Model: model}
}

func (c *QiuQiuProClient) Run(ctx context.Context, systemPrompt, userTask string) (string, error) {
	a, err := agent.New(c.APIKey, c.Model, false)
	if err != nil {
		return "", err
	}
	a.SetSystemPrompt(systemPrompt)
	a.MaxMessages = 5000
	return a.Run(ctx, userTask)
}
