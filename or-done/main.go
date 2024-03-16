package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	done := make(chan interface{})
	// defer close(done)
	defer func() {
		fmt.Println("gracefully shutting down the done channel")
		close(done)
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go seeStars(done, &wg, generate(done, "*"))

	wg.Add(1)
	go seeLines(done, &wg, generate(done, "-"))

	wg.Wait()

	fmt.Println(time.Since(start))
}

func generate(done <-chan interface{}, e interface{}) <-chan interface{} {
	stream := make(chan interface{}, 100)

	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				return
			case stream <- fmt.Sprintf("%s%s%s%s%s%s", e, e, e, e, e, e):
			}
		}
	}()

	return stream
}

func seeStars(done <-chan interface{}, wg *sync.WaitGroup, stars <-chan interface{}) {
	defer wg.Done()

	for star := range orDone(done, stars) {
		fmt.Println(star)
	}
}

func seeLines(done <-chan interface{}, wg *sync.WaitGroup, lines <-chan interface{}) {
	defer wg.Done()

	for line := range orDone(done, lines) {
		fmt.Println(line)
	}
}

func orDone(done <-chan interface{}, stream <-chan interface{}) <-chan interface{} {
	relay := make(chan interface{})

	go func() {
		defer close(relay)

		for {
			select {
			case <-done:
				return
			case s, ok := <-stream:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case relay <- s:
				}
			}
		}
	}()

	return relay
}
