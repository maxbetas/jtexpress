package jtexpress

import (
	"net/http"
	"testing"
	"time"
)

// TestNewClient 测试客户端初始化
// 验证客户端是否正确初始化，包括：
// - API账号设置
// - 默认基础URL
// - HTTP客户端
// - 物流服务
func TestNewClient(t *testing.T) {
	tests := []struct {
		name       string // 测试用例名称
		apiAccount string // API账号
		privateKey string // API私钥
		wantErr    bool   // 是否期望错误
	}{
		{
			name:       "正常初始化",
			apiAccount: "602773536896983053",
			privateKey: "test_key",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.apiAccount, tt.privateKey)

			// 验证 API 账号是否正确设置
			if client.apiAccount != tt.apiAccount {
				t.Errorf("apiAccount = %v, want %v", client.apiAccount, tt.apiAccount)
			}

			// 验证默认基础 URL 是否正确
			if client.baseURL != "https://openapi.jtexpress.com.cn" {
				t.Errorf("baseURL = %v, want %v", client.baseURL, "https://openapi.jtexpress.com.cn")
			}

			// 验证 HTTP 客户端是否已初始化
			if client.httpClient == nil {
				t.Error("httpClient should not be nil")
			}

			// 验证物流服务是否已初始化
			if client.Logistics == nil {
				t.Error("Logistics service should not be nil")
			}
		})
	}
}

// TestWithOptions 测试客户端配置选项
// 验证自定义配置是否生效，包括：
// - 自定义 HTTP 客户端超时
// - 自定义基础 URL
func TestWithOptions(t *testing.T) {
	customTimeout := 5 * time.Second
	customURL := "https://test.api.com"

	client := NewClient(
		"test_account",
		"test_key",
		WithHTTPClient(&http.Client{Timeout: customTimeout}),
		WithBaseURL(customURL),
	)

	// 验证自定义超时设置是否生效
	if client.httpClient.Timeout != customTimeout {
		t.Errorf("timeout = %v, want %v", client.httpClient.Timeout, customTimeout)
	}

	// 验证自定义 URL 是否生效
	if client.baseURL != customURL {
		t.Errorf("baseURL = %v, want %v", client.baseURL, customURL)
	}
}

// TestGetAPIAccount 测试 API 账号转换
// 验证字符串格式的 API 账号是否能正确转换为数字格式
func TestGetAPIAccount(t *testing.T) {
	tests := []struct {
		name       string // 测试用例名称
		apiAccount string // 输入的 API 账号
		want       int64  // 期望的数字结果
	}{
		{
			name:       "正常数字",
			apiAccount: "602773536896983053",
			want:       602773536896983053,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.apiAccount, "test_key")
			if got := client.GetAPIAccount(); got != tt.want {
				t.Errorf("GetAPIAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}
