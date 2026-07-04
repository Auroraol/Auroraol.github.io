[buraksezer/consistent: Consistent hashing with bounded loads in Golang](https://github.com/buraksezer/consistent)

带有限负载的一致性哈希库，它能够帮助分布式系统实现数据的均匀分布和高效查询

**实际场景：分布式本地缓存**

假设你有 3 台缓存节点：

```
cache-node-1
cache-node-2
cache-node-3
```

业务请求来了一个用户 ID：

```
user:1001
```

你希望它固定落到某一台缓存节点：

```
user:1001 -> cache-node-2
```

这样下次查 `user:1001`，还会去 `cache-node-2`，缓存命中率更高。

代码大概是：

```
package main

import (
	"fmt"

	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
)

type CacheNode string

func (n CacheNode) String() string {
	return string(n)
}

type Hasher struct{}

func (h Hasher) Sum64(data []byte) uint64 {
	return xxhash.Sum64(data)
}

func main() {
	cfg := consistent.Config{
		PartitionCount:    271,
		ReplicationFactor: 20,
		Load:              1.25,
		Hasher:            Hasher{},
	}

	c := consistent.New(nil, cfg)

	c.Add(CacheNode("cache-node-1"))
	c.Add(CacheNode("cache-node-2"))
	c.Add(CacheNode("cache-node-3"))

	key := []byte("user:1001")

	node := c.LocateKey(key)
	fmt.Println(node.String())
}
```

这段代码的意义是：

```
user:1001 -> 找到负责它的缓存节点
```

然后你就可以把请求发给这个节点：

```
node := c.LocateKey([]byte("user:1001"))
// 向 node.String() 对应的服务发请求
```

比如：

```
cache-node-2.Get("user:1001")
```

**扩容时怎么用**

原来：

```
cache-node-1
cache-node-2
cache-node-3
```

现在加一台：

```
c.Add(CacheNode("cache-node-4"))
```

一致性哈希会让一部分 key 改归属：

```
部分 key 从 cache-node-1/2/3 改到 cache-node-4
大部分 key 仍然留在原节点
```

这就是它的价值：**扩容时不是所有 key 都重新分配，只迁移一小部分映射关系。**

但是注意：这个库不会帮你搬数据。它只告诉你：

```
这个 key 现在应该属于谁
```

如果你存的是本地缓存，通常不需要主动迁移，新的请求自然会落到新节点，缓存慢慢重建。

如果你存的是持久数据，就需要你自己做迁移逻辑：

```
计算哪些 partition/key 归属变化
从老节点读数据
写到新节点
删除老节点副本
```

**适合的真实应用**

常见用法有：

```
1. 分布式本地缓存路由
2. 多个后端 worker 的任务分配
3. 分布式爬虫 URL 分配
4. 多个存储节点之间的数据归属计算
5. 多个 RPC 节点之间按用户 ID 做稳定路由
```

比如任务分配：

```
order:1001 -> worker-1
order:1002 -> worker-3
order:1003 -> worker-2
```

worker 扩容时，只让少部分订单任务换 worker，避免全量打乱。

一句话总结：

**`buraksezer/consistent` 适合你自己实现“key 到节点”的稳定路由。它解决的是分配问题，不负责网络请求，也不负责数据迁移。**