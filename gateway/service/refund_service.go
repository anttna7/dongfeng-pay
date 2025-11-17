package service

import (
	"fmt"
	"gateway/common"
	"gateway/models/accounts"
	"gateway/models/order"
	"gateway/models/refund"
	"github.com/beego/beego/v2/core/logs"
	"github.com/rs/xid"
	"time"
)

// RefundService 退款服务
type RefundService struct{}

// RefundRequest 退款请求
type RefundRequest struct {
	MerchantOrderId string  // 商户订单号
	BankOrderId     string  // 系统订单号
	RefundAmount    float64 // 退款金额
	RefundReason    string  // 退款原因
	RefundType      string  // 退款类型：full-全额，partial-部分
	OperatorName    string  // 操作人
	NotifyUrl       string  // 退款回调地址
	MerchantUid     string  // 商户UID
}

// RefundResponse 退款响应
type RefundResponse struct {
	RefundUid    string  // 退款单号
	RefundAmount float64 // 退款金额
	Status       string  // 退款状态
	Message      string  // 退款消息
}

// CreateRefund 创建退款
func (rs *RefundService) CreateRefund(req *RefundRequest) (*RefundResponse, error) {
	logger := common.NewLogger().
		WithMerchant(req.MerchantUid).
		WithOrder(req.BankOrderId)

	logger.Info("开始处理退款请求", map[string]interface{}{
		"merchant_order_id": req.MerchantOrderId,
		"refund_amount":     req.RefundAmount,
		"refund_type":       req.RefundType,
	})

	// 1. 查询订单信息
	var orderInfo *order.OrderInfo
	var err error

	if req.BankOrderId != "" {
		orderInfo, err = order.GetOrderByOrderId(req.BankOrderId)
	} else if req.MerchantOrderId != "" {
		orderInfo, err = order.GetOrderByMerchantOrderId(req.MerchantOrderId)
	} else {
		return nil, common.NewBusinessError(common.PARAMS_ERROR)
	}

	if err != nil {
		logger.Error("查询订单失败", err)
		return nil, common.NewBusinessError(common.ORDER_NOT_EXIST)
	}

	// 2. 验证订单状态
	if orderInfo.Status != "success" {
		return nil, common.NewBusinessErrorWithDetail(common.ORDER_STATUS_ERROR, "只能退款成功的订单")
	}

	// 3. 验证退款金额
	if req.RefundAmount <= 0 {
		return nil, common.NewBusinessError(common.REFUND_AMOUNT_ERROR)
	}

	if req.RefundAmount > orderInfo.FactAmount {
		return nil, common.NewBusinessErrorWithDetail(common.REFUND_AMOUNT_ERROR,
			fmt.Sprintf("退款金额不能大于订单金额，订单金额: %.2f", orderInfo.FactAmount))
	}

	// 4. 检查是否已经退款
	existingRefunds, err := refund.GetRefundByOrderId(orderInfo.BankOrderId)
	if err == nil && len(existingRefunds) > 0 {
		var totalRefunded float64
		for _, r := range existingRefunds {
			if r.Status == "success" || r.Status == "processing" {
				totalRefunded += r.RefundAmount
			}
		}

		if totalRefunded+req.RefundAmount > orderInfo.FactAmount {
			return nil, common.NewBusinessErrorWithDetail(common.REFUND_AMOUNT_ERROR,
				fmt.Sprintf("已退款金额: %.2f，本次退款: %.2f，超过订单金额: %.2f",
					totalRefunded, req.RefundAmount, orderInfo.FactAmount))
		}
	}

	// 5. 创建退款记录
	refundInfo := &refund.RefundInfo{
		RefundUid:       xid.New().String(),
		MerchantOrderId: orderInfo.MerchantOrderId,
		BankOrderId:     orderInfo.BankOrderId,
		BankTransId:     orderInfo.BankTransId,
		RefundAmount:    req.RefundAmount,
		RefundReason:    req.RefundReason,
		RefundType:      req.RefundType,
		Status:          "pending",
		MerchantUid:     orderInfo.MerchantUid,
		MerchantName:    orderInfo.MerchantName,
		AgentUid:        orderInfo.AgentUid,
		RoadUid:         orderInfo.RoadUid,
		OperatorName:    req.OperatorName,
		NotifyUrl:       req.NotifyUrl,
		NotifyStatus:    "pending",
		NotifyCount:     0,
	}

	if err := refundInfo.Insert(); err != nil {
		logger.Error("创建退款记录失败", err)
		return nil, common.NewBusinessError(common.SYSTEM_ERROR)
	}

	logger.Info("退款记录创建成功", map[string]interface{}{
		"refund_uid": refundInfo.RefundUid,
	})

	// 6. 异步处理退款
	go rs.ProcessRefund(refundInfo)

	return &RefundResponse{
		RefundUid:    refundInfo.RefundUid,
		RefundAmount: req.RefundAmount,
		Status:       "pending",
		Message:      "退款申请已提交，正在处理中",
	}, nil
}

// ProcessRefund 处理退款
func (rs *RefundService) ProcessRefund(refundInfo *refund.RefundInfo) {
	logger := common.NewLogger().
		WithMerchant(refundInfo.MerchantUid).
		WithOrder(refundInfo.BankOrderId).
		WithExtra("refund_uid", refundInfo.RefundUid)

	logger.Info("开始处理退款")

	// 更新状态为处理中
	refundInfo.Status = "processing"
	refundInfo.Update("Status", "UpdateTime")

	// 1. 调用上游退款接口（这里需要根据实际的支付通道实现）
	success, resultMsg := rs.callUpstreamRefund(refundInfo)

	if success {
		// 2. 退款成功，更新账户余额
		if err := rs.updateAccountBalance(refundInfo); err != nil {
			logger.Error("更新账户余额失败", err)
			// 即使余额更新失败，退款仍然算成功（上游已经退款）
			// 需要人工介入处理
		}

		// 3. 更新订单退款状态
		if err := rs.updateOrderRefundStatus(refundInfo); err != nil {
			logger.Error("更新订单退款状态失败", err)
		}

		// 4. 更新退款状态为成功
		refund.UpdateRefundStatus(refundInfo.RefundUid, "success", resultMsg)

		logger.Info("退款处理成功", map[string]interface{}{
			"refund_amount": refundInfo.RefundAmount,
		})

		// 5. 发送退款通知
		go rs.sendRefundNotify(refundInfo)

	} else {
		// 退款失败
		refund.UpdateRefundStatus(refundInfo.RefundUid, "failed", resultMsg)

		logger.Error("退款处理失败", nil, map[string]interface{}{
			"result": resultMsg,
		})
	}
}

// callUpstreamRefund 调用上游退款接口
func (rs *RefundService) callUpstreamRefund(refundInfo *refund.RefundInfo) (bool, string) {
	// TODO: 根据实际的支付通道实现退款逻辑
	// 这里提供一个模拟实现

	logs.Info("调用上游退款接口: refund_uid=%s, road_uid=%s, amount=%.2f",
		refundInfo.RefundUid, refundInfo.RoadUid, refundInfo.RefundAmount)

	// 模拟退款处理
	time.Sleep(2 * time.Second)

	// 实际应该调用对应通道的退款接口
	// supplier := GetSupplierByRoadUid(refundInfo.RoadUid)
	// result := supplier.Refund(refundInfo)

	// 这里模拟成功
	return true, "退款成功"
}

// updateAccountBalance 更新账户余额
func (rs *RefundService) updateAccountBalance(refundInfo *refund.RefundInfo) error {
	// 退款金额应该从商户账户余额中扣除
	accountInfo := accounts.GetAccountByUid(refundInfo.MerchantUid)
	if accountInfo == nil {
		return fmt.Errorf("账户不存在")
	}

	// 扣除余额
	return accounts.UpdateAccountBalanceByUid(
		refundInfo.MerchantUid,
		"REFUND",
		-refundInfo.RefundAmount,
		fmt.Sprintf("订单退款: %s", refundInfo.BankOrderId),
	)
}

// updateOrderRefundStatus 更新订单退款状态
func (rs *RefundService) updateOrderRefundStatus(refundInfo *refund.RefundInfo) error {
	orderInfo, err := order.GetOrderByOrderId(refundInfo.BankOrderId)
	if err != nil {
		return err
	}

	orderInfo.Refund = "YES"
	orderInfo.RefundTime = time.Now().Format("2006-01-02 15:04:05")

	return order.UpdateOrderRefund(orderInfo.BankOrderId, "YES")
}

// sendRefundNotify 发送退款通知
func (rs *RefundService) sendRefundNotify(refundInfo *refund.RefundInfo) {
	if refundInfo.NotifyUrl == "" {
		return
	}

	logger := common.NewLogger().
		WithExtra("refund_uid", refundInfo.RefundUid).
		WithExtra("notify_url", refundInfo.NotifyUrl)

	logger.Info("发送退款通知")

	// TODO: 实现退款通知逻辑
	// 1. 构建通知参数
	// 2. 签名
	// 3. 发送HTTP请求
	// 4. 更新通知状态

	// 这里提供一个简单的实现框架
	// notifyParams := map[string]string{
	//     "refund_uid":    refundInfo.RefundUid,
	//     "order_id":      refundInfo.MerchantOrderId,
	//     "refund_amount": fmt.Sprintf("%.2f", refundInfo.RefundAmount),
	//     "status":        refundInfo.Status,
	// }
	//
	// // 发送通知
	// success := sendHttpNotify(refundInfo.NotifyUrl, notifyParams)
	//
	// if success {
	//     refundInfo.NotifyStatus = "success"
	// } else {
	//     refundInfo.NotifyCount++
	// }
	// refundInfo.Update("NotifyStatus", "NotifyCount")

	logger.Info("退款通知发送完成")
}

// QueryRefund 查询退款状态
func (rs *RefundService) QueryRefund(refundUid string) (*refund.RefundInfo, error) {
	return refund.GetRefundByUid(refundUid)
}

// GetRefundList 获取退款列表
func (rs *RefundService) GetRefundList(page, pageSize int, merchantUid, status string) ([]*refund.RefundInfo, int64, error) {
	filters := make(map[string]interface{})

	if merchantUid != "" {
		filters["merchant_uid"] = merchantUid
	}

	if status != "" {
		filters["status"] = status
	}

	return refund.GetRefundList(page, pageSize, filters)
}
