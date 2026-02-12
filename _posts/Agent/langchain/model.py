from langchain_openai import ChatOpenAI
from langchain_ollama import ChatOllama
from pydantic import SecretStr

# 本地模型
llm = ChatOllama(
    model="deepseek-r1:1.5b",
    base_url="http://localhost:11434",  # 默认地址
    temperature=0.7,
)

# kimi
KimiLlm = ChatOpenAI(
    model="kimi-k2-turbo-preview",
    base_url="https://api.moonshot.cn/v1",
    api_key=SecretStr("sk-32lqdlnpX9Cgc7iwARGYCAaRPa7KiMdeoepSaIO6B79BqbyL"),
    temperature=0.7,
)

# 单轮
response = KimiLlm.invoke("北京的首都是什么？")
print(response.content)