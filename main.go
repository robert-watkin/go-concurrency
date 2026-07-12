package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	goroutineBasics()

	goroutineChannels()

	// goroutineWorkerPools()
}

// examples of goroutines with waitgroups
func goroutineBasics() {
	// start a waitgroup to track how many goroutines are running
	// wg.Add increments, wg.Done decrements, go.Wait blocks until counter is 0
	var wg sync.WaitGroup

	// to start a goroutine, use the go command with a function
	// you can use a function literal (anonymous function)
	wg.Add(2) // track the 2 goroutine we are going to run
	go func() {
		defer wg.Done()
		response := "Hello, World!"
		time.Sleep(time.Millisecond * 1000)
		fmt.Println(response)
	}()

	// or a function defined elsewhere
	go greet("Robert", &wg)

	wg.Wait() // will wait until all wg.Done
}

// greet function used within goroutineBasics
func greet(name string, wg *sync.WaitGroup) {
	defer wg.Done()

	response := "Hello, " + name + "!"
	time.Sleep(time.Millisecond * 500)
	fmt.Println(response)
}

// this function demonstrates how to use a channel
func goroutineChannels() {
	// messages channel that accepts a string initialised
	messages := make(chan string)

	// call greetViaChannel as a goroutine passing through the channel we have created
	// if we had used a function literal then we could access messages directly instead of passing
	go greetViaChannel("John", messages)

	// single receive - idiomatic when you expect exactly one value
	msg := <-messages
	fmt.Println(msg)

	// Multi-receive - runs until the sender closes the channel
	// for msg := range messages {
	//     fmt.Println(msg)
	// }
}

func greetViaChannel(name string, ch chan<- string) {
	response := "Hello, " + name + "!"
	time.Sleep(time.Millisecond * 500)

	// write the response message to the channel
	// we will wait here until there is a matching read from the channel
	ch <- response
	close(ch) // close channel at sender
}
