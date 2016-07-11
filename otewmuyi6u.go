package main

import &#34;fmt&#34;

const ( // iota is reset to 0
	c0 = iota // c0 == 0
	c1 = iota // c1 == 1
	c2 = iota // c2 == 2
)

const (
	a = 1 &lt;&lt; iota // a == 1 (iota has been reset)
	b = 1 &lt;&lt; iota // b == 2
	c = 1 &lt;&lt; iota // c == 4
)

const (
	u         = iota * 42 // u == 0     (untyped integer constant)
	v float64 = iota * 42 // v == 42.0  (float64 constant)
	w         = iota * 42 // w == 84    (untyped integer constant)
)

const x = iota // x == 0 (iota has been reset)
const y = iota // y == 0 (iota has been reset)
/**
Within an ExpressionList, the value of each iota is the same
because it is only incremented after each ConstSpec:
**/

const (
	bit0, mask0 = 1 &lt;&lt; iota, 1&lt;&lt;iota - 1 // bit0 == 1, mask0 == 0
	bit1, mask1                          // bit1 == 2, mask1 == 1
	_, _                                 // skips iota == 2
	bit3, mask3                          // bit3 == 8, mask3 == 7
)

func main() {

	if x == y &amp;&amp; y == u &amp;&amp; u == c0 {
		fmt.Println(&#34;In each Const Declaration iota reset to 0&#34;)
		fmt.Println(&#34;x, y, u, c0 == 0&#34;)
	}

	fmt.Println(&#34;For each ConstSpec the IOTA gets incremented only by one&#34;)
	fmt.Println(&#34;Within an ExpressionList,\nthe value of each iota is the same&#34;)

	fmt.Println(&#34;\t\tbit0, mask0 = 1 &lt;&lt; iota, 1 &lt;&lt; iota - 1&#34;)
	fmt.Println(&#34;Notice the first ConstSpec of last Const expression&#34;)
	fmt.Println(&#34;The first ConstSpec is defined as \n\t&lt; identifier-1 &gt; &lt; identifier2 &gt; = 1 &lt;&lt; iota, 1 &lt;&lt; iota - 1&#34;)
	fmt.Println(&#34;Therefore, each successive ConstSpec will\n\t1. Increase the iota only once. \n\t2. Use the same expression to evaluate the values of Identifiers&#34;)
	fmt.Println(&#34;\n\tLets see the same behavior in the last consts  &#34;)
	fmt.Println(&#34;bit0, mask0 = &#34;, bit1, mask1, &#34; &amp;&amp; IOTA is 0&#34;)
	fmt.Println(&#34;bit1, mask1 = &#34;, bit1, mask1, &#34; &amp;&amp; IOTA is 1&#34;)
	fmt.Println(&#34;bit2 and mask2 are defined as _, _,\n\t but since it is a next ConstSpec, IOTA gets incremented to 2&#34;)
	fmt.Println(&#34;bit2, mask2 = UNDEfined in Const Spec as _,_&#34;, &#34;&amp;&amp; IOTA is 2&#34;)
	fmt.Println(&#34;bit3, mask3 = &#34;, bit3, mask3, &#34; &amp;&amp; IOTA is 3&#34;)

}
