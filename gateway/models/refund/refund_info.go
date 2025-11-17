package refund

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"time"
)

// RefundInfo 退款信息表
type RefundInfo struct {
	Id                int       `orm:"column(id);auto"`
	RefundUid         string    `orm:"column(refund_uid);size(32)"`          // 退款单号
	MerchantOrderId   string    `orm:"column(merchant_order_id);size(64)"`   // 商户订单号
	BankOrderId       string    `orm:"column(bank_order_id);size(64)"`       // 系统订单号
	BankTransId       string    `orm:"column(bank_trans_id);size(128)"`      // 上游流水号
	RefundAmount      float64   `orm:"column(refund_amount);digits(10);decimals(2)"` // 退款金额
	RefundReason      string    `orm:"column(refund_reason);size(255)"`      // 退款原因
	RefundType        string    `orm:"column(refund_type);size(20)"`         // 退款类型：full-全额，partial-部分
	Status            string    `orm:"column(status);size(20)"`              // 退款状态：pending-待处理，processing-处理中，success-成功，failed-失败
	Result            string    `orm:"column(result);type(text)"`            // 退款结果（JSON）
	MerchantUid       string    `orm:"column(merchant_uid);size(32)"`        // 商户UID
	MerchantName      string    `orm:"column(merchant_name);size(128)"`      // 商户名称
	AgentUid          string    `orm:"column(agent_uid);size(32)"`           // 代理UID
	RoadUid           string    `orm:"column(road_uid);size(32)"`            // 通道UID
	OperatorName      string    `orm:"column(operator_name);size(64)"`       // 操作人
	NotifyUrl         string    `orm:"column(notify_url);size(255)"`         // 退款回调地址
	NotifyStatus      string    `orm:"column(notify_status);size(20)"`       // 通知状态
	NotifyCount       int       `orm:"column(notify_count);default(0)"`      // 通知次数
	CreateTime        time.Time `orm:"column(create_time);type(datetime);auto_now_add"` // 创建时间
	UpdateTime        time.Time `orm:"column(update_time);type(datetime);auto_now"`     // 更新时间
	FinishTime        time.Time `orm:"column(finish_time);type(datetime);null"`         // 完成时间
}

func init() {
	orm.RegisterModel(new(RefundInfo))
}

// TableName 指定表名
func (r *RefundInfo) TableName() string {
	return "refund_info"
}

// Insert 插入退款记录
func (r *RefundInfo) Insert() error {
	o := orm.NewOrm()
	_, err := o.Insert(r)
	if err != nil {
		logs.Error("插入退款记录失败: %v", err)
		return err
	}
	return nil
}

// Update 更新退款记录
func (r *RefundInfo) Update(fields ...string) error {
	o := orm.NewOrm()
	_, err := o.Update(r, fields...)
	if err != nil {
		logs.Error("更新退款记录失败: %v", err)
		return err
	}
	return nil
}

// GetRefundByUid 根据退款单号查询
func GetRefundByUid(refundUid string) (*RefundInfo, error) {
	o := orm.NewOrm()
	refund := &RefundInfo{RefundUid: refundUid}
	err := o.Read(refund, "RefundUid")
	if err != nil {
		return nil, err
	}
	return refund, nil
}

// GetRefundByOrderId 根据订单号查询退款记录
func GetRefundByOrderId(bankOrderId string) ([]*RefundInfo, error) {
	o := orm.NewOrm()
	var refunds []*RefundInfo
	_, err := o.QueryTable(new(RefundInfo)).
		Filter("bank_order_id", bankOrderId).
		OrderBy("-create_time").
		All(&refunds)
	if err != nil {
		return nil, err
	}
	return refunds, nil
}

// GetRefundList 分页查询退款列表
func GetRefundList(page, pageSize int, filters map[string]interface{}) ([]*RefundInfo, int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(RefundInfo))

	// 应用过滤条件
	for key, value := range filters {
		if value != "" && value != nil {
			qs = qs.Filter(key, value)
		}
	}

	// 总数
	total, err := qs.Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	var refunds []*RefundInfo
	offset := (page - 1) * pageSize
	_, err = qs.OrderBy("-create_time").Limit(pageSize, offset).All(&refunds)
	if err != nil {
		return nil, 0, err
	}

	return refunds, total, nil
}

// GetRefundTotalAmount 获取某商户的总退款金额
func GetRefundTotalAmount(merchantUid string, startDate, endDate string) (float64, error) {
	o := orm.NewOrm()

	var maps []orm.Params
	_, err := o.Raw(`
		SELECT IFNULL(SUM(refund_amount), 0) as total
		FROM refund_info
		WHERE merchant_uid = ?
		  AND status = 'success'
		  AND create_time >= ?
		  AND create_time <= ?
	`, merchantUid, startDate, endDate).Values(&maps)

	if err != nil {
		return 0, err
	}

	if len(maps) > 0 {
		if total, ok := maps[0]["total"].(string); ok {
			totalFloat, _ := strconv.ParseFloat(total, 64)
			return totalFloat, nil
		}
	}

	return 0, nil
}

// UpdateRefundStatus 更新退款状态
func UpdateRefundStatus(refundUid, status, result string) error {
	o := orm.NewOrm()

	refund := &RefundInfo{RefundUid: refundUid}
	if err := o.Read(refund, "RefundUid"); err != nil {
		return err
	}

	refund.Status = status
	refund.Result = result

	if status == "success" || status == "failed" {
		refund.FinishTime = time.Now()
	}

	_, err := o.Update(refund, "Status", "Result", "FinishTime", "UpdateTime")
	return err
}

// GetPendingRefunds 获取待处理的退款
func GetPendingRefunds(limit int) ([]*RefundInfo, error) {
	o := orm.NewOrm()
	var refunds []*RefundInfo

	_, err := o.QueryTable(new(RefundInfo)).
		Filter("status", "pending").
		OrderBy("create_time").
		Limit(limit).
		All(&refunds)

	if err != nil {
		return nil, err
	}

	return refunds, nil
}
