package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)
	done := make(chan struct{})

	go func() {
		i := 0
		for {
			// вроде как хочу __правильно__ закончить работу горутины
			select {
			case <-done:
				fmt.Println("goroutine ended")
				return
			default:
				ch <- i
				i++
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	timeout := time.After(5 * time.Second)

	for {
		select {
		case <-timeout:
			fmt.Println("your time has passed")
			return
		case val := <-ch:
			fmt.Printf("got message %d\n", val)
		}
	}
}
