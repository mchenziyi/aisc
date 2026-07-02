你被邀请参加第 1 轮需求评审会议。

这是 PRD 的**首次全量深度评审**。你的目标是尽可能一次性暴露所有问题，不要保留到下一轮。

{role_hint}

# 评审要求
1. 完整阅读 PRD
2. 从你的专业视角逐一审查
3. 每个意见标注严重程度：
   - 🔴 blocker：如果不解决，下一阶段（API 设计）必然失败或产出的 API 不可用
   - 🟡 important：应该改，但不阻塞下一阶段启动
   - 🟢 suggestion：建议改，仅供参考
   - ✅ approve：没有问题

4. **关键规则**：
   - blocker 必须是"不改就无法进入 API 设计"的问题，而不是"可以更好"
   - 功能范围争议（如 X 功能该不该进 MVP）不属于 blocker——由 PM 决定
   - 尽量一次性列出所有问题，本轮结束后不再接受新的 blocker（除非后续发现真正阻断级问题）

5. 输出格式：

[REVIEW_COMMENT]
## {role} 的评审意见

### 🔴 Blocker
- <具体问题>：<为什么是 blocker，不改会导致什么后果>

### 🟡 Important
- <具体问题>

### 🟢 Suggestion
- <建议>

### ✅ Approved
- <确认没问题的部分>

## 总体评价
- [ ] Approve（没有 blocker，同意冻结）
- [ ] Needs Revision（有 blocker 或 important 问题）
- [ ] Reject（方向完全不对）
