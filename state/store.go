package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ─── 目录常量 ────────────────────────────────────────────────

const (
	DirAISC     = ".aisc"
	DirStages   = ".aisc/stages"
	DirMeetings = ".aisc/meetings"
	DirMemory   = ".aisc/memory"
	DirDocs     = "docs"
)

// ─── 数据模型 ────────────────────────────────────────────────

// Project 项目顶层状态 (project.json)
type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"` // active, paused, completed
}

// Stage 单个 Stage 的状态 (stages/01-requirement/stage.json)
type Stage struct {
	ID             string   `json:"id"`
	Type           string   `json:"type"`
	Status         string   `json:"status"` // drafting, in_review, revising, frozen
	Order          int      `json:"order"`
	OwnerAgent     string   `json:"owner_agent"`
	ReviewerAgents []string `json:"reviewer_agents"`
	ArtifactID     string   `json:"artifact_id"`
	CurrentVersion int      `json:"current_version"`
	MeetingIDs     []string `json:"meeting_ids"`
	MeetingCounter int      `json:"meeting_counter,omitempty"`
}

// Meeting 会议记录 (meetings/01-requirement/meeting-{id}.md)
type Meeting struct {
	ID      string      `json:"id"`
	Meta    MeetingMeta `json:"meta"`
	Body    string      `json:"body"`
	Reviews []Review    `json:"reviews"`
}

type MeetingMeta struct {
	ID             string `json:"id"`
	Round          int    `json:"round"`
	PRDVersion     int    `json:"prd_version"`
	Type           string `json:"type"`
	Stage          string `json:"stage"`
	TargetArtifact string `json:"target_artifact"`
	Moderator      string `json:"moderator"`
	Participants   string `json:"participants"`
	Status         string `json:"status"`
	Decision       string `json:"decision,omitempty"`
	CreatedAt      string `json:"created_at"`
}

type Review struct {
	AgentID string `json:"agent_id"`
	Content string `json:"content"`
}

// Memory Agent 长期记忆
type Memory struct {
	Type        string     `json:"type"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Relations   []Relation `json:"relations"`
	Tags        []string   `json:"tags"`
}

type Relation struct {
	Type       string `json:"type"`
	TargetType string `json:"target_type"`
	TargetID   string `json:"target_id"`
}

// ─── Project ─────────────────────────────────────────────────

func LoadProject(root string) (*Project, error) {
	var p Project
	if err := readJSON(filepath.Join(root, DirAISC, "project.json"), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func SaveProject(root string, p *Project) error {
	os.MkdirAll(filepath.Join(root, DirAISC), 0755)
	return writeJSON(filepath.Join(root, DirAISC, "project.json"), p)
}

// ─── Stage ───────────────────────────────────────────────────

// stageDir 根据 stage 的 order 推断目录名（MVP 只有 01-requirement）
func stageDir(root string, stage *Stage) string {
	return filepath.Join(root, DirStages, fmt.Sprintf("%02d-%s", stage.Order, stage.ID))
}

func LoadStage(root, stageID string) (*Stage, error) {
	// 遍历 stages/ 下所有子目录找匹配的 stage.json
	entries, err := os.ReadDir(filepath.Join(root, DirStages))
	if err != nil {
		return nil, fmt.Errorf("read stages dir: %w", err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		var s Stage
		path := filepath.Join(root, DirStages, e.Name(), "stage.json")
		if err := readJSON(path, &s); err != nil {
			continue
		}
		if s.ID == stageID {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("stage %q not found", stageID)
}

func SaveStage(root string, stage *Stage) error {
	dir := stageDir(root, stage)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return writeJSON(filepath.Join(dir, "stage.json"), stage)
}

// ─── Artifact ────────────────────────────────────────────────

func SaveArtifact(root, stageID string, filename string, content string, version int) (string, error) {
	dir := filepath.Join(root, DirStages, fmt.Sprintf("01-%s", stageID), "artifact")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, fmt.Sprintf("%s-v%d.md", filename, version))
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", err
	}
	// 同时写工作副本到 docs/
	docsDir := filepath.Join(root, DirDocs)
	os.MkdirAll(docsDir, 0755)
	os.WriteFile(filepath.Join(docsDir, filename+".md"), []byte(content), 0644)
	return path, nil
}

func ReadArtifact(root, filename string) (string, error) {
	data, err := os.ReadFile(filepath.Join(root, DirDocs, filename+".md"))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func SaveFrozenPRD(root, content string) error {
	docsDir := filepath.Join(root, DirDocs)
	os.MkdirAll(docsDir, 0755)
	return os.WriteFile(filepath.Join(docsDir, "prd-frozen.md"), []byte(content), 0644)
}

// ─── Meeting ─────────────────────────────────────────────────

func SaveMeeting(root string, meeting *Meeting) error {
	dir := filepath.Join(root, DirMeetings, "01-requirement")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	frontmatter := fmt.Sprintf(`id: %s
type: %s
stage: %s
target_artifact: %s
moderator: %s
participants: %s
status: %s
round: %d
prd_version: %d
created_at: %s
`, meeting.Meta.ID, meeting.Meta.Type, meeting.Meta.Stage,
		meeting.Meta.TargetArtifact, meeting.Meta.Moderator,
		meeting.Meta.Participants, meeting.Meta.Status,
		meeting.Meta.Round, meeting.Meta.PRDVersion, meeting.Meta.CreatedAt)

	if meeting.Meta.Decision != "" {
		frontmatter += fmt.Sprintf("decision: %s\n", meeting.Meta.Decision)
	}

	content := fmt.Sprintf("---\n%s---\n\n%s", frontmatter, meeting.Body)
	path := filepath.Join(dir, fmt.Sprintf("meeting-%s.md", meeting.ID))
	return os.WriteFile(path, []byte(content), 0644)
}

// ─── Memory ──────────────────────────────────────────────────

func SaveMemory(root, agentID, memoryID string, mem *Memory) error {
	dir := filepath.Join(root, DirMemory, agentID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return writeJSON(filepath.Join(dir, memoryID+".json"), mem)
}

// ─── Decision Memory (Scope Lock) ────────────────────────────

func SaveDecisionMemory(root string, decision any) error {
	return writeJSON(filepath.Join(root, DirAISC, "decision_memory.json"), decision)
}

func LoadDecisionMemory(root string) (map[string]any, error) {
	var d map[string]any
	if err := readJSON(filepath.Join(root, DirAISC, "decision_memory.json"), &d); err != nil {
		return nil, err
	}
	return d, nil
}

func DeleteDecisionMemory(root string) {
	os.Remove(filepath.Join(root, DirAISC, "decision_memory.json"))
}

// ─── helpers ─────────────────────────────────────────────────

func readJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
