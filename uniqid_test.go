package tsuniqid

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

// TestUniqID_Uniqueness tests that UniqID generates unique string identifiers
// across multiple goroutines to ensure thread safety and uniqueness.
func TestUniqID_Uniqueness(t *testing.T) {
	const (
		numGoroutines      = 10
		numIDsPerGoroutine = 100000
	)

	var (
		counter = make(map[string]int)
		mu      sync.Mutex
		wg      sync.WaitGroup
	)

	// Launch multiple goroutines to generate IDs concurrently
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < numIDsPerGoroutine; i++ {
				id := UniqID()

				mu.Lock()
				counter[id]++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Verify all generated IDs are unique
	totalIDs := 0
	duplicates := 0
	for id, count := range counter {
		totalIDs++
		if count > 1 {
			duplicates++
			t.Errorf("Duplicate ID found: %s appeared %d times", id, count)
		}
	}

	expectedTotal := numGoroutines * numIDsPerGoroutine
	if totalIDs != expectedTotal {
		t.Errorf("Expected %d unique IDs, got %d", expectedTotal, totalIDs)
	}

	t.Logf("Generated %d unique string IDs with %d duplicates", totalIDs, duplicates)
}

// TestUniqUID_Uniqueness tests that UniqUID generates unique uint64 identifiers
// across multiple goroutines to ensure thread safety and uniqueness.
func TestUniqUID_Uniqueness(t *testing.T) {
	const (
		numGoroutines      = 10
		numIDsPerGoroutine = 100000
	)

	var (
		counter = make(map[uint64]int)
		mu      sync.Mutex
		wg      sync.WaitGroup
	)

	// Launch multiple goroutines to generate IDs concurrently
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < numIDsPerGoroutine; i++ {
				id := UniqUID()

				mu.Lock()
				counter[id]++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Verify all generated IDs are unique
	totalIDs := 0
	duplicates := 0
	for id, count := range counter {
		totalIDs++
		if count > 1 {
			duplicates++
			t.Errorf("Duplicate ID found: %d appeared %d times", id, count)
		}
	}

	expectedTotal := numGoroutines * numIDsPerGoroutine
	if totalIDs != expectedTotal {
		t.Errorf("Expected %d unique IDs, got %d", expectedTotal, totalIDs)
	}

	t.Logf("Generated %d unique uint64 IDs with %d duplicates", totalIDs, duplicates)
}

// TestIDGenerator_StringIDFormat tests the format of generated string IDs.
func TestIDGenerator_StringIDFormat(t *testing.T) {
	gen := NewGenerator()

	for i := 0; i < 100; i++ {
		id := gen.GenerateStringID()

		// Check minimum length (hex uint64 + random suffix)
		if len(id) < 16+RandomSuffixLength {
			t.Errorf("String ID too short: %s (length: %d)", id, len(id))
		}

		// Check that the ID contains valid hexadecimal characters and suffix
		hexPart := id[:len(id)-RandomSuffixLength]
		if _, err := strconv.ParseUint(hexPart, 16, 64); err != nil {
			t.Errorf("Invalid hex part in ID %s: %v", id, err)
		}
	}
}

// TestIDGenerator_Uint64IDComponents tests the bit layout of uint64 IDs.
func TestIDGenerator_Uint64IDComponents(t *testing.T) {
	gen := NewGenerator()

	// Generate multiple IDs and verify component extraction
	for i := 0; i < 100; i++ {
		id := gen.GenerateUint64ID()

		// Extract components
		machineID := (id >> MachineIDShift) & MaxMachineID
		instanceID := (id >> InstanceIDShift) & MaxInstanceID
		timestamp := (id >> TimestampShift) & MaxTimestamp
		counter := id & MaxCounter

		// Verify machine ID is within valid range
		if machineID > MaxMachineID {
			t.Errorf("Machine ID out of range: %d > %d", machineID, MaxMachineID)
		}

		// Verify instance ID is within valid range
		if instanceID > MaxInstanceID {
			t.Errorf("Instance ID out of range: %d > %d", instanceID, MaxInstanceID)
		}

		// Verify timestamp is reasonable (within last few seconds and next few seconds)
		now := uint64(time.Now().UnixMilli())
		tolerance := uint64(5000) // 5 seconds tolerance
		if timestamp < now-tolerance || timestamp > now+tolerance {
			t.Errorf("Timestamp out of reasonable range: %d (now: %d, diff: %d)", timestamp, now, int64(timestamp)-int64(now))
		}

		// Verify counter is within valid range
		if counter > MaxCounter {
			t.Errorf("Counter out of range: %d > %d", counter, MaxCounter)
		}
	}
}

// TestIDGenerator_MultipleInstances tests that multiple generator instances
// can work independently without conflicts.
func TestIDGenerator_MultipleInstances(t *testing.T) {
	const numGenerators = 5
	const numIDsPerGenerator = 1000

	generators := make([]*IDGenerator, numGenerators)
	for i := 0; i < numGenerators; i++ {
		generators[i] = NewGenerator()
	}

	allIDs := make(map[uint64]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Generate IDs from multiple generators concurrently
	for i, gen := range generators {
		wg.Add(1)
		go func(genIndex int, generator *IDGenerator) {
			defer wg.Done()

			for j := 0; j < numIDsPerGenerator; j++ {
				id := generator.GenerateUint64ID()

				mu.Lock()
				if allIDs[id] {
					t.Errorf("Duplicate ID %d from generator %d", id, genIndex)
				}
				allIDs[id] = true
				mu.Unlock()
			}
		}(i, gen)
	}

	wg.Wait()

	expectedTotal := numGenerators * numIDsPerGenerator
	if len(allIDs) != expectedTotal {
		t.Errorf("Expected %d unique IDs, got %d", expectedTotal, len(allIDs))
	}
}

// TestIDGenerator_CounterIncrement tests that the counter increments properly
// and handles overflow correctly.
func TestIDGenerator_CounterIncrement(t *testing.T) {
	gen := NewGenerator()

	// Generate several IDs and check counter progression
	var lastCounter uint64
	for i := 0; i < 10; i++ {
		id := gen.GenerateUint64ID()
		counter := id & MaxCounter

		if i > 0 && counter <= lastCounter {
			// Allow for counter overflow
			if counter == 0 && lastCounter == MaxCounter {
				t.Logf("Counter overflow detected at iteration %d", i)
			} else {
				t.Errorf("Counter not incrementing properly: %d -> %d", lastCounter, counter)
			}
		}
		lastCounter = counter
	}
}

// TestIDGenerator_RandomSuffixVariety tests that random suffixes are diverse.
func TestIDGenerator_RandomSuffixVariety(t *testing.T) {
	gen := NewGenerator()
	suffixes := make(map[string]int)

	// Generate many string IDs and collect suffixes
	for i := 0; i < 10000; i++ {
		id := gen.GenerateStringID()
		if len(id) >= RandomSuffixLength {
			suffix := id[len(id)-RandomSuffixLength:]
			suffixes[suffix]++
		}
	}

	// Check that we have good variety in suffixes
	uniqueSuffixes := len(suffixes)
	minExpectedUnique := 9000 // Should have high variety

	if uniqueSuffixes < minExpectedUnique {
		t.Errorf("Insufficient suffix variety: got %d unique suffixes, expected at least %d",
			uniqueSuffixes, minExpectedUnique)
	}

	t.Logf("Generated %d unique suffixes out of 10000 IDs", uniqueSuffixes)
}

// BenchmarkUniqID benchmarks the performance of string ID generation.
func BenchmarkUniqID(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = UniqID()
		}
	})
}

// BenchmarkUniqUID benchmarks the performance of uint64 ID generation.
func BenchmarkUniqUID(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = UniqUID()
		}
	})
}

// BenchmarkIDGenerator_GenerateStringID benchmarks the performance of
// string ID generation using a specific generator instance.
func BenchmarkIDGenerator_GenerateStringID(b *testing.B) {
	gen := NewGenerator()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = gen.GenerateStringID()
		}
	})
}

// BenchmarkIDGenerator_GenerateUint64ID benchmarks the performance of
// uint64 ID generation using a specific generator instance.
func BenchmarkIDGenerator_GenerateUint64ID(b *testing.B) {
	gen := NewGenerator()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = gen.GenerateUint64ID()
		}
	})
}
