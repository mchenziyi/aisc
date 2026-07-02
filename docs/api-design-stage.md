# API Design Stage 设计方案

## 目标

在 Requirement Stage 冻结后，Tech Lead 起草 API Spec → 多角色评审 → 冻结。冻结后 Backend 和 Frontend 可并行开发。

## 核心改动：StageRunner 参数化

当前 `StageRunner.Run()` 硬编码了 Requirement Stage。新增 `StageConfig` 后，同一套 Run 逻辑驱动所有 Stage。

```go
// StageConfig 描述一个 Stage 的元信息
type StageConfig struct {
    StageID        string   // "stage-requirement" | "stage-api-design"
    StageName      string   // "Requirement" | "API Design"
    OwnerAgent     string   // "pm" | "tech-lead"
    ReviewerAgents []string
    ArtifactName   string   // "prd" | "api-spec"
    MeetingType    string   // "requirement_review" | "api_review"
    MemoryTags     []string // ["需求评审","PRD"] | ["API评审","API Spec"]
    MaxRounds      int      // 5
    PromptDraft    string   // prompts.Load("pm","draft") | prompts.Load("tech-lead","api-design")
    PromptRevise   string   // prompts.Load("pm","revise") | prompts.Load("tech-lead","api-design-revise")
    // 输入：如何获取上一阶段的产物
    InputReader    func(root string) (string, error)
}
```

`StageRunner.Run(ctx)` 改为 `StageRunner.Run(ctx, cfg *StageConfig)`，所有 stage-specific 的值从 cfg 取。

## 文件改动清单

```
aisc/
├── orchestration/
│   ├── stage.go        ← 改：StageRunner 参数化，新增 StageConfig
│   └── meeting.go      ← 改：参会人从 cfg.ReviewerAgents 获取
├── state/
│   └── store.go        ← 改：SaveFrozenPRD → SaveFrozenArtifact（通用化）
│                          PRDExists → ArtifactExists（通用化）
│                          SaveMeeting 去掉硬编码 "01-requirement"
├── agents/prompts/
│   └── tech-lead/
│       ├── api-design.md         ← 新：Tech Lead 起草 API Spec 的 prompt
│       └── api-design-revise.md  ← 新：Tech Lead 修订 API Spec 的 prompt
├── .aisc/
│   └── stages/
│       └── 02-api-design/
│           └── stage.json        ← 新：API Design Stage 的元信息
└── main.go              ← 改：先跑 Requirement，后跑 API Design
```

## 数据流

```
runRequirement() → PRD frozen
        ↓
runAPIDesign() → 读取 docs/prd-frozen.md 作为输入
        ↓
     Tech Lead 起草 API Spec v1
        ↓
     Review Meeting（Backend + Frontend + QA + PM 评审）
        ↓
     Moderator 裁决 → revise / freeze
        ↓
     API Spec frozen → docs/api-spec-frozen.md
```

## stage.json（02-api-design）

```json
{
  "id": "stage-api-design",
  "type": "API Design",
  "status": "drafting",
  "order": 2,
  "owner_agent": "tech-lead",
  "reviewer_agents": ["pm-agent", "backend", "frontend", "qa"],
  "artifact_id": "api-spec",
  "current_version": 0,
  "meeting_ids": [],
  "meeting_counter": 0
}
```

注意 reviewer 里没有 `ui-designer`——API 设计和 UI 无关。也没有 `tech-lead` 自己——Owner 不参与评审自己的产物。

## API Spec 格式

用 OpenAPI 3.0 YAML。LLM 输出标准 OpenAPI 文档，可被工具链（Swagger UI、代码生成）直接消费。

## main.go 改动

```go
func main() {
    // 1. Requirement Stage
    if !isFrozen("stage-requirement") {
        runRequirementStage()
    }

    // 2. API Design Stage（依赖 PRD 已冻结）
    if isFrozen("stage-requirement") && !isFrozen("stage-api-design") {
        runAPIDesignStage()
    }
}
```

## 工作量估算

| 文件 | 操作 | 行数 |
|------|------|------|
| stage.go | 重构参数化 | ~40 行改 |
| store.go | 通用化 3 个函数 | ~20 行改 |
| main.go | Stage 串联 | ~20 行改 |
| tech-lead/api-design.md | 新增 prompt | ~50 行 |
| tech-lead/api-design-revise.md | 新增 prompt | ~30 行 |
| stage.json | 新增 | ~10 行 |

总计约 **170 行**改动，无新包，无新依赖。
