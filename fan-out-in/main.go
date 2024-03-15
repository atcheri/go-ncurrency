package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	done := make(chan interface{})
	defer close(done)

	randomNumberSteam := repeatFunc(done, randomNumberGenerator)

	// fan-out
	cpuCount := runtime.NumCPU()
	primeCatcherChans := make([]<-chan int, cpuCount)
	for i := 0; i < cpuCount; i++ {
		primeCatcherChans[i] = primeCatcher(done, randomNumberSteam)
	}

	// fan-in
	primeStream := fanInPrimes(done, primeCatcherChans...)
	for nb := range take(done, primeStream, 10) {
		fmt.Printf("receiving prime number %d from fanned-in stream\n", nb)
	}

	fmt.Println(time.Since(start))
}
func fanInPrimes[T any, D any](done <-chan D, catchers ...<-chan T) <-chan T {
	mergeStream := make(chan T)
	var wg sync.WaitGroup

	transfer := func(c <-chan T) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case mergeStream <- i:
			}
		}

	}

	for _, c := range catchers {
		wg.Add(1)
		go transfer(c)
	}

	go func() {
		wg.Wait()
		close(mergeStream)
	}()

	return mergeStream
}

func randomNumberGenerator() int {
	r := rand.Intn(100000000000)
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
				// fmt.Printf("sending %v in stream\n", r)
			}
		}
	}()

	return stream
}

func take[T any, D any](done <-chan D, source <-chan T, n int) <-chan T {
	taken := make(chan T)
	go func() {
		defer close(taken)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case taken <- <-source:
			}
		}
	}()

	return taken
}

func primeCatcher[T any, D any](done <-chan D, source <-chan T) <-chan T {
	primes := make(chan T)

	go func() {
		defer close(primes)
		for {
			select {
			case <-done:
				return
			case randomNb := <-source:
				if n, ok := any(randomNb).(int); ok {
					if isPrime(n) {
						primes <- randomNb
					}
				}
			}
		}
	}()
	return primes
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	} else if n == 2 {
		return true
	} else if n%2 == 0 {
		return false
	}
	sqrt := int(math.Sqrt(float64(n)))
	for i := 3; i <= sqrt; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}
