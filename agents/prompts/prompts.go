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
// 从 {stage}/review-{scene}.md 加载模板，注入 {role} 的角色视角。
// scene: "exhaustive"（全量评审）或 "verification"（定向复核）
func LoadReviewer(stage, agentID, scene string) (string, error) {
	template, err := Load(stage, "review-"+scene)
	if err != nil {
		// 回退到 _shared
		template, err = Load("_shared", "review-"+scene)
		if err != nil {
			return "", err
		}
	}
	roleHint, err := loadRoleHint(agentID)
	if err != nil {
		return "", err
	}

	result := strings.ReplaceAll(template, "{role_hint}", roleHint)
	result = strings.ReplaceAll(result, "{role}", agentID)
	return result, nil
}

// loadRoleHint 尝试多种 ID 变体加载 role.md
func loadRoleHint(agentID string) (string, error) {
	// 按优先级尝试
	candidates := []string{agentID}
	if strings.HasSuffix(agentID, "-agent") {
		candidates = append(candidates, strings.TrimSuffix(agentID, "-agent"))
	}
	for _, c := range candidates {
		if role, err := Load(c, "role"); err == nil {
			return role, nil
		}
	}
	return "", fmt.Errorf("role hint not found for %s (tried: %v)", agentID, candidates)
}

// ReviewerRoles 返回所有评审角色 ID 列表。
func ReviewerRoles() []string {
	return []string{"tech-lead", "ui-designer", "backend", "frontend", "qa"}
}
