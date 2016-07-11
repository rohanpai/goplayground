package main

import (
	&#34;fmt&#34;
	&#34;math/rand&#34;
)

func main() {
	fmt.Println(&#34;Spaceline Company Days Round-trip Price&#34;)
	fmt.Println(&#34;=======================================&#34;)

	var distance = 57600000
	var count = 0

	for count &lt; 10 {
		var speed = rand.Intn(15) &#43; 16          // 16-30 km/s
		var duration = distance / speed / 86400 // days
		var price = 20.0 &#43; speed                // $ millions

		switch rand.Intn(3) {
		case 0:
			fmt.Print(&#34;Space Adventures  &#34;)
		case 1:
			fmt.Print(&#34;SpaceX            &#34;)
		case 2:
			fmt.Print(&#34;Virgin Galactic   &#34;)
		}

		fmt.Printf(&#34;%4v &#34;, duration)

		if rand.Intn(2) == 1 {
			fmt.Print(&#34;Round-trip &#34;)
			price = price * 2
		} else {
			fmt.Print(&#34;One-way    &#34;)
		}

		fmt.Printf(&#34;$%4v&#34;, price)
		fmt.Println()

		count = count &#43; 1
	}
}
