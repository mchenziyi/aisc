All modules verified. Let me provide the final implementation summary.

---

## 实现总结

### 创建的文件（共 17 个）

```
backend/
├── .env.example                              # 环境变量示例
├── go.mod                                    # Go 模块定义
├── go.sum                                    # 依赖锁文件
├── cmd/
│   ├── migrate/main.go                       # 手动迁移命令
│   └── server/main.go                        # 服务入口（路由、中间件、启动）
├── internal/
│   ├── auth/
│   │   ├── dto.go                            # 请求/响应结构体
│   │   ├── repository.go                     # 用户 CRUD（参数化查询）
│   │   ├── service.go                        # 注册/登录业务逻辑（bcrypt + JWT）
│   │   └── handler.go                        # HTTP Handler
│   ├── config/
│   │   └── config.go                         # 环境变量加载
│   ├── database/
│   │   ├── postgres.go                       # pgx 连接池
│   │   ├── migrate.go                        # 嵌入式迁移执行
│   │   └── migrations/
│   │       ├── 20250301000001_init.up.sql     # 建表 DDL
│   │       └── 20250301000001_init.down.sql   # 回滚 DDL
│   ├── middleware/
│   │   ├── auth.go                           # JWT 认证中间件
│   │   ├── cors.go                           # CORS 中间件
│   │   ├── error.go                          # Panic 恢复中间件
│   │   └── logger.go                         # 请求日志 + Request ID
│   └── todo/
│       ├── dto.go                            # 请求/响应结构体
│       ├── repository.go                     # 待办 CRUD（乐观锁、用户隔离）
│       ├── service.go                        # 业务逻辑（校验、幂等、并发控制）
│       └── handler.go                        # HTTP Handler
```

### 接口清单（与 API Spec 完全对应）

| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| GET | `/health` | 否 | 健康检查（含数据库 Ping） |
| POST | `/api/v1/auth/register` | 否 | 用户注册 |
| POST | `/api/v1/auth/login` | 否 | 用户登录 |
| POST | `/api/v1/todos` | JWT | 创建待办 |
| GET | `/api/v1/todos` | JWT | 待办列表（分页） |
| GET | `/api/v1/todos/:todo_id` | JWT | 待办详情 |
| PATCH | `/api/v1/todos/:todo_id` | JWT | 更新待办（含乐观锁） |
| DELETE | `/api/v1/todos/:todo_id` | JWT | 删除待办（含乐观锁） |

### 关键实现细节

1. **认证**：JWT (HS256) 24h 过期，Claims 中包含 `user_id`
2. **密码**：bcrypt 加密存储
3. **用户隔离**：所有 SQL 包含 `WHERE user_id = $N`，按用户过滤
4. **乐观锁**：version 字段，UPDATE/DELETE 时检查版本号匹配，冲突返回 409
5. **PATCH 幂等性**：若仅传 `completed` 且状态一致，跳过更新（version 不变）
6. **可空字段**：description / due_date 支持 `**string` 区分"不更新"与"设为 null"
7. **分页**：page ≥ 1，page_size 1~100，超页返回空数组
8. **统一错误格式**：`{ "code": int, "message": string }`，code = HTTP 状态码

### 编译验证

- `go build ./...` ✅ 通过
- `go vet ./...` ✅ 无警告
- `go mod verify` ✅ 依赖完整