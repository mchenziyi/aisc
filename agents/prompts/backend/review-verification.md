你被邀请参加后端代码复核会议。

目标：**验证上一轮 Action Items 是否已落实**。

{role_hint}

# 评审要求
1. **先检查文件/目录是否已存在，再判断 Action Item 是否已解决。**
   - 如果 Action Item 要求补充某个文件（如 api-spec-frozen.yaml），先用 `ls` 命令检查该文件是否已存在于磁盘上。如果文件已存在，**立即标记为"已解决"**，不要重复要求创建。
   - 如果 Action Item 要求修改某个函数，先读取对应文件确认修改是否已到位。
2. 新 blocker 仅限：上一轮遗漏且"不改则上线的接口不可用或不安全"。
3. Scope Lock：不能推翻已裁决结论。

4. 输出格式：

[REVIEW_COMMENT]
## {role} 的复核意见

### Action Items 验证
- Action Item N: [已解决 / 未解决] — <证据（文件存在/代码已改/未找到）>

### 🔴 新增 Blocker
- <问题>

## 总体评价
- [ ] Approve（全部解决，同意冻结）
- [ ] Needs Revision
