package main

import (
	&#34;database/sql&#34;
	&#34;fmt&#34;
	_ &#34;github.com/mattn/go-adodb&#34;
	&#34;os&#34;
)

func main() {
	if _, err := os.Stat(&#34;./example.mdb&#34;); err != nil {
		fmt.Println(&#34;put here empty database named &#39;example.mdb&#39;.&#34;)
		return
	}

	db, err := sql.Open(&#34;adodb&#34;, &#34;Provider=Microsoft.Jet.OLEDB.4.0;Data Source=./example.mdb;&#34;)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sqls := []string{
		&#34;DROP TABLE languages&#34;,
		&#34;CREATE TABLE languages (id text not null primary key, name text)&#34;,
	}
	for _, sql := range sqls {
		_, err = db.Exec(sql)
		if err != nil {
			fmt.Printf(&#34;%q: %s\n&#34;, err, sql)
			return
		}
	}

	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}

	stmt, err := tx.Prepare(&#34;insert into languages(id, name) values(?, ?)&#34;)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(&#34;en&#34;, &#34;English&#34;)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = stmt.Exec(&#34;fr&#34;, &#34;French&#34;)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = stmt.Exec(&#34;de&#34;, &#34;German&#34;)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = stmt.Exec(&#34;es&#34;, &#34;Spanish&#34;)
	if err != nil {
		fmt.Println(err)
		return
	}

	tx.Commit()
}