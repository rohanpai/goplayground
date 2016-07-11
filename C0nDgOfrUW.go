// == MySQL
//
// To create the database:
//
//   mysql -p
//   mysql> create database foo_test;
//   mysql> GRANT ALL PRIVILEGES ON modsql_test.* to USER@localhost;
//
// Note: substitute "USER" by your user name.
//
// To remove it:
//
//   mysql> drop database foo_test;

// == PostgreSQL
//
// To create the database:
//
//   sudo -u postgres createuser USER --no-superuser --no-createrole --no-createdb
//   sudo -u postgres createdb foo_test --owner USER
//
// Note: substitute "USER" by your user name.
//
// To remove it:
//
//   sudo -u postgres dropdb foo_test

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/user"
	"time"

	_ "github.com/Go-SQL-Driver/MySQL"
	_ "github.com/bmizerany/pq"
	_ "github.com/mattn/go-sqlite3"
)

const DB_NAME = "foo_test"

const (
	MySQL = iota
	Postgres
	SQLite
)

var engine = map[int]string{
	0: "MySQL",
	1: "PostgreSQL",
	2: "SQLite",
}

var dbs = make([]*sql.DB, 3)

var datetime = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

var CREATE = map[int]string{
	MySQL:    "CREATE TABLE times (id INT, datetime TIMESTAMP)",
	Postgres: "CREATE TABLE times (id integer, datetime timestamp without time zone)",
	SQLite:   "CREATE TABLE times (id INTEGER, datetime TEXT)",
}

var INSERT = fmt.Sprintf("INSERT INTO times (id, datetime) VALUES(0, '%s')",
	datetime.Format(time.RFC3339))

const (
	DROP   = "DROP TABLE times"
	SELECT = "SELECT * FROM times WHERE id = 0"
)

// For SQL table "times"
type Times struct {
	Id       int
	Datetime time.Time
}

func (t Times) Args1() []interface{} {
	return []interface{}{&t.Id, &t.Datetime}
}

func (t Times) Args2() []interface{} {
	tt := t.Datetime.Format(time.RFC3339)
	return []interface{}{&t.Id, &tt}
}

func main() {
	log.SetFlags(0)

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

	// SQLite3

	filename := DB_NAME + ".db"
	defer os.Remove(filename)

	dbSQLite, err := sql.Open("sqlite3", filename)
	if err != nil {
		panic(err)
	}

	if _, err = dbSQLite.Exec(CREATE[SQLite]); err != nil {
		panic(err)
	}

	dbs[2] = dbSQLite

	// == Insert

	input := Times{0, datetime} // SQL statement in INSERT
	output := Times{}

	for i, db := range dbs {
		if _, err := db.Exec(INSERT); err != nil {
			log.Print(err)
		}

		// Scan output
		rows := db.QueryRow(SELECT)
		if err = rows.Scan(output.Args1()...); err != nil {
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
		if _, err = db.Exec(DROP); err != nil {
			log.Print(err)
		}
		db.Close()
	}
}