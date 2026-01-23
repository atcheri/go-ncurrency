package main

import (
	"fmt"
	"math/rand"
	"time"
)

func orDone(done <-chan struct{}, input <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for {
			select {
			case <-done:
				return
			case v, ok := <-input:
				if !ok {
					return
				}

				select {
				case out <- v:
				case <-done:
				}
			}
		}
	}()

	return out
}

func tee(input <-chan int) (<-chan int, <-chan int) {
	out1 := make(chan int)
	out2 := make(chan int)

	go func() {
		defer close(out1)
		defer close(out2)

		for v := range input {
			select {
			case out1 <- v:
			case out2 <- v:
			}
		}

	}()

	return out1, out2
}

func transactionGenerator(done <-chan struct{}) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for range 200 {
			select {
			case <-done:
				fmt.Println("Transaction Generator: Shutting down.")
				return
			case out <- rand.Intn(1000):
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return out
}

func fraudDetectionPipeline(done <-chan struct{}, input <-chan int) {
	for transaction := range orDone(done, input) {
		if transaction > 700 { // Example condition for "suspicious" activity
			fmt.Println("Fraud Detection: Suspicious transaction detected:", transaction)
		} else {
			fmt.Println("Fraus Detection: Normal transaction:", transaction)
		}

		fmt.Println("Fraud Detection Pipeline shutting down.")
	}
}

func analyticsPipeline(done <-chan struct{}, input <-chan int) {
	total := 0
	count := 0

	for transaction := range orDone(done, input) {
		fmt.Println("Analytics Pipeline: Processing transaction:", transaction)
		total += transaction
		count++
	}
	if count > 0 {
		fmt.Printf("Analytics Pipeline: Average transaction amount: %d\n", total/count)
	}

	fmt.Println("Analytics Pipeline shutting down.")
}

func main() {
	done := make(chan struct{})

	// Generate transactions and split them for both pipelines
	source := transactionGenerator(done)
	fraudCh, analyticsCh := tee(source)

	// Run pipelines concurrently
	go fraudDetectionPipeline(done, fraudCh)
	go analyticsPipeline(done, analyticsCh)

	// Simulate a cancellation after 10 seconds
	time.Sleep(10 * time.Second)
	close(done)

	// Allow time for the pipelines to wrap up
	time.Sleep(time.Second)
	fmt.Println("Main: All pipelines are stopped.")
}
