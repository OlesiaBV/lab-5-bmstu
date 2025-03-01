package main

import (
	"fmt"
	"time"
)

func calculator(firstChan <-chan int, secondChan <-chan int, stopChan <-chan struct{}) <-chan int {
	resultChan := make(chan int)

	go func() {
		defer close(resultChan)

		select {
		case val := <-firstChan:
			resultChan <- val * val
		case val := <-secondChan:
			resultChan <- val * 3
		case <-stopChan:
			return
		}
	}()

	return resultChan
}

func main() {
	firstChan := make(chan int, 1)
	secondChan := make(chan int, 1)
	stopChan := make(chan struct{})
	resultChan := calculator(firstChan, secondChan, stopChan)

	firstChan <- 4
	select {
	case result := <-resultChan:
		if result == 16 {
			fmt.Println("Expected 16, got: ", result)
		}
	case <-time.After(time.Second):
		fmt.Println("Timeout waiting for result")
	}

	secondChan <- 5
	resultChan = calculator(firstChan, secondChan, stopChan)
	select {
	case result := <-resultChan:
		if result == 15 {
			fmt.Println("Expected 15, got: ", result)
		}
	case <-time.After(time.Second):
		fmt.Println("Timeout waiting for result")
	}

	close(stopChan)
	resultChan = calculator(firstChan, secondChan, stopChan)
	select {
	case _, ok := <-resultChan:
		if ok {
			fmt.Println("Expected channel to be closed")
		}
	case <-time.After(time.Second):
		fmt.Println("Timeout waiting for channel closure")
	}
}
