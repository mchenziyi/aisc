package prompts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// dir 提示词根目录，由外部在启动时设置。
var dir = "agents/prompts"

// SetDir 设置提示词文件的根目录路径。
func SetDir(d string) { dir = d }

// Load 读取指定角色和场景的提示词文件。
// 例如 Load("pm", "draft") → pm/draft.md
func Load(role, scene string) (string, error) {
	path := filepath.Join(dir, role, scene+".md")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("load prompt %s/%s: %w", role, scene, err)
	}
	return string(data), nil
}

// LoadReviewer 为评审角色加载组合后的 system prompt。
// 读取 _shared/review-{scene}.md 模板，注入 {role} 的角色视角。
// scene: "exhaustive"（全量评审）或 "verification"（定向复核）
func LoadReviewer(agentID, scene string) (string, error) {
	template, err := Load("_shared", "review-"+scene)
	if err != nil {
		return "", err
	}
	roleHint, err := Load(agentID, "role")
	if err != nil {
		return "", err
	}

	result := strings.ReplaceAll(template, "{role_hint}", roleHint)
	result = strings.ReplaceAll(result, "{role}", agentID)
	return result, nil
}

// ReviewerRoles 返回所有评审角色 ID 列表。
func ReviewerRoles() []string {
	return []string{"tech-lead", "ui-designer", "backend", "frontend", "qa"}
}
