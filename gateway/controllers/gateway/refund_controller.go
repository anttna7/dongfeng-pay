package gateway

import (
	"encoding/json"
	"gateway/common"
	"gateway/service"
	"github.com/beego/beego/v2/server/web"
	"strconv"
)

// RefundController 退款控制器
type RefundController struct {
	web.Controller
}

// CreateRefund 创建退款
// @router /gateway/refund/create [post]
func (c *RefundController) CreateRefund() {
	logger := common.NewLogger().
		WithRequest(c.Ctx.Request.Method, c.Ctx.Request.URL.Path).
		WithIP(c.Ctx.Input.IP())

	logger.Info("收到退款请求")

	// 1. 解析请求参数
	var req service.RefundRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		logger.Error("解析请求参数失败", err)
		c.Data["json"] = common.NewErrorResponse(common.PARAMS_ERROR)
		c.ServeJSON()
		return
	}

	// 2. 参数验证
	if req.BankOrderId == "" && req.MerchantOrderId == "" {
		c.Data["json"] = common.NewErrorResponseWithMsg(common.PARAMS_ERROR, "订单号不能为空")
		c.ServeJSON()
		return
	}

	if req.RefundAmount <= 0 {
		c.Data["json"] = common.NewErrorResponseWithMsg(common.PARAMS_ERROR, "退款金额必须大于0")
		c.ServeJSON()
		return
	}

	if req.RefundReason == "" {
		c.Data["json"] = common.NewErrorResponseWithMsg(common.PARAMS_ERROR, "退款原因不能为空")
		c.ServeJSON()
		return
	}

	// 设置默认退款类型
	if req.RefundType == "" {
		req.RefundType = "full"
	}

	// 3. 调用退款服务
	refundService := &service.RefundService{}
	resp, err := refundService.CreateRefund(&req)

	if err != nil {
		if bizErr, ok := err.(*common.BusinessError); ok {
			c.Data["json"] = common.NewErrorResponse(bizErr.ErrorCode)
		} else {
			c.Data["json"] = common.NewErrorResponse(common.SYSTEM_ERROR)
		}
		c.ServeJSON()
		return
	}

	logger.Info("退款请求处理成功", map[string]interface{}{
		"refund_uid": resp.RefundUid,
	})

	// 4. 返回成功响应
	c.Data["json"] = common.NewSuccessResponse(resp)
	c.ServeJSON()
}

// QueryRefund 查询退款
// @router /gateway/refund/query [get]
func (c *RefundController) QueryRefund() {
	logger := common.NewLogger().
		WithRequest(c.Ctx.Request.Method, c.Ctx.Request.URL.Path).
		WithIP(c.Ctx.Input.IP())

	// 1. 获取退款单号
	refundUid := c.GetString("refund_uid")
	if refundUid == "" {
		c.Data["json"] = common.NewErrorResponseWithMsg(common.PARAMS_ERROR, "退款单号不能为空")
		c.ServeJSON()
		return
	}

	logger.WithExtra("refund_uid", refundUid).Info("查询退款")

	// 2. 查询退款信息
	refundService := &service.RefundService{}
	refundInfo, err := refundService.QueryRefund(refundUid)

	if err != nil {
		c.Data["json"] = common.NewErrorResponse(common.REFUND_NOT_EXIST)
		c.ServeJSON()
		return
	}

	// 3. 返回退款信息
	c.Data["json"] = common.NewSuccessResponse(refundInfo)
	c.ServeJSON()
}

// RefundList 退款列表
// @router /gateway/refund/list [get]
func (c *RefundController) RefundList() {
	logger := common.NewLogger().
		WithRequest(c.Ctx.Request.Method, c.Ctx.Request.URL.Path).
		WithIP(c.Ctx.Input.IP())

	// 1. 获取分页参数
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 20)
	merchantUid := c.GetString("merchant_uid")
	status := c.GetString("status")

	logger.Info("查询退款列表", map[string]interface{}{
		"page":         page,
		"page_size":    pageSize,
		"merchant_uid": merchantUid,
		"status":       status,
	})

	// 2. 查询退款列表
	refundService := &service.RefundService{}
	refunds, total, err := refundService.GetRefundList(page, pageSize, merchantUid, status)

	if err != nil {
		logger.Error("查询退款列表失败", err)
		c.Data["json"] = common.NewErrorResponse(common.SYSTEM_ERROR)
		c.ServeJSON()
		return
	}

	// 3. 返回分页数据
	c.Data["json"] = common.NewPageResponse(refunds, total, page, pageSize)
	c.ServeJSON()
}

// RefundCallback 退款回调（供上游调用）
// @router /gateway/refund/callback [post]
func (c *RefundController) RefundCallback() {
	logger := common.NewLogger().
		WithRequest(c.Ctx.Request.Method, c.Ctx.Request.URL.Path).
		WithIP(c.Ctx.Input.IP())

	logger.Info("收到退款回调", map[string]interface{}{
		"body": string(c.Ctx.Input.RequestBody),
	})

	// TODO: 根据实际的上游通道实现回调处理逻辑
	// 1. 解析回调参数
	// 2. 验证签名
	// 3. 更新退款状态
	// 4. 返回响应给上游

	c.Ctx.WriteString("success")
}

// BatchRefund 批量退款
// @router /gateway/refund/batch [post]
func (c *RefundController) BatchRefund() {
	logger := common.NewLogger().
		WithRequest(c.Ctx.Request.Method, c.Ctx.Request.URL.Path).
		WithIP(c.Ctx.Input.IP())

	logger.Info("收到批量退款请求")

	// 1. 解析请求参数
	var requests []service.RefundRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &requests); err != nil {
		logger.Error("解析请求参数失败", err)
		c.Data["json"] = common.NewErrorResponse(common.PARAMS_ERROR)
		c.ServeJSON()
		return
	}

	if len(requests) == 0 {
		c.Data["json"] = common.NewErrorResponseWithMsg(common.PARAMS_ERROR, "退款列表不能为空")
		c.ServeJSON()
		return
	}

	if len(requests) > 100 {
		c.Data["json"] = common.NewErrorResponseWithMsg(common.PARAMS_ERROR, "单次最多支持100笔退款")
		c.ServeJSON()
		return
	}

	// 2. 批量处理退款
	refundService := &service.RefundService{}
	results := make([]map[string]interface{}, 0, len(requests))

	successCount := 0
	failCount := 0

	for i, req := range requests {
		result := map[string]interface{}{
			"index":      i,
			"order_id":   req.BankOrderId,
			"merchant_order_id": req.MerchantOrderId,
		}

		resp, err := refundService.CreateRefund(&req)
		if err != nil {
			result["success"] = false
			result["message"] = err.Error()
			failCount++
		} else {
			result["success"] = true
			result["refund_uid"] = resp.RefundUid
			result["message"] = resp.Message
			successCount++
		}

		results = append(results, result)
	}

	logger.Info("批量退款处理完成", map[string]interface{}{
		"total":   len(requests),
		"success": successCount,
		"fail":    failCount,
	})

	// 3. 返回结果
	c.Data["json"] = common.NewSuccessResponse(map[string]interface{}{
		"total":        len(requests),
		"success_count": successCount,
		"fail_count":    failCount,
		"results":       results,
	})
	c.ServeJSON()
}

// RefundStatistics 退款统计
// @router /gateway/refund/statistics [get]
func (c *RefundController) RefundStatistics() {
	logger := common.NewLogger().
		WithRequest(c.Ctx.Request.Method, c.Ctx.Request.URL.Path).
		WithIP(c.Ctx.Input.IP())

	merchantUid := c.GetString("merchant_uid")
	startDate := c.GetString("start_date")
	endDate := c.GetString("end_date")

	logger.Info("查询退款统计", map[string]interface{}{
		"merchant_uid": merchantUid,
		"start_date":   startDate,
		"end_date":     endDate,
	})

	// TODO: 实现退款统计逻辑
	// 1. 统计总退款金额
	// 2. 统计退款笔数
	// 3. 统计各状态退款数量
	// 4. 统计退款成功率

	statistics := map[string]interface{}{
		"total_amount":    0,
		"total_count":     0,
		"success_count":   0,
		"failed_count":    0,
		"processing_count": 0,
		"success_rate":    "0%",
	}

	c.Data["json"] = common.NewSuccessResponse(statistics)
	c.ServeJSON()
}
