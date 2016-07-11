package main

import (
	&#34;database/sql&#34;
	&#34;fmt&#34;
	//_ &#34;github.com/gwenn/gosqlite&#34;
	_ &#34;code.google.com/p/go-sqlite/go1/sqlite3&#34;
)

type Test struct {
	Id int64
	// ... other fields
	Data []float64
}

const schema = `CREATE TABLE test (
  id INTEGER PRIMARY KEY NOT NULL
-- ... other columns
);

CREATE TABLE test_data (
  test_id INTEGER NOT NULL,
  data REAL NOT NULL,
  FOREIGN KEY (test_id) REFERENCES test(id) ON DELETE CASCADE
);

-- tests
INSERT INTO test VALUES (1);
INSERT INTO test_data VALUES (1, 273.15);
INSERT INTO test_data VALUES (1, 459.67);
INSERT INTO test VALUES (2);
INSERT INTO test_data VALUES (2, 0.239);
INSERT INTO test_data VALUES (2, 2.7778);
`
const query = `SELECT t.id/*, t...other columns */, d.data
FROM test t
LEFT OUTER JOIN test_data d ON d.test_id = t.id
`

func findTests(db *sql.DB) []Test {
	rows, err := db.Query(query)
	checkNoError(err, &#34;error querying tests: %s&#34;)
	defer checkSqlRowsClose(rows)

	tests := make([]Test, 0, 20)
	indexById := make(map[int64]int)
	for rows.Next() {
		var id int64
		var dummy float64
		//err = rows.Scan(&amp;id) // expected 2 destination arguments in Scan, not 1
		err = rows.Scan(&amp;id, &amp;dummy)
		checkNoError(err, &#34;error scanning: %s&#34;)
		index, ok := indexById[id]
		if !ok {
			index = len(tests)
			tests = append(tests, Test{})
			indexById[id] = index
		}
		test := &amp;tests[index]

		var data sql.NullFloat64
		if ok {
			//err = rows.Scan(&amp;_ /*, &amp;_... */, &amp;data) // cannot use _ as value
			//err = rows.Scan(nil /*, nil...*/, &amp;data) // Scan error on column index 1: destination not a pointer
			err = rows.Scan(&amp;dummy /*, &amp;dummy*/, &amp;data)
			checkNoError(err, &#34;error scanning: %s&#34;)
		} else {
			err = rows.Scan(&amp;test.Id /*, &amp;test...other fields */, &amp;data)
			checkNoError(err, &#34;error scanning: %s&#34;)
		}
		if data.Valid {
			if test.Data == nil {
				test.Data = make([]float64, 0)
			}
			test.Data = append(test.Data, data.Float64)
		}
	}
	return tests
}

func main() {
	db, err := sql.Open(&#34;sqlite3&#34;, &#34;:memory:&#34;)
	checkNoError(err, &#34;error opening connection: %s&#34;)
	defer checkSqlDbClose(db)

	_, err = db.Exec(schema)
	checkNoError(err, &#34;error creating table: %s&#34;)
	tests := findTests(db)
	for i, test := range tests {
		fmt.Printf(&#34;&gt; Tests[%d]: %v\n&#34;, i, test)
	}
}

func checkNoError(err error, format string) {
	if err != nil {
		panic(fmt.Sprintf(format, err.Error()))
	}
}
func checkSqlDbClose(db *sql.DB) {
	checkNoError(db.Close(), &#34;error closing connection: %s&#34;)
}
func checkSqlRowsClose(rows *sql.Rows) {
	checkNoError(rows.Close(), &#34;error closing rows: %s&#34;)
}
