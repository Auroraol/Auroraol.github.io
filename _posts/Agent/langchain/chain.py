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