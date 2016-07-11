package main

import (
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

func setDB() *sqlx.DB {
	db, err := sqlx.Connect("pgx","postgresql://user:pass@localhost:5433/mydb")
	if err != nil {
		fmt.Println(err)
	}
	return db
}

func main() {
	db := setDB()
	rows, err := db.Query(`SELECT coordinates FROM init_locations`)
	......
}
