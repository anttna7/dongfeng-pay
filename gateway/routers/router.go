package routers

import (
	"gateway/controllers/gateway"
	"gateway/middleware"
	"gateway/supplier/third_party"
	"github.com/beego/beego/v2/server/web"
)

func init() {
	// 初始化限流器
	middleware.InitRateLimiters()

	// 应用全局限流中间件到所有请求
	web.InsertFilter("/*", web.BeforeRouter, middleware.GlobalRateLimitFilter)

	// 应用商户级别限流到网关相关接口
	web.InsertFilter("/gateway/*", web.BeforeRouter, middleware.MerchantRateLimitFilter)

	// 应用支付级别限流到支付和代付接口（最严格）
	web.InsertFilter("/gateway/scan", web.BeforeRouter, middleware.PaymentRateLimitFilter)
	web.InsertFilter("/gateway/payfor", web.BeforeRouter, middleware.PaymentRateLimitFilter)

	//网关处理函数
	web.Router("/gateway/scan", &gateway.ScanController{}, "*:Scan")
	web.Router("/err/params", &gateway.ErrorGatewayController{}, "*:ErrorParams")
	//代付相关的接口
	web.Router("/gateway/payfor", &gateway.PayForGateway{}, "*:PayFor")
	web.Router("/gateway/payfor/query", &gateway.PayForGateway{}, "*:PayForQuery")
	web.Router("/gateway/balance", &gateway.PayForGateway{}, "*:Balance")
	// 接收银行回调
	web.Router("/daili/notify", &third_party.DaiLiImpl{}, "*:PayNotify")

	web.Router("/gateway/supplier/order/query", &gateway.OrderController{}, "*:OrderQuery")
	web.Router("/gateway/update/order", &gateway.OrderController{}, "*:OrderUpdate")
	web.Router("/gateway/supplier/payfor/query", &gateway.PayForGateway{}, "*:QuerySupplierPayForResult")
	web.Router("/solve/payfor/result", &gateway.PayForGateway{}, "*:SolvePayForResult")

	// 退款相关接口
	web.Router("/gateway/refund/create", &gateway.RefundController{}, "post:CreateRefund")
	web.Router("/gateway/refund/query", &gateway.RefundController{}, "get:QueryRefund")
	web.Router("/gateway/refund/list", &gateway.RefundController{}, "get:RefundList")
	web.Router("/gateway/refund/callback", &gateway.RefundController{}, "post:RefundCallback")
	web.Router("/gateway/refund/batch", &gateway.RefundController{}, "post:BatchRefund")
	web.Router("/gateway/refund/statistics", &gateway.RefundController{}, "get:RefundStatistics")
}
