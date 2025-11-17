# Redis 缓存使用说明

## 当前实现

当前使用**内存缓存**作为开发环境的实现。这是一个简单的基于 map 的缓存系统，适用于：
- 开发环境
- 单机部署
- Redis 不可用时的回退方案

## 生产环境升级到真实 Redis

### 1. 安装 Redis 客户端库

选择以下任一 Redis 客户端库：

#### 方案一：使用 go-redis（推荐）

```bash
cd gateway
go get github.com/go-redis/redis/v8
```

#### 方案二：使用 redigo

```bash
cd gateway
go get github.com/gomodule/redigo/redis
```

### 2. 配置 Redis 连接

在 `conf/app.conf` 中添加 Redis 配置：

```ini
[redis]
host = localhost
port = 6379
password =
db = 0
max_idle = 10
max_active = 100
idle_timeout = 300
```

### 3. 实现真实的 Redis 客户端

创建新文件 `cache/redis_client.go`：

```go
package cache

import (
    "context"
    "encoding/json"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/beego/beego/v2/server/web"
)

type RedisClient struct {
    client *redis.Client
    ctx    context.Context
}

func NewRedisClient() (*RedisClient, error) {
    host, _ := web.AppConfig.String("redis::host")
    port, _ := web.AppConfig.String("redis::port")
    password, _ := web.AppConfig.String("redis::password")
    db, _ := web.AppConfig.Int("redis::db")

    client := redis.NewClient(&redis.Options{
        Addr:     host + ":" + port,
        Password: password,
        DB:       db,
    })

    ctx := context.Background()

    // 测试连接
    _, err := client.Ping(ctx).Result()
    if err != nil {
        return nil, err
    }

    return &RedisClient{
        client: client,
        ctx:    ctx,
    }, nil
}

func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
    return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
    return r.client.Get(r.ctx, key).Result()
}

func (r *RedisClient) GetObject(key string, obj interface{}) error {
    val, err := r.Get(key)
    if err != nil {
        return err
    }
    return json.Unmarshal([]byte(val), obj)
}

func (r *RedisClient) Delete(key string) error {
    return r.client.Del(r.ctx, key).Err()
}

func (r *RedisClient) Exists(key string) bool {
    result, _ := r.client.Exists(r.ctx, key).Result()
    return result > 0
}

func (r *RedisClient) SetNX(key string, value interface{}, expiration time.Duration) bool {
    result, _ := r.client.SetNX(r.ctx, key, value, expiration).Result()
    return result
}

func (r *RedisClient) Expire(key string, expiration time.Duration) error {
    return r.client.Expire(r.ctx, key, expiration).Err()
}

func (r *RedisClient) TTL(key string) (time.Duration, error) {
    return r.client.TTL(r.ctx, key).Result()
}

func (r *RedisClient) IncrBy(key string, value int64) (int64, error) {
    return r.client.IncrBy(r.ctx, key, value).Result()
}

func (r *RedisClient) DecrBy(key string, value int64) (int64, error) {
    return r.client.DecrBy(r.ctx, key, value).Result()
}
```

### 4. 修改初始化代码

修改 `cache/redis_cache.go` 中的 `InitCache()` 函数：

```go
func InitCache() RedisCache {
    if globalCache == nil {
        // 尝试连接 Redis
        redisClient, err := NewRedisClient()
        if err != nil {
            logs.Warn("Redis 连接失败，使用内存缓存: %v", err)
            globalCache = InitMemoryCache()
        } else {
            logs.Info("Redis 连接成功")
            globalCache = redisClient
        }
    }
    return globalCache
}
```

## 使用示例

### 缓存商户信息

```go
import "gateway/cache"

// 设置商户缓存
merchantInfo := &merchant.MerchantInfo{
    MerchantUid:  "M001",
    MerchantName: "测试商户",
}
err := cache.SetMerchantCache("M001", merchantInfo)

// 获取商户缓存
var cachedMerchant merchant.MerchantInfo
err = cache.GetMerchantCache("M001", &cachedMerchant)
if err != nil {
    // 缓存未命中，从数据库查询
}

// 删除商户缓存（商户信息更新时）
cache.DeleteMerchantCache("M001")
```

### 分布式锁

```go
import "gateway/cache"

// 获取锁
if cache.AcquireLock("order:123456", 30*time.Second) {
    defer cache.ReleaseLock("order:123456")

    // 执行业务逻辑
    // ...
} else {
    // 获取锁失败，订单正在处理中
}
```

### 自定义缓存

```go
import "gateway/cache"

cache := cache.GetCache()

// 设置缓存
cache.Set("mykey", "myvalue", 10*time.Minute)

// 获取缓存
value, err := cache.Get("mykey")

// 检查键是否存在
if cache.Exists("mykey") {
    // ...
}

// 原子递增
count, err := cache.IncrBy("counter", 1)
```

## 性能优化建议

### 1. 合理设置过期时间

- 商户信息：30分钟（变化不频繁）
- 通道信息：15分钟
- 订单信息：10分钟
- 账户余额：5分钟（变化频繁，缓存时间短）

### 2. 使用缓存穿透保护

```go
// 使用空值缓存防止缓存穿透
func GetMerchantWithCache(merchantUID string) (*merchant.MerchantInfo, error) {
    var cached merchant.MerchantInfo
    err := cache.GetMerchantCache(merchantUID, &cached)
    if err == nil {
        return &cached, nil
    }

    // 从数据库查询
    info, err := merchant.GetMerchantByUid(merchantUID)
    if err != nil {
        // 设置空值缓存（短时间）
        cache.SetMerchantCache(merchantUID, nil)
        return nil, err
    }

    // 更新缓存
    cache.SetMerchantCache(merchantUID, info)
    return info, nil
}
```

### 3. 批量操作

使用 Pipeline 批量操作以减少网络往返：

```go
// 在 RedisClient 中添加 Pipeline 支持
func (r *RedisClient) Pipeline() redis.Pipeliner {
    return r.client.Pipeline()
}
```

## 监控与运维

### Redis 监控指标

- 连接数
- 命令执行速率
- 内存使用率
- 缓存命中率
- 慢查询

### 常见问题

1. **缓存雪崩**：大量缓存同时过期
   - 解决：设置随机过期时间

2. **缓存击穿**：热点数据过期
   - 解决：使用分布式锁或永不过期 + 异步更新

3. **缓存穿透**：查询不存在的数据
   - 解决：布隆过滤器或空值缓存

## 注意事项

1. 缓存更新策略
   - 数据更新时及时删除缓存
   - 使用 Cache Aside 模式

2. 序列化
   - 推荐使用 JSON 序列化
   - 大对象考虑使用 Protocol Buffers

3. 数据一致性
   - 缓存不保证强一致性
   - 关键数据以数据库为准
