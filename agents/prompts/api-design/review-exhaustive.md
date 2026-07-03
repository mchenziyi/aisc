你被邀请参加第 1 轮 API 设计评审会议。

这是 OpenAPI 规范的**首次全量深度评审**。

{role_hint}

# 评审要求
1. 完整阅读 OpenAPI 规范
2. 逐接口检查：
   - 路径和方法是否和 PRD 功能一一对应
   - 请求/响应 schema 是否完整（字段、类型、必填/可选）
   - 错误码是否覆盖 400/401/403/404/409/422/429/500
   - 认证标注是否正确（公开接口 vs Bearer Token）
3. 每个意见标注严重程度：
   - 🔴 blocker：不改则实现阶段必然产生不可用的接口
   - 🟡 important：应该改但不阻塞开发启动
   - 🟢 suggestion：建议优化

4. 输出格式：

[REVIEW_COMMENT]
## {role} 的评审意见

### 🔴 Blocker
- <问题>：<为什么是 blocker>

### 🟡 Important
- <问题>

### 🟢 Suggestion
- <建议>

## 总体评价
- [ ] Approve（无 blocker，同意冻结）
- [ ] Needs Revision
- [ ] Reject
