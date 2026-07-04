[Go每日一库之113：retry-go](https://juejin.cn/post/7289325465809911847)

```
go get https://github.com/avast/retry-go
```

使用例子

```go
	var lastError error
	
	err = retry.Do(
		func() error {
		   // 业务处理逻辑
         // ....
			return nil
		},
		retry.Attempts(3), // 最多执行3次
		retry.Delay(1*time.Second),
		retry.BackoffExponential(2*time.Second),  //了指数退避策略：初始延迟时间为2秒
		retry.DelayType(retry.CombineDelay(retry.BackOffDelay, retry.RandomDelay)),
		retry.OnRetry(func(n uint, err error) {
            // 每当发生重试时，这个函数就会被调用一次
			logger.WithField("shop_id", shopInfo.ID.Hex()).
				WithField("goods_id", goodsID).
				Warningf(ctx, "the %v time to get goods info, error: %v", n+1, err)
		}),
		retry.LastErrorOnly(true), // 只返回最后一次错误
		retry.Context(ctx),     // 使用传入的context控制
	)

	// 如果达到最大重试次数仍然失败
	if err != nil {
		return nil, err
	}
```

retry.BackoffExponential(2*time.Second)
这个选项设置了指数退避策略：

1. 初始延迟时间为2秒
   每次重试的延迟时间按指数增长
   例如：第1次重试等待2秒，第2次等待4秒，第3次等待8秒，以此类推
   这样可以避免在系统繁忙时频繁重试，给系统恢复的时间

retry.DelayType(retry.CombineDelay(retry.BackOffDelay, retry.RandomDelay))
这个选项设置了组合延迟类型：

1. retry.BackOffDelay：使用上面定义的指数退避延迟
   retry.RandomDelay：添加随机延迟，防止惊群效应
   retry.CombineDelay：将两种延迟策略组合使用

组合作用效果
当这两个选项一起使用时，会产生以下效果：

```
组合作用效果
当这两个选项一起使用时，会产生以下效果：
实际延迟时间 = 指数退避时间 + 随机时间

例如：
第1次重试：2秒 + 随机时间(0-几百毫秒)
第2次重试：4秒 + 随机时间(0-几百毫秒)
第3次重试：8秒 + 随机时间(0-几百毫秒)
```

retry.LastErrorOnly(true)

+ 这个选项控制返回的错误信息：
  当设置为 true 时，重试失败后只返回最后一次重试的错误
  当设置为 false 或不设置时，会返回所有重试过程中产生的错误列表
  作用：
  简化错误信息，避免日志中出现重复的错误记录
  减少内存占用，特别是在重试次数较多时
  使错误处理更简洁

retry.Context(ctx)

+ 这个选项用于传递 context 控制：
  将外部的 context 传入重试机制中
  允许通过 context 控制重试的取消或超时
  当 context 被取消或超时时，重试会立即停止
  作用：
  支持超时控制：如果 context 设置了超时，重试不会超过这个时间
  支持取消操作：如果请求被取消，重试也会立即停止
  保持 context 传播：确保 tracing、logging 等上下文信息正确传播





普通错误：会触发重试机制，直到达到最大重试次数
retry.Unrecoverable 错误：立即停止重试，直接返回错误





```go

func (c *RetryClient) GetGoodsInfo(ctx context.Context, shopInfo models.ShopInfoT, goodsID string) (resp *GetGoodsInfoResponse, err error) {
	lastErr := retry.Do(
		func() error {
			resp, err = c.client.GetGoodsInfo(ctx, shopInfo, goodsID)
			if err != nil {
				return err
			}

			// 如果是限流错误，需要重试
			if resp.ErrorResponse.ErrorCode == RateLimitErrCode {
				return errors.Errorf("rate limit error: code:%v, msg:%v", resp.ErrorResponse.ErrorCode, resp.ErrorResponse.ErrorMsg)
			}

			// 其他错误
			if resp.ErrorResponse.ErrorCode != 0 {
				logger.Warningf(ctx, "get goods info failed:code:%v, msg:%v", resp.ErrorResponse.ErrorCode, resp.ErrorResponse.ErrorMsg)
				// 如果是商品不存在错误，直接返回特定错误
				if resp.ErrorResponse.ErrorCode == constant.CODE_GOODS_NOT_EXISTED {
					return retry.Unrecoverable(errors.New(constant.MSG_GOODS_NOT_EXISTED))
				}
				return retry.Unrecoverable(errors.Errorf("get goods info failed:code:%v, msg:%v", resp.ErrorResponse.ErrorCode, resp.ErrorResponse.ErrorMsg))
			}

			// 如果商品已删除
			if resp.GoodsDetailGetResponse.Status == constant.GOODS_STATUS_DELETED {
				return retry.Unrecoverable(errors.New(constant.MSG_GOODS_NOT_EXISTED))
			}

			return nil
		},
		retry.Attempts(uint(c.retryCount)),
		retry.Delay(1*time.Second),
		retry.DelayType(retry.CombineDelay(retry.BackOffDelay, retry.RandomDelay)),
		retry.LastErrorOnly(true),
		retry.OnRetry(func(n uint, err error) {
			logger.WithField("shop_id", shopInfo.ID.Hex()).WithField("goods_id", goodsID).
				Warningf(ctx, "the %v time to get goods info, error: %v", n+1, err)
		}),
	)

	if lastErr != nil {
		return nil, lastErr
	}

	return resp, nil
}
```

重试次数：由 retry.Attempts(uint(c.retryCount)) 控制，最多重试 c.retryCount 次
重试时间间隔：
初始延迟：1秒 (retry.Delay(1*time.Second))
延迟策略：组合了退避延迟和随机延迟 (retry.CombineDelay(retry.BackOffDelay, retry.RandomDelay))
这意味着每次重试的间隔会逐渐增加，并加入随机因素以避免惊群效应
触发重试的条件：
网络请求返回错误
返回限流错误 (RateLimitErrCode = 70031)
不触发重试的条件（使用 retry.Unrecoverable 包装）：
商品不存在
商品已删除
其他业务错误