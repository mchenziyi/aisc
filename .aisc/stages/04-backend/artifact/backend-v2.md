## 修改总结

所有 11 项 Action Item 已全部完成，`go build ./...` 和 `go vet ./...` 均通过。

### 修改清单

| # | Action Item | 修改文件 | 变更说明 |
|---|------------|---------|---------|
| 1 | 补充 API 规范文件 | `docs/api-spec-frozen.yaml` | 更新为 v5.0.0，新增 `/auth/me`、`/health` 端点；DELETE 响应改为 200 + 确认体；details 改为 JSON 对象；todo_id 增加 minimum:1；用户名示例改为原始大小写 |
| 2 | 修复日期解析 | `backend/internal/todo/service.go` | `parseDate` 增加回写校验，拒绝 `2024-02-30` 等无效日期 |
| 3 | 自动迁移解耦 | `backend/internal/config/config.go` + `backend/cmd/server/main.go` | 新增 `RunMigrations` 配置项，通过 `RUN_MIGRATIONS=true` 环境变量控制，生产环境默认不执行 |
| 4 | 日志级别控制 | `backend/internal/middleware/logger.go` | 重写日志中间件，读取 `LOG_LEVEL` 环境变量（debug/info/warn/error），按状态码自动选择日志级别 |
| 5 | NullableString 统一 | `backend/internal/todo/dto.go` + `service.go` | 添加明确注释：空字符串 `""` 与 `null` 同等对待，Create 和 Update 行为一致 |
| 6 | 版本冲突详情结构化 | `backend/internal/errors/errors.go` | `NewVersionConflictError` 的 Details 改为 `{"current_version": 5}` JSON 对象 |
| 7 | 简化健康检查 | `backend/cmd/server/main.go` | 移除冗余 `database` 字段，DB 不可用时 `status` 变为 `"degraded"`，HTTP 状态码 503 |
| 8 | todo_id 正数校验 | `backend/internal/todo/handler.go` | GetTodo/UpdateTodo/DeleteTodo 均增加 `todoID < 1` → 400 |
| 9 | /auth/me 响应统一 | `backend/internal/auth/handler.go` | 直接返回 UserPublic 对象，不再嵌套 `{"user": ...}` |
| 10 | 删除返回 200 | `backend/internal/todo/handler.go` | DELETE 响应改为 200 OK + `{"id": ..., "deleted": true}` |
| 11 | 用户名大小写处理 | `backend/internal/auth/service.go` | 注册时保留用户输入的原始大小写，响应中原样返回；登录仍通过 `LOWER()` 保持大小写不敏感 |