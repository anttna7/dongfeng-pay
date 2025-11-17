package middleware

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"
	"sync"
	"time"
)

// RateLimiter 限流器结构
type RateLimiter struct {
	rate       int           // 每秒允许的请求数
	burst      int           // 令牌桶容量
	mu         sync.Mutex    // 互斥锁
	tokens     map[string]*TokenBucket
	cleanupTTL time.Duration // 清理过期令牌桶的时间
}

// TokenBucket 令牌桶
type TokenBucket struct {
	tokens    float64   // 当前令牌数
	lastCheck time.Time // 上次检查时间
	mu        sync.Mutex
}

// NewRateLimiter 创建限流器
// rate: 每秒允许的请求数
// burst: 令牌桶容量（允许的突发流量）
func NewRateLimiter(rate, burst int) *RateLimiter {
	limiter := &RateLimiter{
		rate:       rate,
		burst:      burst,
		tokens:     make(map[string]*TokenBucket),
		cleanupTTL: 5 * time.Minute,
	}

	// 启动清理协程
	go limiter.cleanup()

	return limiter
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	bucket, exists := rl.tokens[key]
	if !exists {
		bucket = &TokenBucket{
			tokens:    float64(rl.burst),
			lastCheck: time.Now(),
		}
		rl.tokens[key] = bucket
	}
	rl.mu.Unlock()

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(bucket.lastCheck).Seconds()
	bucket.lastCheck = now

	// 添加新令牌
	bucket.tokens += elapsed * float64(rl.rate)
	if bucket.tokens > float64(rl.burst) {
		bucket.tokens = float64(rl.burst)
	}

	// 检查是否有可用令牌
	if bucket.tokens >= 1 {
		bucket.tokens--
		return true
	}

	return false
}

// cleanup 定期清理过期的令牌桶
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, bucket := range rl.tokens {
			bucket.mu.Lock()
			if now.Sub(bucket.lastCheck) > rl.cleanupTTL {
				delete(rl.tokens, key)
			}
			bucket.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

// RateLimitMiddleware 限流中间件
// 全局限流器实例
var (
	globalLimiter     *RateLimiter
	merchantLimiter   *RateLimiter
	paymentLimiter    *RateLimiter
	initLimiterOnce   sync.Once
)

// InitRateLimiters 初始化限流器
func InitRateLimiters() {
	initLimiterOnce.Do(func() {
		// 全局限流：每秒1000个请求，突发2000
		globalLimiter = NewRateLimiter(1000, 2000)

		// 商户级别限流：每秒100个请求，突发200
		merchantLimiter = NewRateLimiter(100, 200)

		// 支付接口限流：每秒50个请求，突发100（更严格）
		paymentLimiter = NewRateLimiter(50, 100)

		logs.Info("Rate limiters initialized successfully")
	})
}

// GlobalRateLimitFilter 全局限流过滤器
func GlobalRateLimitFilter(ctx *context.Context) {
	InitRateLimiters()

	// 使用IP作为限流key
	clientIP := ctx.Input.IP()

	if !globalLimiter.Allow(clientIP) {
		logs.Warn("Global rate limit exceeded for IP: %s", clientIP)
		ctx.Output.SetStatus(429)
		ctx.Output.JSON(map[string]interface{}{
			"code":    "RATE_LIMIT_EXCEEDED",
			"message": "请求过于频繁，请稍后再试",
		}, false, false)
		return
	}
}

// MerchantRateLimitFilter 商户级别限流过滤器
func MerchantRateLimitFilter(ctx *context.Context) {
	InitRateLimiters()

	// 从请求参数中获取商户ID
	merchantUid := ctx.Input.Query("merchant_id")
	if merchantUid == "" {
		merchantUid = ctx.Input.Query("merchantId")
	}

	// 如果没有商户ID，使用IP作为限流key
	if merchantUid == "" {
		merchantUid = ctx.Input.IP()
	}

	key := fmt.Sprintf("merchant:%s", merchantUid)

	if !merchantLimiter.Allow(key) {
		logs.Warn("Merchant rate limit exceeded for: %s", merchantUid)
		ctx.Output.SetStatus(429)
		ctx.Output.JSON(map[string]interface{}{
			"code":    "MERCHANT_RATE_LIMIT_EXCEEDED",
			"message": "商户请求过于频繁，请稍后再试",
		}, false, false)
		return
	}
}

// PaymentRateLimitFilter 支付接口限流过滤器（最严格）
func PaymentRateLimitFilter(ctx *context.Context) {
	InitRateLimiters()

	// 使用商户ID + IP作为限流key
	merchantUid := ctx.Input.Query("merchant_id")
	if merchantUid == "" {
		merchantUid = ctx.Input.Query("merchantId")
	}

	clientIP := ctx.Input.IP()
	key := fmt.Sprintf("payment:%s:%s", merchantUid, clientIP)

	if !paymentLimiter.Allow(key) {
		logs.Warn("Payment rate limit exceeded for merchant: %s, IP: %s", merchantUid, clientIP)
		ctx.Output.SetStatus(429)
		ctx.Output.JSON(map[string]interface{}{
			"code":    "PAYMENT_RATE_LIMIT_EXCEEDED",
			"message": "支付请求过于频繁，请稍后再试",
		}, false, false)
		return
	}
}
