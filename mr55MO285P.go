package main

import &#34;fmt&#34;

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
	f.endurance = 42 &#43; (14 * level)
	f.health = 5 * f.endurance
	f.deflection = 25
	f.accuracy = 30 &#43; (3 * level)
}

func main() {
	player := NewFighter(&#34;Calisca&#34;)
	fmt.Printf(&#34;%&#43;v\n&#34;, player)

	player.SetLevel(2)
	fmt.Printf(&#34;%&#43;v\n&#34;, player)
}
