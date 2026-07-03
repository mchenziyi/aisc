你被邀请参加第 1 轮后端代码评审会议。

这是后端代码的**首次全量深度评审**。你会看到完整的源代码文件。

{role_hint}

# 评审要求
1. 阅读所有代码文件
2. 重点检查：
   - **API 契约一致性**：每个接口是否严格对应 api-spec-frozen.yaml 中的路径、方法、schema
   - **未定义接口**：是否有 API Spec 未定义的端点（如 /auth/me），有则标记 blocker
   - **权限隔离**：用户是否只能操作自己的数据（SQL 是否带 user_id 过滤）
   - **错误码**：是否和 API Spec 定义一致
   - **JWT_SECRET**：配置是否强制必需（不能有默认弱密钥）
   - **技术设计一致性**：目录结构、模块划分是否和 tech-design-frozen.md 一致

3. 每个意见标注严重程度：
   - 🔴 blocker：不改则上线的接口不可用或不安全
   - 🟡 important：应该改但不阻塞上线
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
