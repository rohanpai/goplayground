package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

// fanIn takes zero or more channels and merges the received data to a
// single output channel. For efficiency, the output channel should be
// buffered to the number of inputs to prevent goroutines blocking each
// other.
func fanIn(inputs []chan []byte, output chan []byte, exit chan bool, timeout time.Duration) {
	if len(inputs) == 0 {
		log.Println("zero inputs")
		return
	}

	defer log.Println("cleaning up fanIn")

	// Always signal the exit
	defer func() {
		exit <- true
	}()

	// Used to signal goroutines to exit
	signal := make(chan struct{})

	// Wait group for spawned routines used after exit is signaled
	wg := sync.WaitGroup{}
	wg.Add(len(inputs))

	// Spawn goroutines for each input channel
	for i, input := range inputs {
		log.Println("spawning input", i)

		// Spawn go routine for each input
		go func(input chan []byte, i int) {
			defer log.Println("closing input", i)
			defer wg.Done()

			open := true
			// for-select idiom to constantly receive off the input
			// channel until it is closed on it has been signaled
			// to exit
			for open {
				select {
				case value, open := <-input:
					// Input is closed, break
					if !open {
						log.Println("(closed) input", i)
						break
					}
					output <- value
					log.Printf("input %d -> %d\n", i, value)
				case <-signal:
					log.Println("(signaled) input", i)
					open = false
				default:
					open = false
				}
			}
		}(input, i)
	}

	// The exit channel is expected to send a true value and wait
	// until it receives a response, however if it is closed,
	// immediately signal the goroutines.
	if _, ok := <-exit; !ok {
		log.Println("exit channel closed")
		close(signal)
	} else if timeout > 0 {
		log.Println("timeout of", timeout, "started")
		<-time.After(timeout)
		close(signal)
	}

	// Wait until all routines are done and exit
	log.Println("waiting for goroutines to finish")
	wg.Wait()
}

// Takes int denoting how many inputs are used by the fanIn function
// and a timeout. Use a timeout or 0 to never timeout.
func testFanIn(n int, timeout time.Duration) {
	// Array of `n` inputs channels
	inputs := make([]chan []byte, n)

	// Output channel
	output := make(chan []byte)

	// Exit channel
	exit := make(chan bool)

	// Seeded random number generator for populating the input channels
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Initialize and populate buffered input channels with a few messages
	for i := 0; i < n; i++ {
		inputs[i] = make(chan []byte, n)

		go func(i int) {
			for j := 0; j < n; j++ {
				inputs[i] <- []byte{byte(r.Intn(20))}
			}
		}(i)
	}

	// Spawn fanIn in a goroutine
	go fanIn(inputs, output, exit, timeout)

	// Spawn goroutine to read and log the values from the output channel
	// as they are received.
	go func() {
		for m := range output {
			log.Println("output <-", m)
		}
	}()

	// Request exit
	log.Println("exit signaled")
	if timeout < 0 {
		close(exit)
	} else {
		exit <- true
	}

	// Wait for response from fanIn
	<-exit
	log.Println("exit confirmed")
}

func main() {
	testFanIn(10, 0)
}
