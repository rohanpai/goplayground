package main

import (
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

func main() {
	// connect to localhost, make sure to have redis-server running on the default port
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// add some keys
	if _, err = conn.Do("SET", "k1", "a"); err != nil {
		log.Fatal(err)
	}
	if _, err = conn.Do("SET", "k2", "b"); err != nil {
		log.Fatal(err)
	}
	// for fun, let's leave k3 non-existing

	// get many keys in a single MGET, ask redigo for []string result
	strs, err := redis.Strings(conn.Do("MGET", "k1", "k2", "k3"))
	if err != nil {
		log.Fatal(err)
	}

	// prints [a b ]
	fmt.Println(strs)

	// now what if we want some integers instead?
	if _, err = conn.Do("SET", "k4", "1"); err != nil {
		log.Fatal(err)
	}
	if _, err = conn.Do("SET", "k5", "2"); err != nil {
		log.Fatal(err)
	}

	// get the keys, but ask redigo to give us a []interface{}
	// (it doesn't have a redis.Ints helper).
	vals, err := redis.Values(conn.Do("MGET", "k4", "k5", "k6"))
	if err != nil {
		log.Fatal(err)
	}

	// scan the []interface{} slice into a []int slice
	var ints []int
	if err = redis.ScanSlice(vals, &ints); err != nil {
		log.Fatal(err)
	}

	// prints [1 2 0]
	fmt.Println(ints)
}
