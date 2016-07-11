package main

import "fmt"

type EventType uint8

var nextEventId EventType

func NextEventId() EventType {
	defer func() { nextEventId++ }()
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
	if int(eventType) >= len(e.rCounts) { // Check if we have enough room
		// Resize the handler tables accordingly
		newRCounts := make([]int, eventType+1) // Counts
		newRChans := make([]chan Event, eventType+1) // Channels
		copy(newRCounts, e.rCounts)
		copy(newRChans, e.channels)
		e.rCounts = newRCounts
		e.channels = newRChans
		
		// Create new count and channel
		e.rCounts[eventType] = 1
		e.channels[eventType] = make(chan Event, 1)
	} else {
		e.rCounts[eventType]++
		e.channels[eventType] = make(chan Event, e.rCounts[eventType])
	}

	go handler.EventLoop(e.channels[eventType])
}

// FireEvent fires an event to all receivers receiving
func (e *EventManager) FireEvent(event Event) {
	// No handlers for this event type
	if len(e.rCounts) <= int(event.Type()) {
		return
	}

	for i := 0; i < e.rCounts[event.Type()]; i++ {
		e.channels[event.Type()] <- event
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
		if test, ok := event.(*TestEvent); ok && test.myNum == 42 {
		event.(*TestEvent).ch <- true
		} else {
			event.(*TestEvent).ch <- false
		}
	}
}

// Main #################################################################################################

func main() {
	ch := make(chan bool)
	eventManager := &EventManager{}
	testHandler := &TestEventHandler{}
	eventManager.AddHandler(0, testHandler)
	go eventManager.FireEvent(&TestEvent{42, ch})
	if !<-ch {
		fmt.Println("Fail.")
	} else {
		fmt.Println("Success!")
	}
}