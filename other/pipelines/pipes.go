package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	in := gen(2, 3, 4, 5, 6)

	// Distribute the sq work across two goroutines that both read from in.
	c1 := sq(in)
	c2 := sq(in)
	c3 := sq(in)

	fmt.Printf("starting iterating over merged\n")
	// Consume the merged output from c1 and c2.
	for n := range merge(c1, c2, c3) {
		fmt.Println("final range merge:", n) // 4 then 9, or 9 then 4
	}
}

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan int) {
		for n := range c {
			fmt.Printf("set %d to out from %v\n", n, c)
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		fmt.Printf("outputting %v to out\n", c)
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		fmt.Printf("waiter started\n")
		wg.Wait()
		close(out)
		fmt.Printf("waiter finished job\n")
	}()
	fmt.Printf("main closed\n")
	return out
}

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			time.Sleep(1 * time.Second)
			out <- n * n
		}
		close(out)
	}()
	return out
}
