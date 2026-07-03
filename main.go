package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	goroutineBasics()

	// goroutineChannels()

	// goroutineWorkerPools()
}

// examples of goroutines with waitgroups
func goroutineBasics() {
	// start a waitgroup to track how many goroutines are running
	// wg.Add increments, wg.Done decrements, go.Wait() waits for 0
	var wg sync.WaitGroup

	// to start a goroutine, use the go command with a function
	// you can use a function literal (anonymous function)
	wg.Add(2) // track the 2 goroutine we are going to run
	go func() {
		response := "Hello, World!"
		time.Sleep(time.Millisecond * 1000)
		fmt.Println(response)
		wg.Done()
	}()

	// or a function defined elsewhere
	go greet("Robert", &wg)

	wg.Wait() // will wait until all wg.Done
}

func greet(name string, wg *sync.WaitGroup) {
	response := "Hello, " + name + "!"
	time.Sleep(time.Millisecond * 500)
	fmt.Println(response)
	wg.Done()
}
