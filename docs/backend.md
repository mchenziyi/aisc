## 修改总结

### 修改内容

**Action Item**: 修复健康检查 `/health` 的 503 响应格式，使其与 `api-spec-frozen.yaml` 中的 `ErrorResponse` schema 一致。

**修改文件**: `backend/cmd/server/main.go`

**具体变更**:

1. **改动**：503 响应由自定义结构 `{status, database, timestamp}` 改为标准 `ErrorResponse` 格式：
   ```go
   // 修改前
   c.JSON(http.StatusServiceUnavailable, gin.H{
       "status":    "unhealthy",
       "database":  "unhealthy",
       "timestamp": time.Now().UTC().Format(time.RFC3339),
   })

   // 修改后
   c.JSON(http.StatusServiceUnavailable, gin.H{
       "code":       503,
       "error_code": apperrors.ErrorCodeInternal,
       "message":    "database is unhealthy",
       "request_id": rid,
       "details":    nil,
   })
   ```

2. **新增导入**: 添加 `apperrors "todo-api/internal/errors"` 以复用 `ErrorCodeInternal` 常量，确保 error_code 值与规范一致。

3. **数据来源**: `request_id` 从 Gin context 中读取（由 `LoggerMiddleware` 生成并注入），其他字段硬编码符合规范定义。

### 编译验证

```
cd backend
go build ./...   → OK
go vet ./...     → OK
```