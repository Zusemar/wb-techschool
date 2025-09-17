package main

import (
	"context"
	"fmt"
	"time"
)

// пособие как завершать работу горутин
func main() {
	// через контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() error {
		counter := 0

		select {
		case <-ctx.Done():
			fmt.Println("ctx gorutine end")
			return ctx.Err()
		default:
			fmt.Println(counter)
			counter++
		}
		return nil
	}()

	// через рантайм го

	// через канал завершения

	// естественное завершение

	// waitgroup

	//
}
