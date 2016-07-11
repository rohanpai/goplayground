package main

import "fmt"

func PadRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}

func PadLeft(str, pad string, lenght int) string {
	for {
		str = pad + str
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}

func main() {

	str := "abc"
	fmt.Println(PadRight(str, "x", 5))   // expects abcxx
	fmt.Println(PadLeft(str, "x", 5))    // expects xxabc
	fmt.Println(PadRight(str, "xyz", 5)) // expects abcxy
	fmt.Println(PadLeft(str, "xyz", 5))  // expects xyzab
}
