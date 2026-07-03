## 修改总结

### 1. 修正 go.mod 中 Go 版本为有效版本
- **文件**: `backend/go.mod`
- **修改**: 将 `go` 指令设为 `go 1.25.0`（经 `go mod tidy` 确认，gin v1.12.0 要求 Go >= 1.25.0，此版本在当前 Go 1.25.5 工具链下有效且可编译）

### 2. 修复数据库迁移失败时服务器仍继续启动的阻断问题
- **文件**: `backend/cmd/server/main.go`
- **修改**: 将
  ```go
  log.Printf("Warning: migration error (non-fatal): %v", err)
  ```
  改为
  ```go
  log.Fatalf("Migration failed: %v", err)
  ```
  确保 `RunMigrations` 返回错误时进程直接退出，不会在 schema 不完整的情况下继续运行

### 3. 运行 `go mod tidy` 清理未使用的间接依赖
- 成功执行 `go mod tidy`，移除了未使用的间接依赖（如 go.sum 中不再需要的条目）
- 保留的间接依赖（如 `quic-go`、`mongo-driver/v2` 等）均为 gin v1.12.0 的传递依赖，属于必需依赖

### 验证结果
- `go build ./...` — ✅ 编译通过
- `go vet ./...` — ✅ 无警告