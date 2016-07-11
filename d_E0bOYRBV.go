// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program to show how the program can access a value
// of an unexported identifier from another package.
package main

import (
	&#34;fmt&#34;

	&#34;github.com/ardanlabs/gotraining/topics/exporting/example3/counters&#34;
)

func main() {

	// Create a variable of the unexported type using the exported
	// New function from the package counters.
	counter := counters.New(10)

	fmt.Printf(&#34;Counter: %d\n&#34;, counter)
}
