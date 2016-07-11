package main

import (
        &#34;fmt&#34;
        &#34;regexp&#34;
)

func main() {
        src := []byte(&#34;eo eo eo eo&#34;)
        search := regexp.MustCompile(&#34;e&#34;)
        repl := []byte(&#34;AEI&#34;)

        i := 0
	src = search.ReplaceAllFunc(src, func(s []byte) []byte {
		if i &lt; 2 {
			i &#43;= 1
			return repl
		}
		return s
	})
	
        fmt.Println(string(src))
}