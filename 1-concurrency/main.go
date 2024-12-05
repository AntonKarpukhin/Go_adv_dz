package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	counter := 10
	startNumbersChan := make(chan int, counter)
	doubleNumbersChan := make(chan int, counter)

	var wg sync.WaitGroup
	wg.Add(2)

	go randomNumbers(startNumbersChan, &wg, counter)
	go doubleNumbers(startNumbersChan, doubleNumbersChan, &wg)

	go func() {
		wg.Wait()
	}()

	for square := range doubleNumbersChan {
		fmt.Print(" ", square)
	}
}

func randomNumbers(startNumbersChan chan int, wg *sync.WaitGroup, counter int) {
	defer wg.Done()
	for i := 0; i < counter; i++ {
		num := rand.Intn(101)
		startNumbersChan <- num
	}
	close(startNumbersChan)
}

func doubleNumbers(startNumbersChan chan int, doubleNumbersChan chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range startNumbersChan {
		doubleNumbersChan <- num * num
	}
	close(doubleNumbersChan)
}
