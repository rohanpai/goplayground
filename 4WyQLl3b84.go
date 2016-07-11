package main

import (
	&#34;log&#34;
	&#34;os&#34;
	&#34;sync&#34;
	&#34;time&#34;
)

type Config struct {
	Address string
	Port    uint32

	serial int // known only to Watcher
}

// NewWatcher creates a new Watcher watching filename.
func NewWatcher(filename string) (*Watcher, error) {
	c, err := readConfig(filename)
	if err != nil {
		return nil, err
	}
	c.serial = 0
	w := &amp;Watcher{current: c, broadcast: make(chan struct{})}
	go w.watch(filename)
	return w, nil
}

// Watcher watches a configuration file and
// sends updated Configs to any waiting clients.
type Watcher struct {
	mu        sync.Mutex
	current   *Config
	broadcast chan struct{}
}

// Config returns the current Config.
func (w *Watcher) Config() *Config {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.current
}

// Wait returns a channel on which it will send
// a Config newer than the provided Config.
func (w *Watcher) Wait(c *Config) &lt;-chan *Config {
	ch := make(chan *Config, 1)
	go func() {
		w.mu.Lock()
		if w.current.serial &lt; c.serial {
			panic(&#34;provided config newer than current!&#34;) // shouldn&#39;t happen
		}
		if w.current.serial &gt; c.serial {
			// we have a newer config, deliver immediately
			ch &lt;- w.current
			w.mu.Unlock()
			return
		}
		bc := w.broadcast
		w.mu.Unlock()

		// waiter has current config, wait for update
		&lt;-bc

		w.mu.Lock()
		ch &lt;- w.current // deliver new config
		w.mu.Unlock()
	}()
	return ch
}

const pollInterval = 10 * time.Second

func (w *Watcher) watch(filename string) {
	var last time.Time
	for {
		for {
			time.Sleep(pollInterval)
			fi, err := os.Stat(filename)
			if err != nil {
				log.Printf(&#34;Statting config file %v: %v&#34;, filename, err)
				continue
			}
			if mtime := fi.ModTime(); last.IsZero() || mtime.After(last) {
				last = mtime
				break
			}
		}
		c, err := readConfig(filename)
		if err != nil {
			log.Printf(&#34;Reading config file %v: %v&#34;, filename, err)
			continue
		}
		w.update(c)
	}
}

func (w *Watcher) update(c *Config) {
	w.mu.Lock()
	defer w.mu.Unlock()

	c.serial = w.current.serial &#43; 1   // increment serial of new config
	w.current = c                     // replace config
	close(w.broadcast)                // tell any watchers that the config has been updated
	w.broadcast = make(chan struct{}) // create new broadcast channel
}

func readConfig(filename string) (*Config, error) {
	panic(&#34;not implemented&#34;)
}

func main() {
	// Create one Watcher for your entire process.
	w, err := NewWatcher(&#34;config.json&#34;)
	if err != nil {
		log.Fatal(err)
	}

	worker(w) // presumably many of these
}

func worker(w *Watcher) {
	// The various parts of your system that need a config file
	// should retrieve it at startup with the Config method.
	cfg := w.Config()

	// Then the processing loops for those various parts should
	// create a waiter channel that will receive the next available config.
	next := w.Wait(cfg)
	for {
		// At some convenient point, check whether a new config is available.
		select {
		case cfg = &lt;-next:
			// A new config has been updated,
			// make the necessary changes to the worker&#39;s state.

			// Set up a waiter channel for the next config.
			next = w.Wait(cfg)
		default:
		}

		// Do whatever work the worker does.
	}
}
