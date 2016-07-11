package main

import (
	"os"
)

func main() {

	for i := 0; i < 67757; i++ {
		file, err := os.Open(".")
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}
}
