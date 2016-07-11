package main

import (
	&#34;log&#34;
	&#34;math/rand&#34;
	&#34;sync&#34;
	&#34;time&#34;
)

// fanIn takes zero or more channels and merges the received data to a
// single output channel. For efficiency, the output channel should be
// buffered to the number of inputs to prevent goroutines blocking each
// other.
func fanIn(inputs []chan []byte, output chan []byte, exit chan bool, timeout time.Duration) {
	if len(inputs) == 0 {
		log.Println(&#34;zero inputs&#34;)
		return
	}

	defer log.Println(&#34;cleaning up fanIn&#34;)

	// Always signal the exit
	defer func() {
		exit &lt;- true
	}()

	// Used to signal goroutines to exit
	signal := make(chan struct{})

	// Wait group for spawned routines used after exit is signaled
	wg := sync.WaitGroup{}
	wg.Add(len(inputs))

	// Spawn goroutines for each input channel
	for i, input := range inputs {
		log.Println(&#34;spawning input&#34;, i)

		// Spawn go routine for each input
		go func(input chan []byte, i int) {
			defer log.Println(&#34;closing input&#34;, i)
			defer wg.Done()

			open := true
			// for-select idiom to constantly receive off the input
			// channel until it is closed on it has been signaled
			// to exit
			for open {
				select {
				case value, open := &lt;-input:
					// Input is closed, break
					if !open {
						log.Println(&#34;(closed) input&#34;, i)
						break
					}
					output &lt;- value
					log.Printf(&#34;input %d -&gt; %d\n&#34;, i, value)
				case &lt;-signal:
					log.Println(&#34;(signaled) input&#34;, i)
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
	if _, ok := &lt;-exit; !ok {
		log.Println(&#34;exit channel closed&#34;)
		close(signal)
	} else if timeout &gt; 0 {
		log.Println(&#34;timeout of&#34;, timeout, &#34;started&#34;)
		&lt;-time.After(timeout)
		close(signal)
	}

	// Wait until all routines are done and exit
	log.Println(&#34;waiting for goroutines to finish&#34;)
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
	for i := 0; i &lt; n; i&#43;&#43; {
		inputs[i] = make(chan []byte, n)

		go func(i int) {
			for j := 0; j &lt; n; j&#43;&#43; {
				inputs[i] &lt;- []byte{byte(r.Intn(20))}
			}
		}(i)
	}

	// Spawn fanIn in a goroutine
	go fanIn(inputs, output, exit, timeout)

	// Spawn goroutine to read and log the values from the output channel
	// as they are received.
	go func() {
		for m := range output {
			log.Println(&#34;output &lt;-&#34;, m)
		}
	}()

	// Request exit
	log.Println(&#34;exit signaled&#34;)
	if timeout &lt; 0 {
		close(exit)
	} else {
		exit &lt;- true
	}

	// Wait for response from fanIn
	&lt;-exit
	log.Println(&#34;exit confirmed&#34;)
}

func main() {
	testFanIn(10, 0)
}
