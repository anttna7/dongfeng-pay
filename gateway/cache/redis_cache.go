package cache

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"time"
)

// RedisCache Redis缓存接口
// 注意：本实现提供了接口定义和内存缓存回退方案
// 生产环境建议使用真实的Redis实现（如 github.com/gomodule/redigo 或 github.com/go-redis/redis）

type RedisCache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	GetObject(key string, obj interface{}) error
	Delete(key string) error
	Exists(key string) bool
	SetNX(key string, value interface{}, expiration time.Duration) bool
	Expire(key string, expiration time.Duration) error
	TTL(key string) (time.Duration, error)
	IncrBy(key string, value int64) (int64, error)
	DecrBy(key string, value int64) (int64, error)
}

// MemoryCache 内存缓存实现（用于开发环境或Redis不可用时的回退方案）
type MemoryCache struct {
	data       map[string]*CacheItem
	expiration map[string]time.Time
}

type CacheItem struct {
	Value     interface{}
	ExpiredAt time.Time
}

var memCache *MemoryCache

// InitMemoryCache 初始化内存缓存
func InitMemoryCache() *MemoryCache {
	if memCache == nil {
		memCache = &MemoryCache{
			data:       make(map[string]*CacheItem),
			expiration: make(map[string]time.Time),
		}
		// 启动清理过期数据的协程
		go memCache.cleanExpired()
	}
	return memCache
}

// Set 设置缓存
func (m *MemoryCache) Set(key string, value interface{}, expiration time.Duration) error {
	expiredAt := time.Now().Add(expiration)
	m.data[key] = &CacheItem{
		Value:     value,
		ExpiredAt: expiredAt,
	}
	m.expiration[key] = expiredAt
	return nil
}

// Get 获取缓存
func (m *MemoryCache) Get(key string) (string, error) {
	item, exists := m.data[key]
	if !exists {
		return "", fmt.Errorf("key not found: %s", key)
	}

	if time.Now().After(item.ExpiredAt) {
		delete(m.data, key)
		delete(m.expiration, key)
		return "", fmt.Errorf("key expired: %s", key)
	}

	return fmt.Sprintf("%v", item.Value), nil
}

// GetObject 获取对象
func (m *MemoryCache) GetObject(key string, obj interface{}) error {
	value, err := m.Get(key)
	if err != nil {
		return err
	}

	// 尝试JSON反序列化
	if err := json.Unmarshal([]byte(value), obj); err != nil {
		return fmt.Errorf("failed to unmarshal object: %v", err)
	}

	return nil
}

// Delete 删除缓存
func (m *MemoryCache) Delete(key string) error {
	delete(m.data, key)
	delete(m.expiration, key)
	return nil
}

// Exists 检查键是否存在
func (m *MemoryCache) Exists(key string) bool {
	item, exists := m.data[key]
	if !exists {
		return false
	}

	if time.Now().After(item.ExpiredAt) {
		delete(m.data, key)
		delete(m.expiration, key)
		return false
	}

	return true
}

// SetNX 只有键不存在时才设置
func (m *MemoryCache) SetNX(key string, value interface{}, expiration time.Duration) bool {
	if m.Exists(key) {
		return false
	}

	m.Set(key, value, expiration)
	return true
}

// Expire 设置过期时间
func (m *MemoryCache) Expire(key string, expiration time.Duration) error {
	item, exists := m.data[key]
	if !exists {
		return fmt.Errorf("key not found: %s", key)
	}

	expiredAt := time.Now().Add(expiration)
	item.ExpiredAt = expiredAt
	m.expiration[key] = expiredAt
	return nil
}

// TTL 获取剩余生存时间
func (m *MemoryCache) TTL(key string) (time.Duration, error) {
	item, exists := m.data[key]
	if !exists {
		return 0, fmt.Errorf("key not found: %s", key)
	}

	ttl := time.Until(item.ExpiredAt)
	if ttl < 0 {
		delete(m.data, key)
		delete(m.expiration, key)
		return 0, fmt.Errorf("key expired: %s", key)
	}

	return ttl, nil
}

// IncrBy 增加值
func (m *MemoryCache) IncrBy(key string, value int64) (int64, error) {
	item, exists := m.data[key]
	if !exists {
		m.Set(key, value, 24*time.Hour)
		return value, nil
	}

	currentValue, ok := item.Value.(int64)
	if !ok {
		// 尝试转换
		if intVal, ok := item.Value.(int); ok {
			currentValue = int64(intVal)
		} else {
			return 0, fmt.Errorf("value is not an integer")
		}
	}

	newValue := currentValue + value
	item.Value = newValue
	return newValue, nil
}

// DecrBy 减少值
func (m *MemoryCache) DecrBy(key string, value int64) (int64, error) {
	return m.IncrBy(key, -value)
}

// cleanExpired 清理过期数据
func (m *MemoryCache) cleanExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		for key, expiredAt := range m.expiration {
			if now.After(expiredAt) {
				delete(m.data, key)
				delete(m.expiration, key)
			}
		}
	}
}

// 全局缓存实例
var globalCache RedisCache

// InitCache 初始化缓存
// 注意：这里使用内存缓存，生产环境应该替换为真实的Redis实现
func InitCache() RedisCache {
	if globalCache == nil {
		logs.Info("初始化内存缓存（开发环境）")
		globalCache = InitMemoryCache()
	}
	return globalCache
}

// GetCache 获取缓存实例
func GetCache() RedisCache {
	if globalCache == nil {
		return InitCache()
	}
	return globalCache
}

// 缓存键前缀
const (
	MerchantPrefix = "merchant:"     // 商户信息缓存前缀
	RoadPrefix     = "road:"         // 通道信息缓存前缀
	OrderPrefix    = "order:"        // 订单信息缓存前缀
	AccountPrefix  = "account:"      // 账户信息缓存前缀
	RateLimitPrefix = "ratelimit:"   // 限流缓存前缀
	LockPrefix     = "lock:"         // 分布式锁前缀
)

// 缓存时间
const (
	MerchantCacheTTL = 30 * time.Minute // 商户信息缓存30分钟
	RoadCacheTTL     = 15 * time.Minute // 通道信息缓存15分钟
	OrderCacheTTL    = 10 * time.Minute // 订单信息缓存10分钟
	AccountCacheTTL  = 5 * time.Minute  // 账户信息缓存5分钟（余额变化频繁）
	LockTTL          = 30 * time.Second // 分布式锁30秒
)

// 便捷方法

// SetMerchantCache 设置商户缓存
func SetMerchantCache(merchantUID string, data interface{}) error {
	cache := GetCache()
	key := MerchantPrefix + merchantUID

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return cache.Set(key, string(jsonData), MerchantCacheTTL)
}

// GetMerchantCache 获取商户缓存
func GetMerchantCache(merchantUID string, result interface{}) error {
	cache := GetCache()
	key := MerchantPrefix + merchantUID
	return cache.GetObject(key, result)
}

// DeleteMerchantCache 删除商户缓存
func DeleteMerchantCache(merchantUID string) error {
	cache := GetCache()
	key := MerchantPrefix + merchantUID
	return cache.Delete(key)
}

// SetRoadCache 设置通道缓存
func SetRoadCache(roadUID string, data interface{}) error {
	cache := GetCache()
	key := RoadPrefix + roadUID

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return cache.Set(key, string(jsonData), RoadCacheTTL)
}

// GetRoadCache 获取通道缓存
func GetRoadCache(roadUID string, result interface{}) error {
	cache := GetCache()
	key := RoadPrefix + roadUID
	return cache.GetObject(key, result)
}

// DeleteRoadCache 删除通道缓存
func DeleteRoadCache(roadUID string) error {
	cache := GetCache()
	key := RoadPrefix + roadUID
	return cache.Delete(key)
}

// SetOrderCache 设置订单缓存
func SetOrderCache(orderID string, data interface{}) error {
	cache := GetCache()
	key := OrderPrefix + orderID

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return cache.Set(key, string(jsonData), OrderCacheTTL)
}

// GetOrderCache 获取订单缓存
func GetOrderCache(orderID string, result interface{}) error {
	cache := GetCache()
	key := OrderPrefix + orderID
	return cache.GetObject(key, result)
}

// DeleteOrderCache 删除订单缓存
func DeleteOrderCache(orderID string) error {
	cache := GetCache()
	key := OrderPrefix + orderID
	return cache.Delete(key)
}

// AcquireLock 获取分布式锁
func AcquireLock(lockKey string, expiration time.Duration) bool {
	cache := GetCache()
	key := LockPrefix + lockKey

	if expiration == 0 {
		expiration = LockTTL
	}

	return cache.SetNX(key, "1", expiration)
}

// ReleaseLock 释放分布式锁
func ReleaseLock(lockKey string) error {
	cache := GetCache()
	key := LockPrefix + lockKey
	return cache.Delete(key)
}
