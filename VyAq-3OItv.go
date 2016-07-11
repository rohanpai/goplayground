package main

import "fmt"

type T struct {
	Name  string
	Port  int
	State State
}

type State int

const (
	Running State = iota + 1
	Stopped
	Rebooting
	Terminated
)

func (s State) String() string {
	switch s {
	case Running:
		return "Running"
	case Stopped:
		return "Stopped"
	case Rebooting:
		return "Rebooting"
	case Terminated:
		return "Terminated"
	default:
		return "Unknown"
	}
}

func main() {
	t := T{Name: "example", Port: 6666}

	// prints: "t {Name:example Port:6666 State:Running}"
	fmt.Printf("t %+v\n", t)
}