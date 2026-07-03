package logger

import "agentdemo/agent"

// AgentSink 将 QiuQiuPro Agent 的事件转发为日志。
// 实现 agent.Sink 接口。
type AgentSink struct {
	Logger *Logger
}

func (s *AgentSink) Emit(ev agent.Event) {
	if s.Logger == nil {
		return
	}
	switch ev.Kind {
	case agent.EventToolCall:
		s.Logger.Debug("tool_call", map[string]any{
			"tool": ev.Name,
			"args": ev.Text,
			"id":   ev.ID,
		})
	case agent.EventToolResult:
		s.Logger.Debug("tool_result", map[string]any{
			"tool":   ev.Name,
			"result": ev.Text[:min(len(ev.Text), 200)],
			"id":     ev.ID,
		})
	case agent.EventNotice:
		// 仅记录非 verbose 通知（如错误）
		if !ev.Verbose {
			s.Logger.Debug("agent_notice", map[string]any{"text": ev.Text})
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
