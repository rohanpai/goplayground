package main

import (
	"github.com/rjeczalik/notify"
	"log"
)

func main() {
	var must = func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	var stop = func(c ...chan<- notify.EventInfo) {
		for _, c := range c {
			notify.Stop(c)
		}
	}

	// Make the channels buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.
	dir := make(chan notify.EventInfo, 1)
	file := make(chan notify.EventInfo, 1)
	symlink := make(chan notify.EventInfo, 1)
	all := make(chan notify.EventInfo, 1)

	// Set up a single watchpoint listening for FSEvents-specific events on
	// multiple user-provided channels.
	must(notify.Watch(".", dir, notify.FSEventsIsDir))
	must(notify.Watch(".", file, notify.FSEventsIsFile))
	must(notify.Watch(".", symlink, notify.FSEventsIsSymlink))
	must(notify.Watch(".", all, notify.All))
	defer stop(dir, file, symlink, all)

	// Block until an event is received.
	select {
	case ei := <-dir:
		log.Println("The directory", ei.Path(), "has changed")
	case ei := <-file:
		log.Println("The file", ei.Path(), "has changed")
	case ei := <-symlink:
		log.Println("The symlink", ei.Path(), "has changed")
	case ei := <-all:
		var kind string

		// Investigate underlying *notify.FSEvent struct to access more
		// information about the event.
		switch flags := ei.Sys().(*notify.FSEvent).Flags; {
		case flags&notify.FSEventsIsFile != 0:
			kind = "file"
		case flags&notify.FSEventsIsDir != 0:
			kind = "dir"
		case flags&notify.FSEventsIsSymlink != 0:
			kind = "symlink"
		}

		log.Printf("The %s under path %s has been %sd\n", kind, ei.Path(), ei.Event())
	}
}
