package main

import (
	"database/sql"
	"fmt"
	"log"
	//	"os"
	//"os/user"
	"time"

	//_ "github.com/Go-SQL-Driver/MySQL"
	//_ "github.com/bmizerany/pq"
	_ "github.com/mattn/go-sqlite3"
)

const DB_NAME = "foo_test"

const (
	//MySQL = iota
	//Postgres
	SQLite = 0
)

var engine = map[int]string{
	//0: "MySQL",
	//1: "PostgreSQL",
	0: "SQLite",
}

var dbs = make([]*sql.DB, 1)

//var datetime = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
var datetime = time.Now().UTC()

var CREATE = map[int]string{
	//MySQL:    "CREATE TABLE times (id INT, datetime TIMESTAMP)",
	//Postgres: "CREATE TABLE times (id integer, datetime timestamp with time zone)",
	SQLite: "CREATE TABLE IF NOT EXISTS times (id INTEGER PRIMARY KEY NOT NULL, datetime TEXT)",
}

var INSERT = fmt.Sprintf("INSERT OR REPLACE INTO times (id, datetime) VALUES(0, '%s')",
	datetime.Format(time.RFC3339))

const (
	//DROP   = "DROP TABLE times"
	SELECT = "SELECT * FROM times WHERE id = 0"
)

// For SQL table "times"
type Times struct {
	Id       int
	Datetime time.Time
}

func (t Times) String() string {
	return fmt.Sprintf("{%d, %s}", t.Id, t.Datetime.Format(time.RFC3339))
}

func (t *Times) Args1() []interface{} {
	return []interface{}{&t.Id, &t.Datetime}
}

func (t *Times) Scan(src interface{}) (err error) {
	value, ok := src.(string)
	if ok {
		t.Datetime, err = time.Parse(time.RFC3339, value)
	} else {
		err = fmt.Errorf("Unexpected type: %T", src)
	}
	return
}

func (t *Times) Args2() []interface{} {
	return []interface{}{&t.Id, t}
}

func main() {
	log.SetFlags(0)
	/*
		u, err := user.Current()
		if err != nil {
			panic(err)
		}
		username := u.Username

		// MySQL

		dbMySQL, err := sql.Open("mysql", fmt.Sprintf("%s@unix(%s)/%s",
			username, "/var/run/mysqld/mysqld.sock", DB_NAME))
		if err != nil {
			panic(err)
		}

		if _, err = dbMySQL.Exec(CREATE[MySQL]); err != nil {
			panic(err)
		}

		dbs[0] = dbMySQL

		// PostgreSQL

		dbPostgres, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s host=%s sslmode=disable",
			username, DB_NAME, "/var/run/postgresql"))
		if err != nil {
			panic(err)
		}

		if _, err = dbPostgres.Exec(CREATE[Postgres]); err != nil {
			panic(err)
		}

		dbs[1] = dbPostgres
	*/
	// SQLite3

	filename := DB_NAME + ".db"
	//defer os.Remove(filename)

	dbSQLite, err := sql.Open("sqlite3", filename)
	if err != nil {
		panic(err)
	}

	if _, err = dbSQLite.Exec(CREATE[SQLite]); err != nil {
		panic(err)
	}

	dbs[0] = dbSQLite

	// == Insert

	input := Times{0, datetime} // SQL statement in INSERT
	output := &Times{}

	for i, db := range dbs {
		if _, err := db.Exec(INSERT); err != nil {
			log.Print(err)
		}

		// Scan output
		rows := db.QueryRow(SELECT)
		fmt.Println("== Testing with Args1()")

		if err = rows.Scan(output.Args1()...); err != nil {
			log.Printf("%s: %s\n\n", engine[i], err)
		} else {
			if fmt.Sprintf("%s", input) != fmt.Sprintf("%s", output) {
				log.Printf("%s: got different data\ninput:  %v\noutput: %v\n\n",
					engine[i], input, output)
			}
		}

		rows = db.QueryRow(SELECT)
		fmt.Println("== Testing with Args2()")

		if err = rows.Scan(output.Args2()...); err != nil {
			log.Printf("%s: %s\n\n", engine[i], err)
		} else {
			if fmt.Sprintf("%v", input) != fmt.Sprintf("%v", output) {
				log.Printf("%s: got different data\ninput:  %v\noutput: %v\n\n",
					engine[i], input, output)
			}
		}
	}

	// == Close

	for _, db := range dbs {
		/*if _, err = db.Exec(DROP); err != nil {
			log.Print(err)
		}*/
		db.Close()
	}
}
