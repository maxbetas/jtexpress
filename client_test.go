package jtexpress

import (
	"encoding/json"
	"testing"
)

func TestSignature(t *testing.T) {
	client := NewClient("test_account", "test_key")

	testCases := []struct {
		name    string
		data    map[string]interface{}
		wantErr bool
	}{
		{
			name: "轨迹查询签名",
			data: map[string]interface{}{
				"billCodes": "JT2099306666983",
				"lang":      "zh",
				"timeType":  "1",
			},
			wantErr: false,
		},
		{
			name: "订阅签名",
			data: map[string]interface{}{
				"id": "test_account",
				"list": []map[string]string{
					{
						"traceNode":   "1&2&3&4&5",
						"waybillCode": "JT2099306666983",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 测试业务参数签名
			bizDigest, err := client.getBizContentDigest(tc.data)
			if (err != nil) != tc.wantErr {
				t.Errorf("getBizContentDigest() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if bizDigest == "" {
				t.Error("Business content digest should not be empty")
			}

			// 测试请求头签名
			jsonStr, err := json.Marshal(tc.data)
			if err != nil {
				t.Errorf("Failed to marshal test data: %v", err)
				return
			}
			headerDigest, err := client.getHeaderDigest(string(jsonStr))
			if (err != nil) != tc.wantErr {
				t.Errorf("getHeaderDigest() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if headerDigest == "" {
				t.Error("Header digest should not be empty")
			}
		})
	}
}

func TestQueryTrack(t *testing.T) {
	client := NewClient("test_account", "test_key")

	testCases := []struct {
		name      string
		billCodes string
		wantErr   bool
		wantCode  int
	}{
		{
			name:      "单个运单查询",
			billCodes: "JT2099306666983",
			wantErr:   false,
			wantCode:  1,
		},
		{
			name:      "批量运单查询",
			billCodes: "JT2099306666983,JT2099306666984",
			wantErr:   false,
			wantCode:  1,
		},
		{
			name:      "空运单号",
			billCodes: "",
			wantErr:   true,
			wantCode:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.Logistics.QueryTrack(tc.billCodes)
			if (err != nil) != tc.wantErr {
				t.Errorf("QueryTrack() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {
				if resp.Code != "1" {
					t.Errorf("Expected code '1', got '%s'", resp.Code)
				}
				if resp.Msg == "" {
					t.Error("Response message should not be empty")
				}
				if resp.Data != nil {
					validateTrackData(t, resp.Data)
				}
			}
		})
	}
}

// validateTrackData 验证轨迹数据
func validateTrackData(t *testing.T, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Failed to marshal response data: %v", err)
		return
	}

	var trackInfo struct {
		BillCode    string `json:"billCode"`
		TrackPoints []struct {
			ScanTime    string `json:"scanTime"`
			ScanType    string `json:"scanType"`
			Location    string `json:"location"`
			Description string `json:"description"`
		} `json:"trackPoints"`
	}

	if err := json.Unmarshal(jsonData, &trackInfo); err != nil {
		t.Errorf("Failed to unmarshal track info: %v", err)
	}
}

func TestSubscribe(t *testing.T) {
	client := NewClient("test_account", "test_key")

	testCases := []struct {
		name      string
		billCode  string
		traceNode string
		backUrl   string
		wantErr   bool
		wantCode  int
	}{
		{
			name:      "正常订阅",
			billCode:  "JT2099306666983",
			traceNode: "1&2&3&4&5&6&7&8&9&10&11",
			backUrl:   "https://your-domain.com/callback",
			wantErr:   false,
			wantCode:  1,
		},
		{
			name:      "空运单号",
			billCode:  "",
			traceNode: "1&2&3&4&5",
			backUrl:   "https://your-domain.com/callback",
			wantErr:   true,
			wantCode:  0,
		},
		{
			name:      "空节点",
			billCode:  "JT2099306666983",
			traceNode: "",
			backUrl:   "https://your-domain.com/callback",
			wantErr:   true,
			wantCode:  0,
		},
		{
			name:      "空回调地址",
			billCode:  "JT2099306666983",
			traceNode: "1&2&3",
			backUrl:   "",
			wantErr:   true,
			wantCode:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.Logistics.Subscribe(tc.billCode, tc.traceNode, tc.backUrl)
			if (err != nil) != tc.wantErr {
				t.Errorf("Subscribe() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr && resp != nil {
				if resp.Code != "1" {
					t.Errorf("Expected code '1', got '%s'", resp.Code)
				}
				if resp.Msg == "" {
					t.Error("Response message should not be empty")
				}

				// 将响应数据转换为具体类型
				jsonData, err := json.Marshal(resp.Data)
				if err != nil {
					t.Errorf("Failed to marshal response data: %v", err)
					return
				}

				var subscribeData SubscribeResponseData
				if err := json.Unmarshal(jsonData, &subscribeData); err != nil {
					t.Errorf("Failed to unmarshal subscribe data: %v", err)
					return
				}

				if len(subscribeData.List) == 0 {
					t.Error("Response data list should not be empty")
				}

				for _, result := range subscribeData.List {
					if !result.IsSuccess {
						t.Errorf("Subscribe failed for waybill %s", result.WaybillCode)
					}
					if result.TraceNode == "" {
						t.Error("TraceNode should not be empty")
					}
					if result.ID == "" {
						t.Error("ID should not be empty")
					}
				}
			}
		})
	}
}

func TestSubscribeBatch(t *testing.T) {
	client := NewClient("test_account", "test_key")

	testCases := []struct {
		name      string
		billCodes []string
		traceNode string
		backUrl   string
		wantErr   bool
		wantCode  int
	}{
		{
			name:      "正常批量订阅",
			billCodes: []string{"JT2099306666983", "JT2099306666984"},
			traceNode: "1&2&3&4&5&6&7&8&9&10&11",
			backUrl:   "https://your-domain.com/callback",
			wantErr:   false,
			wantCode:  1,
		},
		{
			name:      "空运单号列表",
			billCodes: []string{},
			traceNode: "1&2&3&4&5",
			backUrl:   "https://your-domain.com/callback",
			wantErr:   true,
			wantCode:  0,
		},
		{
			name:      "空回调地址",
			billCodes: []string{"JT2099306666983"},
			traceNode: "1&2&3",
			backUrl:   "",
			wantErr:   true,
			wantCode:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.Logistics.SubscribeBatch(tc.billCodes, tc.traceNode, tc.backUrl)
			if (err != nil) != tc.wantErr {
				t.Errorf("SubscribeBatch() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr && resp != nil {
				if resp.Code != "1" {
					t.Errorf("Expected code '1', got '%s'", resp.Code)
				}
				if resp.Msg == "" {
					t.Error("Response message should not be empty")
				}
				if resp.Code == "1" && resp.Data != "SUCCESS" {
					t.Errorf("Expected data 'SUCCESS', got %v", resp.Data)
				}
			}
		})
	}
}
