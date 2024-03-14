package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	var wg sync.WaitGroup
	input := []int{3, -2, 5, 0, 15}
	result := make([]int, len(input))

	for i, v := range input {
		wg.Add(1)
		go processDataSlow(&wg, &result[i], v)
	}

	wg.Wait()

	fmt.Printf("%v\n", result)
	fmt.Println(time.Since(start))
}

func processDataSlow(wg *sync.WaitGroup, result *int, v int) {
	defer wg.Done()

	*result = processValueSlow(v)
}

func processValueSlow(v int) int {
	time.Sleep(time.Duration(time.Second * 2))
	return v * v
}
