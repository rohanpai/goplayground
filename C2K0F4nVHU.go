package main

import (
	"fmt"
	"regexp"
)

func main() {
	str := []byte("12:00")
	assigned := regexp.MustCompile("(.*):(.*)")
	group := assigned.FindSubmatch(str)
	fmt.Println(string(group[0]));
	fmt.Println();
	fmt.Println(string(group[1]))
	fmt.Println(string(group[2]))
}
