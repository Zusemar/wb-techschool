package main

import (
	"fmt"
	"os"
	"os/signal"
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

	// создаем канал сообщений
	messages := make(chan int)
	var wg sync.WaitGroup

	// создаем канал для получения сигнала
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	// создаем канал для завершения
	done := make(chan struct{})

	// создаем горутину для получения сигнала
	go func() {
		sig := <-sigCh
		fmt.Println("received signal:", sig)
		close(done)
	}()

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

	// добавляем сообщения в канал из главной горутины, обрабатывая SIGINT в case <-done
	counter := 0
	for {
		select {
		case <-done:
			fmt.Println("done")
			close(messages)
			wg.Wait()
			return
		default:
			counter++
			messages <- counter
			time.Sleep(200 * time.Millisecond)
		}

	}
}
