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
    api_key=SecretStr("sk-32lqdlnpX9Cgc7iwARGYCAaRPa7KiMdeoepSaIO6B79BqbyL"),
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
