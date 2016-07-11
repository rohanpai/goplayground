package main

import (
	&#34;database/sql&#34;
	&#34;fmt&#34;
	&#34;github.com/jinzhu/gorm&#34;
	_ &#34;github.com/mattn/go-sqlite3&#34;
	&#34;os&#34;
	&#34;time&#34;
)

type User struct {
	Id        int64
	Birthday  time.Time
	Age       int64
	Name      string `sql:&#34;size:255&#34;`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	Emails            []Email       // One-To-Many relationship (has many)
	BillingAddress    Address       // One-To-One relationship (has one)
	BillingAddressId  sql.NullInt64 // Foreign key of BillingAddress
	ShippingAddress   Address       // One-To-One relationship (has one)
	ShippingAddressId int64         // Foreign key of ShippingAddress
	IgnoreMe          int64         `sql:&#34;-&#34;`                          // Ignore this field
	Languages         []Language    `gorm:&#34;many2many:user_languages;&#34;` // Many-To-Many relationship, &#39;user_languages&#39; is join table
}

type Email struct {
	Id         int64
	UserId     int64  // Foreign key for User (belongs to)
	Email      string `sql:&#34;type:varchar(100);&#34;` // Set field&#39;s type
	Subscribed bool
}

type Address struct {
	Id       int64
	Address1 string         `sql:&#34;not null;unique&#34;` // Set field as not nullable and unique
	Address2 string         `sql:&#34;type:varchar(100)&#34;`
	Post     sql.NullString `sql:not null`
}

type Language struct {
	Id   int64
	Name string
}

func main() {
	err := os.Remove(&#34;/tmp/sqlite.db&#34;)
	if err != nil {
		fmt.Println(err)
	}
	db, err := gorm.Open(&#34;sqlite3&#34;, &#34;/tmp/sqlite.db&#34;)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	db.Exec(&#34;PRAGMA foreign_keys = ON;&#34;)

	db.AutoMigrate(User{}, Email{}, Address{}, Language{})

	user := User{
		Name:            &#34;jinzhu&#34;,
		BillingAddress:  Address{Address1: &#34;Billing Address - Address 1&#34;},
		ShippingAddress: Address{Address1: &#34;Shipping Address - Address 1&#34;},
		Emails:          []Email{{Email: &#34;jinzhu@example.com&#34;}, {Email: &#34;jinzhu-2@example@example.com&#34;}},
		Languages:       []Language{{Name: &#34;ZH&#34;}, {Name: &#34;EN&#34;}},
	}

	db.Create(&amp;user)

	var u User
	db.Debug().First(&amp;u).Related(&amp;u.Emails)
	fmt.Println(u.Emails)
}
