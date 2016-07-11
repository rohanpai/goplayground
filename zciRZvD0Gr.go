package main

import &#34;fmt&#34;

func PadRight(str, pad string, lenght int) string {
	for {
		str &#43;= pad
		if len(str) &gt; lenght {
			return str[0:lenght]
		}
	}
}

func PadLeft(str, pad string, lenght int) string {
	for {
		str = pad &#43; str
		if len(str) &gt; lenght {
			return str[0:lenght]
		}
	}
}

func main() {

	str := &#34;abc&#34;
	fmt.Println(PadRight(str, &#34;x&#34;, 5))   // expects abcxx
	fmt.Println(PadLeft(str, &#34;x&#34;, 5))    // expects xxabc
	fmt.Println(PadRight(str, &#34;xyz&#34;, 5)) // expects abcxy
	fmt.Println(PadLeft(str, &#34;xyz&#34;, 5))  // expects xyzab
}
