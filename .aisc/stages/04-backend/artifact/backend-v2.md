所有 13 个 Action Item 修改完成，编译和静态分析均通过。下面是完整修改总结。

---

## 修改总结

### 1. 修正 go.mod Go 版本
- **文件**: `backend/go.mod`
- **修改**: 将 `go 1.25.5` 改为 `go 1.22`（经 `go mod tidy` 后自动调整为 `go 1.25.0`，与当前工具链匹配，项目可编译）

### 2. 修复迁移脚本索引创建
- **文件**: `backend/migrations/20250301000001_create_users.up.sql`
- **文件**: `backend/migrations/20250301000002_create_todos.up.sql`
- **修改**: 将 `CREATE UNIQUE INDEX IF NOT EXISTS` 移除，改用 `DO $$ ... END $$;` 匿名块包裹，先检查 `pg_indexes` 中索引是否存在，不存在时再创建，实现幂等迁移

### 3. 统一可空字段语义
- **文件**: `backend/internal/todo/dto.go` — 更新 `UpdateTodoRequest` 注释：`description`/`due_date` 不再说 "use null to clear"，改为 "empty string clears the field (sets to NULL)"
- **文件**: `backend/internal/todo/repository.go` — `UpdateFields` 结构体重构：`Description *string` 改为 `DescriptionVal *string` + `UpdateDescription bool`，与 `DueDateVal`/`UpdateDueDate` 模式一致；nil 值指针映射为 SQL NULL
- **文件**: `backend/internal/todo/service.go` — `Update` 方法中根据 `req.Description` 判断：若为 nil 则不更新；若为 `""` 则 clear（设 DescriptionVal=nil → SQL NULL）；否则设置值

### 4. 修复 parseDate 重复包装
- **文件**: `backend/internal/todo/service.go` — `parseDate` 函数内部不再返回 `*apperrors.AppError`，改为返回 `var errInvalidDateFormat`；调用方统一用 `apperrors.NewValidationError(err.Error())` 单次包装

### 5. AuthMiddleware 统一错误输出
- **文件**: `backend/internal/middleware/auth.go` — 移除所有 `c.AbortWithStatusJSON()` 手动构造，全部改为 `c.Error(apperrors.NewUnauthorizedError("unauthorized"))` + `c.Abort()`，由 ErrorMiddleware 统一格式化
- 移除不再需要的 `"net/http"` 导入

### 6. 优化请求体验证错误信息
- **文件**: `backend/internal/errors/errors.go` — 新增 `NewValidationErrorFromBinding()` 函数和 `AsValidationErrors` 变量，从 `validator.ValidationErrors` 提取字段名和 tag，返回具体错误（如 `"field 'Username' required"`）
- **文件**: `backend/internal/auth/handler.go` — Register/Login 改用 `NewValidationErrorFromBinding(err)`
- **文件**: `backend/internal/todo/handler.go` — CreateTodo/UpdateTodo 改用 `NewValidationErrorFromBinding(err)`

### 7. 改进分页参数错误反馈
- **文件**: `backend/internal/todo/handler.go` — 区分格式错误和范围错误：`parseQueryInt` 解析失败返回 `"must be a valid integer"`，范围检查返回 `"must be >= 1"` 或 `"must not exceed 100"`

### 8. 统一删除操作 version 传递方式
- **文件**: `backend/internal/todo/dto.go` — 新增 `DeleteTodoRequest` 结构体（`version` 字段，`binding:"required"`）
- **文件**: `backend/internal/todo/handler.go` — `DeleteTodo` 改为从请求体解析 `DeleteTodoRequest`，而非 URL 查询参数 `?version=N`

### 9. 添加获取当前用户信息端点
- **文件**: `backend/internal/auth/repository.go` — 新增 `FindByID(ctx, id)` 方法
- **文件**: `backend/internal/auth/service.go` — 新增 `GetMe(ctx, userID)` 方法
- **文件**: `backend/internal/auth/handler.go` — 新增 `Me` handler（`GET /api/v1/auth/me`），返回当前用户信息
- **文件**: `backend/cmd/server/main.go` — 注册 `/auth/me` 路由（受 JWT 中间件保护）

### 10. 优化密码强度验证错误提示
- **文件**: `backend/internal/auth/service.go` — `validatePassword()` 区分三种失败原因返回具体消息：`"password must be at least 8 characters"` / `"password must contain at least one letter"` / `"password must contain at least one digit"` / `"password must contain at least one letter and one digit"`

### 11. 修改 CORS 默认值
- **文件**: `backend/internal/config/config.go` — `CORS_ALLOWED_ORIGINS` 默认值从 `"*"` 改为 `"http://localhost:3000"`

### 12. 配置加载时对默认 JWT 密钥打印警告
- **文件**: `backend/internal/config/config.go` — `Load()` 函数中检测到使用默认密钥时打印 `log.Println("WARNING: ...")`

### 13. 增加安全响应头
- **文件**: `backend/internal/middleware/security.go` — 新增 `SecurityHeadersMiddleware`，设置 `X-Content-Type-Options: nosniff`、`X-Frame-Options: DENY`、`X-XSS-Protection: 1; mode=block`、`Referrer-Policy: strict-origin-when-cross-origin`、`Permissions-Policy`
- **文件**: `backend/cmd/server/main.go` — 注册 `SecurityHeadersMiddleware` 为全局中间件