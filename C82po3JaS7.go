package main

import (
	&#34;fmt&#34;
	&#34;reflect&#34;
	&#34;strconv&#34;
)

var f0 = func() {
	fmt.Println(&#34;Hello 0&#34;)
}

func main() {
	// In order to use a function as a value
	// like in passing as an argument
	// We must use Function Literal
	// or the anonymous function must return something
	/*
	   we can&#39;t do this

	   myqueue.PushBack(func(str string) {
	   	fmt.Println(str)
	   }(&#34;Hello&#34;))

	   fmt.Println(func(str string) {
	   	fmt.Println(str)
	   }(&#34;Hello&#34;))

	   (func literal)(&#34;Hello&#34;) used as value

	   because the anonymous function closure
	   does not return anything
	*/

	// Function Literal
	temp := func(str string) string {
		return str
	}(&#34;Hello 08&#34;)
	fmt.Println(temp)
	// Hello 08

	// Function Literal  http://golang.org/ref/spec#Function_literals
	// This function has a name
	// to define a function {inside block}
	// Here we need this to define inside main function
	// Note that it has no parentheses in the end
	// Think of it like a function variable
	// (Method cannot be used like this)
	f0()                            // Hello 0
	fmt.Println(reflect.TypeOf(f0)) // func()

	f1 := func() {
		fmt.Println(&#34;Hello 09&#34;)
	}
	f1()                            // Hello 09
	fmt.Println(reflect.TypeOf(f1)) // func()

	f2 := func(str string) {
		fmt.Println(str)
	}
	f2(&#34;Hello 10&#34;) // Hello 10

	f3 := func(str string) string {
		return str
	}
	fmt.Println(f3(&#34;Hello 11&#34;)) // Hello 11

	// Just a function named f4
	// It takes an integer as an argument
	// and return a function which returns a string
	f4 := func(num int) func() string {
		// num is valid inside this block
		return func() string {
			s := strconv.Itoa(num)
			return &#34;Hello &#34; &#43; s
		}
	}

	// now ft is a function of type func() string
	// which returns a string
	ft := f4(12)
	fmt.Println(ft())
	fmt.Println(reflect.TypeOf(ft)) // func() string
	// Hello 12
}
