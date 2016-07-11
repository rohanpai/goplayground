package main

import (
	"bitbucket.org/kardianos/table"
	_ "code.google.com/p/odbc"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	connStr := "driver=sql server;server=(local);database=tempdb;trusted_connection=yes"
	if len(os.Args) > 1 && len(os.Args[1]) > 1 {
		connStr = os.Args[1]
	}
	//log.Printf("Connecting with '%s'\n", connStr)
	conn, err := sql.Open("odbc", connStr)
	if err != nil {
		fmt.Println("Connecting Error")
		return
	}
	defer conn.Close()

	connTbl := "sys.databases"
	if len(os.Args) > 2 {
		connTbl = os.Args[2]
	}
	table, err := table.Get(conn, "select * from "+connTbl)
	if err != nil {
		log.Fatal(err)
	}

	dumpTable(table)

	fmt.Fprintf(os.Stderr, "\nFinished correctly\n")
	return
}

func dumpTable(table *table.Buffer) {
	// open the output file
	file, err := os.Create("out.csv")
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
			fmt.Fprintf(file, ",")
		}
		fmt.Fprintf(file, "\"%s\"", element)
	}
	fmt.Fprintf(file, "\n")

	// output body
	const layout = "01/02/2006 15:04:05.999"
	for _, row := range table.Rows {
		for i, colname := range table.ColumnName {
			if i != 0 {
				fmt.Fprintf(file, ",")
			}
			switch x := row.MustGet(colname).(type) {
			case string: // x is a string
				fmt.Fprintf(file, "\"%s\"", x)
			case int: // now x is an int
				fmt.Fprintf(file, "\"%d\"", x)
			case int32: // now x is an int32
				fmt.Fprintf(file, "\"%d\"", x)
			case int64: // now x is an int64
				fmt.Fprintf(file, "\"%d\"", x)
			case float32: // now x is an float32
				fmt.Fprintf(file, "\"%f\"", x)
			case float64: // now x is an float64
				fmt.Fprintf(file, "\"%f\"", x)
			case time.Time: // now x is a time.Time
				fmt.Fprintf(file, "\"%s\"", x.Format(layout))
			default:
				fmt.Fprintf(file, "\"%s\"", x)
			}
		}
		fmt.Fprintf(file, "\n")
		fmt.Fprintf(os.Stderr, ".")
	}
	fmt.Fprintf(os.Stderr, "\n")

}
