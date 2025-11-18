package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	const jobsCount, workerCount = 15, 3

	jobs := make(chan int, 15)
	results := make(chan int, 15)

	for i := 0; i < workerCount; i++ {
		go worker(i+1, jobs, results)
	}

	for i := 0; i < jobsCount; i++ {
		jobs <- i + 1
	}
	close(jobs)

	for i := 0; i < jobsCount; i++ {
		fmt.Printf("result #%d; value = %d\n", i+1, <-results)
	}

	fmt.Printf("time elapsed: %v\n", time.Since(t))
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for v := range jobs {
		time.Sleep(time.Second)
		fmt.Printf("worker #%d finished\n", id)
		results <- v * v
	}

}
