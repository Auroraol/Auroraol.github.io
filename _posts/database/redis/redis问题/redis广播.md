



# 



```go
package sync_worker

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
	"gitlab.xiaoduoai.com/golib/xd_sdk/redisc"
)

type RedisSubscriber struct {
	client  redisc.Client
	channel string
}

// 创建Redis订阅客户端
func NewRedisSubscriber(client redisc.Client, channel string) (*RedisSubscriber, error) {
	return &RedisSubscriber{
		client:  client,
		channel: channel,
	}, nil
}

// 订阅Redis频道
func (r *RedisSubscriber) Subscribe(ctx context.Context, msgChan chan<- GoodsSyncMsg) error {
	for {
		err := r.subscribeAndListen(ctx, msgChan)
		if isRedisConnClosedErr(err) {
			logger.Warnf(ctx, "Redis connection lost, will reconnect after 2s...: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		return err
	}
}

func (r *RedisSubscriber) subscribeAndListen(ctx context.Context, msgChan chan<- GoodsSyncMsg) error {
	conn := r.client.GetConn()
	realConn := conn.GetRealConn()
	psc := redis.PubSubConn{Conn: realConn}
	defer func() {
		psc.Close()
		conn.Close()
	}()

	if err := psc.Subscribe(r.channel); err != nil {
		logger.Errorf(ctx, "failed to subscribe to channel %s: %v", r.channel, err)
		return err
	}

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Infof(ctx, "Redis subscriber context cancelled")
			return ctx.Err()
		case <-ticker.C:
			if err := psc.Ping(""); err != nil {
				logger.Errorf(ctx, "Redis ping failed: %v", err)
				return err
			}
		default:
			if err := r.handleRedisEvent(ctx, &psc, msgChan); err != nil {
				return err
			}
		}
	}
}

func (r *RedisSubscriber) handleRedisEvent(ctx context.Context, psc *redis.PubSubConn, msgChan chan<- GoodsSyncMsg) error {
	switch v := psc.ReceiveWithTimeout(30 * time.Second).(type) {
	case redis.Message:
		var msg GoodsSyncMsg
		if err := json.Unmarshal(v.Data, &msg); err != nil {
			logger.Errorf(ctx, "failed to unmarshal Redis message: %v", err)
			return nil
		}
		select {
		case msgChan <- msg:
		case <-ctx.Done():
			return ctx.Err()
		}
	case redis.Subscription:
		logger.Infof(ctx, "Redis subscription: %s %s %d", v.Kind, v.Channel, v.Count)
	case error:
		return v
	}
	return nil
}

func isRedisConnClosedErr(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "redigo: connection closed") || strings.Contains(errStr, "use of closed network connection") ||
		strings.Contains(errStr, "i/o timeout") || strings.Contains(errStr, "connection reset by peer")
}

func (r *RedisSubscriber) Close() error {
	logger.Infof(context.Background(), "Redis subscriber closed")
	return nil
}

```



```
1. ticker := time.NewTicker(20 * time.Second)
作用：
每20秒主动向Redis发送一次PING，用于保持长连接活跃，防止K8s网络、NAT、LB等中间件或Redis服务端因“空闲”而断开连接。
优点：
能有效防止连接因长时间无数据而被关闭。
20秒属于比较常见的心跳间隔，既不会太频繁增加负载，也不会太久导致连接被回收。
可调建议：
如果你发现连接还是容易断，可以缩短为10秒。
如果Redis和网络都很稳定，可以适当延长到30秒甚至1分钟。
2. psc.ReceiveWithTimeout(30 * time.Second)
作用：
等待Redis消息，超时时间为30秒。如果30秒内没有消息，会返回超时错误，进入下一轮循环。
优点：
能及时响应ctx.Done()和心跳PING，不会因为一直阻塞在Receive而错过心跳或退出信号。
30秒属于比较宽松的等待时间，适合大部分业务场景。
可调建议：
如果你希望更快响应心跳或退出，可以缩短为10~20秒。
如果消息量很大，且对心跳和退出响应不敏感，可以适当延长。
3. 两者配合的合理性
ticker(20s) 比 ReceiveWithTimeout(30s) 更短，保证即使没有消息，心跳也能按时发出。
如果ReceiveWithTimeout超时，下一轮循环会立刻检查ticker，如果到点就发PING。
这样设计可以兼顾心跳及时性和消息处理效率。
4. 典型优化建议
心跳间隔（ticker）：10~30秒都合理，建议20秒是个不错的折中。
接收超时（ReceiveWithTimeout）：20~60秒都合理，建议30秒是常用选择。
两者不要设置成一样，否则可能出现心跳和接收总是“撞车”，建议心跳间隔略短于接收超时。
```

这个 30 * time.Second 是 Redis PubSub 接收消息的超时时间。

## 详细解释

### psc.ReceiveWithTimeout(30 * time.Second) 的作用

1. 等待Redis消息的超时时间

- 如果30秒内没有收到任何消息，会返回超时错误

- 如果30秒内收到消息，立即返回消息内容

1. 防止无限阻塞

- 没有超时的话，Receive() 会无限期等待

- 有超时可以定期检查其他条件（如context取消）































## 测试跨节点发布/订阅

### 1. 在节点1上订阅

*# 在节点1上订阅*

redis-cli -h 10.248.33.114 -p 7001 SUBSCRIBE goods_sync_broadcast

### 2. 在节点2上发布

在另一个终端，在节点2上发布消息：

*# 在节点2上发布*

redis-cli -h 10.248.33.114 -p 7002 PUBLISH goods_sync_broadcast '{"shop_id":"5e7dbfa6e4f3320016e9b7d1","task_id":"3503595164172200701","success":true}'

### 3. 在节点3上发布

在另一个终端，在节点3上发布消息：

*# 在节点3上发布*

redis-cli -h 10.248.33.114 -p 7003 PUBLISH goods_sync_broadcast '{"shop_id":"5e7dbfa6e4f3320016e9b7d1","task_id":"3503595164172200701","success":true}'

## 如果跨节点发布/订阅仍然不工作

如果跨节点发布/订阅仍然不工作，可能的原因有：



### 2. 检查Redis集群状态

bash

Apply to config.go

Run

*# 检查Redis集群状态*

redis-cli -h 10.248.33.114 -p 7001 CLUSTER INFO

### 3. 检查Redis集群配置

bash

Apply to config.go

Run

*# 检查Redis集群配置*

redis-cli -h 10.248.33.114 -p 7001 CONFIG GET cluster-require-full-coverage

### 4. 检查

### 1. 重新设置Redis集群配置

需要重新设置Redis集群的跨节点发布/订阅配置：

```
# 重新设置Redis集群配置
redis-cli -h 10.248.33.114 -p 7001 CONFIG SET cluster-require-full-coverage no
redis-cli -h 10.248.33.114 -p 7002 CONFIG SET cluster-require-full-coverage no
redis-cli -h 10.248.33.114 -p 7003 CONFIG SET cluster-require-full-coverage no

# 重启Redis集群
redis-cli -h 10.248.33.114 -p 7001 CLUSTER RESET
redis-cli -h 10.248.33.114 -p 7002 CLUSTER RESET
redis-cli -h 10.248.33.114 -p 7003 CLUSTER RESET
```



```
/root # redis-cli -h 10.248.33.114 -p 7001 SUBSCRIBE goods_sync_broadcast
Reading messages... (press Ctrl-C to quit)
1) "subscribe"
2) "goods_sync_broadcast"
3) (integer) 1

```





```
/root # redis-cli -h 10.248.33.114 -p 7001 CLUSTER SLOTS                                                    
1) 1) (integer) 5461
   2) (integer) 10922
   3) 1) "10.248.33.114"
      2) (integer) 7002
      3) "51fcb6eb0a0c52a401604a44b4b4e30f1b9c45a0"
2) 1) (integer) 10923
   2) (integer) 16383
   3) 1) "10.248.33.114"
      2) (integer) 7003
      3) "fe8a79175045f08e8b2da659181885d7e8cfea68"
3) 1) (integer) 0
   2) (integer) 5460
   3) 1) "10.248.33.114"
      2) (integer) 7001
      3) "d236ccf97162d1ffa0d68689d428b2936cd50b07"

/root # 
```

从槽位分配可以看出，Redis集群的槽位分配是正常的。





```
/root # redis-cli -h 10.248.33.114 -p 7001 CLUSTER INFO
cluster_state:ok
cluster_slots_assigned:16384
cluster_slots_ok:16384
cluster_slots_pfail:0
cluster_slots_fail:0
cluster_known_nodes:3
cluster_size:3
cluster_current_epoch:3
cluster_my_epoch:1
cluster_stats_messages_ping_sent:856686
cluster_stats_messages_pong_sent:787361
cluster_stats_messages_publish_sent:12
cluster_stats_messages_sent:1644059
cluster_stats_messages_ping_received:787361
cluster_stats_messages_pong_received:856672
cluster_stats_messages_fail_received:5
cluster_stats_messages_publish_received:23659
cluster_stats_messages_received:1667697

/root # redis-cli -h 10.248.33.114 -p 7001 CLUSTER NODES
51fcb6eb0a0c52a401604a44b4b4e30f1b9c45a0 10.248.33.114:7002@17002 master - 0 1751799477761 2 connected 5461-10922
fe8a79175045f08e8b2da659181885d7e8cfea68 10.248.33.114:7003@17003 master - 0 1751799476756 3 connected 10923-16383
d236ccf97162d1ffa0d68689d428b2936cd50b07 10.248.33.114:7001@17001 myself,master - 0 1751799473000 1 connected 0-5460

/root #                                                                                                     
```

