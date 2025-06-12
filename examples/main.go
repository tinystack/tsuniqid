package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/tinystack/tsuniqid"
)

func main() {
	fmt.Println("=== Tsuniqid Package - 全方法调用示例 ===")
	fmt.Println("展示 tsuniqid 包的所有主要方法和功能")
	fmt.Println()

	// 1. 包级别函数调用
	fmt.Println("1. 📦 包级别函数调用")
	packageLevelFunctions()

	fmt.Println()

	// 2. IDGenerator 实例方法调用
	fmt.Println("2. 🏭 IDGenerator 实例方法调用")
	generatorInstanceMethods()

	fmt.Println()

	// 3. 多个生成器实例
	fmt.Println("3. 🔀 多个生成器实例")
	multipleGenerators()

	fmt.Println()

	// 4. ID 位布局分析
	fmt.Println("4. 🔍 ID 位布局分析")
	bitLayoutAnalysis()

	fmt.Println()

	// 5. 并发安全性测试
	fmt.Println("5. 🚀 并发安全性测试")
	concurrencySafety()

	fmt.Println()

	// 6. 性能基准测试
	fmt.Println("6. ⚡ 性能基准测试")
	performanceBenchmarks()

	fmt.Println()

	// 7. ID 格式和特性验证
	fmt.Println("7. ✅ ID 格式和特性验证")
	idFormatValidation()

	fmt.Println()
	fmt.Println("✅ 所有示例执行完成！")
}

// packageLevelFunctions 演示包级别的函数调用
func packageLevelFunctions() {
	fmt.Println("   tsuniqid.UniqID() - 生成字符串 ID:")
	for i := 0; i < 5; i++ {
		id := tsuniqid.UniqID()
		fmt.Printf("     [%d] %s (长度: %d)\n", i+1, id, len(id))
	}

	fmt.Println()
	fmt.Println("   tsuniqid.UniqUID() - 生成 uint64 ID:")
	for i := 0; i < 5; i++ {
		id := tsuniqid.UniqUID()
		fmt.Printf("     [%d] %d (十六进制: 0x%016x)\n", i+1, id, id)
	}
}

// generatorInstanceMethods 演示 IDGenerator 实例的所有方法
func generatorInstanceMethods() {
	// 创建新的生成器实例
	fmt.Println("   tsuniqid.NewGenerator() - 创建新生成器:")
	generator := tsuniqid.NewGenerator()
	fmt.Printf("     ✓ 生成器实例创建成功\n")

	fmt.Println()
	fmt.Println("   generator.GenerateStringID() - 实例方法生成字符串 ID:")
	for i := 0; i < 5; i++ {
		id := generator.GenerateStringID()
		fmt.Printf("     [%d] %s\n", i+1, id)
	}

	fmt.Println()
	fmt.Println("   generator.GenerateUint64ID() - 实例方法生成 uint64 ID:")
	for i := 0; i < 5; i++ {
		id := generator.GenerateUint64ID()
		fmt.Printf("     [%d] %d (0x%016x)\n", i+1, id, id)
	}
}

// multipleGenerators 演示多个生成器实例的独立性
func multipleGenerators() {
	const numGenerators = 3
	generators := make([]*tsuniqid.IDGenerator, numGenerators)

	fmt.Printf("   创建 %d 个独立的生成器实例:\n", numGenerators)
	for i := 0; i < numGenerators; i++ {
		generators[i] = tsuniqid.NewGenerator()
		fmt.Printf("     生成器 %d: 创建成功\n", i+1)
	}

	fmt.Println()
	fmt.Println("   各生成器产生的 ID 样本:")
	for i, gen := range generators {
		fmt.Printf("     生成器 %d:\n", i+1)
		for j := 0; j < 3; j++ {
			stringID := gen.GenerateStringID()
			uint64ID := gen.GenerateUint64ID()
			fmt.Printf("       字符串: %s\n", stringID)
			fmt.Printf("       数字:   %d (0x%016x)\n", uint64ID, uint64ID)
		}
		fmt.Println()
	}
}

// bitLayoutAnalysis 分析 uint64 ID 的位布局
func bitLayoutAnalysis() {
	generator := tsuniqid.NewGenerator()

	fmt.Println("   uint64 ID 位布局分析 (64位总计):")
	fmt.Println("   位 63-60 (4位): 机器ID")
	fmt.Println("   位 59-56 (4位): 实例ID")
	fmt.Println("   位 55-14 (42位): 时间戳 (毫秒)")
	fmt.Println("   位 13-0 (14位): 计数器")
	fmt.Println()

	for i := 0; i < 5; i++ {
		id := generator.GenerateUint64ID()

		// 根据 uniqid.go 中的位布局提取组件
		machineID := (id >> 60) & 0xF           // 前4位
		instanceID := (id >> 56) & 0xF          // 接下来4位
		timestamp := (id >> 14) & 0x3FFFFFFFFFF // 接下来42位
		counter := id & 0x3FFF                  // 最后14位

		fmt.Printf("   ID %d: %d (0x%016x)\n", i+1, id, id)
		fmt.Printf("     机器ID:  %d (二进制: %04b)\n", machineID, machineID)
		fmt.Printf("     实例ID:  %d (二进制: %04b)\n", instanceID, instanceID)
		fmt.Printf("     时间戳:  %d (时间: %s)\n", timestamp, time.UnixMilli(int64(timestamp)).Format("2006-01-02 15:04:05.000"))
		fmt.Printf("     计数器:  %d (二进制: %014b)\n", counter, counter)
		fmt.Println()
	}
}

// concurrencySafety 测试并发安全性
func concurrencySafety() {
	const numGoroutines = 10
	const idsPerGoroutine = 1000

	var wg sync.WaitGroup
	var mu sync.Mutex
	allIDs := make(map[string]bool)

	fmt.Printf("   启动 %d 个协程，每个生成 %d 个ID...\n", numGoroutines, idsPerGoroutine)

	start := time.Now()

	// 测试包级别函数的并发安全性
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			localIDs := make([]string, 0, idsPerGoroutine)

			// 生成ID
			for j := 0; j < idsPerGoroutine; j++ {
				id := tsuniqid.UniqID()
				localIDs = append(localIDs, id)
			}

			// 添加到全局集合
			mu.Lock()
			for _, id := range localIDs {
				allIDs[id] = true
			}
			mu.Unlock()

			fmt.Printf("     协程 %d 完成: 生成 %d 个ID\n", goroutineID, len(localIDs))
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	expectedTotal := numGoroutines * idsPerGoroutine
	actualTotal := len(allIDs)

	fmt.Printf("   结果: 生成 %d 个ID，唯一性 %d 个\n", expectedTotal, actualTotal)
	fmt.Printf("   耗时: %v\n", duration)
	fmt.Printf("   性能: %.2f IDs/秒\n", float64(expectedTotal)/duration.Seconds())

	if actualTotal == expectedTotal {
		fmt.Printf("   ✅ 所有ID都是唯一的！\n")
	} else {
		fmt.Printf("   ❌ 发现 %d 个重复ID！\n", expectedTotal-actualTotal)
	}
}

// performanceBenchmarks 性能基准测试
func performanceBenchmarks() {
	const iterations = 100000

	fmt.Printf("   基准测试 - 各方法生成 %d 个ID的性能:\n", iterations)

	// 测试 UniqID()
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = tsuniqid.UniqID()
	}
	uniqIDDuration := time.Since(start)

	// 测试 UniqUID()
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = tsuniqid.UniqUID()
	}
	uniqUIDDuration := time.Since(start)

	// 测试生成器实例
	generator := tsuniqid.NewGenerator()

	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = generator.GenerateStringID()
	}
	genStringDuration := time.Since(start)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = generator.GenerateUint64ID()
	}
	genUint64Duration := time.Since(start)

	// 输出结果
	fmt.Printf("   tsuniqid.UniqID():               %v (%.2f ns/op)\n",
		uniqIDDuration, float64(uniqIDDuration.Nanoseconds())/float64(iterations))
	fmt.Printf("   tsuniqid.UniqUID():              %v (%.2f ns/op)\n",
		uniqUIDDuration, float64(uniqUIDDuration.Nanoseconds())/float64(iterations))
	fmt.Printf("   generator.GenerateStringID():    %v (%.2f ns/op)\n",
		genStringDuration, float64(genStringDuration.Nanoseconds())/float64(iterations))
	fmt.Printf("   generator.GenerateUint64ID():    %v (%.2f ns/op)\n",
		genUint64Duration, float64(genUint64Duration.Nanoseconds())/float64(iterations))
}

// idFormatValidation 验证ID格式和特性
func idFormatValidation() {
	generator := tsuniqid.NewGenerator()

	fmt.Println("   ID 格式验证:")

	// 验证字符串ID格式
	fmt.Println()
	fmt.Println("   字符串ID格式验证:")
	for i := 0; i < 5; i++ {
		id := generator.GenerateStringID()

		// 检查长度
		fmt.Printf("     ID: %s\n", id)
		fmt.Printf("       长度: %d 字符\n", len(id))

		// 验证前缀是否为有效的十六进制
		suffixLength := 8 // RandomSuffixLength 常量值
		if len(id) >= suffixLength {
			hexPart := id[:len(id)-suffixLength]
			suffix := id[len(id)-suffixLength:]

			if _, err := strconv.ParseUint(hexPart, 16, 64); err != nil {
				fmt.Printf("       ❌ 十六进制部分无效: %s\n", hexPart)
			} else {
				fmt.Printf("       ✅ 十六进制部分有效: %s\n", hexPart)
			}

			fmt.Printf("       随机后缀: %s\n", suffix)
		}
		fmt.Println()
	}

	// 验证uint64 ID的时间戳合理性
	fmt.Println("   uint64 ID 时间戳验证:")
	now := time.Now().UnixMilli()

	for i := 0; i < 5; i++ {
		id := generator.GenerateUint64ID()
		timestamp := (id >> 14) & 0x3FFFFFFFFFF

		timeDiff := int64(timestamp) - now
		timeObj := time.UnixMilli(int64(timestamp))

		fmt.Printf("     ID: %d\n", id)
		fmt.Printf("       时间戳: %d\n", timestamp)
		fmt.Printf("       时间: %s\n", timeObj.Format("2006-01-02 15:04:05.000"))
		fmt.Printf("       与当前时间差: %d 毫秒\n", timeDiff)

		if timeDiff >= -1000 && timeDiff <= 1000 {
			fmt.Printf("       ✅ 时间戳合理\n")
		} else {
			fmt.Printf("       ❌ 时间戳异常\n")
		}
		fmt.Println()
	}
}
