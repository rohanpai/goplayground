package main

import (
	&#34;bitbucket.org/kardianos/table&#34;
	_ &#34;code.google.com/p/odbc&#34;
	&#34;database/sql&#34;
	&#34;fmt&#34;
	&#34;log&#34;
	&#34;os&#34;
	&#34;time&#34;
)

func main() {
	connStr := &#34;driver=sql server;server=(local);database=tempdb;trusted_connection=yes&#34;
	if len(os.Args) &gt; 1 &amp;&amp; len(os.Args[1]) &gt; 1 {
		connStr = os.Args[1]
	}
	//log.Printf(&#34;Connecting with &#39;%s&#39;\n&#34;, connStr)
	conn, err := sql.Open(&#34;odbc&#34;, connStr)
	if err != nil {
		fmt.Println(&#34;Connecting Error&#34;)
		return
	}
	defer conn.Close()

	connTbl := &#34;sys.databases&#34;
	if len(os.Args) &gt; 2 {
		connTbl = os.Args[2]
	}
	table, err := table.Get(conn, &#34;select * from &#34;&#43;connTbl)
	if err != nil {
		log.Fatal(err)
	}

	dumpTable(table)

	fmt.Fprintf(os.Stderr, &#34;\nFinished correctly\n&#34;)
	return
}

func dumpTable(table *table.Buffer) {
	// open the output file
	file, err := os.Create(&#34;out.csv&#34;)
	if err != nil {
		panic(err)
	}
	// close file on exit and check for its returned error
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	// output header
	for i, element := range table.ColumnName {
		if i != 0 {
			fmt.Fprintf(file, &#34;,&#34;)
		}
		fmt.Fprintf(file, &#34;\&#34;%s\&#34;&#34;, element)
	}
	fmt.Fprintf(file, &#34;\n&#34;)

	// output body
	const layout = &#34;01/02/2006 15:04:05.999&#34;
	for _, row := range table.Rows {
		for i, colname := range table.ColumnName {
			if i != 0 {
				fmt.Fprintf(file, &#34;,&#34;)
			}
			switch x := row.MustGet(colname).(type) {
			case string: // x is a string
				fmt.Fprintf(file, &#34;\&#34;%s\&#34;&#34;, x)
			case int: // now x is an int
				fmt.Fprintf(file, &#34;\&#34;%d\&#34;&#34;, x)
			case int32: // now x is an int32
				fmt.Fprintf(file, &#34;\&#34;%d\&#34;&#34;, x)
			case int64: // now x is an int64
				fmt.Fprintf(file, &#34;\&#34;%d\&#34;&#34;, x)
			case float32: // now x is an float32
				fmt.Fprintf(file, &#34;\&#34;%f\&#34;&#34;, x)
			case float64: // now x is an float64
				fmt.Fprintf(file, &#34;\&#34;%f\&#34;&#34;, x)
			case time.Time: // now x is a time.Time
				fmt.Fprintf(file, &#34;\&#34;%s\&#34;&#34;, x.Format(layout))
			default:
				fmt.Fprintf(file, &#34;\&#34;%s\&#34;&#34;, x)
			}
		}
		fmt.Fprintf(file, &#34;\n&#34;)
		fmt.Fprintf(os.Stderr, &#34;.&#34;)
	}
	fmt.Fprintf(os.Stderr, &#34;\n&#34;)

}
