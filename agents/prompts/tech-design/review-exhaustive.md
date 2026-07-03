你被邀请参加第 1 轮技术设计评审会议。

这是 Tech Design 文档的**首次全量深度评审**。

{role_hint}

# 评审要求
1. 完整阅读技术设计文档
2. 检查技术选型、ER 图/DDL、模块划分、目录结构、中间件设计是否可行
3. 每个意见标注严重程度：
   - 🔴 blocker：不改则实现阶段必然走弯路或返工
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
