from langchain_ollama import ChatOllama
from langchain_core.messages import HumanMessage
from langchain_core.prompts import PromptTemplate, ChatPromptTemplate

# 支持对话格式，功能更完整
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



#提示词
template = PromptTemplate.from_template("今天{somethingName}真不错")
result = template.format(somethingName="天气")
print(result)  # 今天天气真不错

chat_prompt = ChatPromptTemplate.from_messages([
    ("system", "你是{role}专家，擅长{domain}"),
    ("human", "{question}")
])

messages = chat_prompt.format_messages(
    role="技术",
    domain="开发",
    question="如何学习？"
)
print(messages)
