package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
	"os"
	"strings"
)

var sqliteFileName = "test_data.db"

type User struct {
	Id   string `sql:"type:varchar(100);primary_key"`
	Name string
}

func main() {
	fmt.Printf("\nTesting gorm string primary key")

	if _, err := os.Stat(sqliteFileName); err == nil {
		err = os.Remove(sqliteFileName)
		if nil != err {
			fmt.Printf("\nError while deleting file:%v %v", sqliteFileName, err.Error())
		}
	}

	db, err := gorm.Open("sqlite3", sqliteFileName)
	if err != nil {
		fmt.Printf("\nError while opening db:%v", err.Error())
		return
	}

	newDBConn := db.DB()
	err = newDBConn.Ping()
	if err != nil {
		fmt.Printf("\nError while pinging new db:%v", err.Error())
		return
	}

	db.AutoMigrate(&User{})

	newId := strings.ToLower(uuid.NewV4().String())
	user1 := User{
		Id:   newId,
		Name: "John Smith",
	}

	jsnStr, err := ConvObjectToJson(user1)
	if nil != err {
		fmt.Printf("Error while converting to json:%v", err.Error())
		return
	}
	fmt.Printf("\nUser before save:%v", jsnStr)

	res := db.Create(&user1)
	if res.Error != nil {
		fmt.Printf("\nError while saving new user:%v", res.Error.Error())
		return
	}

	jsnStr, err = ConvObjectToJson(user1)
	if nil != err {
		fmt.Printf("Error while converting to json:%v", err.Error())
		return
	}
	fmt.Printf("\nUser after save:%v", jsnStr)

}

func ConvObjectToJson(obj interface{}) (string, error) {
	var buf bytes.Buffer
	encdr := json.NewEncoder(&buf)
	err := encdr.Encode(obj)
	if nil == err {
		return buf.String(), nil
	} else {
		return "", err
	}
}
