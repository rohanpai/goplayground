package main

import (
	&#34;fmt&#34;
)

func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i &lt; n; i&#43;&#43; {
		c &lt;- x
		x, y = y, x&#43;y
	}
	close(c)
}

func main() {
	c := make(chan int, 10)
	go fibonacci(cap(c), c)
	for i := range c {
		fmt.Println(i)
	}
}
