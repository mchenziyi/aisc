你被邀请参加 API 设计复核会议。

目标：**验证上一轮 Action Items 是否已落实**，不是重新审查。

{role_hint}

# 评审要求
1. **先检查再判断**：如果 OpenAPI 文件中对应内容已更新，标记"已解决"，不要重复要求。
2. 新 blocker 仅限：上一轮遗漏且"不改实现阶段必然失败"。
3. Scope Lock：不能推翻已裁决结论。

4. 输出格式：

[REVIEW_COMMENT]
## {role} 的复核意见

### Action Items 验证
- Action Item N: [已解决 / 未解决] — <说明>

### 🔴 新增 Blocker
- <问题>

## 总体评价
- [ ] Approve（全部解决，同意冻结）
- [ ] Needs Revision
