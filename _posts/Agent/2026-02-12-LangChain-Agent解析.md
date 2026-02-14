---
title: LangChain-Agent解析
date: 2026-02-13 16:00:00 +0800
categories: [Agent, LangChain-Agent解析]
tags: [LangChain-Agent解析]
---

# LangChain-Agent 解析

## 创建Agent

**Agent（代理）** 是 LangChain 中更高级的机制，允许 LLM 像“智能助手”一样，根据输入动态选择和使用工具，完成复杂任务。代理的核心是“推理 + 行动”，它会分析输入，决定需要调用哪些工具，并循环执行直到任务完成。

**组成**：

- **工具（Tools）**：如搜索引擎、计算器、API。
- **代理执行器**：协调 LLM 和工具，管理执行流程。
- **推理引擎**：LLM 负责决定使用哪个工具。

**应用场景**：需要动态决策的任务，如联网搜索、数学计算。

### LangChain方式

```python
from langchain.agents import create_react_agent, AgentExecutor  # 从 langchain.agents 导入
from langchain_core.prompts import ChatPromptTemplate, MessagesPlaceholder
from langchain_openai import ChatOpenAI
from langchain_core.messages import HumanMessage
from pydantic import BaseModel, Field

# ========== 3. 配置模型 ==========
model = ChatOpenAI(model="gpt-4", temperature=0)

# ========== 4. 创建 PromptTemplate ==========
prompt = ChatPromptTemplate.from_messages([
    ("system", "你是数学专家，必须使用calculator工具计算，然后返回结构化结果"),
    ("human", "{input}"),
    MessagesPlaceholder("agent_scratchpad")  # ReAct 必需
])

# ========== 5. 创建 Agent ==========
agent = create_react_agent(
    llm=model,  # 注意：参数名是 llm 不是 model
    tools=[calculator],
    prompt=prompt
)

# ========== 6. 使用 AgentExecutor 执行 ==========
agent_executor = AgentExecutor(
    agent=agent,
    tools=[calculator],
    verbose=True,
    handle_parsing_errors=True
)

# ========== 7. 执行计算 ==========
result = agent_executor.invoke({"messages": "计算 100+100"})
print(result["messages"][-1].content)
```

### LangGraph方式

```python
# llm
KimiLlm = ChatOpenAI(
    model="kimi-k2-turbo-preview",
    base_url="https://api.moonshot.cn/v1",
    api_key=SecretStr("sk-32lqdlnpX9CgcRPa7KiMdeoepSaIO6B79BqbyL"),
    temperature=1,
)

# 使用LangGraph
agent = create_agent(
    model=KimiLlm,
    tools=[calculator, get_weather],  # 工具列表
)

# 直接 invoke，无需 AgentExecutor
result = agent.invoke({
    "messages": [
        SystemMessage(content="计算 100+100"),
        HumanMessage(content="开始计算")
    ]
})

print(result["messages"][-1].content)
```

### prompt 参数

```python
# 方式1：使用 system_prompt 参数
agent = create_react_agent(
    model=llm,
    tools=[get_weather],
    system_prompt="你是专业的天气分析师..."  # 等价于下面的 SystemMessage
)

# 方式2：手动传入 SystemMessage（等价）
agent = create_react_agent(
    model=llm,
    tools=[get_weather]
)
# 调用时手动添加
result = agent.invoke({
    "messages": [
        SystemMessage(content="你是专业的天气分析师..."),  # 等价于 system_prompt
        HumanMessage(content="北京天气怎么样？")
    ]
})
```

### AgentType

LangChain提供多种Agent类型，每种类型适用于不同的场景和需求。以下是几种主要Agent类型及其特点：

注意:  `langchain`  1.2.10 已废弃AgentType

| 场景                      | 推荐 Agent                                    |
| :------------------------ | :-------------------------------------------- |
| 快速原型验证              | `ZERO_SHOT_REACT_DESCRIPTION`                 |
| 生产环境 + OpenAI         | `OPENAI_FUNCTIONS`                            |
| 生产环境 + 非 OpenAI 模型 | `STRUCTURED_CHAT_ZERO_SHOT_REACT_DESCRIPTION` |
| 需要记住对话历史          | `CONVERSATIONAL_REACT_DESCRIPTION`            |
| 需要自主分解复杂任务      | `AUTO_GPT`（评估稳定性后）                    |
| 严格的数据/表单输出       | `STRUCTURED_CHAT`                             |

## 消息类型

| 类型       | 类名              | 说明                 | 典型使用场景                  |
| :--------- | :---------------- | :------------------- | :---------------------------- |
| `human`    | `HumanMessage`    | 人类/用户输入        | 用户提问、指令                |
| `ai`       | `AIMessage`       | AI/助手回复          | 模型生成内容、工具调用决策    |
| `tool`     | `ToolMessage`     | 工具执行结果         | 计算器、API调用返回值         |
| `system`   | `SystemMessage`   | 系统提示词           | 设定角色、约束条件            |
| `function` | `FunctionMessage` | 函数调用结果（旧版） | 早期 OpenAI 函数调用          |
| `chat`     | `ChatMessage`     | 通用聊天消息         | 指定任意 role（如 assistant） |

```
所有消息类型：
0. type='system'     | class=SystemMessage      ← 您传入的系统提示
1. type='human'      | class=HumanMessage       ← 用户输入
2. type='ai'        | class=AIMessage          ← AI决定调用工具
3. type='tool'       | class=ToolMessage        ← 计算器返回结果
4. type='ai'        | class=AIMessage          ← AI总结最终结果
```

## 工具Tools使用

工具（Tool）是一个封装了特定功能的类，它包含四个核心组成部分：

+ 名称（name）：名称是工具在工具集合中的**唯一标识符**<font style="color:rgb(44, 44, 54);">，必须确保在同一工具集中不重复</font>

+ 描述（description）：描述用于说明工具的功能，为LLM或代理提供上下文信息，**帮助模型理解何时以及如何调用该工具**

+ 参数模式（args_schema）：是使用Pydantic BaseModel定义的输入参数结构，用于验证和解析工具调用的参数

+ 是否直接返回（return_direct）：布尔值属性，当设置为True时，智能体会在调用工具后立即返回结果给用户，而不继续调用其他工具

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
from langgraph.prebuilt import create_react_agent
from langchain_openai import ChatOpenAI
from langchain_core.tools import tool
from langchain_ollama import OllamaLLM

# 1. 定义工具（使用@tool装饰器自动解析参数）
@tool
def calculator(expression: str) -> str:
    """计算数学表达式，如'2+2*5'"""
    try:
        return str(eval(expression))
    except:
        return "计算错误"

@tool
def get_weather(location: str) -> str:
    """获取指定城市的天气"""
    return f"{location}天气：晴天，25°C"

# 2. 工具列表
tools = [calculator, get_weather]

# 3. 选择模型
# OpenAI模型
# model = ChatOpenAI(model="gpt-4", temperature=0)

# 本地Ollama模型（您已安装langchain-ollama==1.0.1）
model = OllamaLLM(model="llama3:8b", temperature=0)

# 4. 创建Agent（一行代码）
agent = create_react_agent(
    model=model,
    tools=tools,
    # 可选：系统提示词
    # prompt="你是一个有用的助手，请用中文回答"
)

# 5. 运行Agent（注意：输入格式为消息列表）
from langchain_core.messages import HumanMessage

result = agent.invoke({
    "messages": [HumanMessage(content="北京和上海温度总和是多少？")]
})

# 输出结果
print(result["messages"][-1].content)
```

## invoke使用

```python
result = agent.invoke({
    "messages": [
        SystemMessage(content="你是数学专家"),      # 系统提示（可选）
        HumanMessage(content="计算 100+100"),       # 用户输入（必需）
        AIMessage(content="..."),               # 历史AI回复（多轮时）
        ToolMessage(content="...", tool_call_id="...")  # 历史工具结果（多轮时）
    ]
})
```

**代码示例**：

```python
from langchain.agents import create_agent
from langchain_openai import ChatOpenAI
from langchain_core.tools import tool
from langchain_core.messages import HumanMessage, SystemMessage, AIMessage
import json
from pydantic import SecretStr


@tool
def calculator(expression: str) -> str:
    """计算数学表达式"""
    print(f"🛠️ 工具执行: {expression}")
    return str(eval(expression))


@tool
def get_weather(city: str) -> str:
    """获取指定城市的天气"""
    print(f"🌤️ 天气工具被调用: city={city}")  # 添加日志
    return f"{city}天气：晴天 25°C"


KimiLlm = ChatOpenAI(
    model="kimi-k2-turbo-preview",
    base_url="https://api.moonshot.cn/v1",
    api_key=SecretStr("sk-32lqdlnpX9Cgc7iwpSaIO6B79BqbyL"),
    temperature=1,
)

# LangGraph
agent = create_agent(
    model=KimiLlm,
    tools=[calculator, get_weather],  # 工具列表
)

# ========== 多轮对话（手动维护历史）==========
print("\n" + "=" * 50)
print("多轮对话示例：")
print("=" * 50)

# 第一轮
result1 = agent.invoke({
    "messages": [HumanMessage(content="计算 100+100")]
})
history = result1["messages"]
print(f"第一轮: {history[-1].content}")
# 第二轮：带上历史 + 新消息
history = result1["messages"]
result2 = agent.invoke({
    "messages": history + [HumanMessage(content="再加 50")]
})
print(f"第二轮: {result2['messages'][-1].content}")
# 第三轮：天气（明确指定城市）
result3 = agent.invoke({
    "messages": [HumanMessage(content="北京天气如何？")]
})
print(f"第三轮: {result3['messages'][-1].content}")


# ========== 带 thread_id（自动持久化）==========
print("\n" + "=" * 50)
print("持久化记忆示例：")
print("=" * 50)

from langgraph.checkpoint.memory import MemorySaver

# 创建带检查点的 Agent
agent_with_memory = create_agent(
    model=KimiLlm,
    tools=[calculator, get_weather],  # 传入所有工具
    checkpointer=MemorySaver()  # 内存存储
)

# 第一轮 - 指定 thread_id
config = {
    "configurable": {
        "thread_id": "user_123",  # 必需：会话ID
        # "user_id": "u_789",           # 可选：用户ID
        # "checkpoint_ns": "tenant_a"   # 可选：命名空间
    }
}
result3 = agent_with_memory.invoke(
    {"messages": [HumanMessage(content="我叫张三, 住在黑龙江")]},
    config=config
)
print(f"第一轮: {result3['messages'][-1].content}")

# 第二轮 - 同一 thread_id，自动加载历史
result4 = agent_with_memory.invoke(
    {"messages": [HumanMessage(content="我叫什么名字？")]},
    config=config
)
print(f"第二轮: {result4['messages'][-1].content}")
# 第3轮 - 同一 thread_id，自动加载历史
result5 = agent_with_memory.invoke(
    {"messages": [HumanMessage(content="天气如何？")]},
    config=config
)
print(f"第三轮: {result5['messages'][-1].content}")

```

```
==================================================
多轮对话示例：
==================================================
🛠️ 工具执行: 100+100
第一轮: 计算结果是：200
🛠️ 工具执行: 200+50
第二轮: 计算结果是：250（200 + 50 = 250）
🌤️ 天气工具被调用: city=北京
第三轮: 北京的天气很好！现在是晴天，气温25摄氏度，非常适合外出活动。

==================================================
持久化记忆示例：
==================================================
第一轮: 你好，张三！很高兴认识你。黑龙江是个美丽的地方，冬天特别冷吧？有什么我可以帮你的吗？
第二轮: 你叫张三。
第三轮: 我现在还不知道你在哪个城市，黑龙江很大，各地天气可能差别不小。请告诉我你想查询黑龙江的哪一座城市（比如哈尔滨、齐齐哈尔、大庆、佳木斯等），我就能帮你查实时天气啦！

```

## 返回体结构化

```python
from langchain.agents import create_agent
#from langgraph.prebuilt import create_react_agent  # 注意：标记为弃用但仍可用
from langchain_openai import ChatOpenAI
from langchain_core.tools import tool
from langchain_core.messages import HumanMessage, SystemMessage
from pydantic import SecretStr

@tool
def calculator(expression: str) -> str:
    """计算数学表达式"""
    return str(eval(expression))

KimiLlm = ChatOpenAI(
    model="kimi-k2-turbo-preview",
    base_url="https://api.moonshot.cn/v1",
    api_key=SecretStr("sk-32lqdlnpX9Cgc7iwAiMdeoepSaIO6B79BqbyL"),
    temperature=0.7,
)

# LangGraph
agent = create_agent(
    model=KimiLlm,
    tools=[calculator],
)

# 直接 invoke，无需 AgentExecutor
result = agent.invoke({
    "messages": [
        SystemMessage(content="计算 100+90"),
        HumanMessage(content="开始计算")
    ]
})

# ========== 5. 详细查看 result 结构 ==========
print("=" * 50)
print("📦 result 类型：", type(result))
print("📦 result 键：", result.keys())
print("=" * 50)

# 查看 messages 列表
print(f"\n📨 messages 列表（共 {len(result['messages'])} 条）：")
print("=" * 50)

for i, msg in enumerate(result["messages"]):
    print(f"\n--- 消息 {i} ---")
    print(f"类型：{msg.type}")
    print(f"角色：{getattr(msg, 'role', 'N/A')}")
    print(f"内容：{msg.content}")

    # 查看额外字段
    if hasattr(msg, 'tool_calls'):
        print(f"工具调用：{msg.tool_calls}")
    if hasattr(msg, 'name'):
        print(f"工具名：{msg.name}")
    if hasattr(msg, 'tool_call_id'):
        print(f"工具调用ID：{msg.tool_call_id}")
    if hasattr(msg, 'additional_kwargs'):
        print(f"额外参数：{msg.additional_kwargs}")

# ========== 6. 最后一条消息详细分析 ==========
print("\n" + "=" * 50)
print("🔍 最后一条消息详细分析：")
print("=" * 50)

last_msg = result["messages"][-1]
print(f"对象类型：{type(last_msg)}")
print(f"消息类型：{last_msg.type}")
print(f"内容：{repr(last_msg.content)}")

# 查看所有属性
print(f"\n所有属性：")
for attr in dir(last_msg):
    if not attr.startswith('_'):
        try:
            val = getattr(last_msg, attr)
            if not callable(val):
                print(f"  {attr}: {repr(val)[:100]}")
        except:
            pass

# ========== 7. 提取结构化数据示例 ==========
# print("\n" + "=" * 50)
# print("📊 提取数据示例：")
# print("=" * 50)
#
# # 提取所有工具消息
# tool_msgs = [m for m in result["messages"] if m.type == "tool"]
# print(f"工具消息数量：{len(tool_msgs)}")
# for tm in tool_msgs:
#     print(f"  - {tm.name}: {tm.content}")
#
# # 提取所有AI消息
# ai_msgs = [m for m in result["messages"] if m.type == "ai"]
# print(f"\nAI消息数量：{len(ai_msgs)}")
# for am in ai_msgs:
#     print(f"  - 内容：{am.content[:200]}...")
#     if hasattr(am, 'tool_calls') and am.tool_calls:
#         print(f"    工具调用：{am.tool_calls}")
```

```
==================================================
📦 result 类型： <class 'dict'>
📦 result 键： dict_keys(['messages'])
==================================================

📨 messages 列表（共 5 条）：
==================================================

--- 消息 0 ---
类型：system
角色：N/A
内容：计算 100+90
工具名：None
额外参数：{}

--- 消息 1 ---
类型：human
角色：N/A
内容：开始计算
工具名：None
额外参数：{}

--- 消息 2 ---
类型：ai
角色：N/A
内容：我来帮您计算 100 + 90。
工具调用：[{'name': 'calculator', 'args': {'expression': '100+90'}, 'id': 'calculator:0', 'type': 'tool_call'}]
工具名：None
额外参数：{'refusal': None}

--- 消息 3 ---
类型：tool
角色：N/A
内容：190
工具名：calculator
工具调用ID：calculator:0
额外参数：{}

--- 消息 4 ---
类型：ai
角色：N/A
内容：计算结果：100 + 90 = 190
工具调用：[]
工具名：None
额外参数：{'refusal': None}

==================================================
🔍 最后一条消息详细分析：
==================================================
对象类型：<class 'langchain_core.messages.ai.AIMessage'>
消息类型：ai
内容：'计算结果：100 + 90 = 190'

所有属性：
  additional_kwargs: {'refusal': None}
  content: '计算结果：100 + 90 = 190'
  content_blocks: [{'type': 'text', 'text': '计算结果：100 + 90 = 190'}]
  id: 'lc_run--019c55e5-2438-74e0-b1a9-8306f388edaf-0'
  invalid_tool_calls: []
  lc_attributes: {'tool_calls': [], 'invalid_tool_calls': []}
  lc_secrets: {}
  model_computed_fields: {}
  model_config: {'extra': 'allow'}
  model_extra: {}
  model_fields: {'content': FieldInfo(annotation=Union[str, list[Union[str, dict]]], required=True), 'additional_kwa
  model_fields_set: {'tool_calls', 'additional_kwargs', 'response_metadata', 'content', 'name', 'usage_metadata', 'id', 
  name: None
  response_metadata: {'token_usage': {'completion_tokens': 11, 'prompt_tokens': 100, 'total_tokens': 111, 'completion_tok
  tool_calls: []
  type: 'ai'
  usage_metadata: {'input_tokens': 100, 'output_tokens': 11, 'total_tokens': 111, 'input_token_details': {}, 'output_t

```

## LangChain内置工具

Tool使用方法：[https://python.langchain.com/docs/integrations/tools/](https://python.langchain.com/docs/integrations/tools/)

内置Tooltiks：[https://api.python.langchain.com/en/latest/community/agent_toolkits.html](https://api.python.langchain.com/en/latest/community/agent_toolkits.html)

### Python REPL

官网：[https://python.langchain.com/docs/integrations/tools/python/](https://python.langchain.com/docs/integrations/tools/python/)

**利用Python REPL开发一个企业官网**

```python
from langchain.agents import initialize_agent, AgentType
from langchain_core.prompts import PromptTemplate
from langchain_experimental.tools.python.tool import PythonREPLTool

from app.common import llm


tools = [PythonREPLTool()]

tool_names = ["PythonREPLTool"]

agent = initialize_agent(
    tools=tools,
    llm=llm,
    agent=AgentType.ZERO_SHOT_REACT_DESCRIPTION,
    verbose=True,
)

prompt_template = PromptTemplate.from_template(template="""尽你所能回答以下问题，你可以使用以下工具：
{tools}
--
请按照以下格式：
问题：你必须回答的输入问题
思考：你应该始终考虑该怎么做
行动：要采取的行动，应该是[{tool_names}]中的一个
行动输入：行动的输入
观察：行动的结果
... (这个思考/行动/行动输入/观察可以重复N次)
思考：我现在知道最终答案了
最终答案：对原始输入问题的最终答案
--
注意：
- PythonREPLTool工具的入参是python代码，不允许添加 ```python 或 ```py 等标记
--
开始吧！
问题：{input}""")

prompt = prompt_template.format(
    tools=tools,
    tool_names=tool_names,
    input="向/Users/sam/llm/.temp目录下写入一个新文件，名称为：index.html，并写一个企业的官网",
)

print(prompt)

agent.invoke(prompt)
```

回答结果：

```python
> Entering new AgentExecutor chain...
我需要使用Python代码来完成文件的创建和写入操作。
Action: Python_REPL
Action Input: 
with open("/Users/sam/llm/.temp/index.html", "w") as file:
    file.write("""
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Our Company - Home</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f4f4f9;
        }
        header {
            background: #2c3e50;
            color: #fff;
            padding: 10px 20px;
            text-align: center;
        }
        nav {
            background: #34495e;
            color: white;
            display: flex;
            justify-content: center;
            padding: 10px;
        }
        nav a {
            color: white;
            margin: 0 15px;
            text-decoration: none;
        }
        section {
            padding: 20px;
            text-align: center;
        }
        footer {
            background: #2c3e50;
            color: white;
            text-align: center;
            padding: 10px 0;
            position: fixed;
            width: 100%;
            bottom: 0;
        }
    </style>
</head>
<body>
    <header>
        <h1>Welcome to Our Company</h1>
    </header>
    <nav>
        <a href="#home">Home</a>
        <a href="#services">Services</a>
        <a href="#about">About Us</a>
        <a href="#contact">Contact</a>
    </nav>
    <section id="home">
        <h2>Your Success is Our Mission</h2>
        <p>We provide cutting-edge solutions tailored to your needs.</p>
    </section>
    <footer>
        &copy; 2023 Our Company. All rights reserved.
    </footer>
</body>
</html>
    """)
Python REPL can execute arbitrary code. Use with caution.
Observation: 
Thought:我现在知道最终答案了。
Final Answer: 已成功在目录 `/Users/sam/llm/.temp` 下创建了一个名为 `index.html` 的文件，并写入了一个企业官网的HTML代码。该页面包含一个欢迎标题、导航栏、主页内容以及页脚信息，展示了基本的企业官网结构。

> Finished chain.
```

### 解析器

LangChain输出解析器（Output Parsers）是将大语言模型（LLM）的原始文本响应转换为结构化、可操作数据的关键组件。输出解析器通过提供格式化指令并解析模型输出，实现**<font style="color:rgb(44, 44, 54);">文本到结构化数据</font>**<font style="color:rgb(44, 44, 54);">的高效转换。</font>

**基础解析器**

基础解析器处理最简单的数据格式转换：

+ **StrOutputParser**：直接提取模型返回的原始文本，不做任何结构化处理。
+ **CommaSeparatedListOutputParser**：将逗号分隔的文本转换为Python列表。例如，将"apples, bananas, oranges"解析为['apples', 'bananas', 'oranges']。
+ **BooleanOutputParser**：解析文本为布尔值（True/False）。模型输出必须是"yes"或"no"（不区分大小写），解析器会统一转为大写后判断。
+ **SimpleJsonOutputParser**：将文本简单处理后转换为JSON格式，通常用于模型已经正确输出JSON的情况。

这些基础解析器适合处理简单数据结构，但缺乏验证和容错机制，适用于对格式要求不严格的场景。

**Pydantic解析器**

Pydantic解析器通过Pydantic模型定义复杂结构，提供更严格的验证和数据校验：

+ **PydanticOutputParser**：将模型输出解析为Pydantic定义的模型对象。需要预先定义Pydantic模型，支持嵌套数据结构和类型验证。
+ **PydanticToolsParser**：专门处理OpenAI工具调用中的Pydantic对象，将工具调用参数解析为Pydantic模型。

Pydantic解析器适合需要严格验证的数据结构，如用户信息、订单详情等，确保输出的完整性和正确性。

**函数调用解析器**

<font style="color:rgb(44, 44, 54);">函数调用解析器处理支持OpenAI函数调用的模型（如GPT-4）的响应：</font>

+ **<font style="color:rgb(44, 44, 54);">OutputFunctionsParser</font>**<font style="color:rgb(44, 44, 54);">：解析模型返回的函数调用参数，通常返回JSON格式的函数参数。</font>
+ **<font style="color:rgb(44, 44, 54);">JsonOutputFunctionsParser</font>**<font style="color:rgb(44, 44, 54);">：提取函数调用中的JSON参数。</font>

<font style="color:rgb(44, 44, 54);">这些解析器专为需要调用特定函数或工具的场景设计，常见于Agent系统或复杂任务处理。</font>

**时间/枚举解析器**

<font style="color:rgb(44, 44, 54);">处理特定类型的数据：</font>

+ **<font style="color:rgb(44, 44, 54);">DatetimeOutputParser</font>**：将文本解析为标准日期时间格式（"%Y-%m-%dT%H:%M:%S.%fZ"）。

```python
from langchain.output_parsers import DatetimeOutputParser
from langchain.prompts import ChatPromptTemplate

from app.common import llm

# 创建DatetimeOutputParser
output_parser = DatetimeOutputParser()

# 获取格式指令
format_instructions = output_parser.get_format_instructions()

# 创建提示模板
prompt = ChatPromptTemplate.from_messages([
    ("system", f"你必须按照以下格式返回日期时间：{format_instructions}"),
    ("human", "请将以下自然语言转换为标准日期时间格式：\n###{text}###")
])


chain = prompt | llm | output_parser

resp = chain.invoke({"text": "转换一下：2025年5月15日"})

print(resp)
print(type(resp))
```

+ **<font style="color:rgb(44, 44, 54);">EnumOutputParser</font>**<font style="color:rgb(44, 44, 54);">：将文本解析为预定义的枚举值，适用于需要严格限制输出选项的场景。</font>

```python
from enum import Enum
from langchain.output_parsers import EnumOutputParser
from langchain_core.prompts import ChatPromptTemplate

from app.common import llm


# 定义枚举类
class Color(Enum):
    RED = "red"
    GREEN = "green"
    BLUE = "blue"

# 创建解析器
output_parser = EnumOutputParser(enum=Color)

# 提示模板
prompt = ChatPromptTemplate.from_messages([
    ("system", "请选择一种颜色："),
    ("human", "选项：{options}，直接返回颜色文本")
])

# 构建链
chain = prompt | llm | output_parser

# 测试输入
result = chain.invoke({"options": "red, green, blue"})
print("解析后的枚举值：", result)
print(type(result))
```

<font style="color:rgb(44, 44, 54);">这些解析器针对特定数据类型进行了优化，提供更精准的解析结果。</font>

**正则解析器**

<font style="color:rgb(44, 44, 54);">通过正则表达式提取特定模式的数据：</font>

+ **<font style="color:rgb(44, 44, 54);">RegexParser</font>**<font style="color:rgb(44, 44, 54);">：单字段正则匹配。</font>
+ **<font style="color:rgb(44, 44, 54);">RegexDictParser</font>**<font style="color:rgb(44, 44, 54);">：多字段正则提取为字典。</font>

<font style="color:rgb(44, 44, 54);">正则解析器适合需要灵活匹配文本模式的场景，如提取电话号码、邮箱地址等。</font>

**复合解析器**

<font style="color:rgb(44, 44, 54);">处理复杂或多部分的输出：</font>

+ **<font style="color:rgb(44, 44, 54);">CombiningOutputParser</font>**<font style="color:rgb(44, 44, 54);">：组合多个解析器，将文本通过分隔符拆分后分别解析。</font>
+ **<font style="color:rgb(44, 44, 54);">RetryWithErrorOutputParser</font>**<font style="color:rgb(44, 44, 54);">：在解析失败时重试并尝试修复输出。</font>

<font style="color:rgb(44, 44, 54);">复合解析器提供了更强大的处理能力，可以应对复杂或多任务的输出场景。</font>

**强校验解析器**

<font style="color:rgb(44, 44, 54);">提供更严格的输出验证：</font>

+ **<font style="color:rgb(44, 44, 54);">GuardrailsOutputParser</font>**<font style="color:rgb(44, 44, 54);">：基于Guardrails库，对输出进行强校验，支持自定义规则（如过滤不当内容、校验数据格式）。</font>

Guardrails解析器特别适合安全敏感场景，确保输出内容符合预设规则和标准。

### 代码示例

| 解析器                           | 用途       | 输入示例     | 输出类型    |
| :------------------------------- | :--------- | :----------- | :---------- |
| `StrOutputParser`                | 纯文本     | "你好"       | `str`       |
| `CommaSeparatedListOutputParser` | 逗号列表   | "a,b,c"      | `list[str]` |
| `JsonOutputParser`               | 通用JSON   | `{"a":1}`    | `dict`      |
| `PydanticOutputParser`           | 强类型结构 | JSON         | `BaseModel` |
| `DatetimeOutputParser`           | 日期时间   | "2025-05-15" | `datetime`  |
| `EnumOutputParser`               | 枚举值     | "positive"   | `Enum`      |
| `BooleanOutputParser`            | 布尔值     | "yes/no"     | `bool`      |

```python
from pydantic import SecretStr
from langchain_core.output_parsers import StrOutputParser
from langchain_openai import ChatOpenAI
from langchain_core.prompts import ChatPromptTemplate

llm = ChatOpenAI(
    model="kimi-k2-turbo-preview",
    base_url="https://api.moonshot.cn/v1",
    api_key=SecretStr("sk-32lqdlnpX9C7KiMdeoepSaIO6B79BqbyL"),
    temperature=1,
)

# 1. StrOutputParser（纯文本）
str_parser = StrOutputParser()

prompt = ChatPromptTemplate.from_messages([
    ("human", "讲一个笑话")
])

chain = prompt | llm | str_parser

result = chain.invoke({})
print(f"结果: {result}")
print(f"类型: {type(result)}")
# 结果: 老师：“为什么地球是圆的，你却总说世界是平的？”  学生：“老师，您不懂，我说的是‘余额’——我的支付宝余额。”
# 类型: <class 'str'>

# 2. CommaSeparatedListOutputParser（逗号分隔列表）
from langchain_core.output_parsers import CommaSeparatedListOutputParser
from langchain_core.prompts import ChatPromptTemplate

# 创建解析器
comma_parser = CommaSeparatedListOutputParser()

# 获取格式指令（自动添加到提示词）
format_instructions = comma_parser.get_format_instructions()
print(f"格式指令: {format_instructions}")

prompt = ChatPromptTemplate.from_messages([
    ("system", f"你是一个列表生成器。{format_instructions}"),
    ("human", "列出3种水果")
])

chain = prompt | llm | comma_parser

result = chain.invoke({})
print(f"结果: {result}")
print(f"类型: {type(result)}")
# 格式指令: Your response should be a list of comma separated values, eg: `foo, bar, baz` or `foo,bar,baz`
# 结果: ['苹果', '香蕉', '橙子']
# 类型: <class 'list'>


# 3. PydanticOutputParser（结构化数据 - 最常用）
from langchain_core.output_parsers import PydanticOutputParser
from langchain_core.prompts import ChatPromptTemplate
from pydantic import BaseModel, Field

# 定义输出结构
class Person(BaseModel):
    name: str = Field(description="姓名")
    age: int = Field(description="年龄")
    hobbies: list[str] = Field(description="爱好列表")
# 创建解析器
parser = PydanticOutputParser(pydantic_object=Person)
raw_instructions = parser.get_format_instructions()
#escaped_instructions = raw_instructions.replace("{T)

# 现在可以安全地使用 f-string
prompt = ChatPromptTemplate.from_messages([
    ("system", f"提取人物信息。{escaped_instructions}"),
    ("human", "{input}")
])

chain = prompt | llm | parser

result = chain.invoke({"input": "张三今年25岁，喜欢游泳和编程"})
print(f"结果: {result}")
# 结果: name='张三' age=25 hobbies=['游泳', '编程']

# 4. JsonOutputParser（通用JSON）
from langchain_core.output_parsers import JsonOutputParser
from langchain_core.prompts import ChatPromptTemplate


json_parser = JsonOutputParser()
system_msg = "提取信息。" + json_parser.get_format_instructions()

prompt = ChatPromptTemplate.from_messages([
    ("system", system_msg),   # 等效于：("system", "提取信息。Return a JSON object.")
    ("human", "{input}")
])

chain = prompt | llm | json_parser

result = chain.invoke({})
print(f"结果: {result}")
print(f"类型: {type(result)}")
# 结果: {'name': '张三', 'age': 25, 'hobbies': ['游泳', '编程']}
# 类型: <class 'dict'>

# 等价
from langchain_core.output_parsers import JsonOutputParser
from langchain_core.prompts import ChatPromptTemplate

json_parser = JsonOutputParser()
# 关键：在 system 中明确指定只返回 name 和 age
prompt = ChatPromptTemplate.from_messages([
    ("system", """
请提取信息并返回JSON格式，只包含以下两个字段：
{
    "name": "姓名",
    "age": 年龄数字
}
不要包含其他字段如 hobbies、address 等。
"""),
    ("human", "张三今年100岁，喜欢游泳和编程")
])

chain = prompt | llm | json_parser

result = chain.invoke({})
print(f"结果: {result}")
print(f"类型: {type(result)}")
# 输出: {'name': '张三', 'age': 100}

```

## 综合代码

```python
from langchain.agents import create_agent
from langchain_core.messages import HumanMessage
from langchain_core.runnables import RunnableLambda
from pydantic import SecretStr
from langchain_core.output_parsers import StrOutputParser
from langchain_openai import ChatOpenAI
from langchain_core.prompts import ChatPromptTemplate

from langgraph.prebuilt import create_react_agent
from langchain_core.output_parsers import PydanticOutputParser
from langchain_core.tools import tool
from langchain_core.messages import HumanMessage
from langchain_core.runnables import RunnableLambda, RunnablePassthrough
from langchain_openai import ChatOpenAI
from pydantic import BaseModel, Field
from typing import List
import json

# ========== 定义数据结构和工具（同上）==========
class WeatherReport(BaseModel):
    city: str = Field(description="城市名称")
    temperature: int = Field(description="当前温度")
    condition: str = Field(description="天气状况")
    humidity: int = Field(description="湿度百分比")
    wind_level: str = Field(description="风力等级")
    advice: List[str] = Field(description="3条具体出行建议")
    warning: str = Field(description="天气预警信息", default="无")

parser = PydanticOutputParser(pydantic_object=WeatherReport)
format_instructions = parser.get_format_instructions()
#escaped_instructions = format_instructions.replace("TODO")


@tool
def get_weather(city: str) -> str:
    """
    获取指定城市的天气信息。

    Args:
        city (str): 城市名称。

    Returns:
        str: 包含天气信息的 JSON 字符串，字段包括温度(temp)、天气状况(condition)、湿度(humidity)和风力(wind)。
    """
    weather_data = {
        "北京": {"temp": 25, "condition": "晴天", "humidity": 45, "wind": "3级"},
        "上海": {"temp": 28, "condition": "多云", "humidity": 60, "wind": "2级"},
        "哈尔滨": {"temp": -5, "condition": "雪天", "humidity": 80, "wind": "5级"},
    }
    data = weather_data.get(city, {"temp": 20, "condition": "晴天", "humidity": 50, "wind": "2级"})
    #.get(key, default) 是 Python 字典的安全取值方法，当 key 不存在时不会报错，而是返回指定的默认值
    return json.dumps(data, ensure_ascii=False)

# ========== Agent ==========
# llm
llm = ChatOpenAI(
    model="kimi-k2-turbo-preview",
    base_url="https://api.moonshot.cn/v1",
    api_key=SecretStr("sk-32lqdiwARGYCAaRPa7KiMdeoepSaIO6B79BqbyL"),
    temperature=1,
)

# 1. 创建 Agent

SYSTEM_PROMPT = f"""你是专业的天气分析师...{escaped_instructions}..."""  # 等价"你是专业的天气分析师...{{city}}: {{temperature}}..."
agent = create_agent(
    model=llm,
    tools=[get_weather],  # 工具列表
    system_prompt=SYSTEM_PROMPT
)

# 2. 定义 chain 组件
# ========== 定义每个步骤的函数 ==========
def extract_content(data: dict) -> str:
    """
    第2步：从 Agent 输出中提取最后一条 AI 消息的内容

    输入示例:
    {
        "messages": [
            HumanMessage(content="北京天气怎么样？"),
            AIMessage(content="", tool_calls=[...]),
            ToolMessage(content="{\"temp\": 25...}"),
            AIMessage(content="```json\n{\"city\": \"北京\"...}\n```")
        ]
    }

    操作: 取最后一个元素的 content
    """
    last_message = data["messages"][-1]
    content = last_message.content

    print(f"【extract_content】输入类型: {type(data)}")
    print(f"【extract_content】输入键: {data.keys()}")
    print(f"【extract_content】提取内容: {content[:50]}...")

    return content  # 输出: 字符串


def clean_json(content: str) -> str:
    """
    第3步：清理 Markdown 标记，提取纯 JSON

    输入示例:
    "```json\n{\n  \"city\": \"北京\",\n  \"temperature\": 25,\n  ...\n}\n```"

    操作: 删除 ```json 和 ``` 标记
    """
    print(f"【clean_json】输入类型: {type(content)}")
    print(f"【clean_json】原始内容:\n{content}")

    # 清理 Markdown
    if "```json" in content:
        content = content.split("```json")[1].split("```")[0]
    elif "```" in content:
        content = content.split("```")[1].split("```")[0]

    result = content.strip()
    print(f"【clean_json】清理后:\n{result}")

    return result  # 输出: 干净 JSON 字符串


def add_metadata(weather: WeatherReport) -> dict:
    """
    第5步：添加元数据，包装为最终输出格式

    输入示例:
    WeatherReport(
        city="北京",
        temperature=25,
        condition="晴天",
        ...
    )

    操作: 包装为 dict 并添加额外信息
    """
    print(f"【add_metadata】输入类型: {type(weather)}")
    print(f"【add_metadata】输入对象: {weather}")

    result = {
        "data": weather,  # 核心数据
        "query_time": "2024-01-01 10:30:00",  # 查询时间
        "source": "weather_api",  # 数据来源
        "version": "v1.0"  # 数据版本
    }

    print(f"【add_metadata】输出类型: {type(result)}")
    print(f"【add_metadata】输出内容: {result}")

    return result  # 输出: dict


# ========== 构建 chain ==========
weather_chain = (
        agent  # 第1步: Runnable (Agent)
        | RunnableLambda(extract_content)  # 第2步: RunnableLambda (函数 → Runnable)
        | RunnableLambda(clean_json)  # 第3步: RunnableLambda
        | parser  # 第4步: PydanticOutputParser (Runnable)  # 使用解析器
        | RunnableLambda(add_metadata)  # 第5步: RunnableLambda
)

# ========== 执行并观察数据流 ==========
print("=" * 60)
print("开始执行 chain")
print("=" * 60)

final_result = weather_chain.invoke({
    "messages": [HumanMessage(content="北京天气怎么样？适合出门吗？")]
})

print("\n" + "=" * 60)
print("最终结果")
print("=" * 60)
print(f"类型: {type(final_result)}")
print(f"内容: {json.dumps(final_result, default=str, ensure_ascii=False, indent=2)}")
# 4. 使用 chain（一行代码完成全部）
result = weather_chain.invoke({
    "messages": [HumanMessage(content="北京天气怎么样？")]
})

print(f"结构化数据: {result['data']}")
print(f"查询时间: {result['query_time']}")
```

````
============================================================
开始执行 chain
============================================================
【extract_content】输入类型: <class 'dict'>
【extract_content】输入键: dict_keys(['messages'])
【extract_content】提取内容: ### 思考过程：

1. **理解用户需求**：用户询问北京的天气情况，并想知道是否适合出门。这需...
【clean_json】输入类型: <class 'str'>
【clean_json】原始内容:
### 思考过程：

1. **理解用户需求**：用户询问北京的天气情况，并想知道是否适合出门。这需要获取北京的实时天气数据，并根据天气条件给出出行建议。

2. **获取天气数据**：使用`get_weather`函数获取北京的天气信息。调用时需要传入城市名称“北京”。

3. **分析天气数据**：从返回的JSON数据中提取关键信息：
   - 温度（temp）
   - 天气状况（condition）
   - 湿度（humidity）
   - 风力（wind）

4. **生成出行建议**：根据天气状况和温度给出具体的建议。例如：
   - 如果温度高且晴朗，建议防晒和补水。
   - 如果下雨或大风，建议带雨具或避免外出。
   - 如果空气质量差（如雾霾），建议戴口罩。

5. **检查天气预警**：如果有极端天气（如暴雨、大风、高温等），需要额外发出警告。

6. **格式化输出**：按照给定的JSON Schema整理数据和建议。

### 实施步骤：

1. 调用`get_weather("北京")`获取天气数据。
2. 解析返回的数据，假设得到：
   - temp: 25
   - condition: "晴"
   - humidity: 30
   - wind: "2级"
3. 根据“晴”和25°C，适合出门，建议：
   - 防晒（如戴帽子或涂防晒霜）
   - 补水（带水瓶）
   - 穿轻薄衣物
4. 无极端天气，警告为“无”。
5. 填充JSON Schema。

### 最终答案：
```json
{
  "city": "北京",
  "temperature": 25,
  "condition": "晴",
  "humidity": 30,
  "wind_level": "2级",
  "advice": [
    "天气晴朗，适合出门，建议做好防晒措施（如戴帽子或涂防晒霜）",
    "温度适中，建议穿轻薄衣物并随身携带水瓶及时补水",
    "风力较小，出行无需特别防风准备"
  ],
  "warning": "无"
}
```
【clean_json】清理后:
{
  "city": "北京",
  "temperature": 25,
  "condition": "晴",
  "humidity": 30,
  "wind_level": "2级",
  "advice": [
    "天气晴朗，适合出门，建议做好防晒措施（如戴帽子或涂防晒霜）",
    "温度适中，建议穿轻薄衣物并随身携带水瓶及时补水",
    "风力较小，出行无需特别防风准备"
  ],
  "warning": "无"
}
【add_metadata】输入类型: <class '__main__.WeatherReport'>
【add_metadata】输入对象: city='北京' temperature=25 condition='晴' humidity=30 wind_level='2级' advice=['天气晴朗，适合出门，建议做好防晒措施（如戴帽子或涂防晒霜）', '温度适中，建议穿轻薄衣物并随身携带水瓶及时补水', '风力较小，出行无需特别防风准备'] warning='无'
【add_metadata】输出类型: <class 'dict'>
【add_metadata】输出内容: {'data': WeatherReport(city='北京', temperature=25, condition='晴', humidity=30, wind_level='2级', advice=['天气晴朗，适合出门，建议做好防晒措施（如戴帽子或涂防晒霜）', '温度适中，建议穿轻薄衣物并随身携带水瓶及时补水', '风力较小，出行无需特别防风准备'], warning='无'), 'query_time': '2024-01-01 10:30:00', 'source': 'weather_api', 'version': 'v1.0'}

============================================================
最终结果
============================================================
类型: <class 'dict'>
内容: {
  "data": "city='北京' temperature=25 condition='晴' humidity=30 wind_level='2级' advice=['天气晴朗，适合出门，建议做好防晒措施（如戴帽子或涂防晒霜）', '温度适中，建议穿轻薄衣物并随身携带水瓶及时补水', '风力较小，出行无需特别防风准备'] warning='无'",
  "query_time": "2024-01-01 10:30:00",
  "source": "weather_api",
  "version": "v1.0"
}
【extract_content】输入类型: <class 'dict'>
【extract_content】输入键: dict_keys(['messages'])
【extract_content】提取内容: ```json
{
  "city": "北京",
  "temperature": 25,
  "...
【clean_json】输入类型: <class 'str'>
【clean_json】原始内容:
```json
{
  "city": "北京",
  "temperature": 25,
  "condition": "晴天",
  "humidity": 45,
  "wind_level": "3级",
  "advice": [
    "天气晴朗，适合户外活动，建议做好防晒措施",
    "温度适宜，建议穿着轻便的夏季服装",
    "湿度较低，注意补充水分，多喝水"
  ]
}
```
【clean_json】清理后:
{
  "city": "北京",
  "temperature": 25,
  "condition": "晴天",
  "humidity": 45,
  "wind_level": "3级",
  "advice": [
    "天气晴朗，适合户外活动，建议做好防晒措施",
    "温度适宜，建议穿着轻便的夏季服装",
    "湿度较低，注意补充水分，多喝水"
  ]
}
【add_metadata】输入类型: <class '__main__.WeatherReport'>
【add_metadata】输入对象: city='北京' temperature=25 condition='晴天' humidity=45 wind_level='3级' advice=['天气晴朗，适合户外活动，建议做好防晒措施', '温度适宜，建议穿着轻便的夏季服装', '湿度较低，注意补充水分，多喝水'] warning='无'
【add_metadata】输出类型: <class 'dict'>
【add_metadata】输出内容: {'data': WeatherReport(city='北京', temperature=25, condition='晴天', humidity=45, wind_level='3级', advice=['天气晴朗，适合户外活动，建议做好防晒措施', '温度适宜，建议穿着轻便的夏季服装', '湿度较低，注意补充水分，多喝水'], warning='无'), 'query_time': '2024-01-01 10:30:00', 'source': 'weather_api', 'version': 'v1.0'}
结构化数据: city='北京' temperature=25 condition='晴天' humidity=45 wind_level='3级' advice=['天气晴朗，适合户外活动，建议做好防晒措施', '温度适宜，建议穿着轻便的夏季服装', '湿度较低，注意补充水分，多喝水'] warning='无'
查询时间: 2024-01-01 10:30:00
````

