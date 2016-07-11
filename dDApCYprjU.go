package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const dbformat = "2006-01-02 15:04:05"

type MysqlReceipt struct {
	Id               int
	Amount           sql.NullFloat64
	Cc_last4         sql.NullString
	Employee_id      sql.NullString
	Employee_name    sql.NullString
	Is_test          byte
	Menu_items       sql.NullString
	Payable          sql.NullFloat64
	Pos_type         sql.NullString
	Pos_version      sql.NullString
	Punchh_key       string
	Receipt_datetime sql.NullString
	Subtotal_amount  sql.NullFloat64
	Transaction_no   sql.NullString
	Business_id      int
	Location_id      int
	Created_at       string
	Updated_at       sql.NullString
	Revenue_code     sql.NullString
	Revenue_id       sql.NullString
	Status           sql.NullString
	Ipv4_addr        sql.NullString
}

type Menu_item struct {
	Id, Name, Family, Major_group, Item_type string
	Qty                                      int
	Amount                                   float64
}

type BigReceipt struct {
	Id                       int
	Amount                   float64
	Cc_last4                 string
	Employee_id              string `json:",omitempty"`
	Employee_name            string `json:",omitempty"`
	Is_test                  byte
	Menu_items               []Menu_item
	Payable                  float64
	Pos_type                 string `json:",omitempty"`
	Pos_version              string `json:",omitempty"`
	Punchh_key               string
	Receipt_datetime         string
	Subtotal_amount          float64
	Transaction_no           string `json:",omitempty"`
	Business_id, Location_id int
	Created_at               time.Time
	Updated_at               time.Time `json:",omitempty"`
	Revenue_code             string    `json:",omitempty"`
	Revenue_id               string    `json:",omitempty"`
	Status                   string    `json:",omitempty"`
	Ipv4_addr                string    `json:",omitempty"`
	Stored_at                int64
}

func (r BigReceipt) AddMenuItems(items []Menu_item) BigReceipt {
	for _, item := range items {
		r.Menu_items = append(r.Menu_items, item)
	}
	return r
}

func (m Menu_item) ValidItem() bool {
	if m.Item_type == "M" || m.Item_type == "D" {
		return true
	} else {
		return false
	}
}

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(host)/dbname")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query(`select id,amount,cc_last4,employee_id,employee_name,is_test,menu_items,payable,pos_type,
    pos_version,punchh_key,receipt_datetime,subtotal_amount,transaction_no,business_id,location_id,created_at,
    updated_at,revenue_code,revenue_id,status,ipv4_addr from receipts`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var mr MysqlReceipt
		err = rows.Scan(&mr.Id, &mr.Amount, &mr.Cc_last4, &mr.Employee_id, &mr.Employee_name, &mr.Is_test, &mr.Menu_items,
			&mr.Payable, &mr.Pos_type, &mr.Pos_version, &mr.Punchh_key, &mr.Receipt_datetime, &mr.Subtotal_amount, &mr.Transaction_no,
			&mr.Business_id, &mr.Location_id, &mr.Created_at, &mr.Updated_at, &mr.Revenue_code, &mr.Revenue_id, &mr.Status, &mr.Ipv4_addr)
		if err != nil {
			log.Fatal(err)
		}
		if !mr.Menu_items.Valid {
			continue
		}
		r := BigReceipt{Id: mr.Id,
			Amount:           mr.Amount.Float64,
			Cc_last4:         mr.Cc_last4.String,
			Employee_id:      mr.Employee_id.String,
			Employee_name:    mr.Employee_name.String,
			Is_test:          mr.Is_test,
			Menu_items:       []Menu_item{},
			Payable:          mr.Payable.Float64,
			Pos_type:         mr.Pos_type.String,
			Pos_version:      mr.Pos_version.String,
			Punchh_key:       mr.Punchh_key,
			Receipt_datetime: mr.Receipt_datetime.String,
			Subtotal_amount:  mr.Subtotal_amount.Float64,
			Transaction_no:   mr.Transaction_no.String,
			Business_id:      mr.Business_id,
			Location_id:      mr.Location_id,
			Revenue_code:     mr.Revenue_code.String,
			Revenue_id:       mr.Revenue_id.String,
			Status:           mr.Status.String,
			Ipv4_addr:        mr.Ipv4_addr.String,
			Stored_at:        time.Now().Unix(),
		}
		r.Created_at = datetimeParse(mr.Created_at)
		if mr.Updated_at.Valid {
			r.Updated_at = datetimeParse(mr.Updated_at.String)
		}
		menuItems := strings.Split(mr.Menu_items.String, "^")
		items := parseMenuItems(menuItems)
		r = r.AddMenuItems(items)
		b, err := json.Marshal(r)
		if err != nil {
			log.Fatal(err)
		}
		var out bytes.Buffer
		json.Compact(&out, b)
		fmt.Println(string(b))
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func datetimeParse(dateStr string) time.Time {
	datetime, err := time.Parse(dbformat, dateStr)
	if err != nil {
		log.Fatal(err)
	}
	return datetime
}

func parseMenuItems(menuItems []string) []Menu_item {
	var items []Menu_item
	var item Menu_item
	for _, v := range menuItems {
		itemParts := strings.Split(v, "|")
		partsLen := len(itemParts)
		if partsLen < 5 {
			continue
		}
		item.Name = itemParts[0]
		item.Qty, _ = strconv.Atoi(itemParts[1])
		item.Amount, _ = strconv.ParseFloat(itemParts[2], 64)
		item.Item_type = strings.ToUpper(itemParts[3])
		item.Id = itemParts[4]
		if partsLen > 5 {
			item.Family = itemParts[5]
		}
		if partsLen > 6 {
			item.Major_group = itemParts[6]
		}
		if item.ValidItem() {
			items = append(items, item)
		} else {
			continue
		}
	}
	return items
}
