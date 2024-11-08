package types

// CommonRequest 通用请求参数
type CommonRequest struct {
	Lang     string `json:"lang,omitempty"`
	TimeType string `json:"timeType,omitempty"`
}

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
	Fail    bool        `json:"fail"`
}
