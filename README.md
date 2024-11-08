# JT Express Go SDK

极兔快递开放平台 API 的 Go 语言封装，提供轨迹查询等功能。

## 功能特性

- 轨迹查询
- 物流订阅（支持单个和批量）
- 完整的错误处理
- 类型安全的 API 封装
- 支持自定义配置

## 安装

```bash
go get github.com/maxbetas/jtexpress
```

## 快速开始

### 初始化客户端

```go
client := jtexpress.NewClient(
    "your_api_account",    // API账号
    "your_private_key",    // API私钥
)
```

### 配置选项
可以通过配置选项自定义客户端行为：

```go
client := jtexpress.NewClient(
    "your_api_account",
    "your_private_key",
    jtexpress.WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
    jtexpress.WithBaseURL("https://custom-api-url.com"),
)
```

## API 文档

### 轨迹查询

支持单个和批量运单查询：

```go
// 单个运单查询
resp, err := client.Logistics.QueryTrack("JT2099306666983")
if err != nil {
    log.Printf("查询失败: %v\n", err)
    return
}

// 多个运单查询（用英文逗号分隔）
resp, err := client.Logistics.QueryTrack("JT2099306666983,JT2099306666984,JT2099306666985")
if err != nil {
    log.Printf("查询失败: %v\n", err)
    return
}

// 处理响应
if resp.Success {
    for _, track := range resp.Data {
        fmt.Printf("运单号: %s\n", track.BillCode)
        for _, point := range track.TrackPoints {
            fmt.Printf("时间：%s\n", point.ScanTime)
            fmt.Printf("类型：%s\n", point.ScanType)
            fmt.Printf("位置：%s\n", point.Location)
            fmt.Printf("描述：%s\n", point.Description)
        }
    }
}
```

注意事项：
- 多个运单号使用英文逗号分隔
- 建议单次查询不超过 10 个运单号
- 响应中的轨迹点按时间倒序排列

### 物流订阅

支持单个和批量运单订阅：

```go
// 单个运单订阅
resp, err := client.Logistics.Subscribe("JT2099306666983", traceNode, backUrl)

// 批量运单订阅
billCodes := []string{"JT2099306666983", "JT2099306666984"}
resp, err := client.Logistics.Subscribe(billCodes, traceNode, backUrl)
```

#### 订阅节点说明
| 节点编号 | 说明 |
|---------|------|
| 1 | 快件揽收 |
| 2 | 入仓扫描（停用）|
| 3 | 发件扫描 |
| 4 | 到件扫描 |
| 5 | 出仓扫描 |
| 6 | 入库扫描 |
| 7 | 代理点收入扫描 |
| 8 | 快件取出扫描 |
| 9 | 出库扫描 |
| 10 | 快件签收 |
| 11 | 问题件扫描 |
| 12 | 安检扫描 |
| 13 | 其他扫描 |
| 14 | 退件扫描 |

## 错误处理

SDK 提供了详细的错误信息：

```go
resp, err := client.Logistics.QueryTrack(billCode)
if err != nil {
    switch e := err.(type) {
    case *errors.APIError:
        fmt.Printf("API错误: [%s] %s\n", e.Code, e.Message)
    default:
        fmt.Printf("其他错误: %v\n", err)
    }
}
```

## 最佳实践

### 1. 配置建议
- 设置合适的超时时间
- 在生产环境使用 HTTPS
- 妥善保管 API 私钥

### 2. 错误处理建议
- 对所有 API 调用进行错误处理
- 实现错误重试机制
- 记录错误日志

### 3. 性能优化建议
- 使用批量接口处理多个运单
- 合理使用缓存机制
- 控制并发请求数量

## 更新日志

### v1.0.0
- 初始版本发布
- 支持轨迹查询
- 支持物流订阅

## 许可证

MIT License 