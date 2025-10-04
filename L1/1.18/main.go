package main

import "sync"

func main() {
	var ctr Counter
	var wg sync.WaitGroup

	// Запускаем N горутин
	for i := 0; i < 5; i++ {
		wg.Go(func() { incrementer(&ctr) })
	}

	wg.Wait()
	println("Final count:", ctr.Data)
}

func incrementer(ctr *Counter) {
	ctr.Mutex.Lock()
	defer ctr.Mutex.Unlock()
	ctr.Data++
}

type Counter struct {
	Data  int
	Mutex sync.RWMutex
}
