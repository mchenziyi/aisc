---
id: meeting-001
type: tech_review
stage: stage-tech-design
target_artifact: .aisc/stages/03-tech-design/artifact/tech-design-v1.md
moderator: project-manager
participants: pm-agent, backend, frontend, qa
status: needs_revision
round: 1
prd_version: 1
created_at: 2026-07-02T12:14:58Z
decision: revise
---

## Decision (revise)

技术设计整体方向正确，但存在多个关键阻塞问题，包括版本号传递机制缺失、分页元数据不完整、用户名小写未强制、迁移策略自相矛盾、注册并发冲突未处理等。需要修订后重审。

```json
{
  "type": "revise",
  "summary": "技术设计整体方向正确，但存在多个关键阻塞问题，包括版本号传递机制缺失、分页元数据不完整、用户名小写未强制、迁移策略自相矛盾、注册并发冲突未处理等。需要修订后重审。",
  "action_items": [
    {
      "description": "在所有待办API响应中强制包含version字段，并在请求体中接收version以支持乐观锁。"
    },
    {
      "description": "在分页响应中添加total、page、page_size等元数据。"
    },
    {
      "description": "统一日期时间格式为RFC3339。"
    },
    {
      "description": "明确更新操作使用PATCH，并定义可修改字段及其可选性。"
    },
    {
      "description": "统一认证Token响应字段名为\"token\"。"
    },
    {
      "description": "添加CORS中间件，配置允许的域名。"
    },
    {
      "description": "定义标准化错误码枚举（如USERNAME_TAKEN、VALIDATION_ERROR、VERSION_CONFLICT等），并在ErrorResponse中增加error_code字段。"
    },
    {
      "description": "在数据库层使用表达式唯一索引（LOWER(username)）强制用户名小写，并确保应用层存储前转换小写。"
    },
    {
      "description": "处理注册并发唯一约束冲突，捕获错误并返回409及标准化错误码。"
    },
    {
      "description": "统一迁移策略：生产环境通过CI/CD执行迁移，开发环境可自动迁移；文档澄清并移除启动自动迁移的歧义。"
    },
    {
      "description": "完善乐观锁实现：使用原子UPDATE/DELETE SQL，通过RowsAffected判断，若为0再查询记录以区分404和409。"
    },
    {
      "description": "在配置中明确JWT_SECRET和JWT_EXPIRATION环境变量，提供默认值或示例。"
    },
    {
      "description": "调整用户名字段长度为50（或符合PRD要求），并在DDL中注明来源。"
    },
    {
      "description": "添加健康检查端点 /health。"
    },
    {
      "description": "在日志中间件中加入request_id，便于链路追踪。"
    },
    {
      "description": "补充环境变量列表示例（.env.example）包含关键配置项。"
    }
  ],
  "conflicts": [
    {
      "topic": "用户名小写强制实现方式",
      "sides": [
        "backend 建议使用 CHECK(username = LOWER(username)) 约束或表达式唯一索引",
        "qa 要求确保小写存储，但未指定具体数据库级措施"
      ],
      "resolution": "采用表达式唯一索引 CREATE UNIQUE INDEX idx_users_username_lower ON users (LOWER(username))，并在应用层统一小写转换，双重保障。",
      "escalate_to_user": false
    },
    {
      "topic": "迁移策略自动与CI/CD矛盾",
      "sides": [
        "文档描述“应用启动时自动执行迁移”",
        "文档又写“生产环境通过 CI/CD 执行迁移”"
      ],
      "resolution": "统一策略：生产环境禁止自动迁移，通过CI/CD执行；开发环境可通过环境变量控制自动迁移。文档修改以消除歧义。",
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