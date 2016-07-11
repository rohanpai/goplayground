package main

import (
	&#34;database/sql&#34;
	&#34;errors&#34;
	&#34;fmt&#34;
	_ &#34;github.com/mattn/go-sqlite3&#34; // init() only
	&#34;log&#34;
	&#34;os&#34;
)

func main() {
	const database = &#34;/tmp/db_test.db&#34;

	os.Remove(database)

	db, err := sql.Open(&#34;sqlite3&#34;, database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sql := `
	create table user(
		id integer not null primary key,
		name text,
		dummy text
	)
	`
	_, err = db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	sql = `insert into user(name) values(?)`
	for _, name := range []string{&#34;taro&#34;, &#34;jiro&#34;, &#34;saburo&#34;} {
		_, err = db.Exec(sql, name)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()

	sql = `select * from user`
	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name, dummy string
		rows.Scan(&amp;id, &amp;name, &amp;dummy)
		fmt.Printf(&#34;id: %d, name: %s, dummy: %s\n&#34;, id, name, dummy)
	}
	fmt.Println()

	sql = `select * from user`
	users, err := query(db, sql)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(&#34;%#&#43;v\n&#34;, users)
	fmt.Println()

	var id = 2
	sql = `select * from user where id = ?`
	users, err = query(db, sql, id)
	if err != nil {
		log.Fatal(err)
	}
	if len(users) &gt; 0 {
		fmt.Printf(`user(id:%d) is &#34;%s&#34;`, id, users[0][&#34;name&#34;].([]byte))
	}
}

func query(db *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
	if db == nil {
		return nil, errors.New(&#34;db is nil&#34;)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	res := make([]map[string]interface{}, 0)

	for rows.Next() {
		container := make([]interface{}, len(cols))
		dest := make([]interface{}, len(cols))
		for i, _ := range container {
			dest[i] = &amp;container[i]
		}
		rows.Scan(dest...)
		r := make(map[string]interface{})
		for i, colname := range cols {
			val := dest[i].(*interface{})
			r[colname] = *val
		}
		res = append(res, r)
	}

	return res, nil
}
