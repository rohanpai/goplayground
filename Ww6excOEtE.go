package main

import (
	&#34;fmt&#34;
	_ &#34;github.com/jackc/pgx/stdlib&#34;
	&#34;github.com/jmoiron/sqlx&#34;
)

func setDB() *sqlx.DB {
	db, err := sqlx.Connect(&#34;pgx&#34;,&#34;postgresql://user:pass@localhost:5433/mydb&#34;)
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
