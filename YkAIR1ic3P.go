package main

import (
	&#34;bytes&#34;
	&#34;fmt&#34;
)

var trivialXml = []byte(`&lt;root&gt;a&lt;foo&gt;b&lt;/foo&gt;c&lt;bar&gt;d&lt;/bar&gt;e&lt;bar&gt;f&lt;/bar&gt;g&lt;/root&gt;`)

func main() {
	node, err := xmlpath.Parse(bytes.NewBuffer(trivialXml))
	if err != nil {
		panic(err)
	}
	path1 := xmlpath.MustCompile(&#34;root&#34;)
	path2 := xmlpath.MustCompile(&#34;foo&#34;)
	result1, ok1 := path1.String(node)
	result2, ok2 := path2.String(node)
	iter := path1.Iter(node)
	if !iter.Next() {
		panic(&#34;must exist&#34;)
	}
	result3, ok3 := path2.String(iter.Node())
	fmt.Printf(&#34;result1: %v %s\n&#34;, ok1, result1)
	fmt.Printf(&#34;result2: %v %s\n&#34;, ok2, result2)
	fmt.Printf(&#34;result3: %v %s\n&#34;, ok3, result3)
}
