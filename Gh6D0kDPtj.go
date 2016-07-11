package main

import (
	&#34;fmt&#34;
	&#34;log&#34;
	&#34;reflect&#34;
)

func shouldEscape(c byte) bool {
	switch c {
	case &#39; &#39;, &#39;?&#39;, &#39;&amp;&#39;, &#39;=&#39;, &#39;#&#39;, &#39;&#43;&#39;, &#39;%&#39;:
		return true
	}
	return false
}

func main() {
	fmt.Println(shouldEscape([]byte(&#34;?&#34;)[0]))     // true
	fmt.Println(shouldEscape([]byte(&#34;abcd#&#34;)[4])) // true
	fmt.Println(shouldEscape([]byte(&#34;abcd#&#34;)[0])) // false

	num := 2
	switch num {
	case 1:
		fmt.Println(1)
	case 2:
		fmt.Println(2)
	case 3:
		fmt.Println(3)
	default:
		panic(&#34;what&#39;s the number?&#34;)
	}
	// 2

	st := &#34;b&#34;
	switch {
	case st == &#34;a&#34;:
		fmt.Println(&#34;a&#34;)
	case st == &#34;b&#34;:
		fmt.Println(&#34;b&#34;)
	case st == &#34;c&#34;:
		fmt.Println(&#34;c&#34;)
	default:
		panic(&#34;what&#39;s the character?&#34;)
	}
	// b

	ts := []interface{}{true, 1, 1.5, &#34;A&#34;}
	for _, t := range ts {
		eval(t)
	}
	/*
	   bool: true is bool
	   int: 1 is int
	   float64: 1.5 is float64
	   string: A is string
	*/

	type temp struct {
		a string
	}
	eval(interface{}(temp{}))
	// 2009/11/10 23:00:00 {} is main.temp
}

func eval(t interface{}) {
	switch typedValue := t.(type) {
	default:
		log.Fatalf(&#34;%v is %v&#34;, typedValue, reflect.TypeOf(typedValue))
	case bool:
		fmt.Println(&#34;bool:&#34;, typedValue, &#34;is&#34;, reflect.TypeOf(typedValue))
	case int:
		fmt.Println(&#34;int:&#34;, typedValue, &#34;is&#34;, reflect.TypeOf(typedValue))
	case float64:
		fmt.Println(&#34;float64:&#34;, typedValue, &#34;is&#34;, reflect.TypeOf(typedValue))
	case string:
		fmt.Println(&#34;string:&#34;, typedValue, &#34;is&#34;, reflect.TypeOf(typedValue))
	}
}
