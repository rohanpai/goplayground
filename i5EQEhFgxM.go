// 123 project main.go
package main

import (
	&#34;database/sql&#34;
	&#34;fmt&#34;
	&#34;html/template&#34;
	&#34;io/ioutil&#34;
	&#34;log&#34;
	&#34;net/http&#34;
	&#34;strings&#34;

	&#34;code.google.com/p/go-charset/charset&#34;
	_ &#34;code.google.com/p/go-charset/data&#34;
	_ &#34;code.google.com/p/odbc&#34;
)

var (
	fam, name, ot string
	otdel         int
	id            int
	a             []string
	x             string
	//DepartamentNAMEUTF string
)

func otdelen(w http.ResponseWriter, r *http.Request) {
	fmt.Println(&#34;method:&#34;, r.Method) //get request method
	if r.Method == &#34;GET&#34; {
		t, _ := template.ParseFiles(&#34;otdel.gtpl&#34;)
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		fmt.Println(&#34;otdel:&#34;, r.Form[&#34;otdel&#34;])
		fmt.Println(r.Form) // print form information in server side
		fmt.Println(&#34;path&#34;, r.URL.Path)
		fmt.Println(&#34;scheme&#34;, r.URL.Scheme)
		a = r.Form[&#34;otdel&#34;]
		x = a[0]
		fmt.Fprintln(w, a)
		if len(r.Form[&#34;otdel&#34;][0]) == 0 {
			// code for empty field
			fmt.Fprintln(w, &#34;Поле пустое!!&#34;)

		}
	}
}

func pacient(w http.ResponseWriter, r *http.Request) {
	//fmt.Print(&#34;Введите № отделения: &#34;)
	//fmt.Scan(&amp;otdel)
	// Подключаемся к базе MSSQL через ODBC драйвер
	db, err := sql.Open(&#34;odbc&#34;, &#34;DSN=DBS0&#34;)
	if err != nil {
		fmt.Println(&#34;Error in connect DB&#34;)
		log.Fatal(err)
	}

	rows, err := db.Query(&#34;select hDED.FAMILY, hDED.Name, hDED.OT, br.rf_departmentid from dbo.stt_medicalhistory hDED inner join  dbo.v_curentmigrationpatient AS m on hDED.medicalhistoryid = m.rf_medicalhistoryid  INNER JOIN dbo.stt_stationarbranch AS br ON br.stationarbranchid = m.rf_stationarbranchid and br.rf_departmentid=?&#34;, &amp;x)

	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&amp;fam, &amp;name, &amp;ot, &amp;id); err != nil {
			log.Fatal(err)
		}
		/*---------------ФАМИЛИЯ------------------*/
		// Перекодируем данные с запроса из cp1251 в UTF8
		r, err := charset.NewReader(&#34;windows-1251&#34;, strings.NewReader(fam))
		if err != nil {
			log.Fatal(err)
		}
		fam, err := ioutil.ReadAll(r)

		/*---------------ИМЯ------------------*/
		r2, err := charset.NewReader(&#34;windows-1251&#34;, strings.NewReader(name))
		if err != nil {
			log.Fatal(err)
		}
		name, err := ioutil.ReadAll(r2)

		/*------------------ОТЧЕСТВО------------*/
		r3, err := charset.NewReader(&#34;windows-1251&#34;, strings.NewReader(ot))
		if err != nil {
			log.Fatal(err)
		}
		ot, err := ioutil.ReadAll(r3)

		if err != nil {
			log.Fatal(err)
		}
		//		DepartamentNAMEUTF = string(name)
		//	fmt.Printf(&#34;%s &#34;, fam)
		//	fmt.Printf(&#34;%s &#34;, name)
		//	fmt.Printf(&#34;%s &#34;, ot)
		//	fmt.Printf(&#34;%d\n&#34;, id)
		fmt.Fprintf(w, &#34;%s %s %s %d\n&#34;, fam, name, ot, id)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc(&#34;/otdelenie&#34;, otdelen)
	http.HandleFunc(&#34;/pacient&#34;, pacient)
	http.ListenAndServe(&#34;:8080&#34;, nil)
}
