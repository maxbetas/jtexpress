package jtexpress

// CommonRequest 通用请求参数
type CommonRequest struct {
	Lang     string `json:"lang,omitempty"`
	TimeType string `json:"timeType,omitempty"`
}

// TrackRequest 轨迹查询请求
type TrackRequest struct {
	CommonRequest
	BillCodes string `json:"billCodes"`
}

// Response 通用响应结构
type Response struct {
	Code    string      `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

// TrackResponse 轨迹查询响应
type TrackResponse struct {
	Response
	Data []TrackInfo `json:"data"`
}

// TrackInfo 轨迹信息
type TrackInfo struct {
	BillCode    string       `json:"billCode"`
	TrackPoints []TrackPoint `json:"details"`
}

// TrackPoint 轨迹节点
type TrackPoint struct {
	ScanTime    string `json:"scanTime"`
	ScanType    string `json:"scanType"`
	Description string `json:"description"`
	Location    string `json:"location"`
}
