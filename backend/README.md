# Todo Management API

轻量级待办事项管理 REST API，基于 Go + Gin + PostgreSQL 构建。

## 技术栈

- **语言**: Go 1.22+
- **框架**: Gin v1.10
- **数据库**: PostgreSQL 15+ (pgx v5)
- **认证**: JWT (golang-jwt v5)
- **密码**: bcrypt
- **验证**: go-playground/validator v10

## 项目结构

```
backend/
├── cmd/
│   └── server/main.go          # 服务入口
├── internal/
│   ├── auth/                   # 用户认证模块
│   │   ├── handler.go          # HTTP 处理器
│   │   ├── service.go          # 业务逻辑
│   │   ├── repository.go       # 数据访问
│   │   ├── types.go            # 请求/响应定义
│   │   └── claims.go           # JWT 声明
│   ├── todo/                   # 待办事项模块
│   │   ├── handler.go          # HTTP 处理器
│   │   ├── service.go          # 业务逻辑
│   │   ├── repository.go       # 数据访问
│   │   └── types.go            # 请求/响应定义
│   ├── middleware/             # 中间件
│   │   ├── auth.go             # JWT 认证
│   │   ├── cors.go             # CORS
│   │   ├── recovery.go         # 异常恢复
│   │   ├── logger.go           # 请求日志
│   │   ├── request_id.go       # 请求追踪 ID
│   │   └── security.go         # 安全头
│   ├── config/config.go        # 配置管理
│   ├── database/
│   │   ├── postgres.go         # 数据库连接池
│   │   └── migrate.go          # 迁移执行器
│   ├── errors/errors.go        # 业务错误处理
│   └── model/types.go          # 共享数据模型
├── migrations/                 # 数据库迁移文件
├── Dockerfile
├── Makefile
├── go.mod
└── .env.example
```

## 快速开始

### 前置要求

- Go 1.22+
- PostgreSQL 15+

### 配置

复制环境变量文件并根据需要修改：

```bash
cp .env.example .env
```

关键配置项：

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `SERVER_PORT` | 服务端口 | 8080 |
| `DATABASE_URL` | 数据库连接串 | postgres://postgres:postgres@localhost:5432/todoapp?sslmode=disable |
| `JWT_SECRET` | JWT 签名密钥（至少32字符） | 必填 |
| `JWT_EXPIRATION` | JWT 过期时间 | 24h |
| `CORS_ALLOWED_ORIGINS` | CORS 允许的域名 | * |
| `LOG_LEVEL` | 日志级别 (debug/info/warn/error) | info |

### 启动

```bash
# 运行数据库迁移
make migrate

# 启动服务
make run

# 构建
make build
```

## API 端点

所有业务接口以 `/api/v1/` 为前缀。

| 方法 | 路径 | 说明 | 认证 | 乐观锁 |
|------|------|------|------|--------|
| GET | `/health` | 健康检查 | 否 | - |
| POST | `/api/v1/auth/register` | 用户注册 | 否 | - |
| POST | `/api/v1/auth/login` | 用户登录 | 否 | - |
| GET | `/api/v1/auth/me` | 获取当前用户 | 是 | - |
| GET | `/api/v1/todos` | 待办列表（分页+过滤） | 是 | - |
| POST | `/api/v1/todos` | 创建待办 | 是 | - |
| GET | `/api/v1/todos/{id}` | 查看待办详情 | 是 | - |
| PATCH | `/api/v1/todos/{id}` | 更新待办（部分更新，含乐观锁） | 是 | 是 (version) |
| DELETE | `/api/v1/todos/{id}` | 删除待办（需 version 参数） | 是 | 是 (version) |

## 请求/响应格式

### 成功响应

直接返回数据对象，无外层信封。

```json
// 单条数据
{
  "id": 1,
  "title": "完成报告",
  "completed": false,
  "version": 1,
  ...
}

// 列表数据
{
  "items": [...],
  "total": 42,
  "page": 1,
  "page_size": 20,
  "total_pages": 3
}

// 认证响应
{
  "user": {
    "id": 1,
    "username": "alice"
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### 错误响应

```json
{
  "error_code": "VALIDATION_ERROR",
  "message": "请求参数错误",
  "request_id": "a1b2c3d4",
  "errors": [
    { "field": "title", "message": "标题不能为空" }
  ]
}
```

## 错误码

| HTTP 状态码 | error_code | 说明 |
|-------------|------------|------|
| 400 | VALIDATION_ERROR | 参数校验失败 |
| 401 | UNAUTHORIZED | 认证失败 |
| 404 | NOT_FOUND | 资源不存在 |
| 409 | CONFLICT | 用户名已存在 |
| 409 | VERSION_CONFLICT | 数据版本冲突 |
| 500 | INTERNAL_ERROR | 服务器内部错误 |

## 认证

所有受保护的端点需要在请求头中携带 JWT Bearer Token：

```
Authorization: Bearer <token>
```

## 分页

`GET /api/v1/todos` 支持分页参数：

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | int | 1 | 页码（从1开始） |
| page_size | int | 20 | 每页条数（最大100） |
| status | string | all | 过滤: all / pending / completed |

## 乐观锁

更新和删除操作使用 `version` 字段实现乐观锁。

- **更新 (PATCH)**: 请求体必须包含 `version` 字段，值为当前数据版本
- **删除 (DELETE)**: 必须携带 `version` 查询参数

如果服务器端版本不匹配，返回 `409 VERSION_CONFLICT`。
