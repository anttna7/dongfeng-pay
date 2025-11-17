package common

// ErrorCode 错误码定义
type ErrorCode struct {
	Code    string // 错误码
	Message string // 错误信息
}

// 系统级错误码 (1xxx)
var (
	SUCCESS              = ErrorCode{"0000", "成功"}
	SYSTEM_ERROR         = ErrorCode{"1000", "系统错误"}
	PARAMS_ERROR         = ErrorCode{"1001", "参数错误"}
	SIGN_ERROR           = ErrorCode{"1002", "签名错误"}
	TIMESTAMP_ERROR      = ErrorCode{"1003", "时间戳错误"}
	TIMESTAMP_EXPIRED    = ErrorCode{"1004", "请求超时"}
	RATE_LIMIT_EXCEEDED  = ErrorCode{"1005", "请求过于频繁"}
	NETWORK_ERROR        = ErrorCode{"1006", "网络错误"}
	DATABASE_ERROR       = ErrorCode{"1007", "数据库错误"}
	REDIS_ERROR          = ErrorCode{"1008", "缓存错误"}
	JSON_PARSE_ERROR     = ErrorCode{"1009", "JSON解析错误"}
)

// 商户相关错误码 (2xxx)
var (
	MERCHANT_NOT_EXIST   = ErrorCode{"2001", "商户不存在"}
	MERCHANT_FROZEN      = ErrorCode{"2002", "商户已冻结"}
	MERCHANT_INVALID     = ErrorCode{"2003", "商户无效"}
	MERCHANT_NOT_CONFIG  = ErrorCode{"2004", "商户未配置"}
	IP_NOT_ALLOWED       = ErrorCode{"2005", "IP地址不在白名单"}
	MERCHANT_SECRET_ERR  = ErrorCode{"2006", "商户密钥错误"}
	BALANCE_NOT_ENOUGH   = ErrorCode{"2007", "余额不足"}
)

// 通道相关错误码 (3xxx)
var (
	ROAD_NOT_EXIST       = ErrorCode{"3001", "通道不存在"}
	ROAD_CLOSED          = ErrorCode{"3002", "通道已关闭"}
	ROAD_MAINTAIN        = ErrorCode{"3003", "通道维护中"}
	NO_AVAILABLE_ROAD    = ErrorCode{"3004", "无可用通道"}
	ROAD_LIMIT_EXCEEDED  = ErrorCode{"3005", "超出通道限额"}
	ROAD_CONFIG_ERROR    = ErrorCode{"3006", "通道配置错误"}
)

// 订单相关错误码 (4xxx)
var (
	ORDER_NOT_EXIST      = ErrorCode{"4001", "订单不存在"}
	ORDER_DUPLICATE      = ErrorCode{"4002", "订单重复"}
	ORDER_STATUS_ERROR   = ErrorCode{"4003", "订单状态错误"}
	ORDER_AMOUNT_ERROR   = ErrorCode{"4004", "订单金额错误"}
	ORDER_TIMEOUT        = ErrorCode{"4005", "订单超时"}
	ORDER_FROZEN         = ErrorCode{"4006", "订单已冻结"}
	ORDER_CREATE_FAIL    = ErrorCode{"4007", "订单创建失败"}
	ORDER_UPDATE_FAIL    = ErrorCode{"4008", "订单更新失败"}
)

// 支付相关错误码 (5xxx)
var (
	PAY_TYPE_NOT_SUPPORT = ErrorCode{"5001", "不支持的支付类型"}
	PAY_AMOUNT_INVALID   = ErrorCode{"5002", "支付金额无效"}
	PAY_FAILED           = ErrorCode{"5003", "支付失败"}
	PAY_PROCESSING       = ErrorCode{"5004", "支付处理中"}
	PAY_CANCELED         = ErrorCode{"5005", "支付已取消"}
	PAY_CLOSED           = ErrorCode{"5006", "支付已关闭"}
)

// 代付相关错误码 (6xxx)
var (
	PAYFOR_NOT_EXIST     = ErrorCode{"6001", "代付订单不存在"}
	PAYFOR_DUPLICATE     = ErrorCode{"6002", "代付订单重复"}
	PAYFOR_FAILED        = ErrorCode{"6003", "代付失败"}
	PAYFOR_PROCESSING    = ErrorCode{"6004", "代付处理中"}
	PAYFOR_BANK_ERROR    = ErrorCode{"6005", "银行卡信息错误"}
	PAYFOR_AMOUNT_ERROR  = ErrorCode{"6006", "代付金额错误"}
	PAYFOR_NOT_CONFIG    = ErrorCode{"6007", "商户未开通代付"}
)

// 账户相关错误码 (7xxx)
var (
	ACCOUNT_NOT_EXIST    = ErrorCode{"7001", "账户不存在"}
	ACCOUNT_FROZEN       = ErrorCode{"7002", "账户已冻结"}
	ACCOUNT_BALANCE_ERR  = ErrorCode{"7003", "账户余额异常"}
	ACCOUNT_OPERATE_FAIL = ErrorCode{"7004", "账户操作失败"}
)

// 通知相关错误码 (8xxx)
var (
	NOTIFY_FAILED        = ErrorCode{"8001", "通知发送失败"}
	NOTIFY_TIMEOUT       = ErrorCode{"8002", "通知超时"}
	NOTIFY_MAX_RETRY     = ErrorCode{"8003", "通知重试次数已达上限"}
)

// 退款相关错误码 (9xxx)
var (
	REFUND_NOT_ALLOWED   = ErrorCode{"9001", "不允许退款"}
	REFUND_AMOUNT_ERROR  = ErrorCode{"9002", "退款金额错误"}
	REFUND_FAILED        = ErrorCode{"9003", "退款失败"}
	REFUND_PROCESSING    = ErrorCode{"9004", "退款处理中"}
	REFUND_DUPLICATE     = ErrorCode{"9005", "重复退款"}
)
