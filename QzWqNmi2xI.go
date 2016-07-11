package main

import &#34;fmt&#34;

func main() {
	// ★★★★★★★★★
	// often used with channels
	func() {
		fmt.Println(&#34;Hello 05&#34;)
	}()
	// Hello 05

	ac := make(chan int, 1)
	func(c chan int) {
		c &lt;- 1
	}(ac)
	fmt.Println(&#34;Anonymous Function Closure for Channel:&#34;, &lt;-ac)
	// Anonymous Function Closure for Channel: 1

	// ★★★★★★★★★
	// For this function to stand alone,
	// it must refer to variables defined in a surrounding function.
	func(str string) {
		fmt.Println(str)
	}(&#34;Hello 06&#34;)
	// Hello 06

	/*
	   We can&#39;t do this
	   func() {
	   	fmt.Println(&#34;Hello&#34;)
	   }

	   func literal evaluated but not used
	*/

	// ★★★★★★★★★
	fmt.Println(func(str string) string {
		return str
	}(&#34;Hello 07&#34;))
	// Hello 07
}
