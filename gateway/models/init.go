/***************************************************
 ** @Desc : Database initialization for PostgreSQL
 ** @Time : 2019/8/9 13:48
 ** @Author : yuebin
 ** @File : init
 ** @Last Modified by : Claude AI
 ** @Last Modified time: 2025-11-17
 ** @Software: GoLand
 ** @Update : Upgraded to PostgreSQL 18
****************************************************/
package models

import (
	"fmt"
	"gateway/conf"
	"gateway/models/accounts"
	"gateway/models/agent"
	"gateway/models/merchant"
	"gateway/models/notify"
	"gateway/models/order"
	"gateway/models/payfor"
	"gateway/models/refund"
	"gateway/models/road"
	"gateway/models/system"
	"gateway/models/user"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/lib/pq"
)

func init() {
	dbHost := conf.DB_HOST
	dbUser := conf.DB_USER
	dbPassword := conf.DB_PASSWORD
	dbBase := conf.DB_BASE
	dbPort := conf.DB_PORT

	// PostgreSQL 连接字符串格式
	// postgresql://username:password@host:port/database?sslmode=disable
	link := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbBase)

	logs.Info("PostgreSQL init.....", link)

	// 注册 PostgreSQL 驱动
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", link)

	// 注册所有模型（包括新增的退款模型）
	orm.RegisterModel(
		new(user.UserInfo),
		new(system.MenuInfo),
		new(system.SecondMenuInfo),
		new(system.PowerInfo),
		new(system.RoleInfo),
		new(system.BankCardInfo),
		new(road.RoadInfo),
		new(road.RoadPoolInfo),
		new(agent.AgentInfo),
		new(merchant.MerchantInfo),
		new(merchant.MerchantDeployInfo),
		new(accounts.AccountInfo),
		new(accounts.AccountHistoryInfo),
		new(order.OrderInfo),
		new(order.OrderProfitInfo),
		new(order.OrderSettleInfo),
		new(notify.NotifyInfo),
		new(merchant.MerchantLoadInfo),
		new(payfor.PayforInfo),
		new(refund.RefundInfo),
	)

	// 开发环境自动同步表结构（生产环境建议注释掉）
	// orm.RunSyncdb("default", false, true)
}
