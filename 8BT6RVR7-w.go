package main

import "fmt"
import "time"

func main() {
	fmt.Println("START 1")
	for i := 0; i < 3; i++ {
		foo(i)
	}
	fmt.Println("END 1")

	fmt.Println("START 2")
	for i := 0; i < 3; i++ {
		go foo(i)
	}
	fmt.Println("END 2")

	time.Sleep(1 * time.Second)
}

func foo(i int) {
	time.Sleep(1)
	fmt.Printf("Call index: %d\n", i)
}
