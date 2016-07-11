package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // init() only
	"log"
	"os"
)

func main() {
	const database = "/tmp/db_test.db"

	os.Remove(database)

	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sql := `
	create table user(
		id integer not null primary key,
		name text,
		dummy text
	)
	`
	_, err = db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	sql = `insert into user(name) values(?)`
	for _, name := range []string{"taro", "jiro", "saburo"} {
		_, err = db.Exec(sql, name)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()

	sql = `select * from user`
	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name, dummy string
		rows.Scan(&id, &name, &dummy)
		fmt.Printf("id: %d, name: %s, dummy: %s\n", id, name, dummy)
	}
	fmt.Println()

	sql = `select * from user`
	users, err := query(db, sql)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#+v\n", users)
	fmt.Println()

	var id = 2
	sql = `select * from user where id = ?`
	users, err = query(db, sql, id)
	if err != nil {
		log.Fatal(err)
	}
	if len(users) > 0 {
		fmt.Printf(`user(id:%d) is "%s"`, id, users[0]["name"].([]byte))
	}
}

func query(db *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	res := make([]map[string]interface{}, 0)

	for rows.Next() {
		container := make([]interface{}, len(cols))
		dest := make([]interface{}, len(cols))
		for i, _ := range container {
			dest[i] = &container[i]
		}
		rows.Scan(dest...)
		r := make(map[string]interface{})
		for i, colname := range cols {
			val := dest[i].(*interface{})
			r[colname] = *val
		}
		res = append(res, r)
	}

	return res, nil
}
