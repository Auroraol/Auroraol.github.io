# 说明

## runner 包 - 优雅退出管理器

主要功能：

- 实现优雅退出机制，处理SIGINT和SIGTERM信号

- 支持多个并发任务的启动和停止

- 当收到终止信号时，等待当前任务完成后再退出

使用方法：

```go
package main

import (
    "context"
    "fmt"
    "time"
    "your-module/runner" // 替换为实际的包路径
)

// 实现Task接口
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
        select {
        case <-ctx.Done():
            return ctx.Err() // 优雅退出
        default:
            // 执行业务逻辑
            fmt.Printf("%s 正在工作...\n", w.Name())
            time.Sleep(1 * time.Second)
        }
    }
}

func main() {
    ctx := context.Background()
    worker := NewWorker("工作线程1")
    
    // 使用runner启动任务
    _ = runner.NewRunner(worker).Run(ctx)
}
```

## gorunner 包 - 并发任务管理器

主要功能：

- 封装goroutine的并发执行

- 防止panic向上传递

- 支持错误收集和等待机制

- 提供灵活的配置选项

```go
package main

import (
    "context"
    "fmt"
    "time"
    "your-module/gorunner" // 替换为实际的包路径
)

func main() {
    ctx := context.Background()
    
    // 方式1：使用全局默认runner
    gorunner.GoRun(ctx, func() error {
        fmt.Println("任务1执行中...")
        time.Sleep(2 * time.Second)
        return nil
    })
    
    gorunner.GoRun(ctx, func() error {
        fmt.Println("任务2执行中...")
        time.Sleep(1 * time.Second)
        return fmt.Errorf("任务2出错")
    }, gorunner.WithWaitErr()) // 收集错误
    
    // 等待所有任务完成
    if err := gorunner.WaitErr(); err != nil {
        fmt.Printf("有任务出错: %v\n", err)
    }
    
    // 方式2：创建独立的runner
    runner := gorunner.NewGoRunner()
    runner.GoRun(ctx, func() error {
        fmt.Println("独立任务执行中...")
        return nil
    })
    runner.Wait()
}
```

配置选项：

- WithSkipWait(): 跳过等待，不阻塞主线程

- WithWaitErr(): 收集并返回错误



# 案例

## Wait等待

```go
func main() {
    runner := NewGoRunner()

    // 启动第一个任务
    runner.GoRun(ctx, func() error {
        // 执行任务1
        return nil
    })

    // 启动第二个任务
    runner.GoRun(ctx, func() error {
        // 执行任务2
        return nil
    })

    // 启动第三个任务
    runner.GoRun(ctx, func() error {
        // 执行任务3
        return nil
    })

    // 等待所有任务完成
    runner.Wait() //没有 Wait() 的话，main函数会立即结束，但goroutine会继续在后台运行
    // 这里会阻塞直到所有GoRun启动的任务执行完毕
}
```

## WaitErr错误

```go
for i := 0; i < 10; i++ {
    tmp := i
    runner.GoRun(ctx, func() error {
        time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
        fmt.Printf("go run: %v", tmp)
        if rand.Intn(3) == 1 {
            return errors.New(strconv.Itoa(tmp) + ":fail")
        }
        return nil
    }, gorunner.WithWaitErr())
}


// 等待所有 Goroutine 执行完毕，如果收集到错误，会将所有错误连接成一个字符串返回
err := runner.WaitErr()
```

# 使用

## runner.go 的使用方法 (任务级别)

---

这个文件实现了一个优雅退出管理器，主要用于处理多个并发任务的启动和优雅停止。

### 核心概念

#### 1. 接口定义

```go
type Runner interface {
    Run(ctx context.Context) (err error)
}

type Task interface {
    Name() string
    Run(ctx context.Context) (err error)
}
```

- **Runner**: 运行器接口，负责管理多个任务的执行
- **Task**: 任务接口，每个任务都需要实现这个接口

#### 2. 核心功能

`runner.go` 的主要功能包括：

- **信号处理**: 自动捕获 SIGINT (Ctrl+C) 和 SIGTERM 信号
- **优雅退出**: 收到信号后，等待当前任务完成再退出
- **并发执行**: 支持多个任务同时运行
- **错误处理**: 当任一任务出错时，会取消所有其他任务

### 详细使用方法

#### 步骤1: 实现Task接口

首先需要创建一个结构体并实现 `Task` 接口：

```go
package main

import (
    "context"
    "fmt"
    "time"
    "errors"
    "your-module/runner" // 替换为实际的包路径
)

// 示例任务1: 数据处理任务
type DataProcessor struct {
    name string
}

func NewDataProcessor(name string) *DataProcessor {
    return &DataProcessor{name: name}
}

func (dp *DataProcessor) Name() string {
    return dp.name
}

func (dp *DataProcessor) Run(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("%s: 收到停止信号，正在清理资源...\n", dp.Name())
            time.Sleep(1 * time.Second) // 模拟清理工作
            fmt.Printf("%s: 清理完成，优雅退出\n", dp.Name())
            return ctx.Err()
        default:
            // 执行业务逻辑
            fmt.Printf("%s: 正在处理数据...\n", dp.Name())
            time.Sleep(2 * time.Second)
        }
    }
}

// 示例任务2: 监控任务
type Monitor struct {
    name string
}

func NewMonitor(name string) *Monitor {
    return &Monitor{name: name}
}

func (m *Monitor) Name() string {
    return m.name
}

func (m *Monitor) Run(ctx context.Context) error {
    ticker := time.NewTicker(3 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("%s: 监控任务停止\n", m.Name())
            return ctx.Err()
        case <-ticker.C:
            fmt.Printf("%s: 执行健康检查\n", m.Name())
        }
    }
}
```

#### 步骤2: 使用Runner启动任务

```go
func main() {
    ctx := context.Background()
    
    // 创建任务实例
    dataProcessor := NewDataProcessor("数据处理器")
    monitor := NewMonitor("系统监控")
    
    // 方式1: 启动单个任务
    fmt.Println("=== 启动单个任务 ===")
    runner1 := runner.NewRunner(dataProcessor)  // 异步
    err := runner1.Run(ctx)
    if err != nil {
        fmt.Printf("任务执行出错: %v\n", err)
    }
    
    // 方式2: 启动多个任务
    fmt.Println("\n=== 启动多个任务 ===")
    runner2 := runner.NewRunner(dataProcessor, monitor) //异步
    err = runner2.Run(ctx)
    if err != nil {
        fmt.Printf("任务执行出错: %v\n", err)
    }
}
```

### 关键特性说明

#### 1. 信号处理机制

```go
ctx, cancel := context.WithCancel(ctx)

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
defer func() { signal.Stop(quit); close(quit) }()
go func() { <-quit; cancel() }()
```

- 自动捕获 `SIGINT` (Ctrl+C) 和 `SIGTERM` 信号
- 收到信号后通过 `context.WithCancel` 通知所有任务停止

#### 2. 并发执行和等待

```go
var wg sync.WaitGroup
for _, t := range r.tasks {
    task := t
    wg.Add(1)
    go func() {
        defer wg.Done()
        logger.Infof(ctx, "task(%s) is starting", task.Name())
        if err := task.Run(ctx); err != nil {
            logger.Errorf(ctx, "task(%s) run with error(%v)", task.Name(), err)
        }
        logger.Infof(ctx, "task(%s) is stopped", task.Name())
        // Cancel the context to stop all other tasks.
        cancel()
    }()
}
wg.Wait()
```

- 使用 `sync.WaitGroup` 等待所有任务完成
- 每个任务在独立的 goroutine 中运行
- 当任一任务出错时，调用 `cancel()` 停止所有其他任务

## gorunner 包 (函数级别)

```go
newCtx := octrace.CopySpanToNewCtx(ctx)
gorunner.GoRun(newCtx, func() error {
    // 业务逻辑
    return nil
})
```



# 案例

## 例子1

 执行流程：

```
商品1  商品2  商品3  ...  商品20
  ↓      ↓      ↓            ↓
[=====并发控制池=10个=====]
  ↓      ↓      ↓            ↓
同时最多10个API请求在执行
```

```go
// 商品列表
func (h *SdkHandler) GoodsListHdl(ctx context.Context, req *libproto.GetGoodsListReq, rsp *libproto.GetGoodsListRsp) error {
	shopInfo, err := h.GetShopInfoByID(req.ShopId)
	if err != nil {
		return errors.Wrapf(err, "failed to GetShopInfoByID, %+v", req)
	}
	retryCount := getRetryCount(ctx)

	// 调用GetGoodsListTODO获取商品列表
	goodsListRes, err := h.gateClient.GetGoodsListTODO(ctx, shopInfo, req.IsOnsale, int64(req.Page), int64(req.PageSize),
		req.Keyword, req.StartTime, req.EndTime, retryCount)
	if err != nil {
		return errors.Wrap(err, "failed to get goods list")
	}

	rsp.Total = goodsListRes.GoodsListGetResponse.TotalCount
	goodsList := goodsListRes.GoodsListGetResponse.GoodsList
	if len(goodsList) == 0 {
		return nil
	}

	// 优化：使用worker pool模式控制并发
	const maxConcurrency = 10
	type workItem struct {
		index int
		goods gate.GoodsListGoods
	}
	workCh := make(chan workItem, len(goodsList))
	results := make([]libproto.GetGoodsInfoRsp, len(goodsList))

	// 启动worker goroutines
	runner := gorunner.NewGoRunner()
	for i := 0; i < maxConcurrency; i++ {
		runner.GoRun(ctx, func() error {
			for work := range workCh {
				newCtx := octrace.CopySpanToNewCtx(ctx)

				goodsInfo := libproto.GetGoodsInfoRsp{}
				err := h.GoodsInfoHdl(newCtx, &libproto.GetGoodsInfoReq{
					ShopBaseReq: libproto.ShopBaseReq{ShopId: req.ShopId},
					PlatGoodsId: fmt.Sprintf("%v", work.goods.GoodsID),
				}, &goodsInfo)
				if err != nil {
					logger.Errorf(newCtx, "goodsListHdl call api get goods err: %v, goodsId: %v", err, work.goods.GoodsID)
					continue
				}

				createTime, timeErr := libtime.TimestampToTime(strconv.FormatInt(work.goods.CreatedAt, 10))
				if timeErr != nil {
					logger.Errorf(newCtx, "parse start time err: %v", timeErr)
				} else {
					goodsInfo.CreateTime = createTime
				}

				results[work.index] = goodsInfo
			}
			return nil
		})
	}

	// 发送工作任务
	for i, goods := range goodsList {
		workCh <- workItem{index: i, goods: goods}
	}
	close(workCh)

	runner.Wait()

	// 处理有效结果
	for _, result := range results {
		if len(result.SkuList) > 0 {
			rsp.List = append(rsp.List, result)
		}
	}

	return nil
}
```

## 例子2

```go
func (w *Worker) collectSyncResult(ctx context.Context, taskFields logger.Fields, pf platform.TPlatform, task *task.Task, total int64) (*SyncResult, error) {
	var onSaleCh <-chan IterGoodsResult
	var err error
	if pf == platform.EPlatformWXXD {
		onSaleCh, err = Instance(pf).IterOnsaleGoodsByCursor(ctx, w, pf, task.ShopId, WXXDMIT)
	} else {
		onSaleCh, err = Instance(pf).IterOnsaleGoodsByPage(ctx, w, pf, task.ShopId, InitialPageNo, DefaultLIMIT, total)
	}

	if err != nil {
		logger.WithFields(taskFields).Errorf(ctx, "%s iter onsale goods failed, err: %v", pf.String(), err)
		return nil, err
	}
	goodsType := w.container.GetGoodsType(pf)
	result := &SyncResult{
		SuccSpuIds:         make(map[string]struct{}),
		SuccSkuIds:         make(map[string]struct{}),
		FailedPlatGoodsIds: make(map[string]struct{}),
		Abort:              false,
		ActualTotal:        0,
	}

	syncM := sync.Mutex{}
	runner := gorunner.NewGoRunner()
	runner.GoRun(ctx, func() error {
		for iterResult := range onSaleCh {
			if iterResult.Err != nil {
				result.Abort = true
				return iterResult.Err
			}

			var err error
			if goodsType == config.EGoodsTypesSpu {
				tempSpu := iterResult.Spu
				_, err = w.container.Spu().AddSpu(ctx, &tempSpu)
				if err != nil {
					logger.Infof(ctx, "%s add spu failed, goods is %v, err: %v", pf.String(), tempSpu.PlatGoodsID, err)
					syncM.Lock()
					result.FailedPlatGoodsIds[tempSpu.PlatGoodsID] = struct{}{}
					syncM.Unlock()
					continue
				}
				syncM.Lock()
				result.SuccSpuIds[tempSpu.PlatGoodsID] = struct{}{}
				result.ActualTotal++
				syncM.Unlock()
			} else if goodsType == config.EGoodsTypesSku {
				tempSku := iterResult.Sku
				_, err = w.container.Sku().AddSku(ctx, &tempSku)
				if err != nil {
					logger.Infof(ctx, "%s add sku failed, sku is %v, err: %v", pf.String(), tempSku.PlatSkuID, err)
					syncM.Lock()
					result.FailedPlatGoodsIds[tempSku.PlatSkuID] = struct{}{}
					syncM.Unlock()
					continue
				}
				syncM.Lock()
				result.SuccSkuIds[tempSku.PlatSkuID] = struct{}{}
				result.ActualTotal++
				syncM.Unlock()
			}
		}
		return nil
	})

	// 定时上报进度
	closeProcessCh := make(chan struct{})
	closeProcessAckCh := make(chan struct{})
	runner.GoRun(ctx, func() error {
		defer func() {
			closeProcessAckCh <- struct{}{}
		}()

		t := time.NewTicker(time.Second * 5)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				syncM.Lock()
				succ := int64(len(result.SuccSpuIds) + len(result.SuccSkuIds))
				failed := int64(len(result.FailedPlatGoodsIds))
				syncM.Unlock()

				fixTotal := max(total, succ+failed)
				logger.WithFields(taskFields).WithFields(logger.Fields{
					"trueTotal": result.ActualTotal,
					"total":     fixTotal,
					"succeed":   succ,
					"failed":    failed,
				}).Infof(ctx, "%s update task process", pf.String())
				err = w.TaskProcessReport(fixTotal, succ, failed)
				if err != nil {
					logger.WithFields(taskFields).WithError(err).Errorf(ctx, "%s report process failed", pf.String())
				}
			case <-closeProcessCh:
				return nil
			}
		}
	}, gorunner.WithSkipWait())

	runner.Wait()
	closeProcessCh <- struct{}{}
	<-closeProcessAckCh

	// 最终一次进度上报
	fixTotal := max(total, int64(len(result.SuccSpuIds)+len(result.SuccSkuIds)+len(result.FailedPlatGoodsIds)))
	succ := int64(len(result.SuccSpuIds) + len(result.SuccSkuIds))
	failed := int64(len(result.FailedPlatGoodsIds))
	err = w.TaskProcessReport(fixTotal, succ, failed)
	if err != nil {
		logger.WithFields(taskFields).WithError(err).Errorf(ctx, "%s report process failed", pf.String())
	}

	return result, nil
}

```

## 例子3

见: D:\Github\go\问题\pulsar+task\worker