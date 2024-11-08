package services

import (
	"fmt"

	"github.com/maxbetas/jtexpress/types"
)

// Package services 提供各种 API 服务的实现

// LogisticsService 提供物流相关的 API 服务
type LogisticsService struct {
	client APIClient
}

// NewLogisticsService 创建物流服务实例
func NewLogisticsService(client APIClient) *LogisticsService {
	return &LogisticsService{client: client}
}

// QueryTrack 查询物流轨迹
// billCodes: 运单号，支持多个运单号，用英文逗号分隔，如："JT1234,JT5678"
func (s *LogisticsService) QueryTrack(billCodes string) (*types.TrackResponse, error) {
	// 参数验证
	if billCodes == "" {
		return nil, fmt.Errorf("billCodes cannot be empty")
	}

	req := types.TrackRequest{
		CommonRequest: types.CommonRequest{
			Lang:     "zh",
			TimeType: "1",
		},
		BillCodes: billCodes, // 支持多个运单号，用逗号分隔
	}

	var resp types.TrackResponse
	err := s.client.Post(req, "/api/logistics/trace", &resp)
	return &resp, err
}

// Subscribe 订阅物流轨迹
func (s *LogisticsService) Subscribe(billCodes interface{}, traceNode, backUrl string) (*types.SubscribeResponse, error) {
	// 参数验证
	if traceNode == "" {
		return nil, fmt.Errorf("traceNode cannot be empty")
	}
	if backUrl == "" {
		return nil, fmt.Errorf("backUrl cannot be empty")
	}

	var list []types.SubscribeTrace
	switch v := billCodes.(type) {
	case string:
		if v == "" {
			return nil, fmt.Errorf("billCode cannot be empty")
		}
		list = []types.SubscribeTrace{{
			TraceNode:   traceNode,
			WaybillCode: v,
		}}
	case []string:
		if len(v) == 0 {
			return nil, fmt.Errorf("billCodes list cannot be empty")
		}
		list = make([]types.SubscribeTrace, len(v))
		for i, code := range v {
			if code == "" {
				return nil, fmt.Errorf("billCode at index %d cannot be empty", i)
			}
			list[i] = types.SubscribeTrace{
				TraceNode:   traceNode,
				WaybillCode: code,
			}
		}
	default:
		return nil, fmt.Errorf("billCodes must be string or []string")
	}

	req := types.SubscribeRequest{
		ID:   s.client.GetAPIAccount(),
		List: list,
	}

	var resp types.SubscribeResponse
	err := s.client.Post(req, "/api/trace/subscribe", &resp)
	return &resp, err
}

// SubscribeBatch 批量订阅物流轨迹
// billCodes: 运单号列表
// traceNode: 订阅节点，如 "1&2&3&4&5"
// backUrl: 回调通知地址
func (s *LogisticsService) SubscribeBatch(billCodes []string, traceNode, backUrl string) (*types.SubscribeResponse, error) {
	// 构建订阅列表
	list := make([]types.SubscribeTrace, len(billCodes))
	for i, billCode := range billCodes {
		list[i] = types.SubscribeTrace{
			TraceNode:   traceNode,
			WaybillCode: billCode,
		}
	}

	req := types.SubscribeRequest{
		ID:   s.client.GetAPIAccount(),
		List: list,
	}

	var resp types.SubscribeResponse
	err := s.client.Post(req, "/api/trace/subscribe", &resp)
	return &resp, err
}

// 这里可以继续添加其他物流相关的接口
// func (s *LogisticsService) CreateOrder(...) {...}
// func (s *LogisticsService) CancelOrder(...) {...}
