// go build -race

// Sample program to show how to use a read/write mutex to define critical
// sections of code that needs synchronous access.
package main

import (
	&#34;fmt&#34;
	&#34;math/rand&#34;
	&#34;runtime&#34;
	&#34;sync&#34;
	&#34;sync/atomic&#34;
	&#34;time&#34;
)

var (
	// data is a slice that will be shared.
	data []string

	// wg is used to wait for the program to finish.
	wg sync.WaitGroup

	// rwMutex is used to define a critical section of code.
	rwMutex sync.RWMutex

	// Number of reads occuring at ay given time.
	readCount int64
)

// init is called before main is executed.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// main is the entry point for all Go programs.
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Add the one goroutines for the writer.
	wg.Add(1)

	// Create the writer goroutine.
	go writer()

	// Create seven reader goroutines.
	for i := 1; i &lt;= 7; i&#43;&#43; {
		go reader(i)
	}

	// Wait for the write goroutine to finish.
	wg.Wait()
	fmt.Println(&#34;Program Complete&#34;)

	// To keep the sample simple we are allowing the runtime to
	// kill the reader goroutines. This is something we should
	// control before allowing main to exit.
}

// writer adds 10 new strings to the slice in random intervals.
func writer() {
	for i := 1; i &lt;= 10; i&#43;&#43; {
		// Only allow one goroutine to read/write to the
		// slice at a time.
		rwMutex.Lock()
		{
			// Capture the current read count.
			// Keep this safe though we can due without this call.
			rc := atomic.LoadInt64(&amp;readCount)

			// Perform some work since we have a full lock.
			fmt.Printf(&#34;****&gt; : Performing Write : RCount[%d]\n&#34;, rc)
			data = append(data, fmt.Sprintf(&#34;String: %d&#34;, i))
		}
		rwMutex.Unlock()
		// Release the lock and allow any waiting goroutines
		// to continue using the slice.

		// Sleep a random amount of time.
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}

	// Tell main we are done.
	wg.Done()
}

// reader wakes up and iterates over the data slice.
func reader(id int) {
	for {
		// Any goroutine can read when no write
		// operation is taking place.
		rwMutex.RLock()
		{
			// Increment the read count value by 1.
			rc := atomic.AddInt64(&amp;readCount, 1)

			// Perform some read work and display values.
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
			fmt.Printf(&#34;%d : Performing Read : Length[%d] RCount[%d]\n&#34;, id, len(data), rc)

			// Decrement the read count value by 1.
			atomic.AddInt64(&amp;readCount, -1)
		}
		rwMutex.RUnlock()
		// Release the read lock.
	}
}
