package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	ch1 := make(chan int)
	ch2 := make(chan int)

	a := make([]int, 0)
	for i := 0; i < 10; i++ {
		a = append(a, i)
	}

	fmt.Println(a)

	//writer to ch1 goroutine
	wg.Go(func() {
		for _, val := range a {
			ch1 <- val
			fmt.Printf("\n%d added to ch1", val)
		}
		close(ch1)
	})

	wg.Go(func() {
		for val := range ch1 {
			ch2 <- val * 2
			fmt.Printf("\n%d added to ch2", val*2)
		}
		close(ch2)
	})
	//writer to ch2 goroutine

	//worker goroutine
	wg.Go(func() {
		for val := range ch2 {
			fmt.Println("\nResult:", val)
		}
	})
	wg.Wait()
}
