package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

/***************************************************
 ** @Desc : 性能分析工具
 ** @Time : 2025-11-17
 ** @Author : Claude AI
 ** @File : performance_analyzer.go
 ****************************************************/

// PerformanceMetrics 性能指标
type PerformanceMetrics struct {
	Timestamp       time.Time              `json:"timestamp"`
	Services        map[string]ServiceInfo `json:"services"`
	OverallHealth   string                 `json:"overall_health"`
	Recommendations []string               `json:"recommendations"`
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	Name         string        `json:"name"`
	Port         int           `json:"port"`
	Status       string        `json:"status"`
	ResponseTime time.Duration `json:"response_time_ms"`
	HealthCheck  bool          `json:"health_check"`
	Endpoints    []EndpointInfo `json:"endpoints"`
}

// EndpointInfo 端点信息
type EndpointInfo struct {
	Path         string        `json:"path"`
	Method       string        `json:"method"`
	ResponseTime time.Duration `json:"response_time_ms"`
	StatusCode   int           `json:"status_code"`
	Success      bool          `json:"success"`
}

// PerformanceAnalyzer 性能分析器
type PerformanceAnalyzer struct {
	services []ServiceConfig
	results  PerformanceMetrics
}

// ServiceConfig 服务配置
type ServiceConfig struct {
	Name      string
	BaseURL   string
	Port      int
	Endpoints []string
}

// NewPerformanceAnalyzer 创建性能分析器
func NewPerformanceAnalyzer() *PerformanceAnalyzer {
	return &PerformanceAnalyzer{
		services: []ServiceConfig{
			{
				Name:    "Boss",
				BaseURL: "http://localhost:12306",
				Port:    12306,
				Endpoints: []string{
					"/login.html",
					"/login_modern.html",
					"/getVerifyImg",
				},
			},
			{
				Name:    "Merchant",
				BaseURL: "http://localhost:12307",
				Port:    12307,
				Endpoints: []string{
					"/login_modern.html",
					"/img.do/test.png",
				},
			},
			{
				Name:    "Agent",
				BaseURL: "http://localhost:12308",
				Port:    12308,
				Endpoints: []string{
					"/login_modern.html",
					"/img.do/test.png",
				},
			},
			{
				Name:    "Gateway",
				BaseURL: "http://localhost:12309",
				Port:    12309,
				Endpoints: []string{
					"/",
					"/health",
					"/ping",
				},
			},
		},
		results: PerformanceMetrics{
			Timestamp: time.Now(),
			Services:  make(map[string]ServiceInfo),
		},
	}
}

// Run 运行性能分析
func (pa *PerformanceAnalyzer) Run() error {
	fmt.Println("🔍 开始性能分析...")
	fmt.Println("==========================================")

	for _, service := range pa.services {
		fmt.Printf("\n📊 分析服务: %s (端口 %d)\n", service.Name, service.Port)

		serviceInfo := ServiceInfo{
			Name:      service.Name,
			Port:      service.Port,
			Endpoints: []EndpointInfo{},
		}

		// 检查服务状态
		start := time.Now()
		available, statusCode := pa.checkServiceAvailability(service.BaseURL)
		responseTime := time.Since(start)

		serviceInfo.ResponseTime = responseTime
		serviceInfo.Status = "offline"
		serviceInfo.HealthCheck = false

		if available {
			serviceInfo.Status = "online"
			serviceInfo.HealthCheck = true
			fmt.Printf("  ✅ 服务在线 (响应时间: %v)\n", responseTime)

			// 测试各个端点
			for _, endpoint := range service.Endpoints {
				endpointInfo := pa.testEndpoint(service.BaseURL, endpoint)
				serviceInfo.Endpoints = append(serviceInfo.Endpoints, endpointInfo)

				if endpointInfo.Success {
					fmt.Printf("  ✅ %s - %v - %d\n", endpointInfo.Path, endpointInfo.ResponseTime, endpointInfo.StatusCode)
				} else {
					fmt.Printf("  ❌ %s - %v - %d\n", endpointInfo.Path, endpointInfo.ResponseTime, endpointInfo.StatusCode)
				}
			}
		} else {
			fmt.Printf("  ❌ 服务离线 (状态码: %d)\n", statusCode)
		}

		pa.results.Services[service.Name] = serviceInfo
	}

	// 生成健康评估
	pa.generateHealthAssessment()

	// 生成建议
	pa.generateRecommendations()

	return nil
}

// checkServiceAvailability 检查服务可用性
func (pa *PerformanceAnalyzer) checkServiceAvailability(baseURL string) (bool, int) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(baseURL)
	if err != nil {
		return false, 0
	}
	defer resp.Body.Close()

	return resp.StatusCode < 500, resp.StatusCode
}

// testEndpoint 测试端点
func (pa *PerformanceAnalyzer) testEndpoint(baseURL, endpoint string) EndpointInfo {
	url := baseURL + endpoint

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	start := time.Now()
	resp, err := client.Get(url)
	responseTime := time.Since(start)

	info := EndpointInfo{
		Path:         endpoint,
		Method:       "GET",
		ResponseTime: responseTime,
	}

	if err != nil {
		info.Success = false
		info.StatusCode = 0
		return info
	}
	defer resp.Body.Close()

	info.StatusCode = resp.StatusCode
	info.Success = resp.StatusCode >= 200 && resp.StatusCode < 400

	return info
}

// generateHealthAssessment 生成健康评估
func (pa *PerformanceAnalyzer) generateHealthAssessment() {
	onlineCount := 0
	totalServices := len(pa.results.Services)

	for _, service := range pa.results.Services {
		if service.Status == "online" {
			onlineCount++
		}
	}

	if onlineCount == totalServices {
		pa.results.OverallHealth = "excellent"
	} else if onlineCount >= totalServices/2 {
		pa.results.OverallHealth = "good"
	} else if onlineCount > 0 {
		pa.results.OverallHealth = "poor"
	} else {
		pa.results.OverallHealth = "critical"
	}
}

// generateRecommendations 生成建议
func (pa *PerformanceAnalyzer) generateRecommendations() {
	pa.results.Recommendations = []string{}

	// 检查离线服务
	for name, service := range pa.results.Services {
		if service.Status == "offline" {
			pa.results.Recommendations = append(pa.results.Recommendations,
				fmt.Sprintf("⚠️ %s 服务离线，请检查服务是否启动", name))
		}
	}

	// 检查响应时间
	for name, service := range pa.results.Services {
		if service.Status == "online" && service.ResponseTime > 1*time.Second {
			pa.results.Recommendations = append(pa.results.Recommendations,
				fmt.Sprintf("🐌 %s 服务响应较慢 (%v)，建议优化", name, service.ResponseTime))
		}
	}

	// 检查端点失败
	for serviceName, service := range pa.results.Services {
		failedCount := 0
		for _, endpoint := range service.Endpoints {
			if !endpoint.Success {
				failedCount++
			}
		}
		if failedCount > 0 {
			pa.results.Recommendations = append(pa.results.Recommendations,
				fmt.Sprintf("❌ %s 服务有 %d 个端点失败，请检查", serviceName, failedCount))
		}
	}

	// 如果一切正常
	if len(pa.results.Recommendations) == 0 {
		pa.results.Recommendations = append(pa.results.Recommendations,
			"✅ 所有服务运行正常，性能良好")
	}
}

// PrintReport 打印报告
func (pa *PerformanceAnalyzer) PrintReport() {
	fmt.Println("\n==========================================")
	fmt.Println("📊 性能分析报告")
	fmt.Println("==========================================")
	fmt.Printf("分析时间: %s\n", pa.results.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("整体健康: %s\n\n", pa.getHealthEmoji(pa.results.OverallHealth))

	fmt.Println("服务状态汇总:")
	fmt.Println("------------------------------------------")

	for _, service := range pa.results.Services {
		status := "🔴"
		if service.Status == "online" {
			status = "🟢"
		}

		fmt.Printf("%s %s (端口 %d)\n", status, service.Name, service.Port)
		fmt.Printf("   响应时间: %v\n", service.ResponseTime)
		fmt.Printf("   端点数: %d\n", len(service.Endpoints))

		successCount := 0
		for _, endpoint := range service.Endpoints {
			if endpoint.Success {
				successCount++
			}
		}
		fmt.Printf("   成功率: %d/%d (%.1f%%)\n",
			successCount, len(service.Endpoints),
			float64(successCount)/float64(len(service.Endpoints))*100)
		fmt.Println()
	}

	fmt.Println("建议:")
	fmt.Println("------------------------------------------")
	for _, rec := range pa.results.Recommendations {
		fmt.Printf("  %s\n", rec)
	}
	fmt.Println()
}

// getHealthEmoji 获取健康状态表情
func (pa *PerformanceAnalyzer) getHealthEmoji(health string) string {
	switch health {
	case "excellent":
		return "🟢 优秀"
	case "good":
		return "🟡 良好"
	case "poor":
		return "🟠 较差"
	case "critical":
		return "🔴 严重"
	default:
		return "⚪ 未知"
	}
}

// SaveReport 保存报告
func (pa *PerformanceAnalyzer) SaveReport(filename string) error {
	data, err := json.MarshalIndent(pa.results, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化报告失败: %v", err)
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("保存报告失败: %v", err)
	}

	fmt.Printf("✅ 报告已保存到: %s\n", filename)
	return nil
}

// SaveMarkdownReport 保存 Markdown 格式报告
func (pa *PerformanceAnalyzer) SaveMarkdownReport(filename string) error {
	md := fmt.Sprintf("# 性能分析报告\n\n")
	md += fmt.Sprintf("**生成时间**: %s\n\n", pa.results.Timestamp.Format("2006-01-02 15:04:05"))
	md += fmt.Sprintf("**整体健康**: %s\n\n", pa.getHealthEmoji(pa.results.OverallHealth))

	md += "## 服务状态\n\n"
	md += "| 服务 | 状态 | 端口 | 响应时间 | 成功率 |\n"
	md += "|------|------|------|----------|--------|\n"

	for _, service := range pa.results.Services {
		status := "🔴 离线"
		if service.Status == "online" {
			status = "🟢 在线"
		}

		successCount := 0
		totalEndpoints := len(service.Endpoints)
		for _, endpoint := range service.Endpoints {
			if endpoint.Success {
				successCount++
			}
		}

		successRate := 0.0
		if totalEndpoints > 0 {
			successRate = float64(successCount) / float64(totalEndpoints) * 100
		}

		md += fmt.Sprintf("| %s | %s | %d | %v | %.1f%% |\n",
			service.Name, status, service.Port, service.ResponseTime, successRate)
	}

	md += "\n## 端点详情\n\n"

	for serviceName, service := range pa.results.Services {
		if len(service.Endpoints) == 0 {
			continue
		}

		md += fmt.Sprintf("### %s\n\n", serviceName)
		md += "| 端点 | 状态码 | 响应时间 | 结果 |\n"
		md += "|------|--------|----------|------|\n"

		for _, endpoint := range service.Endpoints {
			result := "✅ 成功"
			if !endpoint.Success {
				result = "❌ 失败"
			}

			md += fmt.Sprintf("| %s | %d | %v | %s |\n",
				endpoint.Path, endpoint.StatusCode, endpoint.ResponseTime, result)
		}

		md += "\n"
	}

	md += "## 优化建议\n\n"

	for _, rec := range pa.results.Recommendations {
		md += fmt.Sprintf("- %s\n", rec)
	}

	md += "\n---\n\n"
	md += "*报告由 dongfeng-pay 性能分析工具自动生成*\n"

	err := ioutil.WriteFile(filename, []byte(md), 0644)
	if err != nil {
		return fmt.Errorf("保存 Markdown 报告失败: %v", err)
	}

	fmt.Printf("✅ Markdown 报告已保存到: %s\n", filename)
	return nil
}

func main() {
	// 创建性能分析器
	analyzer := NewPerformanceAnalyzer()

	// 运行分析
	err := analyzer.Run()
	if err != nil {
		fmt.Printf("❌ 分析失败: %v\n", err)
		os.Exit(1)
	}

	// 打印报告
	analyzer.PrintReport()

	// 保存报告
	timestamp := time.Now().Format("20060102_150405")

	// 保存 JSON 报告
	jsonFile := fmt.Sprintf("./test-reports/performance_%s.json", timestamp)
	err = analyzer.SaveReport(jsonFile)
	if err != nil {
		fmt.Printf("⚠️ 保存 JSON 报告失败: %v\n", err)
	}

	// 保存 Markdown 报告
	mdFile := fmt.Sprintf("./test-reports/performance_%s.md", timestamp)
	err = analyzer.SaveMarkdownReport(mdFile)
	if err != nil {
		fmt.Printf("⚠️ 保存 Markdown 报告失败: %v\n", err)
	}

	fmt.Println("\n✅ 性能分析完成")
}
