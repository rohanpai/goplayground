package main

import &#34;fmt&#34;
import &#34;time&#34;

func main() {
	fmt.Println(&#34;START 1&#34;)
	for i := 0; i &lt; 3; i&#43;&#43; {
		foo(i)
	}
	fmt.Println(&#34;END 1&#34;)

	fmt.Println(&#34;START 2&#34;)
	for i := 0; i &lt; 3; i&#43;&#43; {
		go foo(i)
	}
	fmt.Println(&#34;END 2&#34;)

	time.Sleep(1 * time.Second)
}

func foo(i int) {
	time.Sleep(1)
	fmt.Printf(&#34;Call index: %d\n&#34;, i)
}
