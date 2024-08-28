package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func doWork(id int, wg *sync.WaitGroup) {
	randomDelay := rand.IntN(10) + 1
	defer wg.Done()

    fmt.Printf("\nWorker %d starting, will sleep for %d seconds", id, randomDelay)
    time.Sleep(time.Duration(randomDelay) * time.Second)
    fmt.Printf("\nWorker %d finished", id)
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go doWork(i, &wg)
	}

	wg.Wait()
	fmt.Println("\nAll Workers completed")

}
