package types

// TrackRequest 轨迹查询请求
type TrackRequest struct {
	CommonRequest
	BillCodes string `json:"billCodes"`
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
	ScanTime    string `json:"scanTime"`                 // 扫描时间
	ScanType    string `json:"scanType"`                 // 扫描类型
	Description string `json:"desc"`                     // 描述信息，注意这里改为 desc
	Location    string `json:"scanNetworkDetailAddress"` // 位置信息，改为实际的地址字段

	// 可选的额外信息
	NetworkName  string `json:"scanNetworkName"`     // 网点名称
	NetworkType  string `json:"networkType"`         // 网点类型
	StaffName    string `json:"staffName"`           // 工作人员姓名
	StaffContact string `json:"staffContact"`        // 工作人员联系方式
	NextStopName string `json:"nextStopName"`        // 下一站名称
	Province     string `json:"scanNetworkProvince"` // 省份
	City         string `json:"scanNetworkCity"`     // 城市
	Area         string `json:"scanNetworkArea"`     // 区域
}

// SubscribeRequest 订阅请求
type SubscribeRequest struct {
	ID   int64            `json:"id"`   // 接入方在平台的api账户标识
	List []SubscribeTrace `json:"list"` // 订阅列表
}

// SubscribeTrace 订阅轨迹信息
type SubscribeTrace struct {
	TraceNode   string `json:"traceNode"`   // 订阅节点：1&2&3&4&5&6&7&8&9&10&11
	WaybillCode string `json:"waybillCode"` // 极兔运单编号，最大32位
	BackUrl     string `json:"backUrl"`     // 回调通知地址
}

// SubscribeResponse 订阅响应
type SubscribeResponse struct {
	Code string `json:"code"` // 响应码
	Msg  string `json:"msg"`  // 响应信息
	Data struct {
		List []interface{} `json:"list"` // 响应数据列表
	} `json:"data"`
}
