package orchestration

import (
	"context"

	"agentdemo/agent"
	"agentdemo/tool"
)

// ─── AgentClient ──────────────────────────────────────────────

// AgentClient 封装对底层 Agent 运行时的调用。
type AgentClient interface {
	Run(ctx context.Context, systemPrompt, userTask string) (string, error)
}

// AgentClientWithTools 扩展：支持工具调用的 Agent。
// Draft/Revise 阶段需要写文件、执行命令时使用此接口。
type AgentClientWithTools interface {
	AgentClient
	RunWithTools(ctx context.Context, systemPrompt, userTask string, tools []tool.Tool) (string, error)
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

func (c *QiuQiuProClient) RunWithTools(ctx context.Context, systemPrompt, userTask string, tools []tool.Tool) (string, error) {
	a, err := agent.New(c.APIKey, c.Model, false)
	if err != nil {
		return "", err
	}
	a.SetSystemPrompt(systemPrompt)
	a.MaxMessages = 5000
	a.RegisterTools(tools)
	a.SetGate(agent.AllowAllGate{})
	return a.Run(ctx, userTask)
}
