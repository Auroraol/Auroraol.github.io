# runner

### 什么是优雅退出
让我们看个简单的例子：
```golang
package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
)

type Worker struct {
	name string
}

func NewWorker(name string) *Worker {
	return &Worker{name: name}
}

func (w *Worker) Name() string {
	return w.name
}

func (w *Worker) Run(ctx context.Context) error {
	for {
		w.Report()
	}
}

func (w *Worker) Report() {
	for counter := 0; counter < 5; counter++ {
		fmt.Printf("%s report number %v\n", w.Name(), counter)
		time.Sleep(3 * time.Second)
	}
}

func main() {
	ctx := context.Background()
	w := NewWorker("Alice")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := w.Run(ctx)
		logger.Infof("%s stopped: error = %v", w.Name(), err)
	}()
	wg.Wait()
}
```
这个程序的功能很简单，就是启动一个worker不停地执行报数，每一轮都是从0到4然后下一轮。输出是这样的：
```shell
Alice report number 0
Alice report number 1
Alice report number 2

Process finished with exit code 130 (interrupted by signal 2: SIGINT)
```
这里面有一个问题，当我们按下Ctrl+C的一瞬间，程序收到SIGINT信号就终止了。然而一轮报数都还没结束，甚至连后面的那行日志都没打印出来。当然这里只是报数，在真实的业务中有可能是接收完消息还没处理完毕。为了尽量减少这种情况的发生，我们一般建议是优先保证本轮操作结束再停掉程序。所以我们建议程序应该写成如下这个样子(只列出变更部分)：
```golang
func (w *Worker) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "daemon context done")
		default:
			w.Report()
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer func() { signal.Stop(quit); close(quit) }()
	go func() { <-quit; cancel() }()

	w := NewWorker("Alice")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := w.Run(ctx)
		logger.Infof("%s stopped: error = %v", w.Name(), err)
	}()
	wg.Wait()
}
```
这里面有几行核心的代码变更，目的是为了主动拦截操作系统的终止信号，然后通知具体的worker停下来，这个便给了worker足够的"善后"时间。当程序启动后我立马按下Ctrl+C，但是程序并不会立即退出，最终的输出是这样的：
```shell
Alice report number 0
Alice report number 1
Alice report number 2
Alice report number 3
Alice report number 4
time="2019-08-12T16:25:22+08:00" level=info msg="Alice stopped: error = daemon context done: context canceled"

Process finished with exit code 0
```

### 如何使用这个包
以上的这个优化就叫做优雅退出，但是为每个程序都写这么一段代码确实是一件比较繁琐的事情，而且你还要考虑假设起多个worker的时候会更复杂些。因此这个包本质上就是对上面的这个优化进行了简单封装，引用了这个包以后，你的代码可以是这样的(只列出变更部分)：
```golang
func main() {
	ctx := context.Background()
	w := NewWorker("Alice")
	_ = runner.NewRunner(w).Run(ctx)
}
```
最终的输出如下：
```shell
time="2019-08-12T16:45:10+08:00" level=info msg="task(Alice) is starting"
Alice report number 0
Alice report number 1
Alice report number 2
Alice report number 3
Alice report number 4
time="2019-08-12T16:45:25+08:00" level=error msg="task(Alice) run with error(daemon context done: context canceled)"
time="2019-08-12T16:45:25+08:00" level=info msg="task(Alice) is stopped"

Process finished with exit code 0
```
可以看出这样也同样实现了优雅退出，而且代码量少了很多，让你从底层细节中解脱出来，更多的时间可以用来关注业务逻辑。
