package main

import (
	&#34;fmt&#34;
	&#34;math/rand&#34;
	&#34;time&#34;
)

type Event int

func NewEvent() Event                   { return 0 }
func (e Event) Merge(other Event) Event { return e &#43; other }

func produce(out chan&lt;- Event) {
	for {
		delay := time.Duration(rand.Intn(5)&#43;1) * time.Second
		nMessages := rand.Intn(10) &#43; 1
		time.Sleep(delay)
		for i := 0; i &lt; nMessages; i&#43;&#43; {
			e := Event(rand.Intn(10))
			fmt.Println(&#34;Producing:&#34;, e)
			out &lt;- e
		}
	}
}

func coalesce(in &lt;-chan Event, out chan&lt;- Event) {
	event := NewEvent()
	timer := time.NewTimer(0)

	var timerCh &lt;-chan time.Time
	var outCh chan&lt;- Event

	for {
		select {
		case e := &lt;-in:
			event = event.Merge(e)
			if timerCh == nil {
				timer.Reset(500 * time.Millisecond)
				timerCh = timer.C
			}
		case &lt;-timerCh:
			outCh = out
			timerCh = nil
		case outCh &lt;- event:
			event = NewEvent()
			outCh = nil
		}
	}
}

func slowReceive(in &lt;-chan Event) {
	for {
		time.Sleep(1500 * time.Millisecond)
		fmt.Println(&#34;Received:&#34;, &lt;-in)
	}
}

func main() {
	source := make(chan Event)
	output := make(chan Event)

	go produce(source)
	go coalesce(source, output)
	slowReceive(output)
}
