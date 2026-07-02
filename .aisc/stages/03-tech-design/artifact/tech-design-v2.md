[ARTIFACT]
type: TechDesign
version: 2
status: draft
changes: |
  - 所有待办API响应中强制包含version字段，请求体中接收version以支持乐观锁。
  - 分页响应添加total、page、page_size元数据。
  - 日期时间格式统一为RFC3339。
  - 明确更新操作使用PATCH，并定义可修改字段（title, description, due_date, completed）及其可选性。
  - 认证Token响应字段名统一为"token"。
  - 新增CORS中间件，配置允许的域名。
  - 定义标准化错误码枚举（USERNAME_TAKEN, VALIDATION_ERROR, VERSION_CONFLICT等），ErrorResponse增加error_code字段。
  - 数据库层使用表达式唯一索引（LOWER(username)）强制用户名小写，应用层存储前转换小写。
  - 处理注册并发唯一约束冲突，捕获错误并返回409及标准错误码。
  - 迁移策略澄清：生产环境通过CI/CD执行迁移，开发环境可自动迁移；移除启动自动迁移歧义。
  - 完善乐观锁实现：使用原子UPDATE/DELETE SQL，通过RowsAffected判断，若为0再查询记录以区分404和409。
  - 配置中明确JWT_SECRET和JWT_EXPIRATION环境变量，提供默认值或示例。
  - 用户名字段长度调整为VARCHAR(50)，DDL中注明来源（PRD要求）。
  - 添加健康检查端点/health。
  - 日志中间件中加入request_id，便于链路追踪。
  - 补充环境变量列表示例（.env.example）包含关键配置项。

# 技术设计文档

## 1. 技术选型

| 组件 | 选型 | 说明 |
|------|------|------|
| 后端语言 | Go 1.21+ | 高性能、编译型、适合 API 服务 |
| Web 框架 | Gin v1.9+ | 轻量、高性能、路由/中间件丰富 |
| 数据库 | PostgreSQL 15+ | 符合 ACID，支持 JSON、索引、事务 |
| 缓存/消息队列 | 无 | 当前业务无需引入，通过数据库优化满足性能基线 |
| JWT 库 | github.com/golang-jwt/jwt/v5 | 提供 JWT 签名和验证（HS256） |
| 密码哈希 | golang.org/x/crypto/bcrypt | bcrypt 加密存储密码 |
| 参数校验 | github.com/go-playground/validator/v10 | 结构体标签校验（request body 和 query） |
| 数据库驱动 | github.com/jackc/pgx/v5 | 高性能 PostgreSQL 驱动，支持连接池 |
| 数据库迁移 | github.com/golang-migrate/migrate/v4 | 文件化版本迁移，支持 up/down |
| 配置管理 | 环境变量 + viper（可选） | 读取环境变量即可，暂不强制依赖 |

## 2. 数据库设计

### 2.1 ER 图（Mermaid）

```mermaid
erDiagram
    User ||--o{ Todo : has
    User {
        bigint id PK
        varchar username UK
        varchar password_hash
        timestamptz created_at
        timestamptz updated_at
    }
    Todo {
        bigint id PK
        bigint user_id FK
        varchar title
        text description
        date due_date
        boolean completed
        integer version
        timestamptz created_at
        timestamptz updated_at
    }
```

### 2.2 DDL

```sql
-- users 表
CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    username      VARCHAR(50) NOT NULL,             -- 长度50，来自PRD要求
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- username 唯一约束（使用表达式索引强制小写）
CREATE UNIQUE INDEX idx_users_username_lower ON users (LOWER(username));

-- trigger 自动更新 updated_at（可选，也可在代码中维护）

-- todos 表
CREATE TABLE todos (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    due_date    DATE,
    completed   BOOLEAN NOT NULL DEFAULT FALSE,
    version     INTEGER NOT NULL DEFAULT 1,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 加速按用户查询
CREATE INDEX idx_todos_user_id ON todos (user_id);

-- 加速按用户和创建时间倒序排序
CREATE INDEX idx_todos_user_created_desc ON todos (user_id, created_at DESC);
```

### 2.3 Migration 策略

- 使用 `golang-migrate/migrate` 管理迁移文件。
- 迁移文件位于 `migrations/` 目录，命名规则 `{timestamp}_{title}.up.sql` 和 `{timestamp}_{title}.down.sql`。
- 首次迁移包含上述两张建表语句。
- **迁移执行策略**：
  - **生产环境**：通过 CI/CD 手动触发执行迁移命令（例如 `migrate up`），确保与代码版本一致。
  - **开发环境**：可在应用启动时自动执行迁移（通过 `cmd/server` 的 init 逻辑），但需明确该行为仅限开发环境，并在文档中注明，避免生产环境自动迁移。
  - 统一文档澄清：生产环境禁止启动时自动迁移，开发环境可选自动迁移以提升效率。

## 3. 模块划分

| 模块 | 职责 | 依赖 |
|------|------|------|
| `auth` | 用户注册、登录、JWT 生成 | database, bcrypt, jwt |
| `todo` | 待办 CRUD、列表分页、乐观锁控制 | database, jwt (从 context 取 user_id) |
| `middleware` | 认证中间件（JWT 验证）、CORS中间件、错误处理中间件、请求日志中间件（含request_id） | jwt, logger |
| `database` | 数据库连接池、迁移执行 | pgx, migrate |
| `config` | 加载配置（环境变量） | 无（或 viper） |

**模块间依赖关系**：
- `auth` 和 `todo` 通过 `database` 包访问数据，不直接访问其他业务模块。
- `middleware` 在路由层注入，不依赖具体业务模块（通过 context 传递 user_id）。
- `config` 是全局独例，供所有模块读取。

## 4. 目录结构

```
backend/
├── cmd/
│   ├── server/
│   │   └── main.go          # 入口：加载配置、初始化数据库、启动 HTTP Server
│   └── migrate/
│       └── main.go          # 手动迁移命令（可选）
├── internal/
│   ├── auth/
│   │   ├── handler.go       # 注册/登录 HTTP Handler
│   │   ├── service.go       # 业务逻辑：密码验证、创建用户（含小写转换、唯一冲突捕获）
│   │   ├── repository.go    # 数据库查询（用户查询、插入）
│   │   └── dto.go           # 请求/响应结构体
│   ├── todo/
│   │   ├── handler.go       # 待办 CRUD Handler
│   │   ├── service.go       # 业务逻辑：权限校验、版本号乐观锁、幂等判断
│   │   ├── repository.go    # 数据库查询（待办增删改查）
│   │   └── dto.go           # 请求/响应结构体（包含version字段）
│   ├── middleware/
│   │   ├── auth.go          # JWT 认证中间件
│   │   ├── cors.go          # CORS 中间件，配置允许的域名
│   │   ├── error.go         # 统一错误响应中间件（恢复 & 自定义错误，含error_code）
│   │   └── logger.go        # 请求日志中间件（含request_id生成与记录）
│   ├── database/
│   │   ├── postgres.go      # 数据库连接初始化 & 连接池配置
│   │   └── migrate.go       # 迁移执行逻辑
│   └── config/
│       └── config.go        # 环境变量读取结构体（含JWT_SECRET、JWT_EXPIRATION等）
├── migrations/
│   ├── 20250301000001_create_users.up.sql
│   ├── 20250301000001_create_users.down.sql
│   └── 20250301000002_create_todos.up.sql
│   └── 20250301000002_create_todos.down.sql
├── .env.example              # 环境变量列表示例（含关键配置项）
├── go.mod
├── go.sum
├── Dockerfile
└── Makefile
```

**关键说明**：
- `repository` 层负责 SQL 交互，`service` 层负责业务规则和校验，`handler` 层负责 HTTP 输入输出。
- 所有 handler 均不直接访问数据库，通过 service 调用 repository。
- `middleware/error.go` 捕获所有 panic 和业务错误，统一输出包含 `error_code` 的 `ErrorResponse`。

## 5. 中间件设计

### 5.1 认证中间件（JWT 校验）

- **路径**：`internal/middleware/auth.go`
- **流程**：
  1. 从 `Authorization` 头解析 `Bearer {token}`。
  2. 若缺失或格式错误，直接返回 401。
  3. 使用 HS256 密钥解析 JWT（`golang-jwt/jwt/v5`）。
  4. 若 Token 过期或签名无效，返回 401。
  5. 从 Claims 中提取 `user_id` 字段，并写入 Gin 上下文 `c.Set("user_id", userID)`。
  6. 后续 handler 从上下文取 `user_id`。
- **使用**：在路由分组“/todos”上应用该中间件，使所有待办接口必须认证。同时可全局应用（但 auth 自身的两个端点无需认证）。

### 5.2 CORS 中间件

- **路径**：`internal/middleware/cors.go`
- **功能**：允许跨域请求，允许的域名通过环境变量 `CORS_ALLOWED_ORIGINS` 配置（多个域名用逗号分隔，支持通配符 `*`）。
- **配置要点**：
  - 允许的方法：`GET, POST, PATCH, DELETE, OPTIONS`。
  - 允许的 Header：`Authorization, Content-Type, X-Request-ID`。
  - 预检请求缓存时间：300秒。
- **默认值**：开发环境可设为 `"*"`，生产环境需指定具体域名。

### 5.3 错误处理中间件

- **路径**：`internal/middleware/error.go`
- **职责**：
  - 捕获 Handler 中 panic，返回 500 及标准错误结构。
  - 统一拦截业务错误（可自定义错误类型），返回对应 HTTP 状态码和标准错误响应。
- **统一错误结构**：
  ```go
  type ErrorResponse struct {
      Code      int    `json:"code"`
      ErrorCode string `json:"error_code"` // 标准化错误码，如 "VALIDATION_ERROR"
      Message   string `json:"message"`
  }
  ```
- **标准化错误码枚举**（定义在 `internal/errors/errors.go` 或类似位置）：
  - `VALIDATION_ERROR` （400，参数校验失败）
  - `UNAUTHORIZED` （401，认证失败）
  - `FORBIDDEN` （403，无权限）
  - `NOT_FOUND` （404，资源不存在）
  - `USERNAME_TAKEN` （409，用户名已占用）
  - `VERSION_CONFLICT` （409，版本号冲突）
  - `INTERNAL_ERROR` （500，内部服务错误）
  - `DUPLICATE_ENTRY` （409，唯一约束冲突的通用码）
- **实现方式**：在路由组顶层使用 `gin.Recovery()` + 自定义 Middleware 对 `c.Errors` 进行格式化。更好的方案：封装 `c.JSON` 并约定业务方法返回错误，由 middleware 统一渲染（也可在 handler 中直接 return 错误并通过 `c.AbortWithStatusJSON` 实现，但统一中间件可以简化代码）。
- **设计原则**：所有业务错误（400/401/403/404/409/500）均以 `ErrorResponse` 格式返回，`code` 等于 HTTP 状态码，`error_code` 为标准化字符串。

### 5.4 日志中间件

- **路径**：`internal/middleware/logger.go`
- **功能**：
  - 为每个请求生成唯一的 `request_id`（UUID v4），并注入到 Gin 上下文 `c.Set("request_id", requestID)`。
  - 记录请求的方法、路径、状态码、耗时、`request_id`。
  - 建议使用 `slog` 或 `logrus` 结构化输出日志，包含 `request_id` 字段。
- **实现**：使用 Gin 自定义中间件，在请求开始时生成 `request_id`，写入上下文和响应头（可选），在结束前记录日志。

## 6. 非功能性设计

### 6.1 并发控制策略（乐观锁）

- **背景**：PRD 要求对同一待办的并发修改（更新、删除）使用乐观锁，版本号不匹配时返回 409。
- **实现机制**：
  - `todos` 表包含 `version` 字段，每次成功更新（或删除）后 `version` +1。
  - 更新/删除时，SQL 的 `WHERE` 条件中必须包含 `version = :reqVersion`。
  - **优化流程**（原子操作）：
    1. 直接执行带版本条件的 UPDATE 或 DELETE SQL。
    2. 检查 `RowsAffected`：
        - 若为 1，操作成功。
        - 若为 0，则再次查询记录是否存在：
            - 若记录存在，说明版本冲突，返回 409（`VERSION_CONFLICT`）。
            - 若记录不存在，返回 404（`NOT_FOUND`）。
  - **幂等性注意**（PATCH）：
    - 若请求仅包含 `completed: true` 且当前已完成，且版本号匹配，则直接返回数据而不修改（version 不变）。
    - 根据 API Spec 细致实现幂等逻辑，该判定在 service 层完成。
- **API 响应和请求中 version 字段**：
  - 所有待办对象的响应中强制包含 `version` 字段。
  - 更新、删除请求体中必须包含 `version` 字段（用于乐观锁校验）。
  - 创建待办不要求客户端传入 version，由服务端初始化为 1 并在响应中返回。

### 6.2 数据库连接池配置

- 使用 `pgx/v5/pgxpool` 管理连接池。
- 建议配置：
  - `MaxConns`: 根据 CPU 核数和并发期望设定，初始 25（在 1000 用户，性能基线场景下足够）。
  - `MinConns`: 2 ~ 5（避免冷启动）。
  - `MaxConnLifetime`: 30 分钟（自动回收老旧连接）。
  - `HealthCheckInterval`: 1 分钟。
- 配置通过环境变量可调整。

### 6.3 API 限流方案

- **当前阶段**：未强制要求，但由于性能基线条件（1000 用户，每个用户最多 1000 条待办），默认负载不高，暂不引入限流。
- **未来扩展**：如果遇到突发流量，可考虑在网关层或中间件中使用令牌桶（如 `github.com/ulule/limiter`）对 /api/v1 全局限流（例如每秒 1000 请求）。
- 本设计不包含限流代码，但预留接口便于后续添加。

### 6.4 性能优化要点

- **索引**：已创建 `user_id + created_at` 复合索引，满足列表查询排序和过滤。
- **查询优化**：分页查询使用 `OFFSET + LIMIT`（数据量不大时可接受，若后续超过百万可改为游标分页）。
  - 分页响应必须包含 `total`（总条数）、`page`（当前页）、`page_size`（每页条数）字段。
- **Response 体量**：避免不必要字段序列化，`Todo` 结构体与数据库字段一致，减少转换开销。
- **压力测试**：上线前使用 k6 或 wrk 模拟 PRD 中定义的基线场景（1000 用户、1000 待办/用户），验证平均响应时间 <200ms。若不符合，考虑增加连接池大小或引入 Redis 缓存（但当前业务简单，无缓存必要）。

## 7. 配置管理

### 7.1 环境变量

| 环境变量 | 必需 | 默认值 / 示例 | 说明 |
|---------|------|----------------|------|
| `DATABASE_URL` | 是 | `postgres://user:pass@localhost:5432/todoapp` | 数据库连接字符串 |
| `JWT_SECRET` | 是 | `your-256-bit-secret`（示例） | JWT 签名密钥，至少 32 字符 |
| `JWT_EXPIRATION` | 是 | `24h`（示例） | Token 过期时间，遵循 Go duration 格式 |
| `SERVER_PORT` | 否 | `8080` | HTTP 监听端口 |
| `CORS_ALLOWED_ORIGINS` | 否 | `http://localhost:3000` | 允许的跨域来源域名，多个用逗号分隔 |
| `LOG_LEVEL` | 否 | `info` | 日志级别（debug, info, warn, error） |
| `DB_MAX_CONNS` | 否 | `25` | 数据库连接池最大连接数 |
| `DB_MIN_CONNS` | 否 | `5` | 连接池最小连接数 |

### 7.2 `.env.example` 示例

```bash
# 必填配置
DATABASE_URL=postgres://user:pass@localhost:5432/todoapp
JWT_SECRET=your-256-bit-secret
JWT_EXPIRATION=24h

# 可选配置
SERVER_PORT=8080
CORS_ALLOWED_ORIGINS=http://localhost:3000
LOG_LEVEL=info
DB_MAX_CONNS=25
DB_MIN_CONNS=5
```

## 8. API 接口（补充说明）

### 8.1 日期时间格式

所有请求和响体中的日期时间字段统一使用 **RFC3339** 格式（例如 `2025-03-15T10:30:00Z`），通过 Go 标准库 `time.Format(time.RFC3339)` 序列化。

### 8.2 更新操作

- **HTTP 方法**：PATCH
- **可修改字段**：`title`（可选）、`description`（可选）、`due_date`（可选）、`completed`（可选）
- **请求体示例**：
  ```json
  {
    "title": "新的标题",
    "description": "可选描述",
    "due_date": "2025-04-01",
    "completed": true,
    "version": 2
  }
  ```
- **强制字段**：`version`（必传）。其余字段至少提供一个。

### 8.3 认证响应

注册和登录成功后返回 JWT Token，响应字段统一命名为 `"token"`：
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### 8.4 健康检查端点

- **路径**：`GET /health`
- **用途**：用于负载均衡或监控检测服务状态。
- **响应**：
  ```json
  {
    "status": "ok",
    "timestamp": "2025-03-15T10:30:00Z"
  }
  ```

### 8.5 分页响应结构

```json
{
  "data": [ ... ],
  "total": 100,
  "page": 1,
  "page_size": 20
}
```

### 8.6 注册并发冲突处理

在 `auth/service.go` 中插入用户时，若 PostgreSQL 返回唯一约束冲突错误（`23505`），应用层捕获该错误并：
- 返回 HTTP 409 状态码。
- ErrorResponse 中 `error_code` 为 `"USERNAME_TAKEN"`。

### 8.7 用户名小写实现

- **应用层**：注册时，在 `service` 层将用户传入的 `username` 强制转换为小写后再存储。
- **数据库层**：唯一索引使用 `CREATE UNIQUE INDEX idx_users_username_lower ON users (LOWER(username));`，确保数据库层面也强制小写唯一性（防御性）。

---

> 本设计遵循 PRD 和 API Spec 所有要求，包括统一错误格式（含 error_code）、乐观锁、严格校验规则、权限隔离、分页元数据、用户名小写、日期格式、CORS、健康检查、request_id 链路追踪等。所有边界条件均明确实现路径。