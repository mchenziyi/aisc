---
id: meeting-002
type: api_review
stage: stage-api-design
target_artifact: .aisc/stages/02-api-design/artifact/api-spec-v2.yaml
moderator: project-manager
participants: pm-agent, backend, frontend, qa
status: needs_revision
round: 2
prd_version: 2
created_at: 2026-07-02T09:59:05Z
decision: revise
---

## Decision (revise)

上一轮 action items 已全部解决，但 qa 发现 PATCH 幂等性规则在混合更新场景存在阻断级歧义，需要澄清后才可冻结。同时还有多个重要问题需要修正。决定继续修订后重审。

```json
{
  "type": "revise",
  "summary": "上一轮 action items 已全部解决，但 qa 发现 PATCH 幂等性规则在混合更新场景存在阻断级歧义，需要澄清后才可冻结。同时还有多个重要问题需要修正。决定继续修订后重审。",
  "action_items": [
    {
      "description": "澄清 PATCH 端点的幂等性规则：明确仅当请求体中除 version 外只包含 completed:true 且当前已完成且版本号匹配时，才跳过数据修改（version 不变）；若同时包含其他字段，则正常更新所有字段并递增 version 号。"
    },
    {
      "description": "补充 completed 设置为 false 时的对称幂等性规则：若请求设置 completed:false 且当前 completed=false 且版本号匹配，则直接返回成功不修改数据。"
    },
    {
      "description": "在 UpdateTodoRequest 的描述或 schema 中明确“至少提供一个要更新的字段（除 version 外）”，若仅传 version 应返回 400，并在 400 响应示例中体现。"
    },
    {
      "description": "在 UpdateTodoRequest 的 description 属性中添加 nullable:true，以允许通过 null 清除描述，与 due_date 的处理保持一致。"
    }
  ],
  "conflicts": [
    {
      "topic": "PRD 是否达到冻结标准",
      "sides": [
        "pm-agent、backend、frontend 认为 action items 全部解决，无新增 blocker，同意冻结",
        "qa 认为 PATCH 幂等性规则在混合更新场景存在严重歧义，属于 blocker，需要修订"
      ],
      "resolution": "采纳 qa 意见，该歧义可能导致实现时数据丢失或不一致，必须在冻结前澄清。决定继续修订。",
      "escalate_to_user": false
    }
  ],
  "freeze_check": {
    "all_blockers_resolved": false,
    "all_conflicts_resolved": true,
    "ready_for_next_stage": false
  }
}
```