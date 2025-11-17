package middleware

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"
)

// Metrics 性能指标
type Metrics struct {
	TotalRequests    int64         // 总请求数
	CurrentRequests  int64         // 当前请求数
	SuccessRequests  int64         // 成功请求数
	FailedRequests   int64         // 失败请求数
	TotalLatency     int64         // 总延迟(纳秒)
	MinLatency       time.Duration // 最小延迟
	MaxLatency       time.Duration // 最大延迟
	StartTime        time.Time     // 启动时间
	RequestsPerRoute map[string]int64 // 每个路由的请求数
	mu               sync.RWMutex
}

// GlobalMetrics 全局指标
var GlobalMetrics = &Metrics{
	StartTime:        time.Now(),
	RequestsPerRoute: make(map[string]int64),
	MinLatency:       time.Hour,
}

// Update 更新指标
func (m *Metrics) Update(route string, latency time.Duration, success bool) {
	atomic.AddInt64(&m.TotalRequests, 1)
	atomic.AddInt64(&m.TotalLatency, int64(latency))

	if success {
		atomic.AddInt64(&m.SuccessRequests, 1)
	} else {
		atomic.AddInt64(&m.FailedRequests, 1)
	}

	// 更新最小/最大延迟
	m.mu.Lock()
	if latency < m.MinLatency {
		m.MinLatency = latency
	}
	if latency > m.MaxLatency {
		m.MaxLatency = latency
	}
	// 更新路由统计
	m.RequestsPerRoute[route]++
	m.mu.Unlock()
}

// GetStats 获取统计信息
func (m *Metrics) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	totalReq := atomic.LoadInt64(&m.TotalRequests)
	totalLat := atomic.LoadInt64(&m.TotalLatency)
	successReq := atomic.LoadInt64(&m.SuccessRequests)
	failedReq := atomic.LoadInt64(&m.FailedRequests)
	currentReq := atomic.LoadInt64(&m.CurrentRequests)

	avgLatency := time.Duration(0)
	if totalReq > 0 {
		avgLatency = time.Duration(totalLat / totalReq)
	}

	uptime := time.Since(m.StartTime)
	qps := float64(totalReq) / uptime.Seconds()

	// 内存统计
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return map[string]interface{}{
		"总请求数":   totalReq,
		"成功请求数":  successReq,
		"失败请求数":  failedReq,
		"当前请求数":  currentReq,
		"平均延迟":   avgLatency.String(),
		"最小延迟":   m.MinLatency.String(),
		"最大延迟":   m.MaxLatency.String(),
		"运行时间":   uptime.String(),
		"QPS":     fmt.Sprintf("%.2f", qps),
		"成功率":    fmt.Sprintf("%.2f%%", float64(successReq)/float64(totalReq)*100),
		"内存使用":   fmt.Sprintf("%.2f MB", float64(memStats.Alloc)/(1024*1024)),
		"总内存":    fmt.Sprintf("%.2f MB", float64(memStats.TotalAlloc)/(1024*1024)),
		"协程数":    runtime.NumGoroutine(),
		"路由统计":   m.RequestsPerRoute,
	}
}

// Reset 重置指标
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	atomic.StoreInt64(&m.TotalRequests, 0)
	atomic.StoreInt64(&m.CurrentRequests, 0)
	atomic.StoreInt64(&m.SuccessRequests, 0)
	atomic.StoreInt64(&m.FailedRequests, 0)
	atomic.StoreInt64(&m.TotalLatency, 0)

	m.MinLatency = time.Hour
	m.MaxLatency = 0
	m.StartTime = time.Now()
	m.RequestsPerRoute = make(map[string]int64)
}

// PerformanceMonitorFilter 性能监控过滤器
func PerformanceMonitorFilter(ctx *context.Context) {
	startTime := time.Now()
	route := ctx.Request.URL.Path

	// 增加当前请求数
	atomic.AddInt64(&GlobalMetrics.CurrentRequests, 1)

	// 请求完成后的处理
	defer func() {
		// 减少当前请求数
		atomic.AddInt64(&GlobalMetrics.CurrentRequests, -1)

		// 计算延迟
		latency := time.Since(startTime)

		// 判断是否成功（2xx, 3xx）
		success := ctx.ResponseWriter.Status >= 200 && ctx.ResponseWriter.Status < 400

		// 更新指标
		GlobalMetrics.Update(route, latency, success)

		// 记录慢请求
		if latency > 1*time.Second {
			logs.Warn("慢请求: %s, 耗时: %v, 状态: %d",
				route, latency, ctx.ResponseWriter.Status)
		}
	}()
}

// HealthCheckFilter 健康检查过滤器
func HealthCheckFilter(ctx *context.Context) {
	if ctx.Request.URL.Path == "/health" || ctx.Request.URL.Path == "/ping" {
		stats := GlobalMetrics.GetStats()

		ctx.Output.JSON(map[string]interface{}{
			"status":  "ok",
			"metrics": stats,
		}, false, false)

		// 阻止继续执行
		panic("healthcheck") // 这里用 panic 是为了阻止继续执行
	}
}

// MetricsReporter 定时报告性能指标
func MetricsReporter(interval time.Duration) {
	ticker := time.Ticker(interval)
	defer ticker.Stop()

	for range ticker.C {
		stats := GlobalMetrics.GetStats()

		logs.Info("=== 性能指标报告 ===")
		for key, value := range stats {
			if key == "路由统计" {
				logs.Info("%s:", key)
				if routes, ok := value.(map[string]int64); ok {
					for route, count := range routes {
						logs.Info("  %s: %d", route, count)
					}
				}
			} else {
				logs.Info("%s: %v", key, value)
			}
		}
		logs.Info("==================")
	}
}

// 启动性能监控报告器
func init() {
	// 每5分钟报告一次性能指标
	go MetricsReporter(5 * time.Minute)
}
