package main

import &#34;fmt&#34;

func main() {
	g0, g1 := generate(&#34;counting the number of some arbitrary values&#34;)
	even, odd := count(g0, &#39;e&#39;), count(g1, &#39;o&#39;)
	for i := 0; i &lt; 2; i&#43;&#43; {
		select {
		case x := &lt;-even:
			fmt.Printf(&#34;%d even &#39;e&#39;s\n&#34;, x)
		case x := &lt;-odd:
			fmt.Printf(&#34;%d odd &#39;o&#39;s\n&#34;, x)
		}
	}
}

// count returns number of x&#39;s in found in c.
func count(c &lt;-chan int, x int) (&lt;-chan int) {
	r := make(chan int)
	go func() {
		sum := 0
		for i := range c {
			if i == x {
				sum&#43;&#43;
			}
		}
		r &lt;- sum
	}()
	return r
}

// generate sends all even-indexed runes in s to the first
// returned channel, and all odd-indexed runes to the other
// returned channel.
func generate(s string) (&lt;-chan int, &lt;-chan int) {
	even := make(chan int)
	odd := make(chan int)
	go func() {
		i := 0
		for _, r := range s {
			if i % 2 == 0 {
				even &lt;- r
			}else{
				odd &lt;- r
			}
			i&#43;&#43;
		}
		close(even)
		close(odd)
	}()
	return even, odd
}
