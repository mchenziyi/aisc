## 修改总结

所有 10 项 ActionItem 均已实施，`go build ./...` 和 `go vet ./...` 均通过。

### 修改清单

| # | ActionItem | 修改的文件 | 改动说明 |
|---|-----------|-----------|----------|
| 1 | DELETE 版本传递方式 | `internal/todo/handler.go`, `internal/todo/dto.go` | 移除 `DeleteTodoRequest`，改用查询参数 `?version=N` 传递版本号；删除 `dto.go` 中已无用的 `DeleteTodoRequest` 结构体 |
| 2 | 修复 date/timestamptz 扫描 | `internal/auth/dto.go`, `internal/todo/repository.go` | User.CreatedAt/UpdatedAt 改为 `time.Time`；todo 仓库扫描改用 `time.Time`/`*time.Time` 中间变量再格式化为字符串（`formatTime` 辅助函数），避免 pgx 无法将 timestamptz 扫描为 string 的运行时崩溃 |
| 3 | JWT_SECRET 强制配置 | `internal/config/config.go` | 移除 `JWT_SECRET` 默认值，新增 `getEnvRequired()` 函数，环境变量未设置时 panic |
| 4 | 添加 `/auth/me` 端点 | `internal/auth/service.go`, `internal/auth/handler.go`, `cmd/server/main.go` | 新增 `GetMe` service 方法；新增 `Me` handler；在 `/auth` 分组添加 `GET /me` 路由（需 JWT 认证） |
| 5 | 修复 user_id 类型断言 | `internal/middleware/auth.go` | 增加逗号断言模式检查：`userIDVal.(float64)` 失败时返回 401，避免 panic |
| 6 | 统一错误处理 | `internal/auth/handler.go`, `internal/todo/handler.go` | 区分客户端验证错误（返回 400）与非预期服务器错误（返回 500）；新增 `isValidationError()` 辅助函数按错误消息关键词判断；Login 同样区分凭据错误(401)和服务器错误(500) |
| 7 | 简化 nullable 字段 | `internal/todo/dto.go`, `internal/todo/service.go` | 用 `OptionalString` 结构体（含 Set/Null/Value 字段）替代 `**string`，提供 `UnmarshalJSON` 方法区分"未提供/设为 null/设置值"三种状态 |
| 8 | 请求体大小限制 | `internal/middleware/bodylimit.go` (新建), `cmd/server/main.go` | 新增 `MaxBodySize` 中间件（基于 `http.MaxBytesReader`），对 `/api/v1` 下所有路由限制 1MB |
| 9 | 密码策略说明 | `internal/auth/service.go` | 更新 `passwordRegex` 注释，说明当前允许特殊符号；错误消息补充 `(special characters allowed)` |
| 10 | 编译验证 | — | `go build ./...` ✅, `go vet ./...` ✅, `go mod verify` ✅ |