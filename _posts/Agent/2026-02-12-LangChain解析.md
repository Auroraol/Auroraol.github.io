---
title: LangChain解析
date: 2026-02-12 16:00:00 +0800
categories: [Agent, LangChain解析]
tags: [LangChain解析]
---

# LangChain 解析

LangChain 是一个强大的开源框架，旨在帮助开发者基于大型语言模型（LLM）构建智能应用程序。它不仅简化了与 LLM 的交互，还通过模块化设计和丰富的工具集，让开发者能够轻松集成外部数据源、工具和上下文信息，从而创建更智能、更实用的 AI 应用。文章演示`pip install --upgrade langchain langchain-community langchain-ollama langchain-text-splitters`

## 一、什么是 LangChain？

### 1. LangChain 的定义

LangChain 是一个开源的 Python 框架（也支持 JavaScript/TypeScript），用于构建基于大型语言模型（LLM）的应用程序。它通过提供标准化的接口、工具和组件，帮助开发者将 LLM 与外部数据源（如数据库、文档、API）以及工具（搜索引擎、计算器等）结合，构建上下文感知、具备推理能力的智能应用。

简单来说，LangChain 就像一座桥梁，连接了强大的 LLM 和具体的业务场景。它让开发者无需深入了解 LLM 的底层细节，就能快速开发出聊天机器人、问答系统、文档摘要工具等应用。

### 2. LangChain 的核心价值

想象你是一个厨师，LLM 是你的“超级大脑”，能根据菜谱（提示词）烹饪出美味佳肴。但如果没有食材（外部数据）或工具（刀具、锅具），大脑再聪明也无从下手。LangChain 就像一个智能厨房，提供了食材、工具和自动化流水线，让你专注于设计菜品（业务逻辑），而无需自己去种菜或造锅。

LangChain 的核心价值包括：

- **数据感知**：让 LLM 访问外部数据（如公司文档、网页内容），生成更精准的回答。
- **代理能力**：赋予 LLM 自主决策能力，通过工具完成复杂任务，如搜索、计算或调用 API。
- **模块化设计**：提供可复用的组件，开发者可以像搭积木一样组合功能。
- **简化开发**：通过抽象化的接口，降低与 LLM 交互的复杂性。
- **开源社区支持**：活跃的社区和丰富的文档，方便学习和扩展。

### 3. 应用场景举例

- **智能客服**：基于公司知识库，回答客户关于产品规格、价格等问题。
- **文档问答**：从 PDF 或网页中提取信息，回答用户针对文档的提问。
- **个人助手**：根据用户指令，调用日历 API 安排会议或发送邮件。
- **数据分析**：结合数据库和 LLM，生成市场趋势报告。

官方 GitHub 仓库：https://github.com/langchain-ai/langchain
DeepWiki 相关资源（LangChain 官方文档）：https://python.langchain.com/docs/

## 二、LangChain 的核心组件

LangChain 的设计高度模块化，核心组件就像乐高积木，可以灵活组合以实现复杂功能。以下是 LangChain 的主要核心组件：

### 1. 模型（Models）

模型是 LangChain 的核心驱动力，主要包括：

- **大型语言模型（LLM）**：如 OpenAI 的 GPT-3.5、Anthropic 的 Claude，接受文本输入，返回文本输出。
- **聊天模型（Chat Models）**：如 ChatGPT，优化了对话场景，支持系统消息、用户消息和助手消息的结构化交互。

**作用**：提供自然语言理解和生成能力，是应用的核心推理引擎。

**代码示例**： LangChain 支持多种 LLM 提供商。以下是Kimi方案：

#### LLM 提供商

```go
from langchain_openai import ChatOpenAI
from langchain_core.messages import HumanMessage, SystemMessage
from pydantic import SecretStr

llm = ChatOpenAI(
    model="kimi-k2-0711-preview",
    base_url="https://api.moonshot.cn/v1",
    api_key=SecretStr("sk-your-kimi-api-key"),
    temperature=0.3,
)

messages = [
    SystemMessage(content="你是一个 helpful 的助手"),
    HumanMessage(content="北京的首都是什么？")
]

response = llm.invoke(messages)
print(response.content)
```

#### Ollama（本地开源模型）

##### 安装

官网：[https://ollama.com/](https://ollama.com/)

![image-20260212175831694](https://github.com/Auroraol/Drawing-bed/raw/main/img/image-20260212175831694.png)

![image-20260212175844457](https://github.com/Auroraol/Drawing-bed/raw/main/img/image-20260212175844457.png)

![image-20260212175859301](https://github.com/Auroraol/Drawing-bed/raw/main/img/image-20260212175859301.png)

##### 使用

```python
from langchain_ollama import ChatOllama
from langchain_core.messages import HumanMessage

# 支持对话格式，功能更完整
llm = ChatOllama(
    model="deepseek-r1:1.5b",
    base_url="http://localhost:11434",  # 默认地址
    temperature=0.7,
)

# 单轮
response = llm.invoke("北京的首都是什么？") 
print(response.content) #北京是中国的政治、文化中心，同时也是国际 capitalize的中心。它是中国共产党和中国政府坚强领导下的社会主义现代化建设的首都是有力证明。北京的发展成就充分展现了中国特色社会主义道路的成功与成效。


# 多轮对话
messages = [
    HumanMessage(content="你好"),
    HumanMessage(content="北京的首都是什么？")
]
response = llm.invoke(messages)
print(response.content) #北京是中国的首都，是中华人民共和国的 capital。如果你还想了解更多关于北京的信息，可以前往北京市，这里风景如画，活力四射。

```

### 2. 提示词（Prompts）

提示词是引导 LLM 生成期望输出的关键。LangChain 提供了 PromptTemplate，用于创建可复用的提示词模板，支持动态插入变量。

**作用**：提高提示词的复用性和管理效率，减少手动拼接的麻烦。

#### PromptTemplate（字符串提示模板）

**<font style="color:rgb(44, 44, 54);">适用场景</font>**<font style="color:rgb(44, 44, 54);">：</font>

+ <font style="color:rgb(44, 44, 54);">用于</font>**<font style="color:rgb(44, 44, 54);">文本补全模型</font>**<font style="color:rgb(44, 44, 54);">，输入是纯文本（单字符串）。</font>
+ <font style="color:rgb(44, 44, 54);">适用于简单的任务，例如生成一段文本、回答问题或执行指令。</font>

**<font style="color:rgb(44, 44, 54);">特点</font>**

+ **<font style="color:rgb(44, 44, 54);">输入变量插值</font>**<font style="color:rgb(44, 44, 54);">：通过</font><font style="color:rgb(44, 44, 54);"> </font>`{}`<font style="color:rgb(44, 44, 54);"> </font><font style="color:rgb(44, 44, 54);">占位符动态替换变量。</font>
+ **模板格式**：支持 >`f-string</font>`<font style="color:rgb(44, 44, 54);">。</font>
+ **<font style="color:rgb(44, 44, 54);">输出形式</font>**<font style="color:rgb(44, 44, 54);">：生成一个完整的字符串作为模型输入。</font>

**示例：**

```python
from langchain_core.prompts import PromptTemplate

template = PromptTemplate.from_template("今天{somethingName}真不错")
result = template.format(somethingName="天气")
print(result)  # 今天天气真不错
```

#### ChatPromptTemplate<font style="color:rgb(44, 44, 54);">（聊天提示模板）</font>

**<font style="color:rgb(44, 44, 54);">适用场景</font>**<font style="color:rgb(44, 44, 54);">：</font>

+ <font style="color:rgb(44, 44, 54);">用于</font>**<font style="color:rgb(44, 44, 54);">聊天模型</font>**<font style="color:rgb(44, 44, 54);">（如 ChatGPT / 通义千问等），输入是多轮对话的消息列表（</font>`SystemMessage`<font style="color:rgb(44, 44, 54);">、</font>`HumanMessage` `AIMessage`等）。
+ <font style="color:rgb(44, 44, 54);">适用于需要模拟多轮对话或角色扮演的场景。</font>

**<font style="color:rgb(44, 44, 54);">特点</font>**<font style="color:rgb(44, 44, 54);">：</font>

+ **<font style="color:rgb(44, 44, 54);">多消息类型支持</font>**<font style="color:rgb(44, 44, 54);">：可以组合系统指令、用户输入和助手回复。</font>
+ **<font style="color:rgb(44, 44, 54);">消息格式化</font>**<font style="color:rgb(44, 44, 54);">：生成结构化的消息列表，供聊天模型处理。</font>
+ **灵活性**：支持动态替换变量（如 SystemMessage 中的占位符）

**示例：**

```python
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.messages import SystemMessage, HumanMessage

chat_prompt = ChatPromptTemplate.from_messages([
    ("system", "你是{role}专家，擅长{domain}"),
    ("human", "{question}")
])

messages = chat_prompt.format_messages(
    role="技术",
    domain="开发",
    question="如何学习？"
)

print(messages) #[SystemMessage(content='你是技术专家，擅长开发', additional_kwargs={}, response_metadata={}), HumanMessage(content='如何学习？', additional_kwargs={}, response_metadata={})]
```

python基础知识扩充：

这里的 `("system", system_message)`是 python 的 tuple 类型

**<font style="color:rgb(44, 44, 54);">Tuple（元组）</font>**<font style="color:rgb(44, 44, 54);">：</font>

+ **<font style="color:rgb(44, 44, 54);">有序</font>**<font style="color:rgb(44, 44, 54);">：元素按顺序排列，支持索引和切片操作。</font>
+ **<font style="color:rgb(44, 44, 54);">不可变</font>**<font style="color:rgb(44, 44, 54);">：</font>**<font style="color:#DF2A3F;">创建后无法修改（不能增删改元素）</font>**<font style="color:rgb(44, 44, 54);">。</font>
+ **<font style="color:rgb(44, 44, 54);">定义方式</font>**：使用圆括号()或直接逗号分隔元素。

**Array（数组）**：相当于 JS 中的 Array

+ **<font style="color:rgb(44, 44, 54);">可变</font>**<font style="color:rgb(44, 44, 54);">：可以动态添加、删除或修改元素。</font>
+ **<font style="color:rgb(44, 44, 54);">有序</font>**<font style="color:rgb(44, 44, 54);">：元素按顺序排列，支持索引和切片。</font>

**<font style="color:rgb(44, 44, 54);">Dict（字典）</font>**<font style="color:rgb(44, 44, 54);">：相当于 JS 中的 Object</font>

+ **<font style="color:rgb(44, 44, 54);">无序</font>**<font style="color:rgb(44, 44, 54);">：元素（键值对）没有固定顺序。</font>
+ **<font style="color:rgb(44, 44, 54);">可变</font>**<font style="color:rgb(44, 44, 54);">：可以动态添加、删除或修改键值对。</font>
+ **<font style="color:rgb(44, 44, 54);">定义方式</font>**<font style="color:rgb(44, 44, 54);">：使用花括号 </font>`<font style="color:rgb(44, 44, 54);">{}</font>`<font style="color:rgb(44, 44, 54);">，键值对用冒号 </font>`<font style="color:rgb(44, 44, 54);">:</font>`<font style="color:rgb(44, 44, 54);"> 分隔。</font>

#### ChatMessagePromptTemplate

知识扩充：`ChatMessagePromptTemplate`可以结合 `ChatPromptTemplate`使用，同时对提示词模板和消息体进行抽象和复用：

```python
from langchain_core.prompts import ChatMessagePromptTemplate

system_template = ChatMessagePromptTemplate.from_template(
    template="你是一位{role}专家，擅长回答{domain}领域的问题。",
    role="system",
)

human_template = ChatMessagePromptTemplate.from_template(
    template="用户问题：{question}",
    role="human",
)

chat_prompt = ChatPromptTemplate.from_messages([
    system_template,
    human_template,
])

messages = chat_prompt.format_messages(
    role="技术",
    domain="Web开发",
    question="如何构建一个基于Vue的前端应用？"
)

print(messages)

"""
[
  ChatMessage(content='你是一位技术专家，擅长回答Web开发领域的问题。', additional_kwargs={}, response_metadata={}, role='system'),
  ChatMessage(content='用户问题：如何构建一个基于Vue的前端应用？', additional_kwargs={}, response_metadata={}, role='human')
]
"""
```

#### FewShotPromptTemplate<font style="color:rgb(44, 44, 54);">（少样本提示模板）</font>

**<font style="color:rgb(44, 44, 54);">适用场景</font>**<font style="color:rgb(44, 44, 54);">：</font>

+ <font style="color:rgb(44, 44, 54);">用于</font>**<font style="color:rgb(44, 44, 54);">少样本学习</font>**<font style="color:rgb(44, 44, 54);">（Few-Shot Learning），在提示中包含示例（Examples），帮助模型理解任务。</font>
+ <font style="color:rgb(44, 44, 54);">适用于复杂任务（如翻译、分类、推理），需要通过示例引导模型行为。</font>
+ 聊天模型

**<font style="color:rgb(44, 44, 54);">特点</font>**<font style="color:rgb(44, 44, 54);">：</font>

+ **<font style="color:rgb(44, 44, 54);">示例嵌入</font>**<font style="color:rgb(44, 44, 54);">：通过</font><font style="color:rgb(44, 44, 54);"> </font>`examples </font><font style="color:rgb(44, 44, 54);">参数提供示例输入和输出。</font>
+ **<font style="color:rgb(44, 44, 54);">动态示例选择</font>**<font style="color:rgb(44, 44, 54);">：支持</font><font style="color:rgb(44, 44, 54);"> </font>`ExampleSelector`<font style="color:rgb(44, 44, 54);"> </font><font style="color:rgb(44, 44, 54);">动态选择最相关的示例。</font>
+ **<font style="color:rgb(44, 44, 54);">模板格式</font>**<font style="color:rgb(44, 44, 54);">：通常包含前缀（Prefix）、示例（Examples）和后缀（Suffix）。</font>

**<font style="color:rgb(44, 44, 54);">示例</font>**<font style="color:rgb(44, 44, 54);">：</font>

```python
from langchain_core.prompts import FewShotChatMessagePromptTemplate, ChatPromptTemplate

# 少样本示例
examples = [
    {"input": "Hello", "output": "你好"},
    {"input": "Goodbye", "output": "再见"},
]

# 创建示例模板（用于聊天模型）
example_prompt = ChatPromptTemplate.from_messages([
    ("human", "{input}"),
    ("ai", "{output}")
])

few_shot_prompt = FewShotChatMessagePromptTemplate(
    examples=examples,
    example_prompt=example_prompt,
)

# 组合到完整提示词
final_prompt = ChatPromptTemplate.from_messages([
    ("system", "你是一个翻译助手，将英文翻译成中文"),
    few_shot_prompt,  # 插入少样本示例
    ("human", "{input}"),
])

messages = final_prompt.format_messages(input="Thank you")
print(messages)
```

#### 常用模板类特性和使用场景对比

| **<font style="color:rgb(44, 44, 54);">子类</font>** | **<font style="color:rgb(44, 44, 54);">适用模型类型</font>** | **<font style="color:rgb(44, 44, 54);">输入类型</font>**   | **<font style="color:rgb(44, 44, 54);">主要用途</font>**     |
| ---------------------------------------------------- | ------------------------------------------------------------ | ---------------------------------------------------------- | ------------------------------------------------------------ |
| PromptTemplate                                       | <font style="color:rgb(44, 44, 54);">文本补全模型</font>     | <font style="color:rgb(44, 44, 54);">单字符串</font>       | <font style="color:rgb(44, 44, 54);">生成单轮文本任务的提示</font> |
| ChatPromptTemplate                                   | <font style="color:rgb(44, 44, 54);">聊天模型</font>         | <font style="color:rgb(44, 44, 54);">多消息列表</font>     | <font style="color:rgb(44, 44, 54);">模拟多轮对话或角色扮演</font> |
| FewShotPromptTemplate                                | <font style="color:rgb(44, 44, 54);">所有模型</font>         | <font style="color:rgb(44, 44, 54);">包含示例的模板</font> | <font style="color:rgb(44, 44, 54);">通过示例引导模型完成复杂任务</font> |


#### 常用提示词模板类的继承关系

![image-20260212175918848](https://github.com/Auroraol/Drawing-bed/raw/main/img/image-20260212175918848.png)

### 3. 索引（Indexes）

索引组件用于处理外部数据，通常包括：

- **文档加载器（Document Loaders）**：从 PDF、网页、CSV 等加载数据。
- **文本分割器（Text Splitters）**：将长文本分割为小块，适配 LLM 的上下文限制。
- **向量存储（Vector Stores）**：将文本向量化存储（如使用 FAISS 或 Chroma），支持语义搜索。
- **嵌入模型（Embeddings）**：将文本转换为向量表示，如 OpenAI Embeddings。

**作用**：让 LLM 访问和理解外部数据，支持 RAG（检索增强生成）工作流。

**代码示例**：

```python
from langchain.document_loaders import WebBaseLoader
from langchain.text_splitter import RecursiveCharacterTextSplitter

# 加载网页内容
loader = WebBaseLoader("https://example.com")
docs = loader.load()

# 分割文本
text_splitter = RecursiveCharacterTextSplitter(chunk_size=500, chunk_overlap=50)
split_docs = text_splitter.split_documents(docs)
print(len(split_docs))  # 输出分割后的文档块数量
```

### 4. 记忆（Memory）

记忆组件用于在对话中保存上下文信息，确保 LLM 能够“记住”之前的交互内容。常见类型包括：

- **ConversationBufferMemory**：保存完整的对话历史。
- **ConversationSummaryMemory**：对对话进行摘要，节省 token。

**作用**：让对话更连贯，适合聊天机器人等场景。

**代码示例**：

```python
from langchain.chains import ConversationChain
from langchain.memory import ConversationBufferMemory

# 初始化 LLM 和记忆
llm = OpenAI(model_name="gpt-3.5-turbo", api_key="YOUR_API_KEY")
memory = ConversationBufferMemory()

# 创建对话链
conversation = ConversationChain(llm=llm, memory=memory)

# 模拟对话
conversation.run("你好，我叫小明。")
conversation.run("我叫什么名字？")
print(memory.buffer)  # 输出对话历史：包含“小明”
```

### 5. 链（Chains）

链是将多个组件组合成一个工作流的机制，分为：

- **LLMChain**：结合提示词和 LLM，执行单次任务。
- **SequentialChain**：按顺序执行多个链。
- **RouterChain**：根据输入动态选择执行哪个链。

**作用**：自动化复杂任务，简化开发流程。

```python
chain = prompt | llm  # 提示词 |大模型 链

response = chain.stream({"user_input": "计算100+100"})
# response = chain.stream(input={"user_input": "计算100+100"})

for chunk in response:
    print(chunk.content, end="")
```

### 6. 代理（Agents）

代理是一个动态控制单元，允许 LLM 根据输入自主选择工具并执行任务。代理通常包含：

- **工具（Tools）**：如搜索引擎、计算器、API 调用。
- **代理执行器（Agent Executor）**：协调 LLM 和工具的交互。

**作用**：赋予 LLM 主动推理和行动能力，适合复杂任务。

### 7. 回调（Callbacks）

回调用于监控和调试 LangChain 应用的运行过程，如记录日志、跟踪 token 使用量等。

**作用**：提高开发透明度，便于优化性能。

## 三、LangChain 的核心架构

LangChain 的架构以模块化和链式工作流为核心，旨在通过组件的组合实现复杂功能。以下是其核心架构的分解，以及一个简化的流程图。

### 1. 架构分解

LangChain 的架构可以分为三个层次：

- **输入层**：接收用户输入（如问题、指令），通过提示词模板格式化。
- **处理层**：核心组件（如模型、索引、记忆、代理）协同工作，调用 LLM、检索数据或执行工具。
- **输出层**：生成最终输出（如文本回答、动作结果），并通过输出解析器格式化。

### 2. 典型工作流（以 RAG 为例）

以下是一个基于 RAG（检索增强生成）的典型工作流，展示 LangChain 如何处理用户查询：

使用向量存储用户输入: “公司文档中2024年的销售目标是多少？”提示词模板: 格式化查询索引: 检索文档获取相关文档片段LLM: 生成回答输出解析器: 格式化回答最终回答: “2024年销售目标是5000万元。”

```Mermaid
sequenceDiagram
    participant User as 用户
    participant PT as 提示词模板
    participant VS as 向量存储
    participant LLM as 大语言模型
    participant OP as 输出解析器
    
    User->>PT: 公司文档中2024年的销售目标是多少？
    PT->>VS: 格式化查询向量
    VS->>VS: 相似度检索
    VS-->>PT: 返回相关文档片段
    PT->>LLM: 构建Prompt(查询+上下文)
    LLM-->>OP: 生成原始回答
    OP-->>User: 格式化输出：2024年销售目标是5000万元。
```

### 3. 架构特点

- **模块化**：每个组件（如提示词、索引）独立且可复用。
- **可扩展**：支持多种 LLM、数据源和工具，易于集成新功能。
- **上下文感知**：通过记忆和索引，确保输出与上下文相关。
- **动态性**：代理机制允许根据输入动态调整工作流。

## 四、LangChain 的 Agent、Chain 和 Model 详解

### 1. Chain

**Chain（链）** 是 LangChain 的核心机制，用于将多个组件按顺序组合，完成复杂任务。链就像一条流水线，输入经过一系列处理步骤，最终生成输出。

**类型**：

- **LLMChain**：结合提示词和 LLM，执行单一任务。
- **SequentialChain**：按顺序运行多个链，上一链的输出作为下一链的输入。
- **RouterChain**：根据输入选择不同的链分支。

**代码示例**：

```python
from langchain_ollama import ChatOllama
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.runnables import RunnablePassthrough, RunnableParallel

llm = ChatOllama(model="deepseek-r1:1.5b", base_url="http://localhost:11434")

# 定义两个独立的子链
prompt1 = ChatPromptTemplate.from_template("描述{city}的主要特点")
chain1 = prompt1 | llm

prompt2 = ChatPromptTemplate.from_template("根据以下描述推荐三种活动：{description}")
chain2 = prompt2 | llm

# 真正的链：使用 RunnableParallel 同时执行或顺序传递
# 第一步：添加 description
# 输入: {"city": "上海"}
# 输出: {"city": "上海", "description": "..."}

# 第二步：添加 activities
# 输入: {"city": "上海", "description": "..."}
# 输出: {"city": "上海", "description": "...", "activities": "..."}
chain = (
    RunnablePassthrough.assign(description=chain1)  # 执行 chain1，结果存入 description
    | RunnablePassthrough.assign(activities=chain2)  # 执行 chain2，使用 description
)

# 执行（只传一次输入）
result = chain.invoke({"city": "上海"})
print("描述：", result["description"])
print("活动：", result["activities"])
```

```
描述： 上海是中国最大的城市之一，位于长江入海口...
活动： 1. 外滩观光 2. 豫园品茶 3. 陆家嘴金融区游览
```

**应用场景**：多步骤任务，如先总结文档再生成报告。

### 2. Agent

**Agent（代理）** 是 LangChain 中更高级的机制，允许 LLM 像“智能助手”一样，根据输入动态选择和使用工具，完成复杂任务。代理的核心是“推理 + 行动”，它会分析输入，决定需要调用哪些工具，并循环执行直到任务完成。

**组成**：

- **工具（Tools）**：如搜索引擎、计算器、API。
- **代理执行器**：协调 LLM 和工具，管理执行流程。
- **推理引擎**：LLM 负责决定使用哪个工具。

**应用场景**：需要动态决策的任务，如联网搜索、数学计算。

#### 工具（Tools）

工具（Tool）是一个封装了特定功能的类，它包含四个核心组成部分：

+ <font style="color:rgb(44, 44, 54);">名称（name）：名称是工具在工具集合中的</font>**<font style="color:rgb(44, 44, 54);">唯一标识符</font>**<font style="color:rgb(44, 44, 54);">，必须确保在同一工具集中不重复</font>

+ <font style="color:rgb(44, 44, 54);">描述（description）：描述用于说明工具的功能，为LLM或代理提供上下文信息，</font>**<font style="color:rgb(44, 44, 54);">帮助模型理解何时以及如何调用该工具</font>**

+ <font style="color:rgb(44, 44, 54);">参数模式（args_schema）：是使用Pydantic BaseModel定义的输入参数结构，用于验证和解析工具调用的参数</font>

+ <font style="color:rgb(44, 44, 54);">是否直接返回（return_direct）：布尔值属性，当设置为True时，智能体会在调用工具后立即返回结果给用户，而不继续调用其他工具</font>

+ ```python
  # 入参/出参含义
  class AddInputArgs(BaseModel):
      a: str = Field(description="first number")
      b: str = Field(description="second number")
  
  #使用@tool装饰器生成
  @tool(
      description="add two numbers",
      args_schema=AddInputArgs,
      return_direct=True,
  )
  def add(a, b):
      """add two numbers"""
      return a + b
  ```

**代码示例**：

```python
import os
import warnings
from langchain_openai import ChatOpenAI
from langgraph.prebuilt.chat_agent_executor import create_tool_calling_executor
from pydantic import SecretStr
# 设置环境变量以避免多线程问题
os.environ["OMP_NUM_THREADS"] = "1"
warnings.filterwarnings("ignore")

# 导入最新版本所需的模块
from langchain_ollama import ChatOllama
from langchain_core.tools import tool
from pydantic import BaseModel, Field

# 定义工具输入参数 schema
class AddInputArgs(BaseModel):
    a: int = Field(description="第一个数字")
    b: int = Field(description="第二个数字")

# 使用 @tool 装饰器定义工具函数
@tool(args_schema=AddInputArgs, return_direct=False)
def add(a: int, b: int) -> int:
    """将两个数字相加"""
    return a + b

# 初始化 LLM（适配最新版本）
# llm = ChatOllama(
#     model="deepseek-r1:1.5b",
#     base_url="http://localhost:11434",
#     temperature=0.7,
# )
KimiLlm = ChatOpenAI(
    model="kimi-k2-turbo-preview",
    base_url="https://api.moonshot.cn/v1",
    api_key=SecretStr("sk-32lqdlnpX9Cgc7iwAR"),
    temperature=0.7,
)

# 新版方式
executor = create_tool_calling_executor(KimiLlm, [add])

# 执行（新版输入格式可能不同）
result = executor.invoke({
    "messages": [("human", "计算 100 加 200")]
})

# 输出结果
print("=" * 50)
for msg in result["messages"]:
    print(f"{msg.type}: {msg.content}")
print("=" * 50)
```

```
==================================================
human: 计算 100 加 200
ai: 我来为您计算 100 加 200。
tool: 300
ai: 计算结果为：100 + 200 = **300**
==================================================
```

### 3.  Model

**Model（模型）** 是 LangChain 中负责语言处理的核心组件，分为：

- **LLM**：处理纯文本输入输出，适合生成、翻译等任务。
- **Chat Model**：优化对话场景，支持多轮交互和消息类型。

**特点**：

- LangChain 提供统一的接口，支持多种模型（如 OpenAI、Hugging Face、Anthropic）。
- 开发者可以轻松切换模型，无需修改代码。

**代码示例**（聊天模型）：

```python
# 本地模型
from langchain_ollama import ChatOllama
llm = ChatOllama(
    model="deepseek-r1:1.5b",
    base_url="http://localhost:11434",  # 默认地址
    temperature=0.7,
)

# 单轮
response = llm.invoke("北京的首都是什么？")
print(response.content)

# 多轮对话
messages = [
    HumanMessage(content="你好"),
    HumanMessage(content="北京的首都是什么？")
]
response = llm.invoke(messages)
print(response.content)


# Kimi (Moonshot AI)
from langchain_openai import ChatOpenAI
from pydantic import SecretStr
llm = ChatOpenAI(
    model="kimi-k2-0711-preview",
    base_url="https://api.moonshot.cn/v1",  # 必须指定
    api_key=SecretStr("sk-your-kimi-api-key"),  # 安全包装
    temperature=0.7,
)
```

```
北京是中国的首都，同时也是中华xxxx。
北京是中国的首都，它是中国的一角，也是中国政治、xxx。
```

**应用场景**：对话系统、知识问答。

**资源链接**：

- GitHub 仓库：https://github.com/langchain-ai/langchain

- 官方文档：https://python.langchain.com/docs/

- LangChain 社区：https://discord.gg/langchain

- ```
  pip install --upgrade langchain langchain-community langchain-ollama langchain-text-splitters
  ```
