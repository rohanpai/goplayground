package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-adodb"
	"os"
)

func main() {
	if _, err := os.Stat("./example.mdb"); err != nil {
		fmt.Println("put here empty database named 'example.mdb'.")
		return
	}

	db, err := sql.Open("adodb", "Provider=Microsoft.Jet.OLEDB.4.0;Data Source=./example.mdb;")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sqls := []string{
		"DROP TABLE languages",
		"CREATE TABLE languages (id text not null primary key, name text)",
	}
	for _, sql := range sqls {
		_, err = db.Exec(sql)
		if err != nil {
			fmt.Printf("%q: %s\n", err, sql)
			return
		}
	}

	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}

	stmt, err := tx.Prepare("insert into languages(id, name) values(?, ?)")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec("en", "English")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = stmt.Exec("fr", "French")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = stmt.Exec("de", "German")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = stmt.Exec("es", "Spanish")
	if err != nil {
		fmt.Println(err)
		return
	}

	tx.Commit()
}