---
title: LangChain-Agent解析
date: 2026-02-12 16:00:00 +0800
categories: [Agent, LangChain-Agent解析]
tags: [LangChain-Agent解析]
---

# LangChain-Agent 解析

## Agent

**Agent（代理）** 是 LangChain 中更高级的机制，允许 LLM 像“智能助手”一样，根据输入动态选择和使用工具，完成复杂任务。代理的核心是“推理 + 行动”，它会分析输入，决定需要调用哪些工具，并循环执行直到任务完成。

**组成**：

- **工具（Tools）**：如搜索引擎、计算器、API。
- **代理执行器**：协调 LLM 和工具，管理执行流程。
- **推理引擎**：LLM 负责决定使用哪个工具。

**应用场景**：需要动态决策的任务，如联网搜索、数学计算。

## 工具（Tools）

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

