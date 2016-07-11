package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pborman/uuid"
)

type User struct {
	Id        string `gorm:"primary_key;uuid"`
	Birthday  time.Time
	Age       int64
	Name      string `sql:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	Emails            []Email    // One-To-Many relationship (has many)
	BillingAddress    Address    // One-To-One relationship (has one)
	BillingAddressId  string     // Foreign key of BillingAddress
	ShippingAddress   Address    // One-To-One relationship (has one)
	ShippingAddressId string     // Foreign key of ShippingAddress
	IgnoreMe          int64      `sql:"-"`                          // Ignore this field
	Languages         []Language `gorm:"many2many:user_languages;"` // Many-To-Many relationship, 'user_languages' is join table
}

type Email struct {
	Id         string `gorm:"primary_key;uuid"`
	UserId     string // Foreign key for User (belongs to)
	Email      string `sql:"type:varchar(100);"` // Set field's type
	Subscribed bool
}

type Address struct {
	Id       string         `gorm:"primary_key;uuid"`
	Address1 string         `sql:"not null;unique"` // Set field as not nullable and unique
	Address2 string         `sql:"type:varchar(100)"`
	Post     sql.NullString `sql:not null`
}

type Language struct {
	Id   string `gorm:"primary_key;uuid"`
	Name string
}

// create UUID
func beforeCreate(scope *gorm.Scope) {
	reflectValue := reflect.Indirect(reflect.ValueOf(scope.Value))
	if strings.Contains(string(reflectValue.Type().Field(0).Tag), "uuid") {
		uuid.SetClockSequence(-1)
		scope.SetColumn("id", uuid.NewUUID().String())
	}
}

func main() {
	err := os.Remove("/tmp/sqlite.db")
	if err != nil {
		fmt.Println(err)
	}
	db, err := gorm.Open("sqlite3", "/tmp/sqlite.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	db.Callback().Create().Before("gorm:create").Register("my_plugin:before_create", beforeCreate)

	db.Exec("PRAGMA foreign_keys = ON;")

	db.AutoMigrate(&User{}, &Email{}, &Address{}, &Language{})

	user := User{
		Name:            "jinzhu",
		BillingAddress:  Address{Address1: "Billing Address - Address 1"},
		ShippingAddress: Address{Address1: "Shipping Address - Address 1"},
		Emails:          []Email{{Email: "jinzhu@example.com"}, {Email: "jinzhu-2@example@example.com"}},
		Languages:       []Language{{Name: "ZH"}, {Name: "EN"}},
	}

	db.Create(&user)

	var u User
	db.First(&u).Related(&u.Emails)
	js, err := json.MarshalIndent(u.Emails, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(js))

}
