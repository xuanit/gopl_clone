package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	done := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func(in <-chan struct{}, out chan<- struct{}) {
		defer wg.Done()
		count := 0
	loop:
		for {
			select {
			case <-in:
				out <- struct{}{}
				count++
			case <-done:
				close(out)
				break loop
			}
		}
		for range in {
			// drain in channel

		}
		fmt.Printf("execute %d in 1 second", count)
	}(ch1, ch2)

	wg.Add(1)
	go func(in <-chan struct{}, out chan<- struct{}) {
		defer wg.Done()
		for range in {
			out <- struct{}{}
		}
		close(out)
	}(ch2, ch1)

	ch1 <- struct{}{}
	time.Sleep(1 * time.Second)
	close(done)
	// close(ch2)
	wg.Wait()
}
