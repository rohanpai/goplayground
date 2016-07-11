package main

import (
	&#34;fmt&#34;
	&#34;github.com/mikespook/gearman-go/worker&#34;
	&#34;sync&#34;
)

func main() {
	// An example of worker
	w := worker.New(worker.Unlimited)
	defer w.Close()
	// Add a gearman job server
	if err := w.AddServer(worker.Network, &#34;127.0.0.1:4730&#34;); err != nil {
		fmt.Println(err)
		return
	}
	// A function for handling jobs
	foobar := func(job worker.Job) ([]byte, error) {
		// Do nothing here
		return nil, nil
	}
	// Add the function to worker
	if err := w.AddFunc(&#34;foobar&#34;, foobar, 0); err != nil {
		fmt.Println(err)
		return
	}
	var wg sync.WaitGroup
	// A custome handler, for handling other results, eg. ECHO, dtError.
	w.JobHandler = func(job worker.Job) error {
		if job.Err() == nil {
			fmt.Println(string(job.Data()))
		} else {
			fmt.Println(job.Err())
		}
		wg.Done()
		return nil
	}
	// An error handler for handling worker&#39;s internal errors.
	w.ErrorHandler = func(e error) {
		fmt.Println(e)
		// Ignore the error or shutdown the worker
	}
	// Tell Gearman job server: I&#39;m ready!
	if err := w.Ready(); err != nil {
		fmt.Println(err)
		return
	}
	// Running main loop
	go w.Work()
	wg.Add(1)
	// calling Echo
	w.Echo([]byte(&#34;Hello&#34;))
	// Waiting results
	wg.Wait()
}
