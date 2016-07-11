package main

import &#34;fmt&#34;

func main() {
	const (
		Enone  = 0
		Eio    = 2
		Einval = 5
	)

	a := [...]string{Enone: &#34;no error&#34;, Eio: &#34;Eio&#34;, Einval: &#34;invalid argument&#34;}
	s := []string{Enone: &#34;no error&#34;, Eio: &#34;Eio&#34;, Einval: &#34;invalid argument&#34;}
	m := map[int]string{Enone: &#34;no error&#34;, Eio: &#34;Eio&#34;, Einval: &#34;invalid argument&#34;}

	fmt.Println(a)
	fmt.Println(s)
	fmt.Println(m)
}
