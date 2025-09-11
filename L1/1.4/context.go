package main

import (
	"context"
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

	workersCount, err := strconv.Atoi(os.Args[1])
	if err != nil || workersCount <= 0 {
		fmt.Println("workers must be a positive integer")
		return
	}

	// Контекст, который отменится при SIGINT или SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	messages := make(chan int)
	var wg sync.WaitGroup

	// Запускаем воркеров
	for i := 1; i <= workersCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					// Контекст отменён — выходим
					fmt.Printf("worker %d: завершаюсь\n", id)
					return
				case msg, ok := <-messages:
					if !ok {
						return
					}
					fmt.Printf("worker %d: %d\n", id, msg)
				}
			}
		}(i)
	}

	// Генерируем сообщения, пока контекст не отменён
	counter := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Println("main: получен сигнал, закрываю messages")
			close(messages) // закрываем канал, чтобы воркеры завершили range
			wg.Wait()       // ждём всех воркеров
			fmt.Println("main: все воркеры завершились")
			return
		default:
			counter++
			messages <- counter
			time.Sleep(200 * time.Millisecond)
		}
	}
}
