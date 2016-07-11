package main

import &#34;fmt&#34;

func pop(list *[]int, c chan int, done chan bool) {
	for len(*list) != 0 {
		result := (*list)[0]
		*list = (*list)[1:]
		fmt.Println(&#34;about to send &#34;, result)
		c &lt;- result
	}
	close(c)
	done &lt;- true
}

func receiver(c chan int) {
	for result := range c {
		fmt.Println(&#34;received &#34;, result)
	}
}

var list = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

func main() {
	c := make(chan int)
	done := make(chan bool)
	go pop(&amp;list, c, done)
	go receiver(c)
	go receiver(c)
	go receiver(c)
	&lt;-done
	fmt.Println(&#34;done&#34;)
}
