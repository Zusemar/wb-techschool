package main

import (
	"fmt"
	"sync"
)

func main() {
	// через конкурентную мапу
	count := 1000
	concurrent := writeWithSyncMap(count)
	fmt.Printf("sync.Map size: %d\n", len(concurrent))

	// через мьютекс
	withMutex := writeWithMutexMap(count)
	fmt.Printf("mutex map size: %d\n", len(withMutex))
}

func writeWithSyncMap(n int) map[int]int {
	var concurrentMap sync.Map
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Go(func() {
			concurrentMap.Store(i, i)
		})
	}

	wg.Wait()

	result := make(map[int]int, n)
	concurrentMap.Range(func(key, value any) bool {
		k := key.(int)
		v := value.(int)
		result[k] = v
		return true
	})

	return result
}

func writeWithMutexMap(n int) map[int]int {
	result := make(map[int]int, n)
	var mu sync.RWMutex
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(value int) {
			defer wg.Done()
			mu.Lock()
			result[value] = value
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	return result
}
