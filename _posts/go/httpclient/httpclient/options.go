package httpclient

import (
	"time"

	"gitlab.xiaoduoai.com/golib/xd_sdk/httpclient/hooks"
)

type Options struct {
	Address            string
	Timeout            time.Duration
	RetryCount         int           // 重试次数
	RetryWaitTime      time.Duration // 重试间隔等待时间
	RetryMaxWaitTime   time.Duration // 重试间隔最大等待时间
	BeforeRequestHooks []hooks.BeforeRequestHook
	PreRequestHooks    []hooks.PreRequestHook
	AfterResponseHooks []hooks.AfterResponseHook

	EnableMetric      bool                           // 是否开启接口调用监控
	ParseMetricPath   func(originPath string) string // 解析访问路径，用于修改需要监控的路径 path
	ParseSvcCodeFuncs []hooks.ParseSvcCodeFunc       // 从 Response Body 中解析解析业务code
}

func newOptions(opts ...Option) Options {
	options := Options{
		Address:            "",
		Timeout:            3 * time.Second,
		RetryCount:         0,
		RetryWaitTime:      time.Duration(100) * time.Millisecond,
		RetryMaxWaitTime:   time.Duration(2000) * time.Millisecond,
		BeforeRequestHooks: []hooks.BeforeRequestHook{},
		PreRequestHooks:    []hooks.PreRequestHook{},
		AfterResponseHooks: []hooks.AfterResponseHook{},
		EnableMetric:       true, // todo: 后续改为默认 false
	}
	for _, opt := range opts {
		opt(&options)
	}
	options.PreRequestHooks = append(options.PreRequestHooks, hooks.DNSTrace())
	if options.EnableMetric {
		options.AfterResponseHooks = append(options.AfterResponseHooks, hooks.LatencyMetrics(options.ParseMetricPath, options.ParseSvcCodeFuncs...))
	}
	return options
}

type Option func(*Options)

func WithAddress(address string) Option {
	return func(options *Options) {
		options.Address = address
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.Timeout = timeout
	}
}

func WithRetryCount(retryCount int) Option {
	return func(options *Options) {
		options.RetryCount = retryCount
	}
}

func WithRetryWaitTime(retryWaitTime time.Duration) Option {
	return func(options *Options) {
		options.RetryWaitTime = retryWaitTime
	}
}

func WithRetryMaxWaitTime(retryMaxWaitTime time.Duration) Option {
	return func(options *Options) {
		options.RetryMaxWaitTime = retryMaxWaitTime
	}
}

func WithBeforeRequestHooks(hs ...hooks.BeforeRequestHook) Option {
	return func(options *Options) {
		options.BeforeRequestHooks = hs
	}
}

func WithPreRequestHooks(hs ...hooks.PreRequestHook) Option {
	return func(options *Options) {
		options.PreRequestHooks = hs
	}
}

func WithAfterResponseHooks(hs ...hooks.AfterResponseHook) Option {
	return func(options *Options) {
		options.AfterResponseHooks = hs
	}
}

// WithMetrics 是否打开接口调用监控
func WithMetrics(enable bool) Option {
	return func(options *Options) {
		options.EnableMetric = enable
	}
}

// WithParseSvcCodeFuncs 设置解析业务code的方式
func WithParseSvcCodeFuncs(fs ...hooks.ParseSvcCodeFunc) Option {
	return func(o *Options) {
		o.ParseSvcCodeFuncs = fs
	}
}

// WithParseMetricPathF 指定解析监控参数 path 的自定义方式，用于带参数的路径
func WithParseMetricPathF(f func(originPath string) string) Option {
	return func(o *Options) {
		o.ParseMetricPath = f
	}
}
