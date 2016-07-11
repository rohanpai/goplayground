package main

import "fmt"

func pop(list *[]int, c chan int, done chan bool) {
	for len(*list) != 0 {
		result := (*list)[0]
		*list = (*list)[1:]
		fmt.Println("about to send ", result)
		c <- result
	}
	close(c)
	done <- true
}

func receiver(c chan int) {
	for result := range c {
		fmt.Println("received ", result)
	}
}

var list = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

func main() {
	c := make(chan int)
	done := make(chan bool)
	go pop(&list, c, done)
	go receiver(c)
	go receiver(c)
	go receiver(c)
	<-done
	fmt.Println("done")
}
