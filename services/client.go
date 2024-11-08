package services

// APIClient 定义API客户端接口
type APIClient interface {
	Post(data interface{}, apiPath string, result interface{}) error
	GetAPIAccount() int64
}
