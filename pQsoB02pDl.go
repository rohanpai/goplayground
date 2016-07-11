// Sample program to show how to compose maps of maps.
package main

import (
	&#34;fmt&#34;
)

// Data defines a key/value store.
type Data map[string]string

// main is the entry point for the application.
func main() {
	// Declare and make a map of Data type values.
	users := make(map[string]Data)

	// Intialize some data into our map of maps.
	users[&#34;clients&#34;] = Data{&#34;123&#34;: &#34;Henry&#34;, &#34;654&#34;: &#34;Joan&#34;}
	users[&#34;admins&#34;] = Data{&#34;398&#34;: &#34;Bill&#34;, &#34;076&#34;: &#34;Steve&#34;}

	// Iterate over the map of maps.
	for key, data := range users {
		fmt.Println(key)
		for key, value := range data {
			fmt.Println(&#34;\t&#34;, key, value)
		}
	}
}
