package agents

// ReviewerRolePrompt 返回 agentID 对应的角色视角提示
func ReviewerRolePrompt(agentID string) string {
	switch agentID {
	case "tech-lead":
		return RoleTechLead
	case "ui-designer":
		return RoleUIDesigner
	case "backend":
		return RoleBackend
	case "frontend":
		return RoleFrontend
	case "qa":
		return RoleQA
	default:
		return ""
	}
}
