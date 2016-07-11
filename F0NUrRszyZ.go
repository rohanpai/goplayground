package main

import (
	&#34;database/sql&#34;
	&#34;fmt&#34;
	&#34;log&#34;
	//	&#34;os&#34;
	//&#34;os/user&#34;
	&#34;time&#34;

	//_ &#34;github.com/Go-SQL-Driver/MySQL&#34;
	//_ &#34;github.com/bmizerany/pq&#34;
	_ &#34;github.com/mattn/go-sqlite3&#34;
)

const DB_NAME = &#34;foo_test&#34;

const (
	//MySQL = iota
	//Postgres
	SQLite = 0
)

var engine = map[int]string{
	//0: &#34;MySQL&#34;,
	//1: &#34;PostgreSQL&#34;,
	0: &#34;SQLite&#34;,
}

var dbs = make([]*sql.DB, 1)

//var datetime = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
var datetime = time.Now().UTC()

var CREATE = map[int]string{
	//MySQL:    &#34;CREATE TABLE times (id INT, datetime TIMESTAMP)&#34;,
	//Postgres: &#34;CREATE TABLE times (id integer, datetime timestamp with time zone)&#34;,
	SQLite: &#34;CREATE TABLE IF NOT EXISTS times (id INTEGER PRIMARY KEY NOT NULL, datetime TEXT)&#34;,
}

var INSERT = fmt.Sprintf(&#34;INSERT OR REPLACE INTO times (id, datetime) VALUES(0, &#39;%s&#39;)&#34;,
	datetime.Format(time.RFC3339))

const (
	//DROP   = &#34;DROP TABLE times&#34;
	SELECT = &#34;SELECT * FROM times WHERE id = 0&#34;
)

// For SQL table &#34;times&#34;
type Times struct {
	Id       int
	Datetime time.Time
}

func (t Times) String() string {
	return fmt.Sprintf(&#34;{%d, %s}&#34;, t.Id, t.Datetime.Format(time.RFC3339))
}

func (t *Times) Args1() []interface{} {
	return []interface{}{&amp;t.Id, &amp;t.Datetime}
}

func (t *Times) Scan(src interface{}) (err error) {
	value, ok := src.(string)
	if ok {
		t.Datetime, err = time.Parse(time.RFC3339, value)
	} else {
		err = fmt.Errorf(&#34;Unexpected type: %T&#34;, src)
	}
	return
}

func (t *Times) Args2() []interface{} {
	return []interface{}{&amp;t.Id, t}
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
	*/
	// SQLite3

	filename := DB_NAME &#43; &#34;.db&#34;
	//defer os.Remove(filename)

	dbSQLite, err := sql.Open(&#34;sqlite3&#34;, filename)
	if err != nil {
		panic(err)
	}

	if _, err = dbSQLite.Exec(CREATE[SQLite]); err != nil {
		panic(err)
	}

	dbs[0] = dbSQLite

	// == Insert

	input := Times{0, datetime} // SQL statement in INSERT
	output := &amp;Times{}

	for i, db := range dbs {
		if _, err := db.Exec(INSERT); err != nil {
			log.Print(err)
		}

		// Scan output
		rows := db.QueryRow(SELECT)
		fmt.Println(&#34;== Testing with Args1()&#34;)

		if err = rows.Scan(output.Args1()...); err != nil {
			log.Printf(&#34;%s: %s\n\n&#34;, engine[i], err)
		} else {
			if fmt.Sprintf(&#34;%s&#34;, input) != fmt.Sprintf(&#34;%s&#34;, output) {
				log.Printf(&#34;%s: got different data\ninput:  %v\noutput: %v\n\n&#34;,
					engine[i], input, output)
			}
		}

		rows = db.QueryRow(SELECT)
		fmt.Println(&#34;== Testing with Args2()&#34;)

		if err = rows.Scan(output.Args2()...); err != nil {
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
		/*if _, err = db.Exec(DROP); err != nil {
			log.Print(err)
		}*/
		db.Close()
	}
}
