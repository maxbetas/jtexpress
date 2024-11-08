package jtexpress

// Response 通用响应结构
type Response struct {
	Code    string      `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	Success bool        `json:"success,omitempty"`
	Fail    bool        `json:"fail,omitempty"`
}

// TrackInfo 轨迹信息
type TrackInfo struct {
	BillCode string       `json:"billCode"`
	Details  []TrackPoint `json:"details"`
}

// TrackPoint 轨迹点信息
type TrackPoint struct {
	ScanTime                 string `json:"scanTime"`
	Desc                     string `json:"desc"`
	ScanType                 string `json:"scanType"`
	ScanNetworkTypeName      string `json:"scanNetworkTypeName"`
	ScanNetworkName          string `json:"scanNetworkName"`
	ScanNetworkId            int    `json:"scanNetworkId"`
	ScanNetworkProvince      string `json:"scanNetworkProvince"`
	ScanNetworkCity          string `json:"scanNetworkCity"`
	ScanNetworkArea          string `json:"scanNetworkArea"`
	ScanNetworkDetailAddress string `json:"scanNetworkDetailAddress"`
}

// SubscribeTrace 订阅信息
type SubscribeTrace struct {
	TraceNode   string `json:"traceNode"`
	WaybillCode string `json:"waybillCode"`
	BackUrl     string `json:"backUrl"`
}

// SubscribeRequest 订阅请求
type SubscribeRequest struct {
	ID   string           `json:"id"`
	List []SubscribeTrace `json:"list"`
}

// SubscribeResult 订阅结果
type SubscribeResult struct {
	ID          string `json:"id"`
	WaybillCode string `json:"waybillCode"`
	IsSuccess   bool   `json:"isSuccess"`
	TraceNode   string `json:"traceNode"`
}

// SubscribeResponseData 订阅响应数据
type SubscribeResponseData struct {
	List []SubscribeResult `json:"list"`
}

// SubscribeResponse 订阅响应
type SubscribeResponse struct {
	Code    string                `json:"code"`
	Msg     string                `json:"msg"`
	Data    SubscribeResponseData `json:"data"`
	Success bool                  `json:"success,omitempty"`
	Fail    bool                  `json:"fail,omitempty"`
}

// TrackResponse 轨迹查询响应
type TrackResponse struct {
	Code    string      `json:"code"`
	Msg     string      `json:"msg"`
	Data    []TrackInfo `json:"data"`
	Success bool        `json:"success,omitempty"`
	Fail    bool        `json:"fail,omitempty"`
}
