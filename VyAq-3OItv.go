package main

import &#34;fmt&#34;

type T struct {
	Name  string
	Port  int
	State State
}

type State int

const (
	Running State = iota &#43; 1
	Stopped
	Rebooting
	Terminated
)

func (s State) String() string {
	switch s {
	case Running:
		return &#34;Running&#34;
	case Stopped:
		return &#34;Stopped&#34;
	case Rebooting:
		return &#34;Rebooting&#34;
	case Terminated:
		return &#34;Terminated&#34;
	default:
		return &#34;Unknown&#34;
	}
}

func main() {
	t := T{Name: &#34;example&#34;, Port: 6666}

	// prints: &#34;t {Name:example Port:6666 State:Running}&#34;
	fmt.Printf(&#34;t %&#43;v\n&#34;, t)
}