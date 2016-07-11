package main

import (
	&#34;bytes&#34;
	&#34;encoding/json&#34;
	&#34;fmt&#34;
	&#34;github.com/jinzhu/gorm&#34;
	_ &#34;github.com/mattn/go-sqlite3&#34;
	&#34;github.com/satori/go.uuid&#34;
	&#34;os&#34;
	&#34;strings&#34;
)

var sqliteFileName = &#34;test_data.db&#34;

type User struct {
	Id   string `sql:&#34;type:varchar(100);primary_key&#34;`
	Name string
}

func main() {
	fmt.Printf(&#34;\nTesting gorm string primary key&#34;)

	if _, err := os.Stat(sqliteFileName); err == nil {
		err = os.Remove(sqliteFileName)
		if nil != err {
			fmt.Printf(&#34;\nError while deleting file:%v %v&#34;, sqliteFileName, err.Error())
		}
	}

	db, err := gorm.Open(&#34;sqlite3&#34;, sqliteFileName)
	if err != nil {
		fmt.Printf(&#34;\nError while opening db:%v&#34;, err.Error())
		return
	}

	newDBConn := db.DB()
	err = newDBConn.Ping()
	if err != nil {
		fmt.Printf(&#34;\nError while pinging new db:%v&#34;, err.Error())
		return
	}

	db.AutoMigrate(&amp;User{})

	newId := strings.ToLower(uuid.NewV4().String())
	user1 := User{
		Id:   newId,
		Name: &#34;John Smith&#34;,
	}

	jsnStr, err := ConvObjectToJson(user1)
	if nil != err {
		fmt.Printf(&#34;Error while converting to json:%v&#34;, err.Error())
		return
	}
	fmt.Printf(&#34;\nUser before save:%v&#34;, jsnStr)

	res := db.Create(&amp;user1)
	if res.Error != nil {
		fmt.Printf(&#34;\nError while saving new user:%v&#34;, res.Error.Error())
		return
	}

	jsnStr, err = ConvObjectToJson(user1)
	if nil != err {
		fmt.Printf(&#34;Error while converting to json:%v&#34;, err.Error())
		return
	}
	fmt.Printf(&#34;\nUser after save:%v&#34;, jsnStr)

}

func ConvObjectToJson(obj interface{}) (string, error) {
	var buf bytes.Buffer
	encdr := json.NewEncoder(&amp;buf)
	err := encdr.Encode(obj)
	if nil == err {
		return buf.String(), nil
	} else {
		return &#34;&#34;, err
	}
}
