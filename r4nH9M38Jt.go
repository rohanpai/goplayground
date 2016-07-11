package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Event int

func NewEvent() Event                   { return 0 }
func (e Event) Merge(other Event) Event { return e + other }

func produce(out chan<- Event) {
	for {
		delay := time.Duration(rand.Intn(5)+1) * time.Second
		nMessages := rand.Intn(10) + 1
		time.Sleep(delay)
		for i := 0; i < nMessages; i++ {
			e := Event(rand.Intn(10))
			fmt.Println("Producing:", e)
			out <- e
		}
	}
}

func coalesce(in <-chan Event, out chan<- Event) {
	event := NewEvent()
	timer := time.NewTimer(0)

	var timerCh <-chan time.Time
	var outCh chan<- Event

	for {
		select {
		case e := <-in:
			event = event.Merge(e)
			if timerCh == nil {
				timer.Reset(500 * time.Millisecond)
				timerCh = timer.C
			}
		case <-timerCh:
			outCh = out
			timerCh = nil
		case outCh <- event:
			event = NewEvent()
			outCh = nil
		}
	}
}

func slowReceive(in <-chan Event) {
	for {
		time.Sleep(1500 * time.Millisecond)
		fmt.Println("Received:", <-in)
	}
}

func main() {
	source := make(chan Event)
	output := make(chan Event)

	go produce(source)
	go coalesce(source, output)
	slowReceive(output)
}
