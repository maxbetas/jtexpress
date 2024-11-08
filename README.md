# JT Express Go SDK

极兔快递开放平台 API 的 Go 语言封装，提供轨迹查询和订阅等功能。

## 功能特性

- 轨迹查询（支持单个和批量查询）
- 物流订阅（支持单个和批量订阅）
- 完整的错误处理
- 类型安全的 API 封装

## 安装

```bash
go get -u github.com/maxbetas/jtexpress@latest
```

## 快速开始

### 初始化客户端

```go
import "github.com/maxbetas/jtexpress"

client := jtexpress.NewClient(
    "your_api_account",    // API账号
    "your_private_key",    // API私钥
)
```

### 轨迹查询

```go
// 单个运单查询
resp, err := client.Logistics.QueryTrack("JT2099306666983")
if err != nil {
    log.Printf("查询失败: %v\n", err)
    return
}

// 批量运单查询（用英文逗号分隔）
resp, err := client.Logistics.QueryTrack("JT2099306666983,JT2099306666984")
if err != nil {
    log.Printf("查询失败: %v\n", err)
    return
}

// 处理响应
if resp.Code == "1" {
    for _, track := range resp.Data {
        fmt.Printf("运单号: %s\n", track.BillCode)
        for _, detail := range track.Details {
            fmt.Printf("时间：%s\n", detail.ScanTime)
            fmt.Printf("描述：%s\n", detail.Desc)
            fmt.Printf("类型：%s\n", detail.ScanType)
            fmt.Printf("地点：%s (%s)\n", detail.ScanNetworkName, detail.ScanNetworkProvince)
        }
    }
}
```

### 物流订阅

```go
// 单个运单订阅
resp, err := client.Logistics.Subscribe(
    "JT2099306666983",
    "1&2&3&4&5",                        // 订阅节点
    "https://your-domain.com/callback", // 回调地址
)
if err != nil {
    log.Printf("订阅失败: %v\n", err)
    return
}

// 批量运单订阅
resp, err := client.Logistics.SubscribeBatch(
    []string{"JT2099306666983", "JT2099306666984"},
    "1&2&3&4&5",                        // 订阅节点
    "https://your-domain.com/callback", // 回调地址
)
if err != nil {
    log.Printf("批量订阅失败: %v\n", err)
    return
}

// 处理响应
if resp.Code == "1" {
    fmt.Println("订阅成功")
} else {
    fmt.Printf("订阅失败: [%s] %s\n", resp.Code, resp.Msg)
}
```

### 订阅节点说明
| 节点编号 | 说明 |
|---------|------|
| 1 | 快件揽收 |
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
    // 处理错误
    fmt.Printf("错误: %v\n", err)
    return
}

// 检查业务响应
if resp.Code != "1" {
    fmt.Printf("业务错误: [%s] %s\n", resp.Code, resp.Msg)
    return
}
```

## 最佳实践

### 1. 错误处理建议
- 对所有 API 调用进行错误处理
- 记录详细的错误日志
- 实现错误重试机制

### 2. 订阅回调建议
- 确保回调地址可以正常访问
- 回调地址需要能处理 POST 请求
- 建议订阅必要的节点，避免无用推送

### 3. 性能优化建议
- 使用批量接口处理多个运单
- 控制并发请求数量
- 合理设置超时时间

## API 文档

完整的 API 文档请参考：[极兔开放平台文档](https://open.jtexpress.com.cn/)

## 许可证

MIT License 