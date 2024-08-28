package main

import (
	"fmt"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("Worker ", id, "started job ", j)
		time.Sleep(time.Second)
		fmt.Println("Workder ", id, "finished job ", j)
		results <- j * 2
	}
}

func main() {
	const numJobs = 5
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for w := range 3 {
		go worker(w, jobs, results)
	}

	for i := range numJobs {
		jobs <- i
	}

	close(jobs)

	for a := range numJobs {
		a++
		fmt.Println(<-results)
	}

}
