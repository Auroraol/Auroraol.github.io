# httpclient

`httpclient` 是一个基于 [resty](https://gopkg.in/resty.v1) 封装的 HTTP 客户端库，为 Go 开发者提供便捷、功能丰富的 HTTP 请求发送能力。该库内置了日志记录、DNS 追踪、指标监控、自动重试、OpenTelemetry 追踪等功能。

## 特性

- 🚀 **简单易用**：提供简洁的 API，快速上手
- 🔄 **自动重试**：支持可配置的重试机制
- 📊 **内置监控**：自动收集请求延迟和计数指标（Prometheus）
- 📝 **日志记录**：内置请求和响应日志记录功能
- 🔍 **DNS 追踪**：支持 DNS 解析过程追踪
- 🔗 **OpenTelemetry 集成**：自动支持分布式追踪
- 🎣 **钩子机制**：灵活的请求前后处理钩子
- 📦 **Metadata 传递**：自动处理上下文 metadata 传递

## 安装

```shell
go get -u gitlab.xiaoduoai.com/golib/xd_sdk/httpclient
```

## 快速开始

### 基础使用

```go
package main

import (
    "context"
    "time"
    
    "gitlab.xiaoduoai.com/golib/xd_sdk/httpclient"
    "gitlab.xiaoduoai.com/golib/xd_sdk/httpclient/hooks"
)

func main() {
    // 创建客户端
    c, err := httpclient.NewClient(
        httpclient.WithAddress("https://api.example.com"),
        httpclient.WithTimeout(5 * time.Second),
    )
    if err != nil {
        panic(err)
    }
    
    ctx := context.Background()
    
    // GET 请求
    var result map[string]interface{}
    err = c.Get(ctx, "/api/users", nil, nil, 0, &result)
    if err != nil {
        panic(err)
    }
    
    // POST JSON 请求
    reqData := map[string]string{"name": "John"}
    err = c.PostJSON(ctx, "/api/users", reqData, nil, 0, &result)
    if err != nil {
        panic(err)
    }
}
```

### 使用 NewRequest 进行更灵活的请求

```go
ctx := context.Background()
response, err := c.NewRequest(ctx).
    SetHeader("Authorization", "Bearer token").
    SetQueryParam("page", "1").
    Get("/api/users")

if err != nil {
    panic(err)
}

// 获取响应数据
var users []User
err = json.Unmarshal(response.Body(), &users)
```

## NewRequest错误类型说明

返回以下类型的错误：

1. **网络连接错误**
   - `*net.OpError`：网络操作失败
     - 连接被拒绝（`connection refused`）：目标服务器未运行或端口未开放
     - 网络不可达（`no route to host`）：无法路由到目标主机
     - 连接超时（`i/o timeout`）：建立连接超时
   - `*net.DNSError`：DNS 解析失败，域名无法解析

2. **超时错误**
   - `context.DeadlineExceeded`：请求超过设定的超时时间
   - 实现了 `Timeout()` 方法的错误：操作超时

3. **TLS/SSL 错误**
   - `*tls.RecordHeaderError`：TLS 握手失败
   - `*x509.UnknownAuthorityError`：证书验证失败（自签名证书或未知 CA）

4. **上下文取消错误**
   - `context.Canceled`：请求被上下文取消

5. **钩子函数错误**
   - `PreRequestHook` 执行失败时返回的错误（会被包装为 `failed to execute pre request hook`）

6. **其他错误**
   - URL 格式错误
   - 请求体序列化错误
   - 其他底层 HTTP 客户端错误

**注意：** HTTP 状态码错误（如 404、500）通常不会作为 `error` 返回，而是需要通过 `response.StatusCode()` 检查。只有在网络层面或请求准备阶段出现问题时才会返回错误。

**错误处理示例：**

```go
ctx := context.Background()
response, err := c.NewRequest(ctx).
    SetHeader("Authorization", "Bearer token").
    Get("/api/users")

if err != nil {
    // 检查是否为连接拒绝错误
    if httpclient.IsConnectRefuseError(err) {
        log.Println("无法连接到服务器，请检查服务是否运行")
        return
    }
    
    // 检查是否为超时错误
    if err == context.DeadlineExceeded {
        log.Println("请求超时")
        return
    }
    
    // 检查是否为 DNS 错误
    var dnsErr *net.DNSError
    if errors.As(err, &dnsErr) {
        log.Printf("DNS 解析失败: %v", dnsErr)
        return
    }
    
    // 其他错误
    log.Printf("请求失败: %v", err)
    return
}

// 检查 HTTP 状态码
if response.StatusCode() >= 400 {
    log.Printf("HTTP 错误: %d, 响应: %s", response.StatusCode(), response.String())
    return
}

// 处理成功响应
var users []User
err = json.Unmarshal(response.Body(), &users)
```

## 详细配置

### 客户端选项

`httpclient` 提供了丰富的配置选项，通过函数式选项模式进行配置：

#### WithAddress

设置基础 URL 地址。

```go
c, err := httpclient.NewClient(
    httpclient.WithAddress("https://api.example.com"),
)
```

#### WithTimeout

设置请求超时时间。

```go
c, err := httpclient.NewClient(
    httpclient.WithTimeout(10 * time.Second),
)
```

#### WithRetryCount

设置重试次数。

```go
c, err := httpclient.NewClient(
    httpclient.WithRetryCount(3), // 失败后重试 3 次
)
```

#### WithRetryWaitTime

设置重试间隔等待时间。

```go
c, err := httpclient.NewClient(
    httpclient.WithRetryWaitTime(500 * time.Millisecond),
)
```

#### WithRetryMaxWaitTime

设置重试间隔最大等待时间。

```go
c, err := httpclient.NewClient(
    httpclient.WithRetryMaxWaitTime(2 * time.Second),
)
```

#### WithMetrics

启用或禁用指标监控（默认启用）。

```go
c, err := httpclient.NewClient(
    httpclient.WithMetrics(true), // 启用 Prometheus 指标
)
```

#### WithParseSvcCodeFuncs

设置从响应中解析业务代码的函数，用于监控告警。

```go
parseCode := func(httpCode int, result interface{}) (svcCode int, ok bool) {
    if resp, ok := result.(*MyResponse); ok {
        return resp.Code, true
    }
    return 0, false
}

c, err := httpclient.NewClient(
    httpclient.WithParseSvcCodeFuncs(parseCode),
)
```

#### WithParseMetricPathF

自定义路径解析函数，用于处理带参数的路径（如 `/users/:id`）。

```go
parsePath := func(originPath string) string {
    // 将 /users/123 转换为 /users/:id
    return regexp.MustCompile(`/\d+`).ReplaceAllString(originPath, "/:id")
}

c, err := httpclient.NewClient(
    httpclient.WithParseMetricPathF(parsePath),
)
```

### 钩子函数

#### BeforeRequestHooks

在请求发送前执行的钩子函数。

```go
beforeHook := func(c *resty.Client, r *resty.Request) error {
    // 在请求前执行的操作
    r.SetHeader("X-Custom-Header", "value")
    return nil
}

c, err := httpclient.NewClient(
    httpclient.WithBeforeRequestHooks(beforeHook),
)
```

#### PreRequestHooks

在请求准备阶段执行的钩子函数（在 BeforeRequest 之后）。

```go
c, err := httpclient.NewClient(
    httpclient.WithPreRequestHooks(
        hooks.LoggingRequest(),  // 记录请求日志
        hooks.DNSTrace(),        // DNS 追踪
    ),
)
```

#### AfterResponseHooks

在收到响应后执行的钩子函数。

```go
c, err := httpclient.NewClient(
    httpclient.WithAfterResponseHooks(
        hooks.LoggingResponse(), // 记录响应日志
    ),
)
```

### 完整配置示例

```go
c, err := httpclient.NewClient(
    httpclient.WithAddress("https://api.example.com"),
    httpclient.WithTimeout(5 * time.Second),
    httpclient.WithRetryCount(3),
    httpclient.WithRetryWaitTime(500 * time.Millisecond),
    httpclient.WithRetryMaxWaitTime(2 * time.Second),
    httpclient.WithMetrics(true),
    httpclient.WithPreRequestHooks(
        hooks.LoggingRequest(),
        hooks.DNSTrace(),
    ),
    httpclient.WithAfterResponseHooks(
        hooks.LoggingResponse(),
    ),
)
```

## API 文档

### Client 接口

#### NewRequest

创建一个新的请求对象，支持链式调用。

```go
request := c.NewRequest(ctx)
response, err := request.
    SetHeader("Content-Type", "application/json").
    SetBody(data).
    Post("/api/endpoint")
```

**参数：**
- `ctx context.Context`：上下文，用于传递 metadata 和取消请求

**返回：**
- `*resty.Request`：resty 请求对象，支持所有 resty 的功能

#### GetKernel

获取底层的 resty.Client，用于访问 resty 的完整功能。

```go
kernel := c.GetKernel()
kernel.SetProxy("http://proxy.example.com:8080")
```

**返回：**
- `*resty.Client`：底层的 resty 客户端

#### Get

发送 GET 请求。

```go
var result MyStruct
err := c.Get(ctx, "/api/users", values, headers, timeoutSeconds, &result)
```

**参数：**
- `ctx context.Context`：上下文
- `path string`：请求路径（相对于基础地址）
- `values url.Values`：查询参数
- `headers http.Header`：请求头
- `timeoutSeconds int`：超时时间（秒），0 表示使用客户端默认超时
- `ret interface{}`：用于接收响应数据的结构体指针

**返回：**
- `error`：错误信息

#### Post

发送 POST 请求（表单格式）。

```go
values := url.Values{}
values.Set("name", "John")
values.Set("email", "john@example.com")

var result MyStruct
err := c.Post(ctx, "/api/users", values, nil, 0, &result)
```

**参数：**
- `ctx context.Context`：上下文
- `path string`：请求路径
- `values url.Values`：表单数据
- `headers http.Header`：请求头
- `timeoutSeconds int`：超时时间（秒）
- `ret interface{}`：响应结构体指针

**返回：**
- `error`：错误信息

#### PostJSON

发送 POST 请求（JSON 格式）。

```go
reqData := map[string]interface{}{
    "name":  "John",
    "email": "john@example.com",
}

var result MyStruct
err := c.PostJSON(ctx, "/api/users", reqData, nil, 0, &result)
```

**参数：**
- `ctx context.Context`：上下文
- `path string`：请求路径
- `values interface{}`：JSON 数据（可以是结构体、map 等）
- `headers http.Header`：请求头
- `timeoutSeconds int`：超时时间（秒）
- `ret interface{}`：响应结构体指针

**返回：**
- `error`：错误信息

#### RawGet

发送 GET 请求并返回原始字节数据。

```go
data, err := c.RawGet(ctx, "/api/data", values, nil, 0)
if err != nil {
    panic(err)
}
// 处理原始数据
```

**参数：**
- `ctx context.Context`：上下文
- `path string`：请求路径
- `values url.Values`：查询参数
- `headers http.Header`：请求头
- `timeoutSeconds int`：超时时间（秒）

**返回：**
- `[]byte`：响应体原始数据
- `error`：错误信息

#### GetWithMock

发送 GET 请求，支持预设请求体（用于某些特殊场景）。

```go
var result MyStruct
err := c.GetWithMock(ctx, "/api/users", values, nil, 0, preSetBody, &result)
```

**参数：**
- `ctx context.Context`：上下文
- `path string`：请求路径
- `values url.Values`：查询参数
- `headers http.Header`：请求头
- `timeoutSeconds int`：超时时间（秒）
- `preSetRet interface{}`：预设的请求体
- `ret interface{}`：响应结构体指针

**返回：**
- `error`：错误信息

### 动态添加钩子

客户端创建后，还可以动态添加钩子：

```go
c.AddBeforeRequestHooks(customBeforeHook)
c.SetPreRequestHooks(customPreHook)
c.AddAfterResponseHooks(customAfterHook)
```

## 内置钩子函数

### LoggingRequest

记录请求日志。

```go
c, err := httpclient.NewClient(
    httpclient.WithPreRequestHooks(hooks.LoggingRequest()),
)
```

### LoggingResponse

记录响应日志。

```go
c, err := httpclient.NewClient(
    httpclient.WithAfterResponseHooks(hooks.LoggingResponse()),
)
```

### DNSTrace

DNS 解析追踪（默认已启用）。

```go
c, err := httpclient.NewClient(
    httpclient.WithPreRequestHooks(hooks.DNSTrace()),
)
```

### LatencyMetrics

延迟和计数指标收集（默认已启用，当 `WithMetrics(true)` 时）。

```go
// 通过 WithMetrics(true) 自动启用
c, err := httpclient.NewClient(
    httpclient.WithMetrics(true),
)
```

## 使用示例

### 示例 1：基础 GET 请求

```go
package main

import (
    "context"
    "fmt"
    "net/url"
    "time"
    
    "gitlab.xiaoduoai.com/golib/xd_sdk/httpclient"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func main() {
    c, err := httpclient.NewClient(
        httpclient.WithAddress("https://api.example.com"),
        httpclient.WithTimeout(5 * time.Second),
    )
    if err != nil {
        panic(err)
    }
    
    ctx := context.Background()
    
    // 带查询参数的 GET 请求
    values := url.Values{}
    values.Set("page", "1")
    values.Set("limit", "10")
    
    var users []User
    err = c.Get(ctx, "/api/users", values, nil, 0, &users)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("获取到 %d 个用户\n", len(users))
}
```

### 示例 2：POST JSON 请求

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "gitlab.xiaoduoai.com/golib/xd_sdk/httpclient"
)

type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserResponse struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    c, err := httpclient.NewClient(
        httpclient.WithAddress("https://api.example.com"),
        httpclient.WithTimeout(5 * time.Second),
    )
    if err != nil {
        panic(err)
    }
    
    ctx := context.Background()
    
    req := CreateUserRequest{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    var resp CreateUserResponse
    err = c.PostJSON(ctx, "/api/users", req, nil, 0, &resp)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("创建用户成功，ID: %d\n", resp.ID)
}
```

### 示例 3：使用 NewRequest 进行复杂请求

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "gitlab.xiaoduoai.com/golib/xd_sdk/httpclient"
)

func main() {
    c, err := httpclient.NewClient(
        httpclient.WithAddress("https://api.example.com"),
        httpclient.WithTimeout(5 * time.Second),
    )
    if err != nil {
        panic(err)
    }
    
    ctx := context.Background()
    
    // 使用 NewRequest 进行更灵活的配置
    response, err := c.NewRequest(ctx).
        SetHeader("Authorization", "Bearer your-token").
        SetHeader("X-Custom-Header", "custom-value").
        SetQueryParam("page", "1").
        SetQueryParam("limit", "20").
        SetBody(map[string]interface{}{
            "filter": "active",
        }).
        Post("/api/users/search")
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("状态码: %d\n", response.StatusCode())
    fmt.Printf("响应: %s\n", response.String())
}
```

### 示例 4：带重试和日志的完整配置

```go
package main

import (
    "context"
    "time"
    
    "gitlab.xiaoduoai.com/golib/xd_sdk/httpclient"
    "gitlab.xiaoduoai.com/golib/xd_sdk/httpclient/hooks"
)

func main() {
    c, err := httpclient.NewClient(
        httpclient.WithAddress("https://api.example.com"),
        httpclient.WithTimeout(5 * time.Second),
        httpclient.WithRetryCount(3),
        httpclient.WithRetryWaitTime(500 * time.Millisecond),
        httpclient.WithRetryMaxWaitTime(2 * time.Second),
        httpclient.WithMetrics(true),
        httpclient.WithPreRequestHooks(
            hooks.LoggingRequest(),
            hooks.DNSTrace(),
        ),
        httpclient.WithAfterResponseHooks(
            hooks.LoggingResponse(),
        ),
    )
    if err != nil {
        panic(err)
    }
    
    ctx := context.Background()
    
    // 使用客户端发送请求
    var result map[string]interface{}
    err = c.Get(ctx, "/api/data", nil, nil, 0, &result)
    if err != nil {
        panic(err)
    }
}
```

### 示例 5：自定义业务代码解析

```go
package main

import (
    "context"
    "time"
    
    "gitlab.xiaoduoai.com/golib/xd_sdk/httpclient"
    "gitlab.xiaoduoai.com/golib/xd_sdk/httpclient/hooks"
)

type APIResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    interface{} `json:"data"`
}

func parseBusinessCode(httpCode int, result interface{}) (svcCode int, ok bool) {
    if resp, ok := result.(*APIResponse); ok {
        return resp.Code, true
    }
    return 0, false
}

func main() {
    c, err := httpclient.NewClient(
        httpclient.WithAddress("https://api.example.com"),
        httpclient.WithTimeout(5 * time.Second),
        httpclient.WithMetrics(true),
        httpclient.WithParseSvcCodeFuncs(parseBusinessCode),
    )
    if err != nil {
        panic(err)
    }
    
    ctx := context.Background()
    
    var resp APIResponse
    err = c.Get(ctx, "/api/users", nil, nil, 0, &resp)
    if err != nil {
        panic(err)
    }
}
```

## 错误处理

### 检查连接拒绝错误

库提供了工具函数来检查是否为连接拒绝错误：

```go
err := c.Get(ctx, "/api/users", nil, nil, 0, &result)
if err != nil {
    if httpclient.IsConnectRefuseError(err) {
        // 处理连接拒绝错误
        fmt.Println("无法连接到服务器")
    } else {
        // 处理其他错误
        fmt.Printf("请求失败: %v\n", err)
    }
}
```

## 指标监控

当启用指标监控时（`WithMetrics(true)`），库会自动收集以下 Prometheus 指标：

- `xd_sdk_httpclient_request_latency`：请求延迟（毫秒）
  - 标签：`host`, `path`
  - 类型：Summary（包含 p50, p95, p99 分位数）

- `xd_sdk_httpclient_request_counter`：请求计数器
  - 标签：`host`, `path`, `method`, `http_code`, `code`（业务代码）
  - 类型：Counter

## 更多信息

- [resty 文档](https://godoc.org/gopkg.in/resty.v1)
- [resty Request 文档](https://godoc.org/gopkg.in/resty.v1#Request)
