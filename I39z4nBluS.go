package main

import (
	&#34;database/sql&#34;
	&#34;fmt&#34;
	&#34;log&#34;
	&#34;time&#34;

	//_ &#34;code.google.com/p/gosqlite/sqlite3&#34;
	_ &#34;github.com/mattn/go-sqlite3&#34;
	//_ &#34;github.com/mxk/go-sqlite/sqlite3&#34;
)

type TimeStamp struct{ *time.Time }

func (t TimeStamp) Scan(value interface{}) error {
	fmt.Printf(&#34;%T\n&#34;, value)
	//...
	return nil
}

func main() {
	db, err := sql.Open(&#34;sqlite3&#34;, &#34;:memory:&#34;)
	if err != nil {
		log.Fatalf(&#34;cannot open an SQLite memory database: %v&#34;, err)
	}
	defer db.Close()

	// sqlite&gt; select strftime(&#39;%J&#39;, &#39;2015-04-13T19:22:19.773Z&#39;), strftime(&#39;%J&#39;, &#39;2015-04-13T19:22:19&#39;);
	_, err = db.Exec(&#34;CREATE TABLE unix_time (time datetime); INSERT INTO unix_time (time) VALUES (strftime(&#39;%Y-%m-%dT%H:%MZ&#39;,&#39;now&#39;))&#34;)
	if err != nil {
		log.Fatalf(&#34;cannot create schema: %v&#34;, err)
	}

	row := db.QueryRow(&#34;SELECT time FROM unix_time&#34;)
	var t time.Time
	err = row.Scan(TimeStamp{&amp;t})
	if err != nil {
		log.Fatalf(&#34;cannot scan time: %v&#34;, err)
	}
	fmt.Println(t)
}
