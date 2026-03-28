---
title: protojson解析库
date: 2025-11-29 11:00:00 +0800
categories: [GO, 解析库]
tags: [protojson解析库]
---

## json 和 protojson使用

**不是直接使用encoding/json库的原因:**

| 序列化器            | 能否用于 Protobuf GO 结构体 | 输出标准 Protobuf JSON | 能正确处理特殊 proto 类型 |
| ------------------- | --------------------------- | ---------------------- | ------------------------- |
| `encoding/json`     | 可以，但有限                | 否                     | 否                        |
| `protojson.Marshal` | 可以                        | 是                     | 是                        |

**解析区别**

```go
MsgId string `protobuf:"bytes,1,opt,name=msg_id,json=msgId,proto3"  json:"msg_id,omitempty"`
```

| 使用场景                         | 使用的标签                     | 使用字段 | 示例              |
| :------------------------------- | :----------------------------- | :------- | :---------------- |
| protojson + UseProtoNames: true  | protobuf 标签中的name=msg_id   | "msg_id" | {"msg_id": "123"} |
| protojson + UseProtoNames: false | protobuf 标签中的   json=msgId | "msgId"  | {"msgId": "123"}  |
| encoding/json                    | json:"msg_id" 标签             | "msg_id" | {"msg_id": "123"} |

### protojson 库

相关的包：

```go
google.golang.org/protobuf/encoding/protojson
```

语法

```go
//
defaultOpts := protojson.MarshalOptions{}
data2, _ := defaultOpts.Marshal(待解析)
// 
protojson.Unmarshal([]byte(待解析proto数据), &结果)
```

### 案例

 **protobuf 定义**

```protobuf
message ActionInfo {
    string action_id = 1;        // protobuf 字段名：action_id
    bool only_client = 10;        // protobuf 字段名：only_client
    ActionType action_type = 5;  // 枚举类型
}

enum ActionType {
  // 未指定的动作类型
  ACTION_TYPE_UNSPECIFIED = 0;
  // 发送消息动作
  ACTION_TYPE_SEND_MSG = 1;
}
```

**生成的go代码**

```go
type ActionInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 动作类型
	ActionType ActionType `protobuf:"varint,1,opt,name=action_type,json=actionType,proto3,enum=xdproto.common_business.msg_adapter.ActionType" json:"action_type,omitempty"`
	// 内部生成的 ID
	ActionId string `protobuf:"bytes,9,opt,name=action_id,json=actionId,proto3" json:"action_id,omitempty"`
	// 是否只能客户端发送
	OnlyClient bool `protobuf:"varint,10,opt,name=only_client,json=onlyClient,proto3" json:"only_client,omitempty"`
}

type ActionType int32
const (
	// 未指定的动作类型
	ActionType_ACTION_TYPE_UNSPECIFIED ActionType = 0
	// 发送消息动作
	ActionType_ACTION_TYPE_SEND_MSG ActionType = 1
)
```

**使用配置（UseProtoNames: true, UseEnumNumbers: true）**

```20:24:component/action_receipt_report/timer.go
var (
	marshalOpts = protojson.MarshalOptions{
		UseProtoNames:  true,
		UseEnumNumbers: true,
	}
)
```

序列化结果：

```json
"action_info": {                 // ✅ 使用 protobuf 原始字段名 action_info
    "action_id": "action_456",      // ✅ 嵌套的 ActionInfo 也使用 action_id
    "only_client": true,          // ✅ 使用 protobuf 原始字段名 only_client
    "action_type": 1             // ✅ 枚举值使用数字 1（不是字符串）
}
```

**使用默认配置（UseProtoNames: false, UseEnumNumbers: false）**

```go
// 默认配置
defaultOpts := protojson.MarshalOptions{}
data, _ := defaultOpts.Marshal(callBackMsg)
```

序列化结果：

```json
"actionInfo": {                   // ❌ 使用 JSON 命名（驼峰式）
    "actionId": "action_456",        // ❌ 嵌套的 ActionInfo 也使用 actionId
    "onlyClient": true,            // ❌ 使用 JSON 命名（驼峰式）
    "actionType": "ACTION_TYPE_TEXT"   // ❌ 枚举值使用字符串名称
} 
```

**总结**

![image-20251129175112804](https://github.com/Auroraol/Drawing-bed/raw/main/img/image-20251129175112804.png)







```
body: { "msg_id": "1774405167106521", "request_id": "", "user_name": "顺然旗舰店", "vender_id": "", "app": "mindflow", "shop_id": "662f5c2dc38b2037c947a2bc", "from_pin": "顺然旗舰店:慕斯", "to_pin": "one_id_299502931", "action_list": [{ "action_type": "ACTION_TYPE_SEND_MSG", "timestamp": "1774405139496", "timeout": 1, "send_msg_action": { "msg_info": { "msg_type": "SEND_MSG_TYPE_TEXT", "msg_sub_type": "SEND_MSG_SUB_TYPE_DEFAULT", "text": { "content": "好的呢~" }, "ref_msg_id": "" }, "delay": 0, "confirm": false, "is_not_clear_timing": false, "is_not_clear_light": false }, "business": "BUSINESS_ROBOT", "action_id": "900322a1-ddf8-483d-b295-2b889e8f76e1", "only_client": false, "extra_data": { "channel": "pulsar", "msg_scenes_source": "7.18.0" } }], "platform": "tb", "plat_shop_id": "662f5c2dc38b2037c947a2bc", "chat_type": "CHAT_SESSION_TYPE_UNKNOWN", "plat_session_id": "" }
```

```
// protobuf tag 中的定义
protobuf:"bytes,1,opt,name=msg_id,json=msgId,proto3"
//                      ↑            ↑
//                   原始名       JSON 映射名
```

protojson.Unmarshal 的智能行为
关键点：protojson 会同时尝试匹配两种格式！
优先匹配 JSON 名称（json=msgId）
如果找不到，会尝试 protobuf 原始名称（name=msg_id）

`rotojson.Unmarshal` 确实会**同时尝试两种名称格式**，并且遵循**JSON 名称优先**的匹配策略。

匹配优先级（从高到低）

| 优先级 | 名称类型              | 来源                             | 示例     |
| :----- | :-------------------- | :------------------------------- | :------- |
| **1**  | **JSON 名称**         | `json_name` 选项 或 驼峰自动转换 | `msgId`  |
| **2**  | **Protobuf 原始名称** | `.proto` 文件中的字段名          | `msg_id` |



# 未知字段

默认配置（当前代码）

err = protojson.Unmarshal(body, req)

DiscardUnknown: false（默认）
行为：会报错，提示遇到未知字段

显式忽略未知字段

unmarshalOpt := protojson.UnmarshalOptions{
    DiscardUnknown: true,  // ← 忽略未知字段
}
err = unmarshalOpt.Unmarshal(body, req)

行为：不会报错，直接忽略 unknown_field

严格模式

unmarshalOpt := protojson.UnmarshalOptions{
    DiscardUnknown: false,  // 不忽略
    AllowPartial:   false,  // 严格要求所有必填字段
}
err = unmarshalOpt.Unmarshal(body, req)

行为：会报错，如果有未知字段或缺少必填字段
