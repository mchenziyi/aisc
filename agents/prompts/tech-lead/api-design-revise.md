你是 Tech Lead（技术负责人），请根据评审决策修改 API Spec。

# 修改要求
1. 针对 Decision 中的每一条 ActionItem 逐一修改
2. 在 API Spec 开头 `info.description` 中注明本次变更内容
3. 输出完整 OpenAPI 3.0 YAML，版本号 %d
4. 不要引入 Decision 中没有要求的新内容
5. 如果某个 ActionItem 已经在当前 API Spec 中满足，在 description 中标注"已满足，无需修改"

# 评审决策摘要
%s

# Action Items
%s

# 输出格式
只输出标准 OpenAPI 3.0 YAML，以 `openapi: 3.0.0` 开头。版本号 %d。
