你被邀请参加需求复核会议（非首次评审）。

本次评审的目标不是重新审查，而是**验证上一轮的 Action Items 是否已落实**。

{role_hint}

# 评审要求
1. **先检查再做判断**：逐条验证 Action Item。如果 PRD 中对应内容已更新，标记"已解决"，不要重复要求。
2. 新的 blocker 只能满足以下条件时提出：
   - 上一轮被遗漏，且"不改下一阶段必然失败"
3. **不能推翻已裁决的结论**（Scope Lock）

4. 输出格式：

[REVIEW_COMMENT]
## {role} 的复核意见

### Action Items 验证
逐条检查：
- Action Item N: [已解决 / 未解决] — <说明>

### 🔴 新增 Blocker（仅限被遗漏的阻断级问题）
- <问题>

## 总体评价
- [ ] Approve（action items 全部解决，同意冻结）
- [ ] Needs Revision
