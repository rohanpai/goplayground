package main

import (
	"log"
	"os"
	"sync"
	"time"
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
	w := &Watcher{current: c, broadcast: make(chan struct{})}
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
func (w *Watcher) Wait(c *Config) <-chan *Config {
	ch := make(chan *Config, 1)
	go func() {
		w.mu.Lock()
		if w.current.serial < c.serial {
			panic("provided config newer than current!") // shouldn't happen
		}
		if w.current.serial > c.serial {
			// we have a newer config, deliver immediately
			ch <- w.current
			w.mu.Unlock()
			return
		}
		bc := w.broadcast
		w.mu.Unlock()

		// waiter has current config, wait for update
		<-bc

		w.mu.Lock()
		ch <- w.current // deliver new config
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
				log.Printf("Statting config file %v: %v", filename, err)
				continue
			}
			if mtime := fi.ModTime(); last.IsZero() || mtime.After(last) {
				last = mtime
				break
			}
		}
		c, err := readConfig(filename)
		if err != nil {
			log.Printf("Reading config file %v: %v", filename, err)
			continue
		}
		w.update(c)
	}
}

func (w *Watcher) update(c *Config) {
	w.mu.Lock()
	defer w.mu.Unlock()

	c.serial = w.current.serial + 1   // increment serial of new config
	w.current = c                     // replace config
	close(w.broadcast)                // tell any watchers that the config has been updated
	w.broadcast = make(chan struct{}) // create new broadcast channel
}

func readConfig(filename string) (*Config, error) {
	panic("not implemented")
}

func main() {
	// Create one Watcher for your entire process.
	w, err := NewWatcher("config.json")
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
		case cfg = <-next:
			// A new config has been updated,
			// make the necessary changes to the worker's state.

			// Set up a waiter channel for the next config.
			next = w.Wait(cfg)
		default:
		}

		// Do whatever work the worker does.
	}
}
