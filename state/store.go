package state

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	ID              string `json:"id"`
	Round           int    `json:"round"`
	ArtifactVersion int    `json:"artifact_version"`
	Type            string `json:"type"`
	Stage           string `json:"stage"`
	TargetArtifact  string `json:"target_artifact"`
	Moderator       string `json:"moderator"`
	Participants    string `json:"participants"`
	Status          string `json:"status"`
	Decision        string `json:"decision,omitempty"`
	CreatedAt       string `json:"created_at"`
}

type Review struct {
	AgentID string `json:"agent_id"`
	Content string `json:"content"`
}

// Memory Agent 长期记忆
type Memory struct {
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Relations []Relation `json:"relations"`
	Tags      []string   `json:"tags"`
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

// stageDir 根据 stage 的 order 和 type 推断目录名（如 01-requirement）
// stageDirName 从 Stage 的 Type 和 Order 推导目录名
func stageDirName(stage *Stage) string {
	dirName := strings.ToLower(stage.Type)
	dirName = strings.ReplaceAll(dirName, " ", "-")
	return fmt.Sprintf("%02d-%s", stage.Order, dirName)
}

func stageDir(root string, stage *Stage) string {
	return filepath.Join(root, DirStages, stageDirName(stage))
}

func LoadStage(root, stageID string) (*Stage, error) {
	// 遍历 stages/ 下所有子目录找匹配的 stage.json
	entries, err := os.ReadDir(filepath.Join(root, DirStages))
	if err != nil {
		return nil, fmt.Errorf("read stages dir: %w", err)
	}
	var found *Stage
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
			if found != nil {
				return nil, fmt.Errorf("stage %q 存在重复目录: %s 和 %s", stageID, stageDirName(found), e.Name())
			}
			found = &s
		}
	}
	if found == nil {
		return nil, fmt.Errorf("stage %q not found", stageID)
	}
	return found, nil
}

func SaveStage(root string, stage *Stage) error {
	dir := stageDir(root, stage)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return writeJSON(filepath.Join(dir, "stage.json"), stage)
}

// ─── Artifact ────────────────────────────────────────────────

// ArtifactExt 根据 artifact 名称推断文件扩展名（导出供 orchestration 使用）
func ArtifactExt(name string) string {
	return artifactExt(name)
}

// artifactExt 内部实现
func artifactExt(name string) string {
	switch strings.ToLower(name) {
	case "api-spec", "openapi", "api":
		return ".yaml"
	default:
		return ".md"
	}
}

// stripMarkdownFence 去掉 LLM 输出中常见的 markdown 代码块包裹。
// 例如 ```yaml\n...\n``` → 纯内容。
func stripMarkdownFence(content string) string {
	text := strings.TrimSpace(content)
	// 尝试匹配 ```lang\n...\n```
	if strings.HasPrefix(text, "```") && strings.HasSuffix(text, "```") {
		// 找到第一个换行和最后一个 ``` 之间的内容
		nl := strings.Index(text, "\n")
		lastBacktick := strings.LastIndex(text, "```")
		if nl != -1 && lastBacktick > nl {
			return strings.TrimSpace(text[nl+1 : lastBacktick])
		}
	}
	return content
}

func SaveArtifact(root, stageID string, filename string, content string, version int) (string, error) {
	content = stripMarkdownFence(content)
	// 加载 stage 确定 order
	stage, err := LoadStage(root, stageID)
	if err != nil {
		return "", err
	}
	dirName := stageDirName(stage)
	dir := filepath.Join(root, DirStages, dirName, "artifact")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	ext := artifactExt(filename)
	path := filepath.Join(dir, fmt.Sprintf("%s-v%d%s", filename, version, ext))
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", err
	}
	docsDir := filepath.Join(root, DirDocs)
	os.MkdirAll(docsDir, 0755)
	os.WriteFile(filepath.Join(docsDir, filename+ext), []byte(content), 0644)
	return path, nil
}

func ReadArtifact(root, filename string) (string, error) {
	// 尝试 .md 和 .yaml
	for _, ext := range []string{".md", ".yaml"} {
		data, err := os.ReadFile(filepath.Join(root, DirDocs, filename+ext))
		if err == nil {
			return string(data), nil
		}
	}
	return "", fmt.Errorf("artifact %q not found", filename)
}

func ArtifactExists(root, artifactName string) bool {
	// 尝试多种扩展名
	for _, ext := range []string{".md", ".yaml"} {
		if _, err := os.Stat(filepath.Join(root, DirDocs, artifactName+ext)); err == nil {
			return true
		}
	}
	return false
}

// PRDExists 向后兼容别名
func PRDExists(root string) bool {
	return ArtifactExists(root, "prd")
}

func ReadRequirement(root string) (string, error) {
	return ReadArtifact(root, "requirement")
}

// SaveFrozenArtifact 通用：保存冻结产物
func SaveFrozenArtifact(root, artifactName, content string) error {
	content = stripMarkdownFence(content)
	docsDir := filepath.Join(root, DirDocs)
	os.MkdirAll(docsDir, 0755)
	// 根据内容或 artifactName 推断扩展名
	ext := ".md"
	if artifactName == "api-spec" {
		ext = ".yaml"
	}
	return os.WriteFile(filepath.Join(docsDir, artifactName+"-frozen"+ext), []byte(content), 0644)
}

// SaveFrozenPRD 向后兼容别名
func SaveFrozenPRD(root, content string) error {
	return SaveFrozenArtifact(root, "prd", content)
}

// ReadFrozenPRD 读取冻结的 PRD 作为 API Design Stage 的输入
func ReadFrozenPRD(root string) (string, error) {
	data, err := os.ReadFile(filepath.Join(root, DirDocs, "prd-frozen.md"))
	if err != nil {
		return "", fmt.Errorf("请先完成 Requirement Stage: %w", err)
	}
	return string(data), nil
}

// ReadFrozenDesignDocs 读取冻结的 PRD + API Spec + Tech Design
func ReadFrozenDesignDocs(root string) (string, error) {
	prd, err := os.ReadFile(filepath.Join(root, DirDocs, "prd-frozen.md"))
	if err != nil {
		return "", fmt.Errorf("请先完成 Requirement Stage: %w", err)
	}
	api, err := os.ReadFile(filepath.Join(root, DirDocs, "api-spec-frozen.yaml"))
	if err != nil {
		return "", fmt.Errorf("请先完成 API Design Stage: %w", err)
	}
	tech, err := os.ReadFile(filepath.Join(root, DirDocs, "tech-design-frozen.md"))
	if err != nil {
		return "", fmt.Errorf("请先完成 Tech Design Stage: %w", err)
	}
	return fmt.Sprintf("## 冻结的 PRD\n\n%s\n\n## 冻结的 API Spec\n\n%s\n\n## 冻结的 Tech Design\n\n%s",
		string(prd), string(api), string(tech)), nil
}

// ReadCodeDir 递归读取目录下所有文件内容，用于代码评审。
func ReadCodeDir(root, dirName string) (string, error) {
	dir := filepath.Join(root, dirName)
	var result strings.Builder
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil // skip unreadable files
		}
		rel, _ := filepath.Rel(root, path)
		result.WriteString(fmt.Sprintf("// ─── %s ──────────────────────────────\n", rel))
		result.Write(data)
		result.WriteString("\n\n")
		return nil
	})
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

// ArchiveDir 将目录打包为 tar.gz 并保存到 .aisc/frozen/{name}.tar.gz。
// 用于 Backend/Frontend Stage 的代码快照。
func ArchiveDir(root, dirName, archiveName string) error {
	dst := filepath.Join(root, DirAISC, "frozen", archiveName+".tar.gz")
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	cmd := fmt.Sprintf("cd %s && tar czf %s %s", root, dst, dirName)
	// Use shell to execute
	return runShell(cmd)
}

// runShell 执行 shell 命令。
func runShell(cmd string) error {
	return exec.Command("sh", "-c", cmd).Run()
}

// BackendSmokeTest 运行 go build + go vet 作为冒烟测试
func BackendSmokeTest(root string) error {
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = filepath.Join(root, "backend")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go build: %w\n%s", err, string(out))
	}
	cmd = exec.Command("go", "vet", "./...")
	cmd.Dir = filepath.Join(root, "backend")
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go vet: %w\n%s", err, string(out))
	}
	return nil
}

// ─── Meeting ─────────────────────────────────────────────────

func SaveMeeting(root string, meeting *Meeting) error {
	// 从 stage metadata 推导目录名
	stage, err := LoadStage(root, meeting.Meta.Stage)
	if err != nil {
		return fmt.Errorf("load stage for meeting: %w", err)
	}
	dir := filepath.Join(root, DirMeetings, stageDirName(stage))
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
artifact_version: %d
created_at: %s
`, meeting.Meta.ID, meeting.Meta.Type, meeting.Meta.Stage,
		meeting.Meta.TargetArtifact, meeting.Meta.Moderator,
		meeting.Meta.Participants, meeting.Meta.Status,
		meeting.Meta.Round, meeting.Meta.ArtifactVersion, meeting.Meta.CreatedAt)

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
