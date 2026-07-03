所有修改完成。以下是修改总结：

## 修改总结

### 逐项修改明细

| # | ActionItem | 状态 | 修改内容 |
|---|-----------|------|---------|
| 1 | 提供完整的 api-spec-frozen.yaml，确认 `/auth/me` 端点 | ✅ 新增 | 创建 `docs/api-spec-frozen.yaml`（OpenAPI 3.0.3），包含全部 8 个端点定义（`/health`、`/auth/register`、`/auth/login`、`/auth/me`、`/todos` GET/POST、`/todos/{todo_id}` GET/PATCH/DELETE），其中 `/auth/me` 标注 `BearerAuth` 安全要求，返回 `UserPublic` schema；同时包含变更日志、标准错误响应、分页参数校验（page ≥1, page_size 1~100）、密码校验说明等契约一致性内容 |
| 2 | 将 Go 版本修正为当前官方支持的稳定版本 | ✅ 修正 | `go.mod`：`go 1.25` → `go 1.22`；`Dockerfile` 已为 `golang:1.22-alpine`，版本一致 |

### 编译验证

```
cd backend
go build ./...   → OK
go vet ./...     → OK
```