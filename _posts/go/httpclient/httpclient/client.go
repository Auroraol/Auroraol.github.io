package httpclient

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	"gitlab.xiaoduoai.com/golib/xd_sdk/httpclient/hooks"
	"gitlab.xiaoduoai.com/golib/xd_sdk/metadata"
	"gitlab.xiaoduoai.com/golib/xd_sdk/octrace"
	"gitlab.xiaoduoai.com/golib/xd_sdk/public"

	"github.com/pkg/errors"
	ochttp "gitlab.xiaoduoai.com/golib/xd_sdk/octrace/integration/http"
	resty "gopkg.in/resty.v1"
)

type Client interface {
	AddBeforeRequestHooks(hs ...hooks.BeforeRequestHook)
	SetPreRequestHooks(hs ...hooks.PreRequestHook)
	AddAfterResponseHooks(hs ...hooks.AfterResponseHook)
	GetKernel() (kernel *resty.Client) // 拿到resty.client,支持更多丰富功能调用

	NewRequest(ctx context.Context) (request *resty.Request)

	PostJSON(ctx context.Context, path string, values interface{}, headers http.Header, timeoutSeconds int, ret interface{}) (err error)
	Post(ctx context.Context, path string, values url.Values, headers http.Header, timeoutSeconds int, ret interface{}) (err error)
	Get(ctx context.Context, path string, values url.Values, headers http.Header, timeoutSeconds int, ret interface{}) (err error)
	RawGet(ctx context.Context, path string, values url.Values, headers http.Header, timeoutSeconds int) (content []byte, err error)
	GetWithMock(ctx context.Context, path string, values url.Values, headers http.Header, timeoutSeconds int, preSetRet interface{}, ret interface{}) (err error)
}

type client struct {
	kernel  *resty.Client
	options Options
}

func NewClient(opts ...Option) (Client, error) {
	return NewClientWithOptions(newOptions(opts...))
}

func NewClientWithOptions(options Options) (Client, error) {
	var c *client

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}
	if octrace.Enabled() {
		c = &client{
			kernel:  resty.NewWithClient(&http.Client{Transport: ochttp.NewTransportWithTransport(transport)}),
			options: options,
		}
	} else {
		c = &client{
			kernel:  resty.NewWithClient(&http.Client{Transport: transport}),
			options: options,
		}
	}

	c.kernel.SetHostURL(c.options.Address)
	c.kernel.SetTimeout(c.options.Timeout)
	c.kernel.SetRetryCount(c.options.RetryCount)
	c.kernel.SetRetryWaitTime(c.options.RetryWaitTime)
	c.kernel.SetRetryMaxWaitTime(c.options.RetryMaxWaitTime)
	c.AddBeforeRequestHooks(c.options.BeforeRequestHooks...)
	c.SetPreRequestHooks(c.options.PreRequestHooks...)
	c.AddAfterResponseHooks(c.options.AfterResponseHooks...)

	info := public.GetCallingInfoFromOS()
	if info.SvcName != "" {
		c.kernel.SetHeader(public.HeaderSvcName, info.SvcName)
	}
	if info.HostName != "" {
		c.kernel.SetHeader(public.HeaderHostName, info.HostName)
	}
	return c, nil
}

func (c *client) AddBeforeRequestHooks(hs ...hooks.BeforeRequestHook) {
	for _, h := range hs {
		c.kernel.OnBeforeRequest(h)
	}
}

func (c *client) SetPreRequestHooks(hs ...hooks.PreRequestHook) {
	c.kernel.SetPreRequestHook(func(c *resty.Client, r *resty.Request) error {
		for _, h := range hs {
			if err := h(c, r); err != nil {
				return errors.Wrap(err, "failed to execute pre request hook")
			}
		}
		return nil
	})
}

func (c *client) AddAfterResponseHooks(hs ...hooks.AfterResponseHook) {
	for _, h := range hs {
		c.kernel.OnAfterResponse(h)
	}
}

func (c *client) NewRequest(ctx context.Context) *resty.Request {
	rr := c.kernel.NewRequest().SetContext(ctx)

	// 设置metadata
	k, v := metadata.Encode(ctx)
	rr.SetHeader(k, v)
	return rr
}

func (c *client) GetKernel() *resty.Client {
	return c.kernel
}

func (c *client) PostJSON(ctx context.Context, path string, values interface{}, headers http.Header, timeoutSeconds int, ret interface{}) (err error) {
	if timeoutSeconds > 0 {
		c.kernel.SetTimeout(time.Duration(timeoutSeconds) * time.Second)
	}
	r := c.kernel.NewRequest().SetContext(ctx)

	if headers != nil {
		r.Header = headers
	}
	// 设置metadata
	k, v := metadata.Encode(ctx)
	r.SetHeader(k, v)
	r.SetHeader("Content-Type", "application/json")

	if values != nil {
		r.SetBody(values)
	}

	if ret != nil {
		r.SetResult(ret)
	}

	resp, err := r.Post(path)

	if err != nil {
		return fmt.Errorf("PostJSON:{%s} param:%+v err: %v", r.URL, values, err)
	}
	httpStatusCode := resp.StatusCode()
	if httpStatusCode < http.StatusOK || httpStatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("PostJSON:{%s} param:%+v failed: %v", resp.Request.URL, values, httpStatusCode)
	}

	return
}

func (c *client) Post(ctx context.Context, path string, values url.Values, headers http.Header, timeoutSeconds int, ret interface{}) (err error) {
	if timeoutSeconds > 0 {
		c.kernel.SetTimeout(time.Duration(timeoutSeconds) * time.Second)
	}
	r := c.kernel.NewRequest().SetContext(ctx)

	if headers != nil {
		r.Header = headers
	}
	// 设置metadata
	k, v := metadata.Encode(ctx)
	r.SetHeader(k, v)
	r.SetHeader("Content-Type", "application/x-www-form-urlencoded")

	if values != nil {
		r.FormData = values
	}

	if ret != nil {
		r.SetResult(ret)
	}

	resp, err := r.Post(path)

	if err != nil {
		return fmt.Errorf("post:{%s} param:%+v err: %v", r.URL, values, err)
	}
	httpStatusCode := resp.StatusCode()
	if httpStatusCode < http.StatusOK || httpStatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("post:{%s} param:%+v failed: %v", resp.Request.URL, values, httpStatusCode)
	}

	return
}

func (c *client) Get(ctx context.Context, path string, values url.Values, headers http.Header, timeoutSeconds int, ret interface{}) (err error) {
	if timeoutSeconds > 0 {
		c.kernel.SetTimeout(time.Duration(timeoutSeconds) * time.Second)
	}
	r := c.kernel.NewRequest().SetContext(ctx)

	if headers != nil {
		r.Header = headers
	}
	// 设置metadata
	k, v := metadata.Encode(ctx)
	r.SetHeader(k, v)

	if values != nil {
		r.QueryParam = values
	}

	if ret != nil {
		r.SetResult(ret)
	}

	resp, err := r.Get(path)
	if err != nil {
		return fmt.Errorf("get:{%s} param:%+v err: %v", r.URL, r.QueryParam, err)
	}
	httpStatusCode := resp.StatusCode()
	if httpStatusCode < http.StatusOK || httpStatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("get:{%s} param:%+v failed: %v", resp.Request.URL, values, httpStatusCode)
	}

	return
}

func (c *client) RawGet(ctx context.Context, path string, values url.Values, headers http.Header, timeoutSeconds int) (content []byte, err error) {
	if timeoutSeconds > 0 {
		c.kernel.SetTimeout(time.Duration(timeoutSeconds) * time.Second)
	}
	r := c.kernel.NewRequest().SetContext(ctx)

	if headers != nil {
		r.Header = headers
	}

	k, v := metadata.Encode(ctx)
	r.SetHeader(k, v)

	if values != nil {
		r.QueryParam = values
	}

	resp, err := r.Get(path)

	if err != nil {
		return nil, fmt.Errorf("get:{%s} param:%+v err: %v", r.URL, r.QueryParam, err)
	}
	httpStatusCode := resp.StatusCode()
	if httpStatusCode < http.StatusOK || httpStatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("get:{%s} param:%+v failed: %v", resp.Request.URL, values, httpStatusCode)
	}

	return resp.Body(), err
}

func (c *client) GetWithMock(ctx context.Context, path string, values url.Values, headers http.Header, timeoutSeconds int, preSetRet interface{}, ret interface{}) (err error) {
	if timeoutSeconds > 0 {
		c.kernel.SetTimeout(time.Duration(timeoutSeconds) * time.Second)
	}
	r := c.kernel.NewRequest().SetContext(ctx)

	if headers != nil {
		r.Header = headers
	}

	k, v := metadata.Encode(ctx)
	r.SetHeader(k, v)

	if values != nil {
		r.QueryParam = values
	}

	if preSetRet != nil {
		r.SetBody(preSetRet)
	}

	if ret != nil {
		r.SetResult(ret)
	}

	resp, err := r.Get(path)
	if err != nil {
		return fmt.Errorf("get:{%s} param:%+v err: %v", r.URL, r.QueryParam, err)
	}
	httpStatusCode := resp.StatusCode()
	if httpStatusCode < http.StatusOK || httpStatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("get:{%s} param:%+v failed: %v", resp.Request.URL, values, httpStatusCode)
	}

	return
}

func IsConnectRefuseError(err error) bool {
	if err == nil {
		return false
	}

	opError, ok := err.(*net.OpError)
	if !ok {
		return false
	}

	const dialString = "dial"
	const netString = "tcp"
	isOk := opError.Op == dialString &&
		strings.Contains(opError.Err.Error(), "connection refused") &&
		opError.Net == netString

	return isOk
}
