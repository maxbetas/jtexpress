package errors

import "fmt"

// APIError 表示API调用错误
type APIError struct {
	Code    string
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: [%s] %s", e.Code, e.Message)
}

// NewAPIError 创建新的API错误
func NewAPIError(code, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}
