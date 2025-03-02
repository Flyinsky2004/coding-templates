package entity

// Response 定义通用的返回结构体
type Response[T any] struct {
	Code    int    `json:"code"`    // 状态码
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 消息内容
	Data    T      `json:"data"`    // 泛型数据
}

// NewResponse 创建一个通用响应
func NewResponse[T any](code int, success bool, message string, data T) Response[T] {
	return Response[T]{
		Code:    code,
		Success: success,
		Message: message,
		Data:    data,
	}
}

// SuccessResponse 快捷生成成功响应
func SuccessResponse[T any](data T) Response[T] {
	return NewResponse(200, true, "OK", data)
}
func SuccessResponseWithMessage[T any](message string, data T) Response[T] {
	return NewResponse(200, true, message, data)
}

// ErrorResponse 快捷生成错误响应
func ErrorResponse[T any](code int, message string) Response[T] {
	var empty T // 空的泛型值
	return NewResponse(code, false, message, empty)
}
