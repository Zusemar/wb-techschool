package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	// ---- Завершение через контекст с таймаутом ----
	fmt.Println("\n---- Завершение через контекст с таймаутом ----")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		counter := 0
		for {
			select {
			case <-ctx.Done():
				fmt.Println("ctx goroutine end")
				return
			default:
				fmt.Println(counter)
				counter++
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()
	wg.Wait()

	// ---- Завершение через runtime.Goexit ----
	fmt.Println("\n---- Завершение через рантайм ----")

	wg.Add(1)
	go func() {
		defer wg.Done()
		goRuntime()
	}()
	wg.Wait()

	// ---- Завершение через канал завершения ----
	fmt.Println("\n---- Завершение через канал завершения ----")

	ch := make(chan int)
	done := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(ch)
		i := 0
		for {
			select {
			case ch <- i:
				i++
				if i == 10 {
					close(done)
				}
			case <-done:
				fmt.Println("ch goroutine end")
				return
			}
		}
	}()

	// получатель сам не горутина — ждём wg ниже
	for value := range ch {
		fmt.Println("Получено:", value)
		time.Sleep(100 * time.Millisecond)
	}
	wg.Wait()

	// ---- Завершение через context.AfterFunc ----
	fmt.Println("\n---- Завершение через context.AfterFunc ----")

	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()

	stop := make(chan struct{})
	context.AfterFunc(ctx2, func() {
		fmt.Println("context.AfterFunc вызван — завершаем горутину")
		close(stop)
	})

	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		for {
			select {
			case <-stop:
				fmt.Println("goroutine (AfterFunc) end")
				return
			default:
				fmt.Println("work", i)
				i++
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()
	wg.Wait()

	// ---- Естественное завершение ----
	fmt.Println("\n---- Завершение естественное ----")

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 2; i++ {
			fmt.Println("some work")
			time.Sleep(200 * time.Millisecond)
		}
		fmt.Println("ending by itself")
	}()
	wg.Wait()

	fmt.Println("enough time passed to job to be done")
}

func goRuntime() {
	for i := 0; i < 5; i++ {
		fmt.Printf("Goexit работа %d\n", i)
		time.Sleep(300 * time.Millisecond)
		if i == 2 {
			fmt.Println("Вызов Goexit()")
			runtime.Goexit() // немедленное завершение
		}
	}
}
