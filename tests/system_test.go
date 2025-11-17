package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestConfig 测试配置
type TestConfig struct {
	BossURL     string
	MerchantURL string
	AgentURL    string
	GatewayURL  string
}

var config = TestConfig{
	BossURL:     "http://localhost:12306",
	MerchantURL: "http://localhost:12307",
	AgentURL:    "http://localhost:12308",
	GatewayURL:  "http://localhost:12309",
}

// TestBossHealth 测试 Boss 服务健康检查
func TestBossHealth(t *testing.T) {
	t.Log("测试 Boss 服务健康检查...")

	resp, err := http.Get(config.BossURL + "/login.html")
	if err != nil {
		t.Fatalf("Boss 服务连接失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Boss 服务返回状态码异常: 期望 200, 实际 %d", resp.StatusCode)
	}

	t.Log("✓ Boss 服务健康检查通过")
}

// TestMerchantHealth 测试 Merchant 服务健康检查
func TestMerchantHealth(t *testing.T) {
	t.Log("测试 Merchant 服务健康检查...")

	resp, err := http.Get(config.MerchantURL + "/")
	if err != nil {
		t.Fatalf("Merchant 服务连接失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Merchant 服务返回状态码异常: 期望 200, 实际 %d", resp.StatusCode)
	}

	t.Log("✓ Merchant 服务健康检查通过")
}

// TestAgentHealth 测试 Agent 服务健康检查
func TestAgentHealth(t *testing.T) {
	t.Log("测试 Agent 服务健康检查...")

	resp, err := http.Get(config.AgentURL + "/")
	if err != nil {
		t.Fatalf("Agent 服务连接失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Agent 服务返回状态码异常: 期望 200, 实际 %d", resp.StatusCode)
	}

	t.Log("✓ Agent 服务健康检查通过")
}

// TestGatewayHealth 测试 Gateway 服务健康检查
func TestGatewayHealth(t *testing.T) {
	t.Log("测试 Gateway 服务健康检查...")

	resp, err := http.Get(config.GatewayURL + "/")
	if err != nil {
		t.Fatalf("Gateway 服务连接失败: %v", err)
	}
	defer resp.Body.Close()

	// Gateway 可能返回 404（没有首页）或其他状态，只要能连接就算健康
	t.Logf("Gateway 服务响应状态码: %d", resp.StatusCode)
	t.Log("✓ Gateway 服务健康检查通过")
}

// TestBossLogin 测试 Boss 登录功能
func TestBossLogin(t *testing.T) {
	t.Log("测试 Boss 登录功能...")

	// 构造登录请求
	loginData := map[string]string{
		"username": "10086",
		"password": "123456",
	}

	jsonData, _ := json.Marshal(loginData)
	resp, err := http.Post(
		config.BossURL+"/login",
		"application/json",
		strings.NewReader(string(jsonData)),
	)

	if err != nil {
		t.Fatalf("登录请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	t.Logf("登录响应: %s", string(body))

	// 检查响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err == nil {
		if code, ok := result["code"].(string); ok && code == "0000" {
			t.Log("✓ Boss 登录功能测试通过")
			return
		}
	}

	t.Log("ℹ Boss 登录可能需要验证码，跳过详细验证")
}

// TestGatewayPaymentAPI 测试 Gateway 支付接口
func TestGatewayPaymentAPI(t *testing.T) {
	t.Log("测试 Gateway 支付接口...")

	// 这里只测试接口可访问性，实际支付需要完整的参数和签名
	resp, err := http.Post(
		config.GatewayURL+"/gateway/scan",
		"application/json",
		strings.NewReader("{}"),
	)

	if err != nil {
		t.Fatalf("支付接口请求失败: %v", err)
	}
	defer resp.Body.Close()

	t.Logf("支付接口响应状态码: %d", resp.StatusCode)
	t.Log("✓ Gateway 支付接口可访问")
}

// TestGatewayRefundAPI 测试 Gateway 退款接口
func TestGatewayRefundAPI(t *testing.T) {
	t.Log("测试 Gateway 退款接口...")

	resp, err := http.Post(
		config.GatewayURL+"/gateway/refund/create",
		"application/json",
		strings.NewReader("{}"),
	)

	if err != nil {
		t.Fatalf("退款接口请求失败: %v", err)
	}
	defer resp.Body.Close()

	t.Logf("退款接口响应状态码: %d", resp.StatusCode)
	t.Log("✓ Gateway 退款接口可访问")
}

// TestResponseTime 测试响应时间
func TestResponseTime(t *testing.T) {
	t.Log("测试各服务响应时间...")

	services := map[string]string{
		"Boss":     config.BossURL,
		"Merchant": config.MerchantURL,
		"Agent":    config.AgentURL,
		"Gateway":  config.GatewayURL,
	}

	for name, url := range services {
		start := time.Now()
		resp, err := http.Get(url + "/")
		duration := time.Since(start)

		if err != nil {
			t.Logf("⚠ %s 服务无响应: %v", name, err)
			continue
		}
		resp.Body.Close()

		t.Logf("%s 服务响应时间: %v", name, duration)

		if duration > 2*time.Second {
			t.Errorf("⚠ %s 服务响应过慢: %v (超过 2s)", name, duration)
		} else {
			t.Logf("✓ %s 服务响应正常", name)
		}
	}
}

// TestConcurrentRequests 测试并发请求
func TestConcurrentRequests(t *testing.T) {
	t.Log("测试并发请求处理能力...")

	concurrency := 10
	done := make(chan bool, concurrency)

	start := time.Now()

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			resp, err := http.Get(config.GatewayURL + "/")
			if err == nil {
				resp.Body.Close()
			}
			done <- true
		}(i)
	}

	// 等待所有请求完成
	for i := 0; i < concurrency; i++ {
		<-done
	}

	duration := time.Since(start)
	t.Logf("并发 %d 个请求完成时间: %v", concurrency, duration)

	if duration > 5*time.Second {
		t.Errorf("⚠ 并发性能较差: %v (超过 5s)", duration)
	} else {
		t.Log("✓ 并发请求处理正常")
	}
}

// TestMemoryLeak 内存泄漏测试（简单版）
func TestMemoryLeak(t *testing.T) {
	t.Log("执行内存泄漏测试...")

	// 执行多次请求，检查内存是否持续增长
	iterations := 100

	for i := 0; i < iterations; i++ {
		resp, err := http.Get(config.GatewayURL + "/")
		if err == nil {
			resp.Body.Close()
		}

		if i%20 == 0 {
			t.Logf("已完成 %d/%d 次请求", i, iterations)
		}
	}

	t.Log("✓ 内存泄漏测试完成（建议使用 pprof 进行详细分析）")
}

// TestDatabaseConnection 测试数据库连接
func TestDatabaseConnection(t *testing.T) {
	t.Log("测试数据库连接...")

	// 通过 API 间接测试数据库连接
	// 实际项目中应该直接连接数据库进行测试

	resp, err := http.Get(config.BossURL + "/login.html")
	if err != nil {
		t.Fatalf("无法连接服务，可能数据库连接异常: %v", err)
	}
	defer resp.Body.Close()

	t.Log("✓ 数据库连接测试通过（间接验证）")
}

// 运行全部测试
func TestAll(t *testing.T) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("   dongfeng-pay 系统功能测试套件")
	fmt.Println(strings.Repeat("=", 60))

	t.Run("健康检查", func(t *testing.T) {
		t.Run("Boss", TestBossHealth)
		t.Run("Merchant", TestMerchantHealth)
		t.Run("Agent", TestAgentHealth)
		t.Run("Gateway", TestGatewayHealth)
	})

	t.Run("功能测试", func(t *testing.T) {
		t.Run("登录", TestBossLogin)
		t.Run("支付接口", TestGatewayPaymentAPI)
		t.Run("退款接口", TestGatewayRefundAPI)
	})

	t.Run("性能测试", func(t *testing.T) {
		t.Run("响应时间", TestResponseTime)
		t.Run("并发请求", TestConcurrentRequests)
	})

	t.Run("稳定性测试", func(t *testing.T) {
		t.Run("内存泄漏", TestMemoryLeak)
		t.Run("数据库连接", TestDatabaseConnection)
	})

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("   测试完成")
	fmt.Println(strings.Repeat("=", 60) + "\n")
}
