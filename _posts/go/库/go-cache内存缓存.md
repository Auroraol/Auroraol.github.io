# go-cache

> 一个适用于单机应用的内存键值存储/缓存库, 支持高并发环境下稳定运行。

## 一、智能过期清理机制：自动释放闲置资源

go-cache的核心优势在于其灵活的过期策略，通过合理设置键值对的生命周期，自动清理不再需要的数据，从源头避免内存无限增长。

```
NoExpiration（-1）：表示键值对永不过期
DefaultExpiration（0）：使用缓存的默认过期时间
```

创建缓存实例时，可通过New函数设置默认过期时间和清理间隔：

```
// 创建默认过期时间为5分钟，清理间隔为10分钟的缓存
c := cache.New(5*time.Minute, 10*time.Minute)
```

过期清理通过后台协程实现，定期扫描并删除过期项。这种机制确保即使应用持续运行，也不会因未清理的过期数据导致内存溢出。

## 二、分片存储设计：降低锁竞争提升性能

对于大型应用，单一缓存实例可能成为性能瓶颈。go-cache提供了分片缓存（sharded cache）解决方案，通过将数据分散到多个子缓存中，显著降低锁竞争。

sharded.go实现了分片缓存的核心逻辑，通过哈希算法将不同的键分配到不同的子缓存：

```go
func (sc *shardedCache) Set(k string, x interface{}, d time.Duration) {
    // 根据键的哈希值选择分片
    shard := sc.shards[sc.hash(k)%sc.numShards]
    shard.Set(k, x, d)
}
```

## 案例

```go
package memory

import (
    "errors"
    "time"

    "github.com/patrickmn/go-cache"
)

const (
    aipassClientVersionCacheExpiration      = time.Minute
    aipassClientVersionCacheCleanupInterval = time.Minute
)

var AipassClientVersionCache = NewCache(CacheConfig{
    Expiration:      aipassClientVersionCacheExpiration,
    CleanupInterval: aipassClientVersionCacheCleanupInterval,
})

type CacheConfig struct {
    Expiration      time.Duration
    CleanupInterval time.Duration
}

type Cache struct {
    cache *cache.Cache
}

func NewCache(config CacheConfig) *Cache {
    return NewMemoryCacheWithConfig(config.Expiration, config.CleanupInterval)
}

func NewMemoryCacheWithConfig(expiration, cleanupInterval time.Duration) *Cache {
    if expiration == 0 {
        expiration = cache.NoExpiration
    }
    if cleanupInterval == 0 {
        cleanupInterval = cache.NoExpiration
    }
    return &Cache{
        cache: cache.New(expiration, cleanupInterval),
    }
}

func (c *Cache) Get(key string) (interface{}, bool) {
    if c == nil || c.cache == nil {
        return nil, false
    }
    return c.cache.Get(key)
}

func (c *Cache) Add(key string, x interface{}) error {
    if c == nil || c.cache == nil {
        return errors.New("memory cache is nil")
    }
    return c.cache.Add(key, x, cache.DefaultExpiration)
}

func (c *Cache) AddWithExpiration(key string, x interface{}, d time.Duration) error {
    if c == nil || c.cache == nil {
        return errors.New("memory cache is nil")
    }
    return c.cache.Add(key, x, d)
}

func (c *Cache) Set(key string, x interface{}) {
    if c == nil || c.cache == nil {
        return
    }
    c.cache.Set(key, x, cache.DefaultExpiration)
}

func (c *Cache) SetWithExpiration(key string, x interface{}, d time.Duration) {
    if c == nil || c.cache == nil {
        return
    }
    c.cache.Set(key, x, d)
}

```

