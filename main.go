package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	goroutineBasics()

	goroutineChannels()

	goroutineWorkerPools()
}

// goroutineBasics shows launching goroutines and waiting for them with a WaitGroup.
func goroutineBasics() {
	// WaitGroup tracks how many goroutines are still running.
	// Add increases the counter, Done decreases it, Wait blocks until it hits 0.
	var wg sync.WaitGroup

	// Start a goroutine with `go`. Function literals (anonymous funcs) work fine.
	wg.Add(2) // we will start 2 goroutines
	go func() {
		defer wg.Done() // always Done when this goroutine exits
		response := "Hello, World!"
		time.Sleep(time.Millisecond * 1000)
		fmt.Println(response)
	}()

	// Named functions work the same way — pass a pointer to the WaitGroup.
	go greet("Robert", &wg)

	wg.Wait() // block until both goroutines call Done
}

// greet is a named function used as a goroutine in goroutineBasics.
func greet(name string, wg *sync.WaitGroup) {
	defer wg.Done()

	response := "Hello, " + name + "!"
	time.Sleep(time.Millisecond * 500)
	fmt.Println(response)
}

//
// CHANNELS
//

// goroutineChannels shows sending a value from a goroutine to main via a channel.
func goroutineChannels() {
	// Unbuffered channel: send blocks until something receives (and vice versa).
	messages := make(chan string)

	// Pass the channel in. A function literal could close over messages instead.
	go greetViaChannel("John", messages)

	// Single receive — fine when you expect exactly one value.
	msg := <-messages
	fmt.Println(msg)

	// Multi-receive alternative: range until the sender closes the channel.
	// for msg := range messages {
	//     fmt.Println(msg)
	// }
}

// greetViaChannel writes one message then closes. Only the sender should close.
func greetViaChannel(name string, ch chan<- string) {
	response := "Hello, " + name + "!"
	time.Sleep(time.Millisecond * 500)

	// Send blocks until main receives (unbuffered channel).
	ch <- response
	close(ch)
}

//
// WORKER POOLS
//

// worker pulls jobs from jobs, does work, and pushes results. jobs is receive-only
// (<-chan) and results is send-only (chan<-) so the direction of data is clear.
func worker(id int, jobs <-chan int, results chan<- int) {
	// range over jobs ends when jobs is closed and drained.
	for j := range jobs {
		fmt.Println("worker", id, "started job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

// goroutineWorkerPools runs a fixed pool of workers over a batch of jobs.
func goroutineWorkerPools() {
	const numJobs = 10
	const poolSize = 4

	// Buffered channels sized to numJobs so we can enqueue all work (and collect
	// all results) without the main goroutine blocking on every send/receive.
	// Without the buffer, main could deadlock: workers might be busy while main
	// is stuck trying to send the next job, or waiting on a result that can't
	// be written because results is full and nobody is reading yet.
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// Start a fixed pool — only poolSize goroutines, not one per job.
	for w := 1; w <= poolSize; w++ {
		go worker(w, jobs, results)
	}

	// Feed all jobs, then close so workers know no more work is coming
	// (range in worker exits after the last value).
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// Collect one result per job. This also keeps main alive until work finishes.
	for a := 1; a <= numJobs; a++ {
		<-results
	}
}
