package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func sleep(duration time.Duration) {
	start := time.Now()
	deadline := start.Add(duration)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	sleep(5 * time.Second)
	fmt.Println("ive waited for 5 seconds")

	fmt.Println(runtime.GOMAXPROCS(0))
}
