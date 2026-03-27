package hooks

import resty "gopkg.in/resty.v1"

type BeforeRequestHook func(*resty.Client, *resty.Request) error

type PreRequestHook func(*resty.Client, *resty.Request) error

type AfterResponseHook func(*resty.Client, *resty.Response) error

// ParseSvcCodeFunc 解析具体调用API的业务 code，用于监控告警
// httpCode 是HTTP状态码，result 是HTTP返回内容
// result 默认是 SetResult 指定的返回值类型；如果没有通过 SetResult 设置返回值类型，则 result 为 []byte 类型的HTTP返回 body
type ParseSvcCodeFunc func(httpCode int, result interface{}) (svcCode int, ok bool)
