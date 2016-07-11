package main

import (
	&#34;github.com/rjeczalik/notify&#34;
	&#34;log&#34;
)

func main() {
	var must = func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	var stop = func(c ...chan&lt;- notify.EventInfo) {
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
	must(notify.Watch(&#34;.&#34;, dir, notify.FSEventsIsDir))
	must(notify.Watch(&#34;.&#34;, file, notify.FSEventsIsFile))
	must(notify.Watch(&#34;.&#34;, symlink, notify.FSEventsIsSymlink))
	must(notify.Watch(&#34;.&#34;, all, notify.All))
	defer stop(dir, file, symlink, all)

	// Block until an event is received.
	select {
	case ei := &lt;-dir:
		log.Println(&#34;The directory&#34;, ei.Path(), &#34;has changed&#34;)
	case ei := &lt;-file:
		log.Println(&#34;The file&#34;, ei.Path(), &#34;has changed&#34;)
	case ei := &lt;-symlink:
		log.Println(&#34;The symlink&#34;, ei.Path(), &#34;has changed&#34;)
	case ei := &lt;-all:
		var kind string

		// Investigate underlying *notify.FSEvent struct to access more
		// information about the event.
		switch flags := ei.Sys().(*notify.FSEvent).Flags; {
		case flags&amp;notify.FSEventsIsFile != 0:
			kind = &#34;file&#34;
		case flags&amp;notify.FSEventsIsDir != 0:
			kind = &#34;dir&#34;
		case flags&amp;notify.FSEventsIsSymlink != 0:
			kind = &#34;symlink&#34;
		}

		log.Printf(&#34;The %s under path %s has been %sd\n&#34;, kind, ei.Path(), ei.Event())
	}
}
