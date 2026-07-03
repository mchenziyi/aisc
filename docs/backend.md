## 实现总结

### 已创建/存在的文件清单（共 20 个 `.go` 文件）

| 路径 | 说明 |
|------|------|
| `cmd/server/main.go` | 服务入口：加载配置、初始化 DB、注册路由、优雅关闭 |
| `cmd/migrate/main.go` | 手动迁移工具 |
| `internal/config/config.go` | 环境变量配置加载 |
| `internal/database/postgres.go` | pgx 连接池初始化 |
| `internal/database/migrate.go` | SQL 迁移执行引擎 |
| `internal/errors/errors.go` | 统一错误类型 + 标准化错误码 |
| `internal/auth/claims.go` | JWT Claims 结构体 |
| `internal/auth/dto.go` | 注册/登录请求响应结构体 |
| `internal/auth/handler.go` | 注册/登录 HTTP Handler |
| `internal/auth/service.go` | 用户注册（含小写转换）、登录、密码验证 |
| `internal/auth/repository.go` | 用户 CRUD（参数化查询） |
| `internal/todo/dto.go` | 待办 CRUD 请求/响应结构体 + NullableString |
| `internal/todo/handler.go` | 待办 CRUD HTTP Handler |
| `internal/todo/service.go` | 待办业务逻辑（乐观锁、幂等性、权限校验） |
| `internal/todo/repository.go` | 待办 CRUD（参数化查询 + 动态 SET 子句） |
| `internal/middleware/auth.go` | JWT 认证中间件 |
| `internal/middleware/cors.go` | CORS 中间件 |
| `internal/middleware/error.go` | 统一错误响应中间件 |
| `internal/middleware/logger.go` | 请求日志中间件（含 request_id UUID） |
| `internal/middleware/security.go` | 安全响应头中间件 |

### 本次修改的内容

| 文件 | 修改内容 |
|------|---------|
| `cmd/server/main.go` | 1) 健康检查 DB 不可用时返回 **503**（而非 200）；2) 移除 `/auth/me` 路由（不在 API Spec 中） |
| `internal/todo/handler.go` | DELETE 操作返回 **204 No Content**（而非 200 JSON） |
| `internal/auth/service.go` | 注册时用户名统一转为 **小写** 存储（符合 PRD + API Spec） |
| `internal/todo/dto.go` | 1) `Version` 字段增加 `binding:"min=1"`；2) `Title` 字段增加 `binding:"min=1,max=255"` |

### 编译验证

```
cd backend && go build ./...  →  OK
cd backend && go vet ./...    →  OK
```

### 与冻结文档的符合度

- **API Spec v4.0.0**：全部 7 个接口（注册、登录、创建待办、列表、详情、更新、删除）均已实现，路径/方法/请求体/响应体/状态码严格对齐。错误响应包含 `code`、`error_code`、`message`、`request_id` 四个字段。
- **PRD**：所有功能点、业务规则、边界条件均已覆盖。
- **Tech Design**：目录结构、模块划分、乐观锁策略、数据隔离（`WHERE user_id = :currentUserId`）、连接池配置均符合设计文档。