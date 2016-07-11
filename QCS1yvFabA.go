package main

import "fmt"

func main() {
	const (
		Enone  = 0
		Eio    = 2
		Einval = 5
	)

	a := [...]string{Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
	s := []string{Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
	m := map[int]string{Enone: "no error", Eio: "Eio", Einval: "invalid argument"}

	fmt.Println(a)
	fmt.Println(s)
	fmt.Println(m)
}
