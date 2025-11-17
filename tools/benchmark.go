package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// BenchmarkConfig 压测配置
type BenchmarkConfig struct {
	URL         string        // 测试URL
	Method      string        // HTTP方法
	Concurrency int           // 并发数
	Duration    time.Duration // 持续时间
	Timeout     time.Duration // 请求超时
	Body        string        // 请求体
}

// BenchmarkResult 压测结果
type BenchmarkResult struct {
	TotalRequests   int64         // 总请求数
	SuccessRequests int64         // 成功请求数
	FailedRequests  int64         // 失败请求数
	TotalDuration   time.Duration // 总耗时
	MinLatency      time.Duration // 最小延迟
	MaxLatency      time.Duration // 最大延迟
	AvgLatency      time.Duration // 平均延迟
	QPS             float64       // 每秒请求数
}

// Benchmarker 压测器
type Benchmarker struct {
	config  BenchmarkConfig
	client  *http.Client
	results BenchmarkResult

	latencies     []time.Duration
	latenciesMux  sync.Mutex
	totalRequests int64
	successCount  int64
	failedCount   int64
}

// NewBenchmarker 创建压测器
func NewBenchmarker(config BenchmarkConfig) *Benchmarker {
	return &Benchmarker{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
		latencies: make([]time.Duration, 0, 10000),
	}
}

// Run 执行压测
func (b *Benchmarker) Run() *BenchmarkResult {
	fmt.Println("\n" + "="*70)
	fmt.Println("   dongfeng-pay 性能压测工具")
	fmt.Println("="*70)
	fmt.Printf("目标URL: %s\n", b.config.URL)
	fmt.Printf("并发数: %d\n", b.config.Concurrency)
	fmt.Printf("持续时间: %v\n", b.config.Duration)
	fmt.Println("="*70)
	fmt.Println("\n开始压测...")

	startTime := time.Now()
	stopChan := make(chan bool)
	var wg sync.WaitGroup

	// 启动worker
	for i := 0; i < b.config.Concurrency; i++ {
		wg.Add(1)
		go b.worker(i, stopChan, &wg)
	}

	// 等待指定时间后停止
	time.Sleep(b.config.Duration)
	close(stopChan)
	wg.Wait()

	endTime := time.Now()
	totalDuration := endTime.Sub(startTime)

	// 计算结果
	b.calculateResults(totalDuration)

	// 打印结果
	b.printResults()

	return &b.results
}

// worker 工作协程
func (b *Benchmarker) worker(id int, stopChan chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-stopChan:
			return
		default:
			b.doRequest()
		}
	}
}

// doRequest 执行单次请求
func (b *Benchmarker) doRequest() {
	atomic.AddInt64(&b.totalRequests, 1)

	start := time.Now()

	var resp *http.Response
	var err error

	if b.config.Method == "POST" {
		resp, err = b.client.Post(
			b.config.URL,
			"application/json",
			strings.NewReader(b.config.Body),
		)
	} else {
		resp, err = b.client.Get(b.config.URL)
	}

	latency := time.Since(start)

	// 记录延迟
	b.latenciesMux.Lock()
	b.latencies = append(b.latencies, latency)
	b.latenciesMux.Unlock()

	if err != nil {
		atomic.AddInt64(&b.failedCount, 1)
		return
	}
	defer resp.Body.Close()

	// 读取响应体（确保完整接收）
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		atomic.AddInt64(&b.successCount, 1)
	} else {
		atomic.AddInt64(&b.failedCount, 1)
	}
}

// calculateResults 计算结果
func (b *Benchmarker) calculateResults(totalDuration time.Duration) {
	b.results.TotalRequests = atomic.LoadInt64(&b.totalRequests)
	b.results.SuccessRequests = atomic.LoadInt64(&b.successCount)
	b.results.FailedRequests = atomic.LoadInt64(&b.failedCount)
	b.results.TotalDuration = totalDuration

	if len(b.latencies) == 0 {
		return
	}

	// 计算延迟统计
	var totalLatency time.Duration
	b.results.MinLatency = b.latencies[0]
	b.results.MaxLatency = b.latencies[0]

	for _, lat := range b.latencies {
		totalLatency += lat
		if lat < b.results.MinLatency {
			b.results.MinLatency = lat
		}
		if lat > b.results.MaxLatency {
			b.results.MaxLatency = lat
		}
	}

	b.results.AvgLatency = totalLatency / time.Duration(len(b.latencies))
	b.results.QPS = float64(b.results.TotalRequests) / totalDuration.Seconds()
}

// printResults 打印结果
func (b *Benchmarker) printResults() {
	fmt.Println("\n" + "="*70)
	fmt.Println("   压测结果")
	fmt.Println("="*70)

	fmt.Printf("总请求数:     %d\n", b.results.TotalRequests)
	fmt.Printf("成功请求数:   %d (%.2f%%)\n",
		b.results.SuccessRequests,
		float64(b.results.SuccessRequests)/float64(b.results.TotalRequests)*100)
	fmt.Printf("失败请求数:   %d (%.2f%%)\n",
		b.results.FailedRequests,
		float64(b.results.FailedRequests)/float64(b.results.TotalRequests)*100)

	fmt.Printf("\n总耗时:       %v\n", b.results.TotalDuration)
	fmt.Printf("QPS:          %.2f 请求/秒\n", b.results.QPS)

	fmt.Printf("\n延迟统计:\n")
	fmt.Printf("  最小:       %v\n", b.results.MinLatency)
	fmt.Printf("  最大:       %v\n", b.results.MaxLatency)
	fmt.Printf("  平均:       %v\n", b.results.AvgLatency)

	// 计算百分位延迟
	if len(b.latencies) > 0 {
		p50 := b.percentile(50)
		p90 := b.percentile(90)
		p95 := b.percentile(95)
		p99 := b.percentile(99)

		fmt.Printf("\n百分位延迟:\n")
		fmt.Printf("  P50:        %v\n", p50)
		fmt.Printf("  P90:        %v\n", p90)
		fmt.Printf("  P95:        %v\n", p95)
		fmt.Printf("  P99:        %v\n", p99)
	}

	fmt.Println("="*70 + "\n")
}

// percentile 计算百分位
func (b *Benchmarker) percentile(p float64) time.Duration {
	if len(b.latencies) == 0 {
		return 0
	}

	// 简单实现：排序并取对应位置
	// 生产环境应使用更高效的算法
	sorted := make([]time.Duration, len(b.latencies))
	copy(sorted, b.latencies)

	// 冒泡排序（简单实现）
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	index := int(float64(len(sorted)) * p / 100)
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	return sorted[index]
}

// 预定义的测试场景
func runBossHealthCheck() {
	config := BenchmarkConfig{
		URL:         "http://localhost:12306/login.html",
		Method:      "GET",
		Concurrency: 50,
		Duration:    10 * time.Second,
		Timeout:     5 * time.Second,
	}

	benchmarker := NewBenchmarker(config)
	benchmarker.Run()
}

func runGatewayStressTest() {
	config := BenchmarkConfig{
		URL:         "http://localhost:12309/",
		Method:      "GET",
		Concurrency: 100,
		Duration:    30 * time.Second,
		Timeout:     5 * time.Second,
	}

	benchmarker := NewBenchmarker(config)
	benchmarker.Run()
}

func runPaymentAPITest() {
	testData := map[string]interface{}{
		"merchant_id": "test",
		"amount":      100.00,
	}

	body, _ := json.Marshal(testData)

	config := BenchmarkConfig{
		URL:         "http://localhost:12309/gateway/scan",
		Method:      "POST",
		Concurrency: 50,
		Duration:    10 * time.Second,
		Timeout:     10 * time.Second,
		Body:        string(body),
	}

	benchmarker := NewBenchmarker(config)
	benchmarker.Run()
}

func main() {
	// 命令行参数
	url := flag.String("url", "", "测试URL")
	method := flag.String("method", "GET", "HTTP方法 (GET/POST)")
	concurrency := flag.Int("c", 50, "并发数")
	duration := flag.Int("d", 10, "持续时间（秒）")
	timeout := flag.Int("t", 5, "请求超时（秒）")
	body := flag.String("body", "{}", "请求体（POST时使用）")
	scenario := flag.String("scenario", "", "预定义场景: boss, gateway, payment")

	flag.Parse()

	// 如果指定了预定义场景
	if *scenario != "" {
		switch *scenario {
		case "boss":
			fmt.Println("运行 Boss 健康检查压测...")
			runBossHealthCheck()
		case "gateway":
			fmt.Println("运行 Gateway 压力测试...")
			runGatewayStressTest()
		case "payment":
			fmt.Println("运行支付接口压测...")
			runPaymentAPITest()
		default:
			fmt.Printf("未知场景: %s\n", *scenario)
			fmt.Println("可用场景: boss, gateway, payment")
			return
		}
		return
	}

	// 自定义压测
	if *url == "" {
		fmt.Println("用法:")
		fmt.Println("  ./benchmark -url <URL> -c <并发数> -d <持续时间>")
		fmt.Println("\n或使用预定义场景:")
		fmt.Println("  ./benchmark -scenario boss      # Boss服务压测")
		fmt.Println("  ./benchmark -scenario gateway   # Gateway压力测试")
		fmt.Println("  ./benchmark -scenario payment   # 支付接口压测")
		fmt.Println("\n示例:")
		fmt.Println("  ./benchmark -url http://localhost:12306/login.html -c 100 -d 30")
		return
	}

	config := BenchmarkConfig{
		URL:         *url,
		Method:      *method,
		Concurrency: *concurrency,
		Duration:    time.Duration(*duration) * time.Second,
		Timeout:     time.Duration(*timeout) * time.Second,
		Body:        *body,
	}

	benchmarker := NewBenchmarker(config)
	benchmarker.Run()
}
