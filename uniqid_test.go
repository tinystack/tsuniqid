package tsuniqid

import (
	"sync"
	"testing"
)

func TestUniqID(t *testing.T) {
	var (
		goroutine = 10
		counter   = make(map[string]int)
		mu        sync.Mutex
		wg        sync.WaitGroup
	)

	for g := 0; g < goroutine; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			maxTimes := 100000
			for maxTimes > 0 {
				uniq := UniqID()
				mu.Lock()
				_, ok := counter[uniq]
				if !ok {
					counter[uniq] = 0
				}
				counter[uniq]++
				mu.Unlock()
				maxTimes--
			}
		}()
	}

	wg.Wait()

	for uniq, num := range counter {
		if num > 1 {
			t.Errorf("uniq = %s, maxTimes = %d", uniq, num)
		}
	}
}

func TestUniqUID(t *testing.T) {
	var (
		goroutine = 10
		counter   = make(map[uint64]int)
		mu        sync.Mutex
		wg        sync.WaitGroup
	)

	for g := 0; g < goroutine; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			maxTimes := 100000
			for maxTimes > 0 {
				uniq := UniqUID()
				mu.Lock()
				_, ok := counter[uniq]
				if !ok {
					counter[uniq] = 0
				}
				counter[uniq]++
				mu.Unlock()
				maxTimes--
			}
		}()
	}

	wg.Wait()

	for uniq, num := range counter {
		if num > 1 {
			t.Errorf("uniq = %d, maxTimes = %d", uniq, num)
		}
	}
}
