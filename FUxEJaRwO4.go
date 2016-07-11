package main

import (
	&#34;os&#34;
)

func main() {

	for i := 0; i &lt; 67757; i&#43;&#43; {
		file, err := os.Open(&#34;.&#34;)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}
}
