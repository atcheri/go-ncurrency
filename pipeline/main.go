package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	start := time.Now()
	done := make(chan interface{})
	defer close(done)

	for nb := range take(done, repeatFunc(done, randomNumberGenerator), 5) {
		fmt.Printf("receiving %d from stream\n", nb)
	}

	fmt.Println(time.Since(start))
}

func randomNumberGenerator() int {
	r := rand.Intn(1000000000)
	return r
}
func repeatFunc[T any, D any](done <-chan D, fn func() T) <-chan T {
	stream := make(chan T)

	go func() {
		defer close(stream)

		for {
			r := fn()
			select {
			case <-done:
				return
			case stream <- r:
				fmt.Printf("sending %v in stream\n", r)
			}
		}
	}()

	return stream
}

func take[T any, D any](done <-chan D, inputStream <-chan T, n int) <-chan T {
	outputStream := make(chan T)
	go func() {
		defer close(outputStream)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case outputStream <- <-inputStream:
			}
		}
	}()

	return outputStream
}
