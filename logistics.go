package jtexpress

import (
	"fmt"
)

// LogisticsService 物流服务
type LogisticsService struct {
	client *Client
}

// NewLogisticsService 创建物流服务实例
func NewLogisticsService(client *Client) *LogisticsService {
	return &LogisticsService{client: client}
}

// QueryTrack 查询物流轨迹
func (s *LogisticsService) QueryTrack(billCodes string) (*Response, error) {
	if billCodes == "" {
		return nil, fmt.Errorf("billCodes cannot be empty")
	}

	data := map[string]interface{}{
		"billCodes": billCodes,
		"lang":      "zh",
		"timeType":  "1",
	}

	return s.client.Post(data, "/api/logistics/trace")
}

// Subscribe 订阅物流轨迹
func (s *LogisticsService) Subscribe(billCode, traceNode, backUrl string) (*Response, error) {
	if err := s.validateSubscribeParams(billCode, traceNode, backUrl); err != nil {
		return nil, err
	}

	req := SubscribeRequest{
		ID: s.client.apiAccount,
		List: []SubscribeTrace{{
			TraceNode:   traceNode,
			WaybillCode: billCode,
			BackUrl:     backUrl,
		}},
	}

	return s.client.Post(req, "/api/trace/subscribe")
}

// SubscribeBatch 批量订阅物流轨迹
func (s *LogisticsService) SubscribeBatch(billCodes []string, traceNode, backUrl string) (*Response, error) {
	if err := s.validateSubscribeParams("", traceNode, backUrl); err != nil {
		return nil, err
	}
	if len(billCodes) == 0 {
		return nil, fmt.Errorf("billCodes cannot be empty")
	}

	traces := make([]SubscribeTrace, len(billCodes))
	for i, code := range billCodes {
		traces[i] = SubscribeTrace{
			TraceNode:   traceNode,
			WaybillCode: code,
			BackUrl:     backUrl,
		}
	}

	req := SubscribeRequest{
		ID:   s.client.apiAccount,
		List: traces,
	}

	return s.client.Post(req, "/api/trace/subscribe")
}

// validateSubscribeParams 验证订阅参数
func (s *LogisticsService) validateSubscribeParams(billCode, traceNode, backUrl string) error {
	if billCode != "" && billCode == "" {
		return fmt.Errorf("billCode cannot be empty")
	}
	if traceNode == "" {
		return fmt.Errorf("traceNode cannot be empty")
	}
	if backUrl == "" {
		return fmt.Errorf("backUrl cannot be empty")
	}
	return nil
}
