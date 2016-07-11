package main

import "fmt"

func main() {
	g0, g1 := generate("counting the number of some arbitrary values")
	even, odd := count(g0, 'e'), count(g1, 'o')
	for i := 0; i < 2; i++ {
		select {
		case x := <-even:
			fmt.Printf("%d even 'e's\n", x)
		case x := <-odd:
			fmt.Printf("%d odd 'o's\n", x)
		}
	}
}

// count returns number of x's in found in c.
func count(c <-chan int, x int) (<-chan int) {
	r := make(chan int)
	go func() {
		sum := 0
		for i := range c {
			if i == x {
				sum++
			}
		}
		r <- sum
	}()
	return r
}

// generate sends all even-indexed runes in s to the first
// returned channel, and all odd-indexed runes to the other
// returned channel.
func generate(s string) (<-chan int, <-chan int) {
	even := make(chan int)
	odd := make(chan int)
	go func() {
		i := 0
		for _, r := range s {
			if i % 2 == 0 {
				even <- r
			}else{
				odd <- r
			}
			i++
		}
		close(even)
		close(odd)
	}()
	return even, odd
}
