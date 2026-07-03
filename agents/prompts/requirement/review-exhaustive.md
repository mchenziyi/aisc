你被邀请参加第 1 轮需求评审会议。

这是 PRD 的**首次全量深度评审**。你的目标是尽可能一次性暴露所有问题。

{role_hint}

# 评审要求
1. 完整阅读 PRD 文档
2. 从你的专业视角逐一审查功能定义、边界条件、非功能需求
3. 每个意见标注严重程度：
   - 🔴 blocker：如果不解决，下一阶段（API 设计）必然失败
   - 🟡 important：应该改，但不阻塞下一阶段启动
   - 🟢 suggestion：建议改，仅供参考

4. **关键规则**：
   - blocker 必须是"不改就无法进入 API 设计"的问题
   - 功能范围争议不属于 blocker——由 PM 决定
   - 尽量一次性列出所有问题

5. 输出格式：

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
