package common

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"runtime"
	"time"
)

// LogLevel 日志级别
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// LogContext 日志上下文
type LogContext struct {
	RequestID   string                 `json:"request_id,omitempty"`   // 请求ID
	MerchantUID string                 `json:"merchant_uid,omitempty"` // 商户ID
	OrderID     string                 `json:"order_id,omitempty"`     // 订单号
	UserIP      string                 `json:"user_ip,omitempty"`      // 用户IP
	Method      string                 `json:"method,omitempty"`       // 请求方法
	Path        string                 `json:"path,omitempty"`         // 请求路径
	Extra       map[string]interface{} `json:"extra,omitempty"`        // 额外信息
}

// StructuredLog 结构化日志
type StructuredLog struct {
	Timestamp  string                 `json:"timestamp"`            // 时间戳
	Level      string                 `json:"level"`                // 日志级别
	Message    string                 `json:"message"`              // 日志消息
	Context    *LogContext            `json:"context,omitempty"`    // 日志上下文
	File       string                 `json:"file,omitempty"`       // 文件名
	Line       int                    `json:"line,omitempty"`       // 行号
	Function   string                 `json:"function,omitempty"`   // 函数名
	Error      string                 `json:"error,omitempty"`      // 错误信息
	Duration   int64                  `json:"duration,omitempty"`   // 耗时（毫秒）
	Extra      map[string]interface{} `json:"extra,omitempty"`      // 额外字段
}

// Logger 日志记录器
type Logger struct {
	context *LogContext
}

// NewLogger 创建新的日志记录器
func NewLogger() *Logger {
	return &Logger{
		context: &LogContext{
			Extra: make(map[string]interface{}),
		},
	}
}

// WithContext 设置日志上下文
func (l *Logger) WithContext(ctx *LogContext) *Logger {
	l.context = ctx
	return l
}

// WithRequestID 设置请求ID
func (l *Logger) WithRequestID(requestID string) *Logger {
	l.context.RequestID = requestID
	return l
}

// WithMerchant 设置商户ID
func (l *Logger) WithMerchant(merchantUID string) *Logger {
	l.context.MerchantUID = merchantUID
	return l
}

// WithOrder 设置订单ID
func (l *Logger) WithOrder(orderID string) *Logger {
	l.context.OrderID = orderID
	return l
}

// WithIP 设置用户IP
func (l *Logger) WithIP(ip string) *Logger {
	l.context.UserIP = ip
	return l
}

// WithRequest 设置请求信息
func (l *Logger) WithRequest(method, path string) *Logger {
	l.context.Method = method
	l.context.Path = path
	return l
}

// WithExtra 添加额外字段
func (l *Logger) WithExtra(key string, value interface{}) *Logger {
	if l.context.Extra == nil {
		l.context.Extra = make(map[string]interface{})
	}
	l.context.Extra[key] = value
	return l
}

// buildLog 构建结构化日志
func (l *Logger) buildLog(level string, message string, err error, extra map[string]interface{}) *StructuredLog {
	log := &StructuredLog{
		Timestamp: time.Now().Format("2006-01-02 15:04:05.000"),
		Level:     level,
		Message:   message,
		Context:   l.context,
		Extra:     extra,
	}

	if err != nil {
		log.Error = err.Error()
	}

	// 获取调用者信息
	if pc, file, line, ok := runtime.Caller(2); ok {
		log.File = file
		log.Line = line
		if fn := runtime.FuncForPC(pc); fn != nil {
			log.Function = fn.Name()
		}
	}

	return log
}

// toJSON 将日志转为JSON
func (l *Logger) toJSON(log *StructuredLog) string {
	jsonData, err := json.Marshal(log)
	if err != nil {
		return fmt.Sprintf(`{"level":"ERROR","message":"日志序列化失败: %v"}`, err)
	}
	return string(jsonData)
}

// Debug 记录调试日志
func (l *Logger) Debug(message string, extra ...map[string]interface{}) {
	var extraData map[string]interface{}
	if len(extra) > 0 {
		extraData = extra[0]
	}
	log := l.buildLog("DEBUG", message, nil, extraData)
	logs.Debug(l.toJSON(log))
}

// Info 记录信息日志
func (l *Logger) Info(message string, extra ...map[string]interface{}) {
	var extraData map[string]interface{}
	if len(extra) > 0 {
		extraData = extra[0]
	}
	log := l.buildLog("INFO", message, nil, extraData)
	logs.Info(l.toJSON(log))
}

// Warn 记录警告日志
func (l *Logger) Warn(message string, extra ...map[string]interface{}) {
	var extraData map[string]interface{}
	if len(extra) > 0 {
		extraData = extra[0]
	}
	log := l.buildLog("WARN", message, nil, extraData)
	logs.Warn(l.toJSON(log))
}

// Error 记录错误日志
func (l *Logger) Error(message string, err error, extra ...map[string]interface{}) {
	var extraData map[string]interface{}
	if len(extra) > 0 {
		extraData = extra[0]
	}
	log := l.buildLog("ERROR", message, err, extraData)
	logs.Error(l.toJSON(log))
}

// Fatal 记录致命错误日志
func (l *Logger) Fatal(message string, err error, extra ...map[string]interface{}) {
	var extraData map[string]interface{}
	if len(extra) > 0 {
		extraData = extra[0]
	}
	log := l.buildLog("FATAL", message, err, extraData)
	logs.Critical(l.toJSON(log))
}

// Performance 记录性能日志
func (l *Logger) Performance(message string, duration int64, extra ...map[string]interface{}) {
	var extraData map[string]interface{}
	if len(extra) > 0 {
		extraData = extra[0]
	}
	log := l.buildLog("PERF", message, nil, extraData)
	log.Duration = duration
	logs.Info(l.toJSON(log))
}

// 全局日志实例
var globalLogger = NewLogger()

// 便捷函数

// LogInfo 记录信息日志（全局）
func LogInfo(message string, extra ...map[string]interface{}) {
	globalLogger.Info(message, extra...)
}

// LogWarn 记录警告日志（全局）
func LogWarn(message string, extra ...map[string]interface{}) {
	globalLogger.Warn(message, extra...)
}

// LogError 记录错误日志（全局）
func LogError(message string, err error, extra ...map[string]interface{}) {
	globalLogger.Error(message, err, extra...)
}

// LogDebug 记录调试日志（全局）
func LogDebug(message string, extra ...map[string]interface{}) {
	globalLogger.Debug(message, extra...)
}

// LogPerformance 记录性能日志（全局）
func LogPerformance(message string, duration int64, extra ...map[string]interface{}) {
	globalLogger.Performance(message, duration, extra...)
}

// PaymentLogger 支付日志记录器
type PaymentLogger struct {
	*Logger
	startTime time.Time
}

// NewPaymentLogger 创建支付日志记录器
func NewPaymentLogger(merchantUID, orderID string) *PaymentLogger {
	logger := NewLogger()
	logger.WithMerchant(merchantUID).WithOrder(orderID)

	return &PaymentLogger{
		Logger:    logger,
		startTime: time.Now(),
	}
}

// LogStart 记录支付开始
func (pl *PaymentLogger) LogStart(payType string, amount float64) {
	pl.Info("支付请求开始", map[string]interface{}{
		"pay_type": payType,
		"amount":   amount,
	})
}

// LogSuccess 记录支付成功
func (pl *PaymentLogger) LogSuccess(transID string) {
	duration := time.Since(pl.startTime).Milliseconds()
	pl.Performance("支付成功", duration, map[string]interface{}{
		"trans_id": transID,
	})
}

// LogFail 记录支付失败
func (pl *PaymentLogger) LogFail(reason string, err error) {
	duration := time.Since(pl.startTime).Milliseconds()
	pl.Error("支付失败", err, map[string]interface{}{
		"reason":   reason,
		"duration": duration,
	})
}

// LogRoadSelect 记录通道选择
func (pl *PaymentLogger) LogRoadSelect(roadUID, roadName string) {
	pl.Info("选择支付通道", map[string]interface{}{
		"road_uid":  roadUID,
		"road_name": roadName,
	})
}
