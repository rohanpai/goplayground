// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// &#43;build ignore

package main

import (
	&#34;database/sql&#34;
	&#34;fmt&#34;
	//_ &#34;github.com/gwenn/gosqlite&#34;
	//_ &#34;github.com/mattn/go-sqlite3&#34;
	_ &#34;code.google.com/p/go-sqlite/go1/sqlite3&#34;
)

func assert(msg string, actual bool) {
	if !actual {
		panic(msg)
	}
}
func assertEquals(format string, expected, actual interface{}) {
	if expected != actual {
		panic(fmt.Sprintf(format, expected, actual))
	}
}

func checkNoError(err error, format string) {
	if err != nil {
		panic(fmt.Sprintf(format, err))
	}
}

func checkSqlDbClose(db *sql.DB) {
	checkNoError(db.Close(), &#34;error closing connection: %s&#34;)
}

func checkSqlStmtClose(stmt *sql.Stmt) {
	checkNoError(stmt.Close(), &#34;error closing statement: %s&#34;)
}

func checkSqlRowsClose(rows *sql.Rows) {
	checkNoError(rows.Close(), &#34;error closing rows: %s&#34;)
}

func insertRecords(db *sql.DB, name1 string, name2 string) (id1, id2 int64) {
	tx, err := db.Begin()
	checkNoError(err, &#34;error beginning transaction&#34;)
	insert := &#34;INSERT into blah (name) values (?)&#34;
	r, err := tx.Exec(insert, name1)
	checkNoError(err, &#34;error inserting record: %s&#34;)
	id1, err = r.LastInsertId()
	checkNoError(err, &#34;error retrieving record id: %s&#34;)
	r, err = tx.Exec(insert, name2)
	checkNoError(err, &#34;error inserting record: %s&#34;)
	id2, err = r.LastInsertId()
	checkNoError(err, &#34;error retrieving record id: %s&#34;)
	println(&#34;id1:&#34;, id1, &#34;id2: &#34;, id2)
	checkNoError(tx.Commit(), &#34;error commiting transaction&#34;)
	return
}

func checkRecords(db *sql.DB, id1 int64, name1 string, id2 int64, name2 string) (rid1, rid2 int64) {
	rows, err := db.Query(&#34;SELECT id, name from blah order by id&#34;)
	checkNoError(err, &#34;error querying records: %s&#34;)
	defer checkSqlRowsClose(rows)

	var name string
	if rows.Next() {
		err = rows.Scan(&amp;rid1, &amp;name)
		checkNoError(err, &#34;error scanning: %s&#34;)
		assertEquals(&#34;expected: %d, actual: %d&#34;, id1, rid1)
		assertEquals(&#34;expected: %s, actual: %s&#34;, name1, name)
	}
	if rows.Next() {
		err = rows.Scan(&amp;rid2, &amp;name)
		checkNoError(err, &#34;error scanning: %s&#34;)
		assertEquals(&#34;expected: %d, actual: %d&#34;, id2, rid2)
		assertEquals(&#34;expected: %s, actual: %s&#34;, name2, name)
	}
	return
}

func updateRecords(db *sql.DB, id1 int64, name1 string, id2 int64, name2 string) {
	tx, err := db.Begin()
	checkNoError(err, &#34;error beginning transaction&#34;)
	stmt, err := tx.Prepare(&#34;update blah set name = ? where id = ?&#34;)
	checkNoError(err, &#34;error updating records&#34;)
	defer checkSqlStmtClose(stmt)

	r, err := stmt.Exec(name1, id1)
	checkNoError(err, &#34;error updating record&#34;)
	changes, err := r.RowsAffected()
	checkNoError(err, &#34;error checking row affected&#34;)
	assert(&#34;fail to update&#34;, changes == 1)

	r, err = stmt.Exec(name2, id2)
	checkNoError(err, &#34;error updating record&#34;)
	changes, err = r.RowsAffected()
	checkNoError(err, &#34;error checking row affected&#34;)
	assert(&#34;fail to update&#34;, changes == 1)
	checkNoError(tx.Commit(), &#34;error commiting transaction&#34;)
}

func deleteRecords(db *sql.DB, id1 int64, id2 int64) {
	r, err := db.Exec(&#34;delete from blah where id in (?, ?)&#34;, id1, id2)
	checkNoError(err, &#34;error deleting records&#34;)
	changes, err := r.RowsAffected()
	checkNoError(err, &#34;error checking row affected&#34;)
	assert(&#34;fail to delete&#34;, changes == 2)
}

func checkNoRecord(db *sql.DB) {
	rows, err := db.Query(&#34;SELECT id, name from blah order by id&#34;)
	checkNoError(err, &#34;error querying records: %s&#34;)
	defer checkSqlRowsClose(rows)
	assert(&#34;No row expected&#34;, !rows.Next())
}

func main() {
	db, err := sql.Open(&#34;sqlite3&#34;, &#34;:memory:&#34;)
	checkNoError(err, &#34;error opening connection: %s&#34;)
	defer checkSqlDbClose(db)

	// create table blah with autogen primary key
	_, err = db.Exec(&#34;CREATE TABLE blah (id INTEGER PRIMARY KEY NOT NULL, name TEXT);&#34;)
	checkNoError(err, &#34;error creating table: %s&#34;)

	// insert 2 records into blah
	name1, name2 := &#34;Bart&#34;, &#34;Lisa&#34;
	id1, id2 := insertRecords(db, name1, name2)

	// select records from blah
	rid1, rid2 := checkRecords(db, id1, name1, id2, name2)
	// verify they have different ids
	assert(fmt.Sprintf(&#34;%d == %d&#34;, rid1, rid2), rid1 != rid2)

	// update both records
	name1, name2 = &#34;El Barto&#34;, &#34;Maggie&#34;
	updateRecords(db, id1, name1, id2, name2)

	// select records from blah
	// verify they have the same value
	checkRecords(db, id1, name1, id2, name2)

	// delete both records from blah
	deleteRecords(db, id1, id2)

	// select records from blah
	// verify no records were found
	checkNoRecord(db)
}
