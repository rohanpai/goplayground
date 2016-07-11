package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"log"
	"os"
	"strings"
)

var conn string

func main() {
	flag.Parse()
	if flag.NArg() >= 1 {
		conn = flag.Arg(0)
	} else {
		conn = "system/123456@XE"
	}

	db, err := sql.Open("oci8", conn)
	if err != nil {
		fmt.Println("can't connect ", conn, err)
		return
	}
	if err = test_conn(db); err != nil {
		fmt.Println("can't connect ", conn, err)
		return
	}
	var in string
	var sqlquery string
	fmt.Print("> ")
	for {
		fmt.Scan(&in)
		if in == "q;" {
			break
		}
		if in[len(in)-1] != ';' {
			sqlquery += in + " "
		} else {
			sqlquery += in[:len(in)-1]
			rows, err := db.Query(sqlquery)
			if err != nil {
				fmt.Println("can't run ", sqlquery, "\n", err)
				fmt.Print("> ")
				sqlquery = ""
				continue
			}
			cols, err := rows.Columns()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(strings.Join(cols, "\t"))
			var result = make([]string, len(cols))
			var s = make([]interface{}, len(result))
			for i, _ := range result {
				s[i] = &result[i]
			}
			for rows.Next() {
				rows.Scan(s...)
				fmt.Println(strings.Join(result, "\t"))
			}
			rows.Close()
			fmt.Print("> ")
			sqlquery = ""
		}
	}

	db.Close()
}

func test_conn(db *sql.DB) (err error) {
	query := "select * from dual"
	_, err = db.Query(query)
	return err
}
