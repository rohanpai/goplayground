package main

import (
	&#34;fmt&#34;
	&#34;log&#34;

	&#34;github.com/garyburd/redigo/redis&#34;
)

func main() {
	// connect to localhost, make sure to have redis-server running on the default port
	conn, err := redis.Dial(&#34;tcp&#34;, &#34;:6379&#34;)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// add some keys
	if _, err = conn.Do(&#34;SET&#34;, &#34;k1&#34;, &#34;a&#34;); err != nil {
		log.Fatal(err)
	}
	if _, err = conn.Do(&#34;SET&#34;, &#34;k2&#34;, &#34;b&#34;); err != nil {
		log.Fatal(err)
	}
	// for fun, let&#39;s leave k3 non-existing

	// get many keys in a single MGET, ask redigo for []string result
	strs, err := redis.Strings(conn.Do(&#34;MGET&#34;, &#34;k1&#34;, &#34;k2&#34;, &#34;k3&#34;))
	if err != nil {
		log.Fatal(err)
	}

	// prints [a b ]
	fmt.Println(strs)

	// now what if we want some integers instead?
	if _, err = conn.Do(&#34;SET&#34;, &#34;k4&#34;, &#34;1&#34;); err != nil {
		log.Fatal(err)
	}
	if _, err = conn.Do(&#34;SET&#34;, &#34;k5&#34;, &#34;2&#34;); err != nil {
		log.Fatal(err)
	}

	// get the keys, but ask redigo to give us a []interface{}
	// (it doesn&#39;t have a redis.Ints helper).
	vals, err := redis.Values(conn.Do(&#34;MGET&#34;, &#34;k4&#34;, &#34;k5&#34;, &#34;k6&#34;))
	if err != nil {
		log.Fatal(err)
	}

	// scan the []interface{} slice into a []int slice
	var ints []int
	if err = redis.ScanSlice(vals, &amp;ints); err != nil {
		log.Fatal(err)
	}

	// prints [1 2 0]
	fmt.Println(ints)
}
