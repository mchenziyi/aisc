# prompts 目录结构

每个角色一个文件夹，`.md` 文件即提示词。`_shared/` 存跨角色共享模板。

## 加载方式

```go
import "github.com/mchenziyi/aisc/agents/prompts"

// 直接加载
p, _ := prompts.Load("pm", "draft")

// 评审角色 = 共享模板 + 角色视角（自动组合）
p, _ := prompts.LoadReviewer("tech-lead", "exhaustive")
```

## 角色清单

| 文件夹 | 角色 | 现有文件 |
|--------|------|---------|
| `pm/` | 产品经理 | draft.md, revise.md |
| `ui-designer/` | UI 设计师 | role.md |
| `tech-lead/` | 技术负责人 | role.md |
| `backend/` | 后端开发 | role.md |
| `frontend/` | 前端开发 | role.md |
| `qa/` | 测试工程师 | role.md |
| `moderator/` | 调度者（Project Manager） | default.md |
| `devops/` | 运维工程师 | （待添加） |
| `documentation/` | 文档工程师 | （待添加） |
| `_shared/` | 共享模板 | review-exhaustive.md, review-verification.md |

## 添加新提示词

在对应角色文件夹下放 `.md` 文件，然后用 `prompts.Load("角色", "文件名不含扩展名")` 加载。无需改代码，无需重新编译。
