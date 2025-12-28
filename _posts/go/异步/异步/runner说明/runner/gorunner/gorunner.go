package gorunner

import (
	"context"
	"errors"
	"strings"
	"sync"

	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
	"gitlab.xiaoduoai.com/golib/xd_sdk/xd_error"
)

var defaultRunner = NewGoRunner()

// GoRun 是全局变量 defaultRunner 的封装
func GoRun(ctx context.Context, f func() error, opts ...OptionFunc) {
	defaultRunner.GoRun(ctx, f, opts...)
}

// Wait 是全局变量 defaultRunner 的封装
func Wait() {
	defaultRunner.Wait()
}

// WaitErr 是全局变量 defaultRunner 的封装
func WaitErr() error {
	return defaultRunner.WaitErr()
}

// GoRunner 封装 goroutine 并发执行及结果等待
type GoRunner interface {
	// GoRun 并发执行 f，且防止 panic 向上传递
	GoRun(ctx context.Context, f func() error, opts ...OptionFunc)

	// Wait 只 Wait GoRun 不带 WithSkipWait 选项的 goroutine
	// Wait 忽略 goroutine 返回的 error
	// (阻塞)
	Wait()

	// WaitErr 等待 GoRun 并发启动的 goroutine 返回 error（前提条件：GoRun 带了 WithWaitErr 且不带 WithSkipWait 选项）
	// WaitErr 只 Wait GoRun 不带 WithSkipWait 选项的 goroutine
	// (阻塞)
	WaitErr() error
}

// NewGoRunner 创建 GoRunner
func NewGoRunner() GoRunner {
	return &goRunner{
		wg:     &sync.WaitGroup{},
		errors: make([]error, 0),
		mutex:  &sync.Mutex{},
	}
}

type goRunner struct {
	wg     *sync.WaitGroup
	errors []error
	mutex  *sync.Mutex
}

func (r *goRunner) appendErr(err error) {
	if err == nil {
		return
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.errors = append(r.errors, err)
}
func (r *goRunner) GoRun(ctx context.Context, f func() error, opts ...OptionFunc) {
	opt := getOptions(opts) // 获取选项，opt 是一个包含并发任务配置的对象
	if !opt.skipWait {      // 如果配置中不跳过等待（skipWait为false），就增加等待计数
		r.wg.Add(1)
	}
	go func() {
		if !opt.skipWait { // 如果配置中不跳过等待（skipWait为false），就减少等待计数
			defer r.wg.Done()
		}
		// 使用defer来处理异常（panic）
		defer func() {
			// Panic 处理
			if err := recover(); err != nil {
				r.appendErr(err.(error))
				logger.Error(ctx, "goRunner.GoRun panic: ", xd_error.Wrap(err.(error)))
			}
		}()

		// 执行传入的函数f，并检查是否发生错误
		if err := f(); err != nil {
			if opt.waitErr { // 如果配置中要求等待错误（waitErr为true），就记录错误
				r.appendErr(err)
			}
			logger.Error(ctx, "goRunner.GoRun exec f fail: ", err)
		}
	}()
}

func (r *goRunner) Wait() {
	r.wg.Wait()
}
func (r *goRunner) WaitErr() error {
	r.wg.Wait()
	if len(r.errors) == 0 {
		return nil
	}
	builder := strings.Builder{}
	for i, e := range r.errors {
		if i != len(r.errors) {
			builder.WriteString(e.Error() + ",")
		} else {
			builder.WriteString(e.Error())
		}
	}
	return errors.New(builder.String())
}

// 配置项
type OptionFunc func(opt *Option)

// 设置当前 GoRun 所执行的 goroutine 跳过 Wait
// 用于指定某个 Goroutine 执行时跳过 `Wait()` 阶段，这对于不需要等待的任务（例如某些无关紧要的异步任务）非常有用。
func WithSkipWait() OptionFunc {
	return func(opt *Option) {
		opt.skipWait = true
	}
}

// 设置当前 GoRun 接收处理函数 f 返回的 error，
// 默认忽略 f 返回的 error
// 用于指定某个 Goroutine 执行时收集其返回的错误。默认情况下，Goroutine 执行的错误是被忽略的。
func WithWaitErr() OptionFunc {
	return func(opt *Option) {
		opt.waitErr = true
	}
}

type Option struct {
	skipWait bool
	waitErr  bool
}

func getOptions(fs []OptionFunc) *Option {
	opt := &Option{}
	for _, f := range fs {
		f(opt)
	}
	return opt
}
