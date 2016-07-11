package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"time"
)

func main() {

	/*
	 * Local server has following settings in postgresql.conf:
	 * timzone = UTC
	 *
	 * Table looks like:
	 *
	 * CREATE TABLE "public"."time" ("t" timestamp(6) WITH TIME ZONE)
	 */

	db, err := sql.Open("postgres", "dbname=postgres host=localhost sslmode=disable")
	if err != nil {
		panic(err)
	}

	if _, err := db.Exec(`TRUNCATE TABLE time`); err != nil {
		panic(err)
	}

	tin := time.Date(2013, 5, 4, 3, 2, 1, 0, time.UTC)

	if _, err := db.Exec(`INSERT into time ("t") VALUES ($1)`, tin); err != nil {
		panic(err)
	}

	var tout time.Time
	if err := db.QueryRow(`SELECT t FROM time`).Scan(&tout); err != nil {
		// we expect a row since we just put it in...
		panic(err)
	}

	fmt.Printf("In %T %v\nOut %T %v", tin, tin, tout, tout)

	// Output:
	// In time.Time 2013-05-04 03:02:01 +0000 UTC
	// Out time.Time 2013-05-04 03:02:01 +0000 +0000
}
