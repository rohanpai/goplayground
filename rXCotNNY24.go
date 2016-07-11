/*
The classical Dining philosophers problem.

Implemented with forks (aka chopsticks) as mutexes.
*/

package main

import (
	&#34;fmt&#34;
	&#34;math/rand&#34;
	&#34;sync&#34;
	&#34;time&#34;
)

type Fork struct{ sync.Mutex }

type Philosopher struct {
	id                  int
	leftFork, rightFork *Fork
}

// Endlessly dine.
// Goes from thinking to hungry to eating and starts over.
// Adapt the pause values to increased or decrease contentions
// around the forks.
func (p Philosopher) dine() {
	say(&#34;thinking&#34;, p.id)
	randomPause(2)

	say(&#34;hungry&#34;, p.id)
	p.leftFork.Lock()
	p.rightFork.Lock()

	say(&#34;eating&#34;, p.id)
	randomPause(5)

	p.rightFork.Unlock()
	p.leftFork.Unlock()

	p.dine()
}

func randomPause(max int) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(max*1000)))
}

func say(action string, id int) {
	fmt.Printf(&#34;#%d is %s\n&#34;, id, action)
}

func init() {
	// Random seed
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	// How many philosophers and forks
	count := 5

	// Create forks
	forks := make([]*Fork, count)
	for i := 0; i &lt; count; i&#43;&#43; {
		forks[i] = new(Fork)
	}

	// Create philospoher, assign them 2 forks and send them to the dining table
	philosophers := make([]*Philosopher, count)
	for i := 0; i &lt; count; i&#43;&#43; {
		philosophers[i] = &amp;Philosopher{
			id: i, leftFork: forks[i], rightFork: forks[(i&#43;1)%count]}
		go philosophers[i].dine()
	}

	// Wait endlessly while they&#39;re dining
	endless := make(chan int)
	&lt;-endless
}