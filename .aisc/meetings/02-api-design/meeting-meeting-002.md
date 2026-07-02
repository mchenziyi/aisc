---
id: meeting-002
type: api_review
stage: stage-api-design
target_artifact: .aisc/stages/02-api-design/artifact/api-spec-v1.yaml
moderator: project-manager
participants: pm-agent, backend, frontend, qa
status: needs_revision
round: 1
prd_version: 1
created_at: 2026-07-02T09:13:32Z
decision: revise
---

## Decision (revise)

存在三个 Blocker（响应结构不一致、UploadTokenResponse.fields 未定义、视频同步校验不现实）及多个重要问题（Tags 类型不统一、缺失用户管理/Token 刷新、搜索/排序定义缺失等）。所有评审人均要求修改（Needs Revision），且无意见冲突。必须完成修正后方可进入下一阶段。

```json
{
  "type": "revise",
  "summary": "存在三个 Blocker（响应结构不一致、UploadTokenResponse.fields 未定义、视频同步校验不现实）及多个重要问题（Tags 类型不统一、缺失用户管理/Token 刷新、搜索/排序定义缺失等）。所有评审人均要求修改（Needs Revision），且无意见冲突。必须完成修正后方可进入下一阶段。",
  "action_items": [
    {
      "description": "修正 VideoListResponse 结构：将 total、page、page_size 移至 data 字段内（如 data.items, data.total, data.page, data.page_size），保持统一响应包裹。"
    },
    {
      "description": "显式定义 UploadTokenResponse.fields 的结构（如 key, policy, signature 等），提供示例或必填字段列表。"
    },
    {
      "description": "重新设计 POST /videos 的校验流程：明确是同步检测（仅读头部元数据）还是异步转码＋状态字段（processing/available/failed），并更新文档与响应。"
    },
    {
      "description": "统一 tags 字段类型：将 CreateVideoRequest 和 UpdateVideoRequest 中的 tags 改为 array of strings，增加 maxItems:5 及 items maxLength:20，并保持与 VideoDetail.tags 一致。"
    },
    {
      "description": "为 CreateVideoRequest 增加 uploadId 对应文件状态的校验（是否已上传、大小/MIME 是否匹配），并补充错误码（如 422）。"
    },
    {
      "description": "补充用户管理接口：至少包括 GET /user/profile（获取个人资料）、PUT /user/profile（更新个人资料）、PUT /user/password（修改密码）和 POST /auth/refresh（刷新 Token）。"
    },
    {
      "description": "将 updateVideo 的 HTTP 方法从 PUT 改为 PATCH，并更新请求体描述以反映部分更新。"
    },
    {
      "description": "取消 GET /videos 和 GET /videos/{id} 的强制认证：允许未登录用户访问，并在文档中说明已登录用户可额外获得收藏状态等信息（需产品确认后实施）。"
    },
    {
      "description": "为列表接口（GET /videos 和 GET /favorites）添加排序参数（sort_by, order），补充 keyword 的搜索范围（标题、描述、标签）及默认排序规则。"
    },
    {
      "description": "明确删除视频的级联行为（是否同步删除存储文件、收藏记录等），补充 /videos/{id} 中 id 的格式约束（如 uuid），统一成功消息为 'ok'，并定义独立业务错误码体系（如 10000+）。"
    }
  ],
  "conflicts": [],
  "freeze_check": {
    "all_blockers_resolved": false,
    "all_conflicts_resolved": true,
    "ready_for_next_stage": false
  }
}
```