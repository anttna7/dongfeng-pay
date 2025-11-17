package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// IntegrityChecker 数据完整性检查器
type IntegrityChecker struct {
	db *sql.DB
}

// NewIntegrityChecker 创建检查器
func NewIntegrityChecker(config DatabaseConfig) (*IntegrityChecker, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库ping失败: %v", err)
	}

	return &IntegrityChecker{db: db}, nil
}

// Close 关闭连接
func (ic *IntegrityChecker) Close() error {
	return ic.db.Close()
}

// CheckTableExists 检查表是否存在
func (ic *IntegrityChecker) CheckTableExists(tableName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = $1
		);
	`

	var exists bool
	err := ic.db.QueryRow(query, tableName).Scan(&exists)
	return exists, err
}

// GetTableRowCount 获取表行数
func (ic *IntegrityChecker) GetTableRowCount(tableName string) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)

	var count int64
	err := ic.db.QueryRow(query).Scan(&count)
	return count, err
}

// CheckAllTables 检查所有表
func (ic *IntegrityChecker) CheckAllTables() error {
	tables := []string{
		"user_info",
		"agent_info",
		"merchant_info",
		"merchant_deploy_info",
		"merchant_load_info",
		"account_info",
		"account_history_info",
		"order_info",
		"order_profit_info",
		"order_settle_info",
		"road_info",
		"road_pool_info",
		"payfor_info",
		"refund_info",
		"notify_info",
		"menu_info",
		"second_menu_info",
		"power_info",
		"role_info",
		"bank_card_info",
	}

	fmt.Println("\n" + "="*60)
	fmt.Println("   数据表完整性检查")
	fmt.Println("="*60)

	totalRows := int64(0)
	existsCount := 0

	for _, table := range tables {
		exists, err := ic.CheckTableExists(table)
		if err != nil {
			fmt.Printf("✗ %-30s 检查失败: %v\n", table, err)
			continue
		}

		if !exists {
			fmt.Printf("⚠ %-30s 不存在\n", table)
			continue
		}

		count, err := ic.GetTableRowCount(table)
		if err != nil {
			fmt.Printf("✗ %-30s 读取行数失败: %v\n", table, err)
			continue
		}

		existsCount++
		totalRows += count
		fmt.Printf("✓ %-30s 存在, 行数: %d\n", table, count)
	}

	fmt.Println("="*60)
	fmt.Printf("总计: %d/%d 个表存在, 总行数: %d\n", existsCount, len(tables), totalRows)
	fmt.Println("="*60 + "\n")

	return nil
}

// CheckOrderIntegrity 检查订单数据完整性
func (ic *IntegrityChecker) CheckOrderIntegrity() error {
	fmt.Println("\n" + "="*60)
	fmt.Println("   订单数据完整性检查")
	fmt.Println("="*60)

	// 1. 检查订单状态分布
	query := `
		SELECT status, COUNT(*) as count
		FROM order_info
		GROUP BY status
		ORDER BY count DESC
	`

	rows, err := ic.db.Query(query)
	if err != nil {
		return fmt.Errorf("查询订单状态失败: %v", err)
	}
	defer rows.Close()

	fmt.Println("\n订单状态分布:")
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			continue
		}
		fmt.Printf("  %-20s: %d\n", status, count)
	}

	// 2. 检查订单金额是否合理
	amountQuery := `
		SELECT
			MIN(order_amount) as min_amount,
			MAX(order_amount) as max_amount,
			AVG(order_amount) as avg_amount
		FROM order_info
		WHERE order_amount > 0
	`

	var minAmount, maxAmount, avgAmount float64
	err = ic.db.QueryRow(amountQuery).Scan(&minAmount, &maxAmount, &avgAmount)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("查询订单金额失败: %v", err)
	}

	fmt.Println("\n订单金额统计:")
	fmt.Printf("  最小金额: %.2f\n", minAmount)
	fmt.Printf("  最大金额: %.2f\n", maxAmount)
	fmt.Printf("  平均金额: %.2f\n", avgAmount)

	// 3. 检查是否有孤立订单（商户不存在）
	orphanQuery := `
		SELECT COUNT(*)
		FROM order_info o
		LEFT JOIN merchant_info m ON o.merchant_uid = m.merchant_uid
		WHERE m.merchant_uid IS NULL
	`

	var orphanCount int
	err = ic.db.QueryRow(orphanQuery).Scan(&orphanCount)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("查询孤立订单失败: %v", err)
	}

	if orphanCount > 0 {
		fmt.Printf("\n⚠ 发现 %d 个孤立订单（商户不存在）\n", orphanCount)
	} else {
		fmt.Println("\n✓ 未发现孤立订单")
	}

	fmt.Println("="*60 + "\n")
	return nil
}

// CheckAccountBalance 检查账户余额完整性
func (ic *IntegrityChecker) CheckAccountBalance() error {
	fmt.Println("\n" + "="*60)
	fmt.Println("   账户余额完整性检查")
	fmt.Println("="*60)

	// 检查是否有负余额账户
	query := `
		SELECT account_uid, account_name, balance
		FROM account_info
		WHERE balance < 0
	`

	rows, err := ic.db.Query(query)
	if err != nil {
		return fmt.Errorf("查询账户余额失败: %v", err)
	}
	defer rows.Close()

	negativeCount := 0
	fmt.Println("\n负余额账户:")
	for rows.Next() {
		var uid, name string
		var balance float64
		if err := rows.Scan(&uid, &name, &balance); err != nil {
			continue
		}
		negativeCount++
		fmt.Printf("  %-20s %-30s 余额: %.2f\n", uid, name, balance)
	}

	if negativeCount == 0 {
		fmt.Println("  ✓ 未发现负余额账户")
	} else {
		fmt.Printf("  ⚠ 发现 %d 个负余额账户\n", negativeCount)
	}

	// 统计总余额
	totalQuery := `SELECT SUM(balance) FROM account_info`
	var totalBalance float64
	err = ic.db.QueryRow(totalQuery).Scan(&totalBalance)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("查询总余额失败: %v", err)
	}

	fmt.Printf("\n系统总余额: %.2f\n", totalBalance)
	fmt.Println("="*60 + "\n")

	return nil
}

// CheckRefundIntegrity 检查退款数据完整性
func (ic *IntegrityChecker) CheckRefundIntegrity() error {
	fmt.Println("\n" + "="*60)
	fmt.Println("   退款数据完整性检查")
	fmt.Println("="*60)

	// 检查退款状态分布
	query := `
		SELECT status, COUNT(*) as count
		FROM refund_info
		GROUP BY status
		ORDER BY count DESC
	`

	rows, err := ic.db.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("\n暂无退款数据")
			fmt.Println("="*60 + "\n")
			return nil
		}
		return fmt.Errorf("查询退款状态失败: %v", err)
	}
	defer rows.Close()

	fmt.Println("\n退款状态分布:")
	hasData := false
	for rows.Next() {
		hasData = true
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			continue
		}
		fmt.Printf("  %-20s: %d\n", status, count)
	}

	if !hasData {
		fmt.Println("  暂无退款数据")
	}

	fmt.Println("="*60 + "\n")
	return nil
}

// GenerateReport 生成完整报告
func (ic *IntegrityChecker) GenerateReport() error {
	fmt.Println("\n" + "="*80)
	fmt.Println("   dongfeng-pay 数据完整性检查报告")
	fmt.Println("   生成时间: " + time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("="*80)

	if err := ic.CheckAllTables(); err != nil {
		return err
	}

	if err := ic.CheckOrderIntegrity(); err != nil {
		return err
	}

	if err := ic.CheckAccountBalance(); err != nil {
		return err
	}

	if err := ic.CheckRefundIntegrity(); err != nil {
		return err
	}

	fmt.Println("="*80)
	fmt.Println("   检查完成")
	fmt.Println("="*80 + "\n")

	return nil
}

func main() {
	// 从环境变量或配置文件读取数据库配置
	config := DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "Kyb^15273031604"),
		DBName:   getEnv("DB_NAME", "juhe_pay"),
	}

	checker, err := NewIntegrityChecker(config)
	if err != nil {
		log.Fatalf("创建检查器失败: %v", err)
	}
	defer checker.Close()

	if err := checker.GenerateReport(); err != nil {
		log.Fatalf("生成报告失败: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
