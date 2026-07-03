package state

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// TechStack 从 Tech Design 中提取的结构化技术选型。
// 下游 Stage（Backend/Frontend/QA/DevOps）据此决定具体行为。
type TechStack struct {
	Language      string `json:"language"`       // "Go", "Python", "TypeScript" 等
	Framework     string `json:"framework"`      // "Gin", "FastAPI", "React+Next.js" 等
	ORM           string `json:"orm"`            // "pgx", "SQLAlchemy", "Prisma" 等
	Database      string `json:"database"`       // "PostgreSQL", "MySQL" 等
	BuildCommand  string `json:"build_command"`  // "go build ./...", "npm run build" 等
	TestCommand   string `json:"test_command"`   // "go test ./...", "npm test" 等
	LintCommand   string `json:"lint_command"`   // "go vet ./...", "ruff check ." 等
	Runtime       string `json:"runtime"`        // Docker 基础镜像，如 "golang:1.22-alpine"
}

// ExtractPrompt 返回让 LLM 从 Tech Design 中提取技术栈的 system prompt。
func ExtractPrompt() string {
	return `你是一个技术文档解析器。请从以下 Tech Design 文档中提取技术栈信息。

输出一个 JSON 对象（不要输出其他内容，不要用 markdown 代码块包裹）：
{
  "language": "Go",
  "framework": "Gin",
  "orm": "pgx",
  "database": "PostgreSQL",
  "build_command": "go build ./...",
  "test_command": "go test ./...",
  "lint_command": "go vet ./...",
  "runtime": "golang:1.22-alpine"
}

规则：
- 如果文档明确写了，用实际值
- build_command/test_command/lint_command 如果文档没有明确写，根据 language 推断合理默认值
- runtime 根据 language 推断 Docker 基础镜像`
}

// ExtractTechStack 读取冻结的 Tech Design 并通过 LLM 提取技术栈。
// client 用于调用 LLM，extracted 写入 stage.json。
func ExtractTechStack(ctx context.Context, root string, runLLM func(systemPrompt, userPrompt string) (string, error)) error {
	data, err := os.ReadFile(filepath.Join(root, DirDocs, "tech-design-frozen.md"))
	if err != nil {
		return fmt.Errorf("read tech-design: %w", err)
	}
	raw, err := runLLM(ExtractPrompt(), string(data))
	if err != nil {
		return fmt.Errorf("llm extract: %w", err)
	}
	var ts TechStack
	if err := json.Unmarshal([]byte(raw), &ts); err != nil {
		// 尝试清理 markdown fence
		raw = stripMarkdownFence(raw)
		if err := json.Unmarshal([]byte(raw), &ts); err != nil {
			return fmt.Errorf("parse tech stack JSON: %w (raw: %.200s)", err, raw)
		}
	}
	return SaveTechStack(root, "stage-tech-design", &ts)
}

// SaveTechStack 将技术栈写入 stage.json 的 tech_stack 字段。
func SaveTechStack(root, stageID string, ts *TechStack) error {
	stage, err := LoadStage(root, stageID)
	if err != nil {
		return err
	}
	stage.TechStack = ts
	return SaveStage(root, stage)
}
