package main

import (
	"fmt"
	"reflect"
)

func removeDuplicates(inputStream <-chan string, outputStream chan<- string) {
	defer close(outputStream)

	var lastVal string
	firstIteration := true

	for val := range inputStream {
		if firstIteration || val != lastVal {
			outputStream <- val
			lastVal = val
			firstIteration = false
		}
	}
}

func main() {
	input := make(chan string, 5)
	output := make(chan string, 5)

	// Заполняем входной канал
	go func() {
		defer close(input)
		input <- "a"
		input <- "a"
		input <- "b"
		input <- "b"
		input <- "c"
	}()

	go removeDuplicates(input, output)

	var results []string
	for val := range output {
		results = append(results, val)
	}

	expected := []string{"a", "b", "c"}
	if reflect.DeepEqual(results, expected) {
		fmt.Println("Expected: ", expected, "Got: ", results)
	}
}
