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
│   ├── server/main.go          # 服务入口
│   └── migrate/main.go         # 数据库迁移入口
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
│   │   ├── error.go            # 统一错误处理
│   │   ├── logger.go           # 请求日志
│   │   └── security.go         # 安全头
│   ├── config/config.go        # 配置管理
│   ├── database/
│   │   ├── postgres.go         # 数据库连接池
│   │   └── migrate.go          # 迁移执行器
│   ├── errors/errors.go        # 业务错误码
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
| `TOKEN_EXPIRY` | Refresh Token 过期时间 | 168h (7天) |
| `CORS_ALLOWED_ORIGINS` | CORS 允许的域名 | * |

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

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/health` | 健康检查 | 否 |
| POST | `/v1/users` | 用户注册 | 否 |
| POST | `/v1/auth/login` | 用户登录 | 否 |
| POST | `/v1/auth/refresh` | 刷新 Token | 是 |
| GET | `/v1/users/me` | 获取当前用户 | 是 |
| GET | `/v1/todos` | 待办列表（分页+过滤） | 是 |
| POST | `/v1/todos` | 创建待办 | 是 |
| GET | `/v1/todos/{id}` | 查看待办详情 | 是 |
| PUT | `/v1/todos/{id}` | 更新待办 | 是 |
| PATCH | `/v1/todos/{id}` | 完成待办 | 是 |
| DELETE | `/v1/todos/{id}` | 删除待办 | 是 |

## 响应格式

### 成功响应

```json
{
  "code": 0,
  "message": "ok",
  "data": { ... }
}
```

### 错误响应

```json
{
  "code": 1001,
  "message": "请求参数错误",
  "errors": [
    { "field": "title", "message": "标题不能为空" }
  ]
}
```

## 错误码

| HTTP 状态码 | 业务码 | 说明 |
|-------------|--------|------|
| 400 | 1001 | 参数校验失败 |
| 409 | 1002 | 用户名已存在 |
| 401 | 2001 | 认证失败 |
| 401 | 2002 | Refresh Token 无效 |
| 404 | 3001 | 资源不存在 |
| 500 | 9999 | 服务器内部错误 |
