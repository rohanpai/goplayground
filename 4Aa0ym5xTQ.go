package main

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"time"
)

type User struct {
	Id        int64
	Birthday  time.Time
	Age       int64
	Name      string `sql:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	Emails            []Email       // One-To-Many relationship (has many)
	BillingAddress    Address       // One-To-One relationship (has one)
	BillingAddressId  sql.NullInt64 // Foreign key of BillingAddress
	ShippingAddress   Address       // One-To-One relationship (has one)
	ShippingAddressId int64         // Foreign key of ShippingAddress
	IgnoreMe          int64         `sql:"-"`                          // Ignore this field
	Languages         []Language    `gorm:"many2many:user_languages;"` // Many-To-Many relationship, 'user_languages' is join table
}

type Email struct {
	Id         int64
	UserId     int64  // Foreign key for User (belongs to)
	Email      string `sql:"type:varchar(100);"` // Set field's type
	Subscribed bool
}

type Address struct {
	Id       int64
	Address1 string         `sql:"not null;unique"` // Set field as not nullable and unique
	Address2 string         `sql:"type:varchar(100)"`
	Post     sql.NullString `sql:not null`
}

type Language struct {
	Id   int64
	Name string
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

	db.Exec("PRAGMA foreign_keys = ON;")

	db.AutoMigrate(User{}, Email{}, Address{}, Language{})

	user := User{
		Name:            "jinzhu",
		BillingAddress:  Address{Address1: "Billing Address - Address 1"},
		ShippingAddress: Address{Address1: "Shipping Address - Address 1"},
		Emails:          []Email{{Email: "jinzhu@example.com"}, {Email: "jinzhu-2@example@example.com"}},
		Languages:       []Language{{Name: "ZH"}, {Name: "EN"}},
	}

	db.Create(&user)

	var u User
	db.Debug().First(&u).Related(&u.Emails)
	fmt.Println(u.Emails)
}
