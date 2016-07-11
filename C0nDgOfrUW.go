// == MySQL
//
// To create the database:
//
//   mysql -p
//   mysql&gt; create database foo_test;
//   mysql&gt; GRANT ALL PRIVILEGES ON modsql_test.* to USER@localhost;
//
// Note: substitute &#34;USER&#34; by your user name.
//
// To remove it:
//
//   mysql&gt; drop database foo_test;

// == PostgreSQL
//
// To create the database:
//
//   sudo -u postgres createuser USER --no-superuser --no-createrole --no-createdb
//   sudo -u postgres createdb foo_test --owner USER
//
// Note: substitute &#34;USER&#34; by your user name.
//
// To remove it:
//
//   sudo -u postgres dropdb foo_test

package main

import (
	&#34;database/sql&#34;
	&#34;fmt&#34;
	&#34;log&#34;
	&#34;os&#34;
	&#34;os/user&#34;
	&#34;time&#34;

	_ &#34;github.com/Go-SQL-Driver/MySQL&#34;
	_ &#34;github.com/bmizerany/pq&#34;
	_ &#34;github.com/mattn/go-sqlite3&#34;
)

const DB_NAME = &#34;foo_test&#34;

const (
	MySQL = iota
	Postgres
	SQLite
)

var engine = map[int]string{
	0: &#34;MySQL&#34;,
	1: &#34;PostgreSQL&#34;,
	2: &#34;SQLite&#34;,
}

var dbs = make([]*sql.DB, 3)

var datetime = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

var CREATE = map[int]string{
	MySQL:    &#34;CREATE TABLE times (id INT, datetime TIMESTAMP)&#34;,
	Postgres: &#34;CREATE TABLE times (id integer, datetime timestamp without time zone)&#34;,
	SQLite:   &#34;CREATE TABLE times (id INTEGER, datetime TEXT)&#34;,
}

var INSERT = fmt.Sprintf(&#34;INSERT INTO times (id, datetime) VALUES(0, &#39;%s&#39;)&#34;,
	datetime.Format(time.RFC3339))

const (
	DROP   = &#34;DROP TABLE times&#34;
	SELECT = &#34;SELECT * FROM times WHERE id = 0&#34;
)

// For SQL table &#34;times&#34;
type Times struct {
	Id       int
	Datetime time.Time
}

func (t Times) Args1() []interface{} {
	return []interface{}{&amp;t.Id, &amp;t.Datetime}
}

func (t Times) Args2() []interface{} {
	tt := t.Datetime.Format(time.RFC3339)
	return []interface{}{&amp;t.Id, &amp;tt}
}

func main() {
	log.SetFlags(0)

	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	username := u.Username

	// MySQL

	dbMySQL, err := sql.Open(&#34;mysql&#34;, fmt.Sprintf(&#34;%s@unix(%s)/%s&#34;,
		username, &#34;/var/run/mysqld/mysqld.sock&#34;, DB_NAME))
	if err != nil {
		panic(err)
	}

	if _, err = dbMySQL.Exec(CREATE[MySQL]); err != nil {
		panic(err)
	}

	dbs[0] = dbMySQL

	// PostgreSQL

	dbPostgres, err := sql.Open(&#34;postgres&#34;, fmt.Sprintf(&#34;user=%s dbname=%s host=%s sslmode=disable&#34;,
		username, DB_NAME, &#34;/var/run/postgresql&#34;))
	if err != nil {
		panic(err)
	}

	if _, err = dbPostgres.Exec(CREATE[Postgres]); err != nil {
		panic(err)
	}

	dbs[1] = dbPostgres

	// SQLite3

	filename := DB_NAME &#43; &#34;.db&#34;
	defer os.Remove(filename)

	dbSQLite, err := sql.Open(&#34;sqlite3&#34;, filename)
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
			log.Printf(&#34;%s: %s\n\n&#34;, engine[i], err)
		} else {
			if fmt.Sprintf(&#34;%v&#34;, input) != fmt.Sprintf(&#34;%v&#34;, output) {
				log.Printf(&#34;%s: got different data\ninput:  %v\noutput: %v\n\n&#34;,
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