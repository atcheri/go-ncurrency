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

	for nb := range repeatFunc(done, randomNumberGenerator) {
		fmt.Printf("receiving %v from stream\n", nb)
	}

	fmt.Println(time.Since(start))
}

func randomNumberGenerator() int {
	return rand.Intn(1000000000)
}
func repeatFunc[T any, D any](done <-chan D, fn func() T) <-chan T {
	stream := make(chan T)

	go func() {
		defer close(stream)

		for {
			select {
			case <-done:
				return
			case stream <- fn():
				fmt.Printf("sending %v in stream\n", fn())
			}
		}
	}()

	return stream
}
