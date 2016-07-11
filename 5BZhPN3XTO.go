// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"database/sql"
	"fmt"
	//_ "github.com/gwenn/gosqlite"
	//_ "github.com/mattn/go-sqlite3"
	_ "code.google.com/p/go-sqlite/go1/sqlite3"
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
	checkNoError(db.Close(), "error closing connection: %s")
}

func checkSqlStmtClose(stmt *sql.Stmt) {
	checkNoError(stmt.Close(), "error closing statement: %s")
}

func checkSqlRowsClose(rows *sql.Rows) {
	checkNoError(rows.Close(), "error closing rows: %s")
}

func insertRecords(db *sql.DB, name1 string, name2 string) (id1, id2 int64) {
	tx, err := db.Begin()
	checkNoError(err, "error beginning transaction")
	insert := "INSERT into blah (name) values (?)"
	r, err := tx.Exec(insert, name1)
	checkNoError(err, "error inserting record: %s")
	id1, err = r.LastInsertId()
	checkNoError(err, "error retrieving record id: %s")
	r, err = tx.Exec(insert, name2)
	checkNoError(err, "error inserting record: %s")
	id2, err = r.LastInsertId()
	checkNoError(err, "error retrieving record id: %s")
	println("id1:", id1, "id2: ", id2)
	checkNoError(tx.Commit(), "error commiting transaction")
	return
}

func checkRecords(db *sql.DB, id1 int64, name1 string, id2 int64, name2 string) (rid1, rid2 int64) {
	rows, err := db.Query("SELECT id, name from blah order by id")
	checkNoError(err, "error querying records: %s")
	defer checkSqlRowsClose(rows)

	var name string
	if rows.Next() {
		err = rows.Scan(&rid1, &name)
		checkNoError(err, "error scanning: %s")
		assertEquals("expected: %d, actual: %d", id1, rid1)
		assertEquals("expected: %s, actual: %s", name1, name)
	}
	if rows.Next() {
		err = rows.Scan(&rid2, &name)
		checkNoError(err, "error scanning: %s")
		assertEquals("expected: %d, actual: %d", id2, rid2)
		assertEquals("expected: %s, actual: %s", name2, name)
	}
	return
}

func updateRecords(db *sql.DB, id1 int64, name1 string, id2 int64, name2 string) {
	tx, err := db.Begin()
	checkNoError(err, "error beginning transaction")
	stmt, err := tx.Prepare("update blah set name = ? where id = ?")
	checkNoError(err, "error updating records")
	defer checkSqlStmtClose(stmt)

	r, err := stmt.Exec(name1, id1)
	checkNoError(err, "error updating record")
	changes, err := r.RowsAffected()
	checkNoError(err, "error checking row affected")
	assert("fail to update", changes == 1)

	r, err = stmt.Exec(name2, id2)
	checkNoError(err, "error updating record")
	changes, err = r.RowsAffected()
	checkNoError(err, "error checking row affected")
	assert("fail to update", changes == 1)
	checkNoError(tx.Commit(), "error commiting transaction")
}

func deleteRecords(db *sql.DB, id1 int64, id2 int64) {
	r, err := db.Exec("delete from blah where id in (?, ?)", id1, id2)
	checkNoError(err, "error deleting records")
	changes, err := r.RowsAffected()
	checkNoError(err, "error checking row affected")
	assert("fail to delete", changes == 2)
}

func checkNoRecord(db *sql.DB) {
	rows, err := db.Query("SELECT id, name from blah order by id")
	checkNoError(err, "error querying records: %s")
	defer checkSqlRowsClose(rows)
	assert("No row expected", !rows.Next())
}

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	checkNoError(err, "error opening connection: %s")
	defer checkSqlDbClose(db)

	// create table blah with autogen primary key
	_, err = db.Exec("CREATE TABLE blah (id INTEGER PRIMARY KEY NOT NULL, name TEXT);")
	checkNoError(err, "error creating table: %s")

	// insert 2 records into blah
	name1, name2 := "Bart", "Lisa"
	id1, id2 := insertRecords(db, name1, name2)

	// select records from blah
	rid1, rid2 := checkRecords(db, id1, name1, id2, name2)
	// verify they have different ids
	assert(fmt.Sprintf("%d == %d", rid1, rid2), rid1 != rid2)

	// update both records
	name1, name2 = "El Barto", "Maggie"
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
