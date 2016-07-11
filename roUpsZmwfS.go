package main

import (
	"fmt"
	"time"
)

func debounce(interval time.Duration, done chan<- struct{}, input <-chan int, f func(arg int)) {
	var (
		item int
		timeout = time.After(interval)
	)
	for {
		select {
		case item = <-input:
			if item % 100000 == 0 {
				fmt.Println("received a send on a spammy channel - might be doing a costly operation if not for debounce")
			}
		case <-timeout:
			f(item)
			done<-struct{}{}
		}
	}
}

func main() {
	spammyChan := make(chan int, 10)
	done := make(chan struct{})
	go debounce(100*time.Millisecond, done, spammyChan, func(arg int) {
		fmt.Println("*****************************")
		fmt.Println("* DOING A COSTLY OPERATION! *")
		fmt.Println("*****************************")
		fmt.Println("In case you were wondering, the value passed to this function is", arg)
		fmt.Println("We could have more args to our \"compiled\" debounced function too, if we wanted.")
	})
	loop:
	for i := 0; i < 10000000; i++ {
		select {
			case spammyChan <- i:
				continue loop
			case <-done:
				fmt.Println("break loop")
				return
		}
	}
	fmt.Println("exited loop")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("done.")
}