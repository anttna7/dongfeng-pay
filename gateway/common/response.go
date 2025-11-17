package common

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
)

// Response 统一响应结构
type Response struct {
	Code      string      `json:"code"`                // 错误码
	Message   string      `json:"message"`             // 错误信息
	Data      interface{} `json:"data,omitempty"`      // 响应数据
	Timestamp int64       `json:"timestamp,omitempty"` // 时间戳
}

// PageResponse 分页响应结构
type PageResponse struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Total     int64       `json:"total,omitempty"`     // 总数
	Page      int         `json:"page,omitempty"`      // 当前页
	PageSize  int         `json:"pageSize,omitempty"`  // 每页数量
	Timestamp int64       `json:"timestamp,omitempty"`
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Code:    SUCCESS.Code,
		Message: SUCCESS.Message,
		Data:    data,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(errorCode ErrorCode) *Response {
	return &Response{
		Code:    errorCode.Code,
		Message: errorCode.Message,
	}
}

// NewErrorResponseWithMsg 创建带自定义消息的错误响应
func NewErrorResponseWithMsg(errorCode ErrorCode, customMsg string) *Response {
	return &Response{
		Code:    errorCode.Code,
		Message: customMsg,
	}
}

// NewPageResponse 创建分页响应
func NewPageResponse(data interface{}, total int64, page, pageSize int) *PageResponse {
	return &PageResponse{
		Code:     SUCCESS.Code,
		Message:  SUCCESS.Message,
		Data:     data,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}

// ToJSON 将响应转换为JSON字符串
func (r *Response) ToJSON() string {
	jsonData, err := json.Marshal(r)
	if err != nil {
		logs.Error("响应转JSON失败: %v", err)
		return `{"code":"1000","message":"系统错误"}`
	}
	return string(jsonData)
}

// ToMap 将响应转换为Map
func (r *Response) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"code":    r.Code,
		"message": r.Message,
	}

	if r.Data != nil {
		result["data"] = r.Data
	}

	if r.Timestamp > 0 {
		result["timestamp"] = r.Timestamp
	}

	return result
}

// BusinessError 业务错误结构
type BusinessError struct {
	ErrorCode ErrorCode
	Message   string
	Detail    string
}

// Error 实现error接口
func (e *BusinessError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.ErrorCode.Message
}

// NewBusinessError 创建业务错误
func NewBusinessError(errorCode ErrorCode) *BusinessError {
	return &BusinessError{
		ErrorCode: errorCode,
		Message:   errorCode.Message,
	}
}

// NewBusinessErrorWithDetail 创建带详情的业务错误
func NewBusinessErrorWithDetail(errorCode ErrorCode, detail string) *BusinessError {
	return &BusinessError{
		ErrorCode: errorCode,
		Detail:    detail,
	}
}

// HandleError 统一错误处理
func HandleError(err error) *Response {
	if err == nil {
		return NewSuccessResponse(nil)
	}

	// 如果是业务错误，直接返回
	if bizErr, ok := err.(*BusinessError); ok {
		if bizErr.Detail != "" {
			return NewErrorResponseWithMsg(bizErr.ErrorCode, bizErr.Detail)
		}
		return NewErrorResponse(bizErr.ErrorCode)
	}

	// 其他错误统一返回系统错误
	logs.Error("系统错误: %v", err)
	return NewErrorResponse(SYSTEM_ERROR)
}
