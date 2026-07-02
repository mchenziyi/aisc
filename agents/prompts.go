package agents

import (
	"fmt"
	"os"
	"path/filepath"
)

// promptsDir 提示词根目录，由外部在启动时设置。
// 默认相对于工作目录的 agents/prompts/。
var promptsDir = "agents/prompts"

// SetPromptsDir 设置提示词文件的根目录路径。
func SetPromptsDir(dir string) {
	promptsDir = dir
}

// Load 读取指定角色和场景的提示词文件。
// 例如 Load("pm", "draft") → agents/prompts/pm/draft.md
func Load(role, scene string) (string, error) {
	path := filepath.Join(promptsDir, role, scene+".md")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("load prompt %s/%s: %w", role, scene, err)
	}
	return string(data), nil
}

// ReviewerRoleHint 返回 reviewer agentID 对应的角色视角提示词文件名。
// 例如 ReviewerRoleHint("tech-lead") → "tech-lead"
func ReviewerRoleHint(agentID string) string {
	switch agentID {
	case "tech-lead", "ui-designer", "backend", "frontend", "qa":
		return agentID
	default:
		return ""
	}
}
