package jtexpress

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client JT Express API 客户端
type Client struct {
	apiAccount string
	privateKey string
	baseURL    string
	httpClient *http.Client
	Logistics  *LogisticsService
}

// NewClient 创建新的客户端
func NewClient(apiAccount, privateKey string) *Client {
	client := &Client{
		apiAccount: apiAccount,
		privateKey: privateKey,
		baseURL:    "https://openapi.jtexpress.com.cn/webopenplatformapi",
		httpClient: &http.Client{},
	}
	client.Logistics = NewLogisticsService(client)
	return client
}

// Post 发送请求
func (c *Client) Post(data interface{}, apiPath string) (*Response, error) {
	// 1. 准备请求数据
	bizContent, err := c.prepareBizContent(data)
	if err != nil {
		return nil, err
	}

	// 2. 发送请求
	resp, err := c.doRequest(bizContent, apiPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 3. 处理响应
	return c.handleResponse(resp)
}

// prepareBizContent 准备业务参数
func (c *Client) prepareBizContent(data interface{}) (string, error) {
	// 序列化数据
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("marshal request error: %w", err)
	}

	return string(jsonBytes), nil
}

// doRequest 执行HTTP请求
func (c *Client) doRequest(bizContent, apiPath string) (*http.Response, error) {
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	formData := url.Values{}
	formData.Set("bizContent", bizContent)

	req, err := http.NewRequest("POST", c.baseURL+apiPath, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}

	// 计算请求头的digest
	headerDigest, err := c.getHeaderDigest(bizContent)
	if err != nil {
		return nil, fmt.Errorf("generate header digest error: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apiAccount", c.apiAccount)
	req.Header.Set("timestamp", timestamp)
	req.Header.Set("digest", headerDigest)

	return c.httpClient.Do(req)
}

// handleResponse 处理响应
func (c *Client) handleResponse(resp *http.Response) (*Response, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("parse response error: %w", err)
	}

	response.Success = response.Code == "1"
	response.Fail = !response.Success

	return &response, nil
}

// getHeaderDigest 计算请求头摘要
func (c *Client) getHeaderDigest(jsonStr string) (string, error) {
	str := jsonStr + c.privateKey
	md5Sum := md5.Sum([]byte(str))
	return base64.StdEncoding.EncodeToString(md5Sum[:]), nil
}

// getBizContentDigest 计算业务参数摘要
func (c *Client) getBizContentDigest(data interface{}) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	str := string(jsonBytes) + c.privateKey
	md5Sum := md5.Sum([]byte(str))
	return base64.StdEncoding.EncodeToString(md5Sum[:]), nil
}
