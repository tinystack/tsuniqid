# tsuniqid

[![English](https://img.shields.io/badge/English-blue)](README.md) [![中文](https://img.shields.io/badge/中文-red)](README_CN.md)

A high-performance unique ID generator for Go that provides both string and uint64 type unique identifiers with excellent concurrency safety and performance characteristics.

[![Go Report Card](https://goreportcard.com/badge/github.com/tinystack/tsuniqid)](https://goreportcard.com/report/github.com/tinystack/tsuniqid)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.18-61CFDD.svg?style=flat-square)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/tinystack/tsuniqid)](https://pkg.go.dev/mod/github.com/tinystack/tsuniqid)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **🚀 High Performance**: String ID ~443 ns/op, uint64 ID ~24 ns/op
- **🔒 Thread Safe**: Fully concurrent safe with atomic operations
- **🎯 Guaranteed Uniqueness**: Tested with 1M+ concurrent IDs with zero duplicates
- **📦 Multiple Formats**: Both string and uint64 ID generation
- **🏭 Multiple Instances**: Support for multiple independent generators
- **🌐 Machine Aware**: Incorporates machine ID for distributed environments
- **⚡ Zero Dependencies**: Pure Go implementation with standard library only

## Installation

```bash
go get -u github.com/tinystack/tsuniqid
```

## Quick Start

### Package-level Functions (Recommended)

```go
package main

import (
    "fmt"
    "github.com/tinystack/tsuniqid"
)

func main() {
    // Generate string ID
    stringID := tsuniqid.UniqID()
    fmt.Println("String ID:", stringID) // e.g., "1a2b3c4d5e6f78901a2b3c4d"

    // Generate uint64 ID
    uint64ID := tsuniqid.UniqUID()
    fmt.Println("Uint64 ID:", uint64ID) // e.g., 1844674407370955161
}
```

### Generator Instances

```go
package main

import (
    "fmt"
    "github.com/tinystack/tsuniqid"
)

func main() {
    // Create independent generator instances
    gen1 := tsuniqid.NewGenerator()
    gen2 := tsuniqid.NewGenerator()

    // Generate IDs from different instances
    id1 := gen1.GenerateStringID()
    id2 := gen2.GenerateUint64ID()

    fmt.Println("Generator 1 String ID:", id1)
    fmt.Println("Generator 2 Uint64 ID:", id2)
}
```

## API Reference

### Package Functions

| Function             | Description               | Return Type | Performance |
| -------------------- | ------------------------- | ----------- | ----------- |
| `tsuniqid.UniqID()`  | Generate unique string ID | `string`    | ~443 ns/op  |
| `tsuniqid.UniqUID()` | Generate unique uint64 ID | `uint64`    | ~24 ns/op   |

### Generator Methods

| Method               | Description                      | Return Type    |
| -------------------- | -------------------------------- | -------------- |
| `NewGenerator()`     | Create new generator instance    | `*IDGenerator` |
| `GenerateStringID()` | Generate string ID from instance | `string`       |
| `GenerateUint64ID()` | Generate uint64 ID from instance | `uint64`       |

## ID Structure

### String ID Format

- **Format**: `{hex_uint64_id}{random_suffix}`
- **Length**: 24 characters (16-char hex + 8-char random suffix)
- **Example**: `"1a2b3c4d5e6f78901a2b3c4d"`

### Uint64 ID Bit Layout (64 bits total)

```
┌─────────────┬─────────────┬──────────────────────────────────────────┬────────────────┐
│   Bits      │    Size     │                Description                │     Range      │
├─────────────┼─────────────┼──────────────────────────────────────────┼────────────────┤
│   63-60     │   4 bits    │            Machine ID                    │      0-15      │
│   59-56     │   4 bits    │           Instance ID                    │      0-15      │
│   55-14     │   42 bits   │     Timestamp (milliseconds)            │   0-4398046511103 │
│   13-0      │   14 bits   │              Counter                     │     0-16383    │
└─────────────┴─────────────┴──────────────────────────────────────────┴────────────────┘
```

## Advanced Usage

### Concurrent Generation

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

    // Launch concurrent ID generation
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
    fmt.Printf("Generated %d unique IDs\n", len(uniqueIDs))
}
```

### ID Component Analysis

```go
package main

import (
    "fmt"
    "time"
    "github.com/tinystack/tsuniqid"
)

func main() {
    id := tsuniqid.UniqUID()

    // Extract components
    machineID := (id >> 60) & 0xF
    instanceID := (id >> 56) & 0xF
    timestamp := (id >> 14) & 0x3FFFFFFFFFF
    counter := id & 0x3FFF

    fmt.Printf("ID: %d (0x%016x)\n", id, id)
    fmt.Printf("Machine ID: %d\n", machineID)
    fmt.Printf("Instance ID: %d\n", instanceID)
    fmt.Printf("Timestamp: %d (%s)\n", timestamp,
        time.UnixMilli(int64(timestamp)).Format("2006-01-02 15:04:05.000"))
    fmt.Printf("Counter: %d\n", counter)
}
```

## Performance Benchmarks

### Benchmark Results

```
BenchmarkUniqID-12                           408.9 ns/op
BenchmarkUniqUID-12                           27.43 ns/op
BenchmarkIDGenerator_GenerateStringID        397.6 ns/op
BenchmarkIDGenerator_GenerateUint64ID         24.85 ns/op
```

### Performance Characteristics

- **String ID Generation**: ~2.4M operations/second
- **Uint64 ID Generation**: ~36M operations/second
- **Memory Allocation**: Minimal heap allocation
- **Concurrency**: Linear scaling with CPU cores

## Testing and Quality

### Test Coverage

- ✅ **Uniqueness Tests**: 1M+ concurrent IDs with zero duplicates
- ✅ **Format Validation**: String ID format and uint64 bit layout
- ✅ **Concurrency Safety**: Multi-goroutine stress testing
- ✅ **Component Verification**: Machine ID, timestamp, counter ranges
- ✅ **Multiple Instances**: Independent generator isolation
- ✅ **Performance Benchmarks**: Comprehensive performance testing

### Running Tests

```bash
# Run all tests
go test -v

# Run benchmarks
go test -bench=. -benchmem

# Test with race detection
go test -race -v

# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Use Cases

- **🔗 Distributed Systems**: Machine-aware IDs for distributed environments
- **📊 Database Records**: Primary keys with embedded timestamp information
- **🌐 Web Applications**: Request IDs, session tokens, API keys
- **📝 Logging Systems**: Trace IDs for distributed tracing
- **🔄 Message Queues**: Message identifiers with ordering information
- **📱 Microservices**: Service instance identification and correlation

## Architecture

### Design Principles

- **High Performance**: Optimized for speed with minimal allocations
- **Thread Safety**: Lock-free atomic operations where possible
- **Uniqueness Guarantee**: Mathematical guarantee of ID uniqueness
- **Flexibility**: Support for multiple formats and generator instances
- **Simplicity**: Clean API with minimal learning curve

### Machine ID Generation

- Based on hostname and local IP address
- SHA1 hash for deterministic generation
- Fallback to random generation if network info unavailable
- 4-bit machine ID supports up to 16 machines

### Instance ID System

- Atomic counter for unique instance identification
- 4-bit instance ID supports up to 16 generators per machine
- Prevents collisions between multiple generator instances

## Examples

Check out the comprehensive examples in the [`examples/`](examples/) directory:

- **Basic Usage**: Simple ID generation examples
- **Concurrent Usage**: Multi-goroutine safe ID generation
- **Web Server**: HTTP service with unique request IDs
- **Microservice**: Distributed service with request tracing
- **Data Storage**: Database-like operations with unique keys

Run examples:

```bash
cd examples
go run main.go
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Requirements

- Go 1.18 or higher
- No external dependencies

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Made with ❤️ for the Go community**
