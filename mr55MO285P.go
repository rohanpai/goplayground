package main

import "fmt"

type Fighter struct {
	name  string
	level int

	endurance, health, deflection, accuracy int
}

func NewFighter(name string) Fighter {
	var f = Fighter{name: name}
	f.SetLevel(1)
	return f
}

func (f *Fighter) SetLevel(level int) {
	f.level = level
	f.endurance = 42 + (14 * level)
	f.health = 5 * f.endurance
	f.deflection = 25
	f.accuracy = 30 + (3 * level)
}

func main() {
	player := NewFighter("Calisca")
	fmt.Printf("%+v\n", player)

	player.SetLevel(2)
	fmt.Printf("%+v\n", player)
}
