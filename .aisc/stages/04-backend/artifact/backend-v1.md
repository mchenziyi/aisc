---

## 实现总结

### 创建的文件清单（共 22 个）

| 文件 | 说明 |
|------|------|
| `cmd/server/main.go` | 入口：加载配置、初始化数据库、路由注册、优雅启停 |
| `cmd/migrate/main.go` | 可选的手动迁移命令 |
| `internal/config/config.go` | 环境变量加载（DATABASE_URL, JWT_SECRET, JWT_EXPIRATION 等） |
| `internal/database/postgres.go` | pgx 连接池初始化（MaxConns/MinConns/MaxConnLifetime） |
| `internal/database/migrate.go` | 按文件名顺序执行 `*.up.sql` 迁移文件 |
| `internal/errors/errors.go` | 标准化错误结构 + 7 种业务错误码（VALIDATION_ERROR, UNAUTHORIZED, NOT_FOUND, VERSION_CONFLICT 等） |
| `internal/auth/dto.go` | RegisterRequest/RegisterResponse/LoginRequest/LoginResponse/UserPublic |
| `internal/auth/repository.go` | `CreateUser`（插入用户）、`FindByUsername`（LOWER 不区分大小写查询） |
| `internal/auth/service.go` | 注册（用户名正则校验 + 密码复杂度校验 + bcrypt + JWT 签发）、登录（小写转换 + bcrypt 比对 + JWT 签发） |
| `internal/auth/handler.go` | `POST /api/v1/auth/register`、`POST /api/v1/auth/login` |
| `internal/todo/dto.go` | CreateTodoRequest/UpdateTodoRequest/TodoResponse/TodoListResponse |
| `internal/todo/repository.go` | CRUD + 动态 SQL 部分更新 + 乐观锁版本控制（`errVersionConflict` 哨兵错误） |
| `internal/todo/service.go` | 创建验证（标题非空+长度、描述长度、日期格式）、列表分页、详情、更新（幂等 completed 处理、字段级验证）、删除 |
| `internal/todo/handler.go` | `GET/POST /todos`、`GET/PATCH/DELETE /todos/:todo_id`（含分页参数校验、version 查询参数） |
| `internal/middleware/auth.go` | JWT 认证中间件（Bearer Token 解析、user_id 注入上下文、401 统一响应） |
| `internal/middleware/cors.go` | CORS 中间件（可配置 allowed origins、支持 OPTIONS 预检） |
| `internal/middleware/error.go` | 统一错误处理中间件（捕获 `c.Error()` 的 `*AppError`，输出标准 JSON 格式） |
| `internal/middleware/logger.go` | 请求日志中间件（UUID v4 request_id + 方法/路径/状态码/耗时） |
| `migrations/20250301000001_create_users.up.sql` | users 建表 + 小写唯一索引 |
| `migrations/20250301000002_create_todos.up.sql` | todos 建表 + user_id 索引 + user_id+created_at 复合降序索引 |
| `go.mod` / `go.sum` | Go 模块定义和依赖锁定 |

### 关键设计决策

1. **API Spec v4.0.0 为准**：跨用户访问返回 404（而非 PRD 中的 403）、注册返回 JWT Token、错误响应包含 `request_id` 和 `error_code`、DELETE version 通过查询参数传递
2. **乐观锁**：UPDATE/DELETE 在 WHERE 中同时校验 `id`、`user_id`、`version`，RowsAffected==0 时二次查询区分「不存在（404）」和「版本冲突（409）」
3. **动态 SQL 部分更新**：PATCH 请求根据提供的字段动态构建 SET 子句，避免 COALESCE 对空字符串的误判
4. **幂等 completed**：当仅传 `version + completed` 且值与当前一致时，直接返回 200 不修改版本号
5. **数据隔离**：所有 SQL 均包含 `WHERE user_id = $N`，确保用户只能操作自己的数据

### 编译验证

```
go build ./...  →  ok
go vet ./...    →  ok
```