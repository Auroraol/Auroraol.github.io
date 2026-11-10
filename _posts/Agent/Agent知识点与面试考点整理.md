# Agent 知识点与面试考点整理

> 目标：把网上已经有人整理的中文教程、面试题和工程实践，收敛成一份可复习的 Agent 知识地图。
>
> 范围：Agent 基础、Agent Loop、查询改写（Query Rewrite）、召回与重排、RAG、GraphRAG、Agent Graph/LangGraph、工具调用、记忆、评测、安全和生产化。
>
> 资料原则：按要求只收录中文教程、中文面试整理和公开项目文档；不下载、不把论文作为学习资料。网页内容会持续更新，面试前建议打开原文复核。

## 0. 先把几个容易混的词分开

你提到的“改写、召回、AskGraph”，在 Agent 面试里通常会落到下面几组概念：

| 名称 | 它解决什么问题 | 所在位置 |
|---|---|---|
| Query Rewrite（查询改写） | 把口语、缩写、指代不清的问题改成更容易检索或执行的查询 | 检索前 |
| Retrieval / 召回 | 从候选库中尽可能找回可能相关的文档、Chunk、实体或工具 | 检索阶段 |
| Rerank（重排） | 在较小候选集里，把“相关”重新排成“最能回答问题” | 召回后 |
| GraphRAG | 用实体、关系、路径或社区摘要补足跨文档、多跳和全局问题 | 知识检索侧 |
| Agent Graph / LangGraph | 用 State、Node、Edge 编排 Agent 的执行、分支、循环、暂停和恢复 | 执行编排侧 |
| AskGraph | 不是中文 Agent 领域中统一的标准术语；可能指某个代码图谱问答项目，也可能是对 LangGraph/GraphRAG 的口误 | 需要结合上下文判断 |

### AskGraph 的特别说明

公开搜索中较明确的 `askgraph` 项目是一个英文的本地代码库问答工具，核心思路是：

```text
代码解析 / 结构图
    + 语义向量检索
    + 标识符或符号名检索
    + 图邻居扩展
    + Git 历史溯源
    -> 给 Agent 提供低 Token、结构化的代码上下文
```

它更接近“代码库知识图谱 + 混合检索 + Agent 工具”，不是 LangGraph，也不是通用 GraphRAG 的中文教程名称。

如果你说的是这个项目，可参考：[askgraph GitHub](https://github.com/bakhliustov/askgraph)。该链接只用于核对项目含义，不作为本文的中文教程来源。

---

## 1. 一张总图：Agent 到底在做什么

一个可落地的 Agent，不是“多写几句 Prompt”，而是一个带状态的决策和执行循环：

```text
用户请求
  -> 意图识别 / 上下文补全 / 查询改写 / 任务拆解
  -> 决定：直接回答、检索、调用工具，还是进入子流程
  -> 执行并读取结果
  -> 更新 State / Memory
  -> 继续、重试、降级、转人工或结束
```

如果走 RAG 分支：

```text
原始 Query
  -> Query Rewrite
  -> Metadata 预过滤
  -> 向量召回 + 关键词召回 + 其他召回
  -> 去重、融合、Rerank
  -> 上下文压缩与排序
  -> LLM 基于证据生成答案
```

如果走工具或工作流分支：

```text
意图判断
  -> 工具选择
  -> 参数校验
  -> 服务端鉴权
  -> 工具执行
  -> 解析成功 / 空结果 / 失败 / 部分成功
  -> 重试、补救、人工确认或结束
```

面试回答任何 Agent 问题，都可以顺着五个问题展开：

1. 这个能力位于链路哪一层？
2. 输入和输出是什么？
3. 哪些步骤由模型决定，哪些步骤由代码控制？
4. 失败、超时、权限错误和空结果怎么处理？
5. 如何通过指标证明改动有效？

---

## 2. Agent 基础知识点

### 2.1 Agent 的组成

常见的面试表达是：

```text
Agent = LLM + Planning + Memory + Tools + Loop + Stop Conditions
```

分别理解为：

- `LLM`：理解意图、做判断、生成结构化动作或答案。
- `Planning`：决定目标、步骤和下一步行动。
- `Memory`：保存当前任务状态、会话信息和可复用的长期事实。
- `Tools`：查询数据库、搜索、调用 API、读写文件或执行代码。
- `Loop`：执行“判断 -> 行动 -> 观察 -> 再判断”的多轮闭环。
- `Stop Conditions`：任务完成、达到最大步数、超时、预算耗尽、需要人工确认或不可恢复失败。

Agent 和普通 Chatbot 的关键区别，不是“有没有大模型”，而是是否会围绕目标持续读取外部结果并采取行动。

### 2.2 Chatbot、Workflow、Agent 的区别

| 类型 | 谁决定路径 | 典型特点 | 适合场景 |
|---|---|---|---|
| Chatbot | 主要由一次模型调用完成 | 对话和内容生成 | FAQ、闲聊、简单问答 |
| 传统程序 | 代码完全决定 | 确定性强、易测试 | 订单状态、扣库存、计费 |
| Workflow | 开发者预先定义节点和分支 | 稳定、可观测、易回放 | 审批、发布、固定处理流程 |
| Agent | 模型参与决定下一步 | 灵活，但有循环和不确定性 | 故障排查、开放式研究、动态工具选择 |
| Agentic Workflow | 主流程由图或代码控制，局部交给 Agent | 兼顾稳定性和动态性 | 企业级长任务、复杂客服流程 |

一个实用判断：

> 执行路径能提前写出来，优先 Workflow；路径依赖中间证据、提前写不死，再引入 Agent。生产系统常见形态是“确定性流程 + 局部 Agent 决策”。

### 2.3 Agent Loop

典型循环：

```text
读取当前 State
  -> LLM 决定：直接回答 / 调工具 / 追问 / 结束
  -> 执行动作
  -> 记录 Observation
  -> 更新 State
  -> 判断是否满足结束条件
```

工程上至少要设置：

- 最大迭代次数；
- 单次和总任务超时；
- Token 或费用预算；
- 重复工具调用检测；
- 无状态进展检测；
- 高风险动作的人工确认；
- 可恢复的 Checkpoint。

### 2.4 常见 Agent 范式

| 范式 | 思路 | 优点 | 风险 / 适用边界 |
|---|---|---|---|
| ReAct | 推理一步、行动一步、读取观察后继续 | 能根据实时结果动态调整 | 调用轮次多、延迟高、容易循环 |
| Plan-and-Execute | 先生成计划，再逐步执行 | 长任务有全局结构 | 计划错了会放大错误，动态调整较弱 |
| Reflection | 生成后检查、批评并修正 | 适合质量要求高的输出 | 增加调用和成本，反思不代表一定正确 |
| Router | 先识别意图，再路由到专用链路 | 适合多业务入口 | 路由错会导致整条链路走偏 |
| Multi-Agent | 不同 Agent 分工协作 | 可并行、职责清晰 | 通信、状态一致性、调试成本上升 |
| Graph Agent | 显式定义节点、边、状态和循环 | 易审计、易暂停恢复 | 图设计和状态建模成本更高 |

面试中不要只背定义，要能回答“为什么选它”：

- 路径未知、依赖实时证据：ReAct。
- 步骤多但总体结构清晰：Plan-and-Execute。
- 结果需要多轮校验：在 ReAct 或计划执行外叠加 Reflection。
- 任务天然按专业角色拆开：Multi-Agent。
- 需要分支、循环、暂停、恢复、人工审批：Graph / LangGraph。

### 2.5 Memory 与 Context Engineering

记忆至少分三类：

| 类型 | 生命周期 | 例子 | 设计重点 |
|---|---|---|---|
| 工作记忆 | 当前任务 | 当前计划、工具结果、待办项 | 结构化、可恢复、控制长度 |
| 会话记忆 | 多轮对话 | 已确认条件、用户偏好 | 摘要、时间、来源、冲突处理 |
| 长期记忆 | 跨会话 | 稳定偏好、历史事实、知识库 | 写入门槛、可删除、权限和版本 |

“有记忆”不等于把全部聊天记录永久塞回 Prompt。更稳的做法是：

```text
最近必要原文
  + 结构化会话摘要
  + 当前任务 State
  + 按需检索的长期记忆
  + 工具和证据的结构化结果
```

`Prompt Engineering` 更关注指令怎么写；`Context Engineering` 管理范围更大，包括规则、工具描述、记忆、会话状态、外部观察、Token 预算和上下文裁剪。

---

## 3. Query Rewrite：改写不是润色，而是让下游可执行

### 3.1 为什么要改写

用户问题经常存在：

- 口语化：`钱能退吗？`
- 指代缺失：`上次那个报错怎么处理？`
- 缩写或术语不一致：用户说“客服没回”，文档写“无应答策略”；
- 多意图混合：`查订单，再看看退款规则；`
- 条件隐藏：时间、平台、版本、租户、权限没有明确抽取。

改写的目标是提高下游有效召回和工具执行成功率，同时保持原始意图不变。

### 3.2 常见改写策略

| 策略 | 做法 | 适用场景 | 主要风险 |
|---|---|---|---|
| 规范化 | 补全同义词、缩写、领域术语 | 口语和业务术语不一致 | 把原意改偏 |
| 上下文补全 | 从必要历史对话中补齐指代 | `它为什么失败` | 带入无关历史 |
| Multi-Query | 生成多个等价查询并合并结果 | 一种说法可能漏召回 | 查询数量和成本上升 |
| Query Decomposition | 拆成独立子问题 | 多意图、多跳问题 | 聚合规则复杂 |
| Step-back | 先退到更通用的背景问题 | 问题太细、缺背景 | 召回过泛 |
| HyDE | 先生成假设性答案或文档，再用其检索 | Query 和文档表达差异大 | 假设答案带来偏差 |
| Self-Query | 抽取过滤条件、排序和结构化字段 | 时间、版本、类别、租户过滤 | 结构化条件解析错误 |

### 3.3 推荐的安全输出结构

生产环境不要只让模型返回一段改写后的字符串，建议保留原始 Query，并把意图、候选查询和过滤条件结构化：

```json
{
  "original_query": "上次那个订单为什么还没退款？",
  "intent": "查询退款状态和原因",
  "rewritten_queries": [
    "查询当前会话关联订单的退款状态和最近一次状态变化原因",
    "查询订单退款处理流程、失败原因和当前处理节点"
  ],
  "filters": {
    "order_id": null,
    "platform": null,
    "time_range": null,
    "version": null
  },
  "missing_slots": ["订单唯一标识"],
  "must_preserve": ["当前订单", "退款", "原因"]
}
```

关键原则：

1. 原始 Query 永远保留，改写失败可以回退。
2. 人名、订单号、错误码、版本号、时间、否定词和权限边界单独校验。
3. 查询文本和过滤条件分开，关键词交给检索器，条件交给 Metadata 过滤器。
4. 改写不是事实补全，缺少信息时标记 `missing_slots`，不要编造。
5. 原始 Query 和改写 Query 一起召回，避免改写模型误解原意。

### 3.4 什么时候不要改写

对于明确的错误码、订单号、SKU、接口路径、代码符号和结构化查询，改写可能损害精确匹配。可以先做意图判断：只有检测到口语、指代、多意图或领域表达不一致时才触发改写。

### 3.5 改写效果怎么评估

不要评价“改写文本是否更通顺”，要看下游结果：

- Recall@K、Hit Rate@K；
- MRR、nDCG；
- 正确证据进入候选池的比例；
- 工具参数校验通过率；
- 最终答案正确率和引用准确率；
- 延迟、Token 和额外模型调用成本；
- 实体、否定词、时间和权限条件保留率。

---

## 4. RAG 与召回：从“能找”到“找得准、能回答”

### 4.1 标准 RAG 链路

```text
文档采集
  -> 解析与清洗
  -> 结构化与 Chunking
  -> Metadata 标注
  -> Embedding / 倒排索引
  -> Query Rewrite 与过滤
  -> 初始召回
  -> 结果融合与去重
  -> Rerank
  -> 上下文压缩、排序和拼接
  -> LLM 生成答案与引用
  -> 评测、回放、灰度和监控
```

RAG 的失败不一定是模型问题，可能发生在任意一环：解析错、Chunk 切断语义、Metadata 过滤过严、关键词没命中、正确证据排得太后、上下文被截断，或者模型没有遵守证据边界。

### 4.2 文档与 Chunking

Chunk 没有万能大小，核心标准是：一个 Chunk 是否包含回答所需的最小完整语义。

常见策略：

- 固定长度 + Overlap：实现简单，适合做基线。
- 按标题、段落、列表、表格和代码边界切分：更保留文档结构。
- Parent-Child Chunk：小子块负责命中，父块负责补上下文。
- 语义切分：根据主题变化确定边界，但离线成本和更新成本更高。
- FAQ / 工单按问答对或会话切分：避免问题和答案分离。
- 代码按文件、类、方法、注释和调用关系切分：不要只按 Token 截断。

每个 Chunk 建议至少保留：

```text
doc_id、chunk_id、title、section_path、source、page、version、language、tenant_id、acl、updated_at
```

表格、图片、PDF 页眉页脚、跨页内容和标题层级如果解析错，后面的 Embedding 和召回都无法完全补救。

### 4.3 Embedding、BM25 和混合召回

| 召回方式 | 强项 | 弱项 | 常见场景 |
|---|---|---|---|
| 向量检索 | 语义相近、同义表达 | 可能漏掉错误码、SKU、版本号 | FAQ、自然语言问题 |
| BM25 / 倒排 | 精确词项、数字和专有名词 | 对同义表达不敏感 | 日志、代码、错误码、型号 |
| 混合召回 | 同时覆盖语义和精确词 | 融合参数更复杂 | 企业知识库的默认方案 |
| 多路召回 | 不同索引或策略并行 | 成本、去重、融合复杂 | 高价值、复杂或多领域问题 |
| 图检索 | 实体关系和路径扩展 | 图构建和维护成本高 | 多跳关系、影响分析 |

一个常见的工程链路是：

```text
Metadata 预过滤
  -> 向量召回 Top-K
  + BM25 召回 Top-K
  + 结构化 / 图检索候选
  -> RRF 或归一化加权融合
  -> 去重、相邻块合并
  -> Rerank
  -> 上下文裁剪
```

RRF 只看候选名次，不需要强行把 BM25 分数和余弦相似度换算到同一尺度。

### 4.4 Rerank 与上下文组装

召回和重排是两个目标：

- 召回追求覆盖率：正确证据尽量进入候选池；
- 重排追求精度：把最能回答问题的候选排到前面。

因此建议区分三个参数：

```text
recall_top_k：粗召回候选数，例如 30～100
rerank_top_n：重排后保留数，例如 5～10
context_top_n：最终送入模型的证据数，例如 3～6
```

`Top-K` 越大不一定越好。候选太少会漏证据，候选太多会提高重排和 Token 成本，并把噪声塞进上下文。

上下文组装至少要做：

- 去重；
- 保持同一文档的相邻片段顺序；
- 优先放最相关、版本最新且有权限的证据；
- 标记来源、页码、版本和时间；
- 对长片段做查询相关压缩或结构化抽取；
- 让 Prompt 明确“只基于上下文回答，证据不足就说明无法判断”。

### 4.5 “召回不到”如何排查

按下面顺序排查，避免一上来就换 Embedding：

```text
文档是否真的入库
  -> 解析是否丢标题、表格、代码和页码
  -> Chunk 是否切断了完整语义
  -> Metadata / ACL 是否过滤过严
  -> Query 是否需要改写、分解或补充上下文
  -> 是否缺少 BM25 / 精确词路由
  -> recall_top_k 是否过小
  -> 正确证据是否已在候选池
       -> 不在：查索引、切分、查询和召回
       -> 在但没进上下文：查融合、Rerank、去重和 Top-N
       -> 已进上下文但答错：查上下文顺序、截断、冲突和生成约束
```

### 4.6 RAG 面试必考指标

| 指标 | 看什么 |
|---|---|
| Recall@K / Hit Rate@K | 正确证据是否进入前 K |
| Precision@K | 前 K 中相关结果的比例 |
| MRR | 第一个正确结果排得多靠前 |
| nDCG | 不同相关等级下的排序质量 |
| Context Recall | 回答所需证据是否找全 |
| Context Precision | 送进模型的上下文是否足够纯 |
| Faithfulness / Groundedness | 答案是否被证据支撑 |
| Answer Relevancy | 答案是否真正回应问题 |
| Citation Accuracy | 引用是否支撑对应结论 |
| P95 延迟 / Token / 成本 | 生产可用性 |

评测一定要区分检索质量和生成质量。只看最终答案，定位不了到底是没找着、排错了，还是模型没用好证据。

---

## 5. GraphRAG、Agent Graph、LangGraph 到底怎么分

### 5.1 GraphRAG 解决什么问题

普通向量 RAG 的主要检索对象是文本 Chunk；GraphRAG 会额外组织实体、关系、路径、社区或主题摘要。

它更适合：

- 从一个实体沿关系找上下游；
- 跨文档做多跳推理；
- 查询影响范围、依赖链、负责人和事故关联；
- 对大量报告、工单或访谈做全局主题归纳；
- 需要结构化过滤、路径解释或关系溯源的场景。

它不一定适合：

- 数据量小、问题只是查一条规则；
- 向量召回本身还没有做稳；
- 文档质量差、实体名称不统一；
- 实时更新要求很高但没有增量图维护能力；
- 团队没有图建模、消歧和评测能力。

### 5.2 GraphRAG 的核心对象

```text
实体 Entity：系统、接口、人员、订单、故障、政策条款
关系 Relationship：依赖、负责、影响、导致、属于、引用
属性 Property：时间、版本、置信度、权限、来源
社区 Community：连接紧密的一组主题或实体
路径 Path：从起点实体到目标实体的关系链
```

构建流程通常是：

```text
文档解析
  -> Chunk / TextUnit
  -> 实体抽取
  -> 关系抽取
  -> 实体消歧与归一
  -> 图存储
  -> 社区发现 / 摘要（可选）
  -> 图索引 + 向量索引 + 全文索引
```

查询时可按问题路由：

| 问题类型 | 更适合的策略 |
|---|---|
| 局部事实 | 向量检索或普通混合检索 |
| 已知实体的上下游关系 | 局部图检索 |
| 多跳问题 | 图遍历 + 向量补充原文证据 |
| 全局主题归纳 | 社区摘要或主题聚合 |
| 精确过滤 | 图查询、结构化查询或固定模板 |

### 5.3 GraphRAG 的工程难点

- 同一个实体可能有多个别名，实体消歧不能全靠模型猜。
- 关系方向错一个字，可能让整条依赖链反过来。
- 每条关系应保存来源文档、原文位置、置信度、版本和抽取模型版本。
- 文档更新可能同时影响实体、关系、社区和摘要，增量失效范围要可追踪。
- 社区摘要可能混入无权限文档内容，不能只在最终答案阶段做权限过滤。
- 图检索结果仍需回到原文证据，不能把图谱抽取结果当成绝对事实。

### 5.4 推荐的 GraphRAG 落地顺序

```text
第一阶段：把文档解析、Chunk、Metadata、混合召回、Rerank 和引用做稳
第二阶段：收集关系型 Badcase，确认是否真的需要图
第三阶段：只抽取少量核心实体和高价值关系
第四阶段：用图做候选扩展，并保留原文证据
第五阶段：有明确全局问题时，再引入社区发现和摘要
第六阶段：按问题类型做向量 / 关键词 / 图 / 结构化查询路由
```

### 5.5 Agent Graph / LangGraph 解决什么问题

Agent Graph 关注的是“任务怎么跑”，不是“知识怎么连”。

| 概念 | 含义 |
|---|---|
| State | 节点之间共享的任务状态，例如 Query、计划、证据、工具结果、重试次数 |
| Node | 执行一个可观察步骤，可以是 LLM、工具或普通代码 |
| Edge | 决定从一个节点流向哪个节点 |
| Conditional Edge | 根据 State 或运行结果选择分支 |
| Loop | 回到前面节点继续处理，但必须有退出边界 |
| Reducer | 并行节点写同一字段时的合并规则 |
| Checkpoint | 保存中间状态，用于暂停、恢复、回放和人工介入 |
| Subgraph | 把一段复杂流程封装成可复用子图 |

以“审核不通过最多重试三次”为例：

```text
生成 -> 审核
         | 通过
         v
        结束
         |
         | 不通过且 retry_count < 3
         v
       修改 -> 审核
         |
         | retry_count >= 3
         v
      转人工 / 返回部分结果
```

LangGraph 的中文文档强调它是一个面向长时间运行、有状态 Agent 的低层编排框架，重点能力包括持久执行、人机协作、短期和长期记忆、流式处理、断点和部署。

### 5.6 GraphRAG 与 LangGraph 的一句话区分

> LangGraph 解决“Agent 如何按状态和节点执行”；GraphRAG 解决“知识如何按实体和关系组织、检索与扩展”。两者可以组合，但不是同一个东西。

一个组合例子：

```text
LangGraph：查询理解 -> Query Rewrite -> 召回路由 -> 图检索 -> 证据审核 -> 回答
GraphRAG：在图检索节点里提供实体、关系、路径和来源证据
```

---

## 6. 工具调用、MCP、Skills 与可靠性

### 6.1 Function Calling

Function Calling 让模型输出结构化的工具调用意图，包括工具名和 JSON Schema 参数；模型负责“建议调用什么、参数是什么”，真正的权限校验和执行仍在服务端。

面试要点：

- 工具名称、用途、输入、必填字段和不适用条件要写清楚；
- 参数必须做 Schema 校验、枚举校验、范围校验和业务校验；
- 工具返回值尽量结构化，区分成功、空结果、失败和部分成功；
- 不能把 Function Calling 当成权限系统；
- 写操作要考虑幂等、二次确认、审计和回滚。

### 6.2 MCP

MCP 更偏工具和上下文能力的统一接入协议。常见原语包括：

| 原语 | 作用 |
|---|---|
| Tools | 可主动调用的函数，例如查库、发消息、执行脚本 |
| Resources | 按需读取的只读资源，例如文件、日志或数据库记录 |
| Prompts | 可复用的提示模板 |

一句话区分：

> Function Calling 解决模型如何表达工具调用；MCP 解决工具如何被 Host / Client 发现和接入。

### 6.3 Skills

Skill 更像“任务手册”：把完成一类任务所需的流程、约束、参考资料和脚本封装起来，按需加载。

```text
Function Calling = 调用语言 / 数据格式
MCP = 工具接入协议 / 工具箱
Skill = 任务流程与操作手册
```

Skills 适合沉淀团队流程和工程经验；工具数量很多时，按需加载还能降低上下文和路由压力。

### 6.4 工具失败、重试和幂等

适合重试的通常是：网络抖动、临时限流、明确的 5xx、短暂不可用。

不应盲目重试的通常是：参数错误、权限不足、业务状态不允许、不可逆写操作和幂等性不明确的请求。

工具返回结果可以统一成：

```json
{
  "ok": false,
  "code": "UPSTREAM_TIMEOUT",
  "data": null,
  "retryable": true,
  "request_id": "...",
  "source": "order-service"
}
```

注意：模型的“重新思考”不等于可以重复执行副作用操作。扣款、下单、发消息、删数据等操作需要业务幂等键和服务端状态机。

---

## 7. 生产化：安全、可观测、可评测

### 7.1 安全考点

重点防护：

- Prompt Injection：把外部文档、网页、邮件、工具结果当作不可信数据，不能让它覆盖系统指令。
- 越权访问：权限由服务端鉴权，不能让模型自己决定“我有权限”。
- RAG 泄露：召回前或召回时做租户、用户、角色和 ACL 过滤。
- 高风险工具：外发、写库、删数据、执行代码、转账和发布操作需要风险分级、确认和审计。
- 记忆污染：长期记忆写入要有门槛、来源、版本、可删除能力和权限边界。
- 敏感数据：日志和 Trace 脱敏，不保存密钥、Token、完整身份证号等不必要数据。

安全边界最好画成：

```text
模型提出意图
  -> 服务端校验身份、权限、参数和业务状态
  -> 风险分级 / 幂等检查 / 人工确认
  -> 执行工具
  -> 记录审计结果
```

### 7.2 可观测性

一次 Agent 任务至少建议关联：

```text
trace_id、session_id、task_id、tenant_id
模型名和版本、Prompt / Skill 版本
Query Rewrite 前后内容及过滤条件
召回文档 ID、分数、Rerank 分数和最终证据
工具名、参数摘要、返回码、耗时、重试次数
Graph 节点、边、State 版本和终止原因
最终答案、引用、质量检查和人工接管原因
```

敏感字段要脱敏，保存足够的结构化证据即可，不要把完整隐私数据原样打进日志。

### 7.3 Agent 评测分层

| 层级 | 评测问题 | 例子 |
|---|---|---|
| 改写层 | 是否保留原意和关键约束 | 实体一致率、条件保留率 |
| 召回层 | 正确证据是否进入候选池 | Recall@K、MRR |
| 重排层 | 正确证据是否排到前面 | Precision@N、nDCG |
| 生成层 | 答案是否正确、忠实、有引用 | Correctness、Faithfulness、Citation Accuracy |
| Agent 层 | 是否完成任务且没有越权 | Task Success、工具成功率、转人工率 |
| 系统层 | 是否满足线上约束 | P95、成本、错误率、SLO |

评测集要覆盖：高频问题、线上 Badcase、口语和缩写、多意图、多跳问题、精确错误码、无答案、冲突版本、权限隔离、工具超时和高风险操作。

每次改动都记录版本：文档解析器、Chunk 策略、Embedding、索引参数、Rewrite Prompt、Rerank、生成 Prompt、Skill 和评测集。

---

## 8. 高频面试题清单

### 8.1 Agent 基础

1. 什么是 Agent？它和 Chatbot、Workflow 的区别是什么？
2. `Agent = LLM + Planning + Memory + Tools` 怎么理解？
3. Agent Loop 的完整流程是什么？
4. ReAct、Plan-and-Execute、Reflection 分别适合什么场景？
5. 什么时候使用单 Agent，什么时候使用 Multi-Agent？
6. 为什么生产环境通常采用“Workflow + 局部 Agent”？
7. Agent 如何避免重复调用、路径震荡和死循环？
8. 长任务中断后如何恢复？
9. 短期记忆、会话记忆和长期记忆有什么区别？
10. Prompt Engineering 和 Context Engineering 有什么区别？

### 8.2 Query Rewrite 与 RAG

11. 为什么要做 Query Rewrite？
12. Query Rewrite 和 Query Decomposition 有什么区别？
13. Multi-Query、Step-back、HyDE、Self-Query 怎么选？
14. 什么情况下不应该改写？
15. 如何防止改写导致语义漂移、否定词丢失和实体被替换？
16. RAG 的完整链路是什么？
17. Chunk 大小、Overlap 和语义边界如何取舍？
18. Parent-Child Chunk 解决什么问题？
19. Embedding、BM25、向量检索和关键词检索的区别是什么？
20. 为什么要做 Hybrid Search？怎么融合？
21. 多路召回具体怎么做？如何去重和合并？
22. Rerank 解决什么问题？为什么不能直接把向量 Top-K 送给模型？
23. `recall_top_k`、`rerank_top_n`、`context_top_n` 为什么要分开？
24. 召回不到正确文档，你的排查顺序是什么？
25. 正确文档已经召回但答案仍然错误，怎么定位？
26. RAG 如何处理多版本文档和知识冲突？
27. RAG 如何做权限过滤？
28. 如何评估 RAG？检索指标和生成指标怎么拆？

### 8.3 GraphRAG 与 LangGraph

29. GraphRAG 解决了普通向量 RAG 的什么问题？
30. 为什么说 Chunk 可能是信息孤岛？
31. GraphRAG 的实体、关系、社区和路径分别是什么？
32. Local Search 和 Global Search 有什么区别？
33. 什么情况下不应该上 GraphRAG？
34. GraphRAG 的实体消歧、关系方向和增量更新怎么处理？
35. LangGraph 的 State、Node、Edge、Checkpoint 分别是什么？
36. 条件边、循环和 Reducer 怎么设计？
37. 如何给 Graph Agent 设置最大步数和安全退出？
38. Agent Graph / LangGraph 与 GraphRAG 的区别是什么？
39. 如果“AskGraph”指代码库问答项目，它和 LangGraph 有什么区别？

### 8.4 工具与生产化

40. Function Calling 的原理是什么？
41. Function Calling 和 MCP 的区别是什么？
42. MCP 的 Host、Client、Server 和 Tools / Resources / Prompts 怎么理解？
43. MCP 和 Agent Skills 的边界是什么？
44. 工具参数错误、超时、限流和权限失败分别怎么处理？
45. 哪些错误可以重试？写操作如何保证幂等？
46. 如何防止 Agent 任务幻觉：没有执行却声称完成？
47. 如何防止 Prompt Injection、越权和敏感数据泄露？
48. 如何记录一条可回放的 Agent Trace？
49. 线上延迟升高时，如何区分模型、召回、重排、工具和排队耗时？
50. 如何证明一次 Agent 优化真的有效？

---

## 9. 综合系统设计题答题模板

题目示例：设计一个企业知识库 Agent，支持问答、关系查询和部分业务操作。

建议按以下顺序回答：

1. **目标和边界**：是问答还是执行操作？哪些高风险动作必须人工确认？
2. **数据链路**：文档采集、解析、Chunk、Metadata、权限和版本。
3. **查询链路**：意图识别、必要时改写、多路召回、过滤、重排、压缩和引用。
4. **图能力**：只有关系型 Badcase 足够多时，才引入实体、关系、图查询或 GraphRAG。
5. **执行编排**：用 State、Node、Edge 表达查询、工具、审核、循环和人工接管。
6. **可靠性**：超时、重试、幂等、空结果、部分成功、降级和 Checkpoint。
7. **安全**：服务端鉴权、最小权限、Prompt Injection 防护、敏感信息处理和审计。
8. **评测指标**：改写、召回、重排、生成、工具和端到端 Task Success 分层评估。
9. **线上闭环**：Trace、Badcase、回放集、灰度、版本对比和回滚。

一个简短的高质量回答可以收束为：

> 先用确定性流程保证边界和安全，再在需要动态决策的节点引入 Agent；RAG 先把解析、Chunk、混合召回、重排和引用做稳；只有失败样本明确需要多跳关系或全局归纳时才引入 GraphRAG；用 LangGraph 一类的图编排保存 State、控制循环并支持暂停恢复；最后用分层指标和 Trace 证明效果，而不是只展示几个 Demo。

---

## 10. 一页复习清单

### Agent

- [ ] 能解释 Agent 的组成和 Loop。
- [ ] 能区分 Chatbot、Workflow、Agent、Multi-Agent、Agent Graph。
- [ ] 能比较 ReAct、Plan-and-Execute、Reflection。
- [ ] 能说明 Memory 与 Context Engineering 的边界。
- [ ] 能处理超时、重试、幂等、循环和人工接管。

### 改写与召回

- [ ] 能说清 Query Rewrite 的目的、触发条件和安全输出。
- [ ] 能区分规范化、上下文补全、多查询、分解、Step-back、HyDE、Self-Query。
- [ ] 能解释向量、BM25、Hybrid Search、多路召回、Rerank 和上下文压缩。
- [ ] 能按“入库 -> 解析 -> Chunk -> Metadata -> Query -> 召回 -> Rerank -> 上下文 -> 生成”排查 Badcase。
- [ ] 能用 Recall@K、MRR、nDCG、Faithfulness 等指标定位问题。

### Graph

- [ ] 能区分 GraphRAG 与 Agent Graph / LangGraph。
- [ ] 能解释 Entity、Relationship、Path、Community、Local / Global Search。
- [ ] 能解释 LangGraph 的 State、Node、Edge、Reducer、Checkpoint。
- [ ] 能设计循环上限、条件分支、暂停恢复和人工审批。
- [ ] 能说明 AskGraph 可能是代码图谱问答项目，而不是 LangGraph 的别名。

### 生产化

- [ ] 能讲 Prompt Injection、越权、ACL、敏感信息和高风险工具治理。
- [ ] 能设计结构化工具结果和服务端权限校验。
- [ ] 能记录 Query、召回证据、工具结果、Graph State、版本和终止原因。
- [ ] 能用离线评测、线上灰度和 Badcase 回放证明优化收益。

---

## 11. 中文教程与面试资料入口

以下资料是本次整理的主要依据，均优先选择中文教程、中文面试题或中文文档，不需要下载论文。

| 来源 | 推荐内容 | 适合怎么用 |
|---|---|---|
| [JavaGuide：AI Agent 面试题总结](https://javaguide.cn/ai/interview-questions/agent-interview-questions.html) | Agent、Memory、Context、MCP、Skills、Workflow、Graph、评测和安全考点 | 面试前快速过题 |
| [JavaGuide：AI Agent 核心概念详解](https://javaguide.cn/ai/agent/agent-basis.html) | Agent Loop、工具注册、ReAct、Plan-and-Execute、Reflection、Multi-Agent、Workflow | 建立知识框架 |
| [JavaGuide：RAG 面试题总结](https://javaguide.cn/ai/interview-questions/rag-interview-questions.html) | RAG 链路、Chunk、Embedding、召回、改写、多路召回、GraphRAG、评测 | RAG 面试题库 |
| [JavaGuide：RAG 优化](https://javaguide.cn/ai/rag/rag-optimization.html) | 数据治理、Chunk、Hybrid Search、Query Rewrite、Rerank、上下文工程、排查路径 | 召回优化实战 |
| [JavaGuide：GraphRAG 详解](https://javaguide.cn/ai/rag/graphrag.html) | GraphRAG 原理、实体关系、社区、局部/全局检索、选型和工程难点 | GraphRAG 选型与面试 |
| [小林面试笔记：Agent 面试题](https://xiaolinnote.com/ai/agent/agent_info.html) | 24 个 Agent 高频问题，覆盖架构、范式、记忆、评测、死循环、延迟和安全 | 看面试官如何追问 |
| [小林面试笔记：RAG 面试题](https://xiaolinnote.com/ai/rag/rag_info.html) | 21 个 RAG 高频问题，重点是 Chunk、Query Rewrite、多路召回、图数据库和知识冲突 | 专门补 RAG 面试 |
| [小林面试笔记：LLM 工具调用面试题](https://xiaolinnote.com/ai/tools/tools_info.html) | Function Calling、MCP、Skills、A2A、通信、Tool Routing 和错误处理 | 工具调用与协议面试 |
| [Datawhale：从零开始构建智能体](https://github.com/datawhalechina/hello-agents) | Agent 基础、ReAct、Plan-and-Solve、Reflection、框架、Memory、RAG、上下文、MCP、评测和综合案例 | 系统学习和动手实践 |
| [LangGraph 中文文档](https://langgraph.com.cn) | State、Node、Edge、持久化、记忆、人机协作、断点、多 Agent 和评估 | LangGraph 入门与复习 |
| [LangGraph 中文翻译仓库](https://github.com/jurnea/LangGraph-Chinese) | LangGraph 中文翻译内容和示例入口 | 需要查 API 或示例时使用 |

### 推荐阅读顺序

```text
JavaGuide Agent 核心概念
  -> JavaGuide Agent 面试题
  -> 小林 Agent 题库
  -> JavaGuide RAG 优化
  -> JavaGuide RAG / GraphRAG 面试题
  -> LangGraph 中文文档
  -> Datawhale Hello-Agents 做系统实践
```

如果时间很少：先看本文件第 0、1、3、4、5、8、10 节，再打开 JavaGuide 和小林对应原文做补充。

---

## 12. 最后记忆口诀

```text
先判断任务边界
再决定是否改写
改写保留原意和约束
召回追求覆盖
重排追求可回答
上下文只放证据
GraphRAG 管知识关系
LangGraph 管执行流程
工具必须服务端鉴权
循环必须有退出条件
效果必须用数据证明
```

