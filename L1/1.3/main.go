package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run main.go <number of workers>")
		return
	}

	// получаем количество воркеров
	workersCount, err := strconv.Atoi(os.Args[1])
	if err != nil || workersCount <= 0 {
		fmt.Println("workers must be a positive integer")
		return
	}

	// создаем канал
	messages := make(chan int)
	var wg sync.WaitGroup

	// создаем воркеров
	for i := 1; i <= workersCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for msg := range messages {
				fmt.Printf("worker %d: %d\n", id, msg)
			}
		}(i)
	}

	// пишем в канал из главной горутины
	counter := 0
	for {
		counter++
		messages <- counter
		time.Sleep(200 * time.Millisecond)
	}
}
