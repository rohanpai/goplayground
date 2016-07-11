package main

import &#34;fmt&#34;

type EventType uint8

var nextEventId EventType

func NextEventId() EventType {
	defer func() { nextEventId&#43;&#43; }()
	return nextEventId
}

type Event interface {
	Type() EventType // Returns the event type for this event
}

type EventHandler interface {
	EventLoop(ch chan Event)
}

type EventManager struct {
	rCounts []int // Receiver counts for each event type
	channels []chan Event // Channels for each event type
}

// AddHandler adds an event handler for a particular event type
func (e *EventManager) AddHandler(eventType EventType, handler EventHandler) {
	if int(eventType) &gt;= len(e.rCounts) { // Check if we have enough room
		// Resize the handler tables accordingly
		newRCounts := make([]int, eventType&#43;1) // Counts
		newRChans := make([]chan Event, eventType&#43;1) // Channels
		copy(newRCounts, e.rCounts)
		copy(newRChans, e.channels)
		e.rCounts = newRCounts
		e.channels = newRChans
		
		// Create new count and channel
		e.rCounts[eventType] = 1
		e.channels[eventType] = make(chan Event, 1)
	} else {
		e.rCounts[eventType]&#43;&#43;
		e.channels[eventType] = make(chan Event, e.rCounts[eventType])
	}

	go handler.EventLoop(e.channels[eventType])
}

// FireEvent fires an event to all receivers receiving
func (e *EventManager) FireEvent(event Event) {
	// No handlers for this event type
	if len(e.rCounts) &lt;= int(event.Type()) {
		return
	}

	for i := 0; i &lt; e.rCounts[event.Type()]; i&#43;&#43; {
		e.channels[event.Type()] &lt;- event
	}
}

// TestHandler ##########################################################################################

type TestEvent struct {
	myNum int
	ch    chan bool
}

func (t *TestEvent) Type() EventType {
	return 0
}

type TestEventHandler struct {
}

func (t *TestEventHandler) EventLoop(ch chan Event) {
	for event := range ch {
		if test, ok := event.(*TestEvent); ok &amp;&amp; test.myNum == 42 {
		event.(*TestEvent).ch &lt;- true
		} else {
			event.(*TestEvent).ch &lt;- false
		}
	}
}

// Main #################################################################################################

func main() {
	ch := make(chan bool)
	eventManager := &amp;EventManager{}
	testHandler := &amp;TestEventHandler{}
	eventManager.AddHandler(0, testHandler)
	go eventManager.FireEvent(&amp;TestEvent{42, ch})
	if !&lt;-ch {
		fmt.Println(&#34;Fail.&#34;)
	} else {
		fmt.Println(&#34;Success!&#34;)
	}
}