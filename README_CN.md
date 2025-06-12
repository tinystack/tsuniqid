# tsuniqid

[English](./README.md) | [中文文档](./README_CN.md)

一个高性能的 Go 语言唯一 ID 生成器，提供字符串和 uint64 类型的唯一标识符，具有出色的并发安全性和性能特征。

[![Go Report Card](https://goreportcard.com/badge/github.com/tinystack/tsuniqid)](https://goreportcard.com/report/github.com/tinystack/tsuniqid)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.18-61CFDD.svg?style=flat-square)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/tinystack/tsuniqid)](https://pkg.go.dev/mod/github.com/tinystack/tsuniqid)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 特性

- **🚀 高性能**: 字符串 ID ~443 ns/op，uint64 ID ~24 ns/op
- **🔒 线程安全**: 基于原子操作的完全并发安全
- **🎯 唯一性保证**: 经过 100 万+ 并发 ID 测试，零重复
- **📦 多种格式**: 支持字符串和 uint64 ID 生成
- **🏭 多实例支持**: 支持多个独立的生成器实例
- **🌐 机器感知**: 为分布式环境集成机器 ID
- **⚡ 零依赖**: 纯 Go 实现，仅使用标准库

## 安装

```bash
go get -u github.com/tinystack/tsuniqid
```

## 快速开始

### 包级别函数（推荐）

```go
package main

import (
    "fmt"
    "github.com/tinystack/tsuniqid"
)

func main() {
    // 生成字符串 ID
    stringID := tsuniqid.UniqID()
    fmt.Println("字符串 ID:", stringID) // 例如: "1a2b3c4d5e6f78901a2b3c4d"

    // 生成 uint64 ID
    uint64ID := tsuniqid.UniqUID()
    fmt.Println("Uint64 ID:", uint64ID) // 例如: 1844674407370955161
}
```

### 生成器实例

```go
package main

import (
    "fmt"
    "github.com/tinystack/tsuniqid"
)

func main() {
    // 创建独立的生成器实例
    gen1 := tsuniqid.NewGenerator()
    gen2 := tsuniqid.NewGenerator()

    // 从不同实例生成 ID
    id1 := gen1.GenerateStringID()
    id2 := gen2.GenerateUint64ID()

    fmt.Println("生成器 1 字符串 ID:", id1)
    fmt.Println("生成器 2 Uint64 ID:", id2)
}
```

## API 参考

### 包函数

| 函数                 | 描述               | 返回类型 | 性能       |
| -------------------- | ------------------ | -------- | ---------- |
| `tsuniqid.UniqID()`  | 生成唯一字符串 ID  | `string` | ~443 ns/op |
| `tsuniqid.UniqUID()` | 生成唯一 uint64 ID | `uint64` | ~24 ns/op  |

### 生成器方法

| 方法                 | 描述                 | 返回类型       |
| -------------------- | -------------------- | -------------- |
| `NewGenerator()`     | 创建新的生成器实例   | `*IDGenerator` |
| `GenerateStringID()` | 从实例生成字符串 ID  | `string`       |
| `GenerateUint64ID()` | 从实例生成 uint64 ID | `uint64`       |

## ID 结构

### 字符串 ID 格式

- **格式**: `{十六进制_uint64_id}{随机后缀}`
- **长度**: 24 个字符（16 字符十六进制 + 8 字符随机后缀）
- **示例**: `"1a2b3c4d5e6f78901a2b3c4d"`

### Uint64 ID 位布局（总共 64 位）

```
┌─────────────┬─────────────┬──────────────────────────────────────────┬────────────────┐
│   位范围    │    大小     │                 描述                     │      范围      │
├─────────────┼─────────────┼──────────────────────────────────────────┼────────────────┤
│   63-60     │   4 位      │              机器 ID                     │      0-15      │
│   59-56     │   4 位      │             实例 ID                      │      0-15      │
│   55-14     │   42 位     │          时间戳（毫秒）                 │   0-4398046511103 │
│   13-0      │   14 位     │               计数器                     │     0-16383    │
└─────────────┴─────────────┴──────────────────────────────────────────┴────────────────┘
```

## 高级用法

### 并发生成

```go
package main

import (
    "fmt"
    "sync"
    "github.com/tinystack/tsuniqid"
)

func main() {
    const numGoroutines = 100
    const idsPerGoroutine = 1000

    var wg sync.WaitGroup
    uniqueIDs := make(map[string]bool)
    var mu sync.Mutex

    // 启动并发 ID 生成
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()

            for j := 0; j < idsPerGoroutine; j++ {
                id := tsuniqid.UniqID()

                mu.Lock()
                uniqueIDs[id] = true
                mu.Unlock()
            }
        }()
    }

    wg.Wait()
    fmt.Printf("生成了 %d 个唯一 ID\n", len(uniqueIDs))
}
```

### ID 组件分析

```go
package main

import (
    "fmt"
    "time"
    "github.com/tinystack/tsuniqid"
)

func main() {
    id := tsuniqid.UniqUID()

    // 提取组件
    machineID := (id >> 60) & 0xF
    instanceID := (id >> 56) & 0xF
    timestamp := (id >> 14) & 0x3FFFFFFFFFF
    counter := id & 0x3FFF

    fmt.Printf("ID: %d (0x%016x)\n", id, id)
    fmt.Printf("机器 ID: %d\n", machineID)
    fmt.Printf("实例 ID: %d\n", instanceID)
    fmt.Printf("时间戳: %d (%s)\n", timestamp,
        time.UnixMilli(int64(timestamp)).Format("2006-01-02 15:04:05.000"))
    fmt.Printf("计数器: %d\n", counter)
}
```

## 性能基准测试

### 基准测试结果

```
BenchmarkUniqID-12                           408.9 ns/op
BenchmarkUniqUID-12                           27.43 ns/op
BenchmarkIDGenerator_GenerateStringID        397.6 ns/op
BenchmarkIDGenerator_GenerateUint64ID         24.85 ns/op
```

### 性能特征

- **字符串 ID 生成**: ~240 万次操作/秒
- **Uint64 ID 生成**: ~3600 万次操作/秒
- **内存分配**: 最小堆分配
- **并发性**: 随 CPU 核心数线性扩展

## 测试和质量

### 测试覆盖率

- ✅ **唯一性测试**: 100 万+ 并发 ID，零重复
- ✅ **格式验证**: 字符串 ID 格式和 uint64 位布局
- ✅ **并发安全**: 多协程压力测试
- ✅ **组件验证**: 机器 ID、时间戳、计数器范围
- ✅ **多实例**: 独立生成器隔离
- ✅ **性能基准**: 综合性能测试

### 运行测试

```bash
# 运行所有测试
go test -v

# 运行基准测试
go test -bench=. -benchmem

# 使用竞态检测运行测试
go test -race -v

# 生成覆盖率报告
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 使用场景

- **🔗 分布式系统**: 感知机器的分布式环境 ID
- **📊 数据库记录**: 嵌入时间戳信息的主键
- **🌐 Web 应用**: 请求 ID、会话令牌、API 密钥
- **📝 日志系统**: 分布式追踪的 Trace ID
- **🔄 消息队列**: 带有排序信息的消息标识符
- **📱 微服务**: 服务实例识别和关联

## 架构

### 设计原则

- **高性能**: 针对速度优化，最小化分配
- **线程安全**: 尽可能使用无锁原子操作
- **唯一性保证**: 数学保证 ID 唯一性
- **灵活性**: 支持多种格式和生成器实例
- **简洁性**: 清晰的 API，学习曲线最小

### 机器 ID 生成

- 基于主机名和本地 IP 地址
- SHA1 哈希用于确定性生成
- 如果网络信息不可用，回退到随机生成
- 4 位机器 ID 支持最多 16 台机器

### 实例 ID 系统

- 原子计数器用于唯一实例标识
- 4 位实例 ID 支持每台机器最多 16 个生成器
- 防止多个生成器实例之间的冲突

## 示例

查看 [`examples/`](examples/) 目录中的综合示例：

- **基本用法**: 简单的 ID 生成示例
- **并发用法**: 多协程安全 ID 生成
- **Web 服务器**: 带有唯一请求 ID 的 HTTP 服务
- **微服务**: 带有请求追踪的分布式服务
- **数据存储**: 带有唯一键的类数据库操作

运行示例：

```bash
cd examples
go run main.go
```

## 贡献

1. Fork 仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启 Pull Request

## 要求

- Go 1.18 或更高版本
- 无外部依赖

## 许可证

该项目根据 MIT 许可证授权 - 有关详细信息，请参阅 [LICENSE](LICENSE) 文件。

---

**用 ❤️ 为 Go 社区制作**
