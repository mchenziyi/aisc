所有修改已完成，编译和静态检查均通过。以下是修改总结：

---

## 修改总结

### Action Item 完成情况

| # | Action Item | 状态 | 说明 |
|---|------------|------|------|
| 1 | 补充 `docs/api-spec-frozen.yaml` 规范文件以完成契约校验 | ✅ 完成 | 修复 PATCH `/todos/{todo_id}` 的 `todo_id` 参数缺少 `minimum: 1` 约束的问题，使其与代码实现及 GET/DELETE 端点的约束保持一致 |

### 详细变更

**文件**: `docs/api-spec-frozen.yaml`

- **PATCH `/todos/{todo_id}` 参数修正**: 在 `todo_id` 的 schema 中补充 `minimum: 1`，与代码中 `UpdateTodo` handler 的 `if todoID < 1` 正数校验逻辑一致，同时也与 GET 和 DELETE 端点的参数约束保持一致。

### 编译验证

```
go build ./... → OK
go vet   ./... → OK
```

### 说明

当前代码（v2）中 11 项功能修改已全部完成，本次仅对 API 规范文件 `docs/api-spec-frozen.yaml` 做了补充修正，使其完整反映代码实现中的参数约束，满足契约校验要求。至此所有 Action Items 已全部关闭，可进入冻结阶段。