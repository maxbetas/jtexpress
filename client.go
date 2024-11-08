package jtexpress

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"jtexpress/services"
)

// Package jtexpress 提供极兔快递 API 的 Go 语言封装

// Client 是极兔快递 API 的客户端
// 它提供了对极兔快递各项服务的访问能力
type Client struct {
	apiAccount string       // API 账号
	signer     Signer       // 签名器
	baseURL    string       // API 基础URL
	httpClient *http.Client // HTTP 客户端

	// 服务
	Logistics *services.LogisticsService // 物流服务

	// 缓存
	cache sync.Map // 用于缓存一些常用数据
}

// ClientOption 定义客户端配置选项
type ClientOption func(*Client)

// WithHTTPClient 设置自定义的 HTTP 客户端
// 可用于配置超时、代理等选项
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL 设置自定义的基础 URL
// 用于测试环境或特殊部署场景
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// NewClient 创建新的极兔快递 API 客户端
// apiAccount: API账号
// privateKey: API私钥
// opts: 可选的配置选项
func NewClient(apiAccount, privateKey string, opts ...ClientOption) *Client {
	c := &Client{
		apiAccount: apiAccount,
		signer:     NewMD5Signer(privateKey),
		baseURL:    "https://openapi.jtexpress.com.cn",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	// 应用配置选项
	for _, opt := range opts {
		opt(c)
	}

	// 初始化服务
	c.Logistics = services.NewLogisticsService(c)

	return c
}

// Post 发送 POST 请求到指定的 API 路径
// data: 请求数据
// apiPath: API路径
// result: 响应结果的指针
func (c *Client) Post(data interface{}, apiPath string, result interface{}) error {
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)

	// 1. 计算业务参数的digest
	bizDigest, err := c.signer.SignStruct(data)
	if err != nil {
		return fmt.Errorf("calculate biz digest error: %w", err)
	}

	// 2. 准备业务参数
	bizContent := map[string]interface{}{
		"digest": bizDigest,
	}

	// 将请求数据合并到bizContent
	if reqData, err := json.Marshal(data); err == nil {
		var reqMap map[string]interface{}
		if err := json.Unmarshal(reqData, &reqMap); err == nil {
			for k, v := range reqMap {
				bizContent[k] = v
			}
		}
	}

	// 3. 准备请求
	jsonBytes, err := json.Marshal(bizContent)
	if err != nil {
		return fmt.Errorf("marshal bizContent error: %w", err)
	}

	// 4. 计算请求头的digest
	headerDigest, err := c.signer.Sign(string(jsonBytes))
	if err != nil {
		return fmt.Errorf("calculate header digest error: %w", err)
	}

	// 5. 发送请求
	formValues := url.Values{}
	formValues.Set("bizContent", string(jsonBytes))

	fullURL := c.baseURL + "/webopenplatformapi" + apiPath
	req, err := http.NewRequest(http.MethodPost, fullURL, strings.NewReader(formValues.Encode()))
	if err != nil {
		return fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apiAccount", c.apiAccount)
	req.Header.Set("timestamp", timestamp)
	req.Header.Set("digest", headerDigest)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request error: %w", err)
	}
	defer resp.Body.Close()

	// 6. 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response error: %w", err)
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("parse response error: %w", err)
	}

	return nil
}

// GetAPIAccount 获取API账号（数字格式）
func (c *Client) GetAPIAccount() int64 {
	account, _ := strconv.ParseInt(c.apiAccount, 10, 64)
	return account
}

// 添加缓存方法
func (c *Client) getCachedValue(key string) (interface{}, bool) {
	return c.cache.Load(key)
}

func (c *Client) setCachedValue(key string, value interface{}) {
	c.cache.Store(key, value)
}
