package main

import (
	&#34;database/sql&#34;
	&#34;flag&#34;
	&#34;fmt&#34;
	_ &#34;github.com/mattn/go-oci8&#34;
	&#34;log&#34;
	&#34;os&#34;
	&#34;strings&#34;
)

var conn string

func main() {
	flag.Parse()
	if flag.NArg() &gt;= 1 {
		conn = flag.Arg(0)
	} else {
		conn = &#34;system/123456@XE&#34;
	}

	db, err := sql.Open(&#34;oci8&#34;, conn)
	if err != nil {
		fmt.Println(&#34;can&#39;t connect &#34;, conn, err)
		return
	}
	if err = test_conn(db); err != nil {
		fmt.Println(&#34;can&#39;t connect &#34;, conn, err)
		return
	}
	var in string
	var sqlquery string
	fmt.Print(&#34;&gt; &#34;)
	for {
		fmt.Scan(&amp;in)
		if in == &#34;q;&#34; {
			break
		}
		if in[len(in)-1] != &#39;;&#39; {
			sqlquery &#43;= in &#43; &#34; &#34;
		} else {
			sqlquery &#43;= in[:len(in)-1]
			rows, err := db.Query(sqlquery)
			if err != nil {
				fmt.Println(&#34;can&#39;t run &#34;, sqlquery, &#34;\n&#34;, err)
				fmt.Print(&#34;&gt; &#34;)
				sqlquery = &#34;&#34;
				continue
			}
			cols, err := rows.Columns()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(strings.Join(cols, &#34;\t&#34;))
			var result = make([]string, len(cols))
			var s = make([]interface{}, len(result))
			for i, _ := range result {
				s[i] = &amp;result[i]
			}
			for rows.Next() {
				rows.Scan(s...)
				fmt.Println(strings.Join(result, &#34;\t&#34;))
			}
			rows.Close()
			fmt.Print(&#34;&gt; &#34;)
			sqlquery = &#34;&#34;
		}
	}

	db.Close()
}

func test_conn(db *sql.DB) (err error) {
	query := &#34;select * from dual&#34;
	_, err = db.Query(query)
	return err
}
