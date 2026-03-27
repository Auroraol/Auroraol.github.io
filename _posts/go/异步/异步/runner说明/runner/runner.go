package runner

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
)

type Runner interface {
	Run(ctx context.Context) (err error)
}

type runner struct {
	tasks []Task
}

type Task interface {
	Name() string
	Run(ctx context.Context) (err error)
}

func NewRunner(tasks ...Task) Runner {
	return &runner{tasks: tasks}
}

func (r *runner) Run(ctx context.Context) error {

	ctx, cancel := context.WithCancel(ctx) // 创建一个可取消的 context, 也会受父 context 的影响

	quit := make(chan os.Signal, 1)                      // 创建一个容量为1的信号通道 quit，用于接收系统信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //注册监听 SIGINT（Ctrl+C）和 SIGTERM（终止信号）两种系统信号
	defer func() { signal.Stop(quit); close(quit) }()
	go func() { <-quit; cancel() }() //启动一个goroutine等待信号，接收到信号后调用 cancel() 取消上下文

	var wg sync.WaitGroup
	for _, t := range r.tasks {
		task := t
		wg.Add(1) // 为每个 task 增加计数，表示有待完成的任务数量
		go func() {
			defer wg.Done() //在每个 goroutine 完成时减少计数
			logger.Infof(ctx, "task(%s) is starting", task.Name())
			if err := task.Run(ctx); err != nil {
				logger.Errorf(ctx, "task(%s) run with error(%v)", task.Name(), err)
			}
			logger.Infof(ctx, "task(%s) is stopped", task.Name())
			// Cancel the context to stop all other tasks.
			cancel()
		}()
	}
	wg.Wait() //阻塞主程序直到所有任务都调用了 Done()

	return nil
}
