package main

import (
	&#34;database/sql&#34;
	&#34;fmt&#34;

	_ &#34;github.com/lib/pq&#34;

	&#34;time&#34;
)

func main() {

	/*
	 * Local server has following settings in postgresql.conf:
	 * timzone = UTC
	 *
	 * Table looks like:
	 *
	 * CREATE TABLE &#34;public&#34;.&#34;time&#34; (&#34;t&#34; timestamp(6) WITH TIME ZONE)
	 */

	db, err := sql.Open(&#34;postgres&#34;, &#34;dbname=postgres host=localhost sslmode=disable&#34;)
	if err != nil {
		panic(err)
	}

	if _, err := db.Exec(`TRUNCATE TABLE time`); err != nil {
		panic(err)
	}

	tin := time.Date(2013, 5, 4, 3, 2, 1, 0, time.UTC)

	if _, err := db.Exec(`INSERT into time (&#34;t&#34;) VALUES ($1)`, tin); err != nil {
		panic(err)
	}

	var tout time.Time
	if err := db.QueryRow(`SELECT t FROM time`).Scan(&amp;tout); err != nil {
		// we expect a row since we just put it in...
		panic(err)
	}

	fmt.Printf(&#34;In %T %v\nOut %T %v&#34;, tin, tin, tout, tout)

	// Output:
	// In time.Time 2013-05-04 03:02:01 &#43;0000 UTC
	// Out time.Time 2013-05-04 03:02:01 &#43;0000 &#43;0000
}
