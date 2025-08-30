package main

import (
	"fmt"
	"sync"
)

func square(n int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(n * n)
}

func main() {
	a := []int{2, 4, 6, 8, 10}

	var wg sync.WaitGroup

	for _, num := range a {
		wg.Add(1)
		go square(num, &wg)
	}

	wg.Wait()
}
