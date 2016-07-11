package main

import (
	"database/sql"
	"fmt"
	//_ "github.com/gwenn/gosqlite"
	_ "code.google.com/p/go-sqlite/go1/sqlite3"
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
	checkNoError(err, "error querying tests: %s")
	defer checkSqlRowsClose(rows)

	tests := make([]Test, 0, 20)
	indexById := make(map[int64]int)
	for rows.Next() {
		var id int64
		var dummy float64
		//err = rows.Scan(&id) // expected 2 destination arguments in Scan, not 1
		err = rows.Scan(&id, &dummy)
		checkNoError(err, "error scanning: %s")
		index, ok := indexById[id]
		if !ok {
			index = len(tests)
			tests = append(tests, Test{})
			indexById[id] = index
		}
		test := &tests[index]

		var data sql.NullFloat64
		if ok {
			//err = rows.Scan(&_ /*, &_... */, &data) // cannot use _ as value
			//err = rows.Scan(nil /*, nil...*/, &data) // Scan error on column index 1: destination not a pointer
			err = rows.Scan(&dummy /*, &dummy*/, &data)
			checkNoError(err, "error scanning: %s")
		} else {
			err = rows.Scan(&test.Id /*, &test...other fields */, &data)
			checkNoError(err, "error scanning: %s")
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
	db, err := sql.Open("sqlite3", ":memory:")
	checkNoError(err, "error opening connection: %s")
	defer checkSqlDbClose(db)

	_, err = db.Exec(schema)
	checkNoError(err, "error creating table: %s")
	tests := findTests(db)
	for i, test := range tests {
		fmt.Printf("> Tests[%d]: %v\n", i, test)
	}
}

func checkNoError(err error, format string) {
	if err != nil {
		panic(fmt.Sprintf(format, err.Error()))
	}
}
func checkSqlDbClose(db *sql.DB) {
	checkNoError(db.Close(), "error closing connection: %s")
}
func checkSqlRowsClose(rows *sql.Rows) {
	checkNoError(rows.Close(), "error closing rows: %s")
}
