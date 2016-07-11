package main

import (
	&#34;bytes&#34;
	&#34;database/sql&#34;
	&#34;encoding/json&#34;
	&#34;fmt&#34;
	&#34;log&#34;
	&#34;strconv&#34;
	&#34;strings&#34;
	&#34;time&#34;

	_ &#34;github.com/go-sql-driver/mysql&#34;
)

const dbformat = &#34;2006-01-02 15:04:05&#34;

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
	Employee_id              string `json:&#34;,omitempty&#34;`
	Employee_name            string `json:&#34;,omitempty&#34;`
	Is_test                  byte
	Menu_items               []Menu_item
	Payable                  float64
	Pos_type                 string `json:&#34;,omitempty&#34;`
	Pos_version              string `json:&#34;,omitempty&#34;`
	Punchh_key               string
	Receipt_datetime         string
	Subtotal_amount          float64
	Transaction_no           string `json:&#34;,omitempty&#34;`
	Business_id, Location_id int
	Created_at               time.Time
	Updated_at               time.Time `json:&#34;,omitempty&#34;`
	Revenue_code             string    `json:&#34;,omitempty&#34;`
	Revenue_id               string    `json:&#34;,omitempty&#34;`
	Status                   string    `json:&#34;,omitempty&#34;`
	Ipv4_addr                string    `json:&#34;,omitempty&#34;`
	Stored_at                int64
}

func (r BigReceipt) AddMenuItems(items []Menu_item) BigReceipt {
	for _, item := range items {
		r.Menu_items = append(r.Menu_items, item)
	}
	return r
}

func (m Menu_item) ValidItem() bool {
	if m.Item_type == &#34;M&#34; || m.Item_type == &#34;D&#34; {
		return true
	} else {
		return false
	}
}

func main() {
	db, err := sql.Open(&#34;mysql&#34;, &#34;user:password@tcp(host)/dbname&#34;)
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
		err = rows.Scan(&amp;mr.Id, &amp;mr.Amount, &amp;mr.Cc_last4, &amp;mr.Employee_id, &amp;mr.Employee_name, &amp;mr.Is_test, &amp;mr.Menu_items,
			&amp;mr.Payable, &amp;mr.Pos_type, &amp;mr.Pos_version, &amp;mr.Punchh_key, &amp;mr.Receipt_datetime, &amp;mr.Subtotal_amount, &amp;mr.Transaction_no,
			&amp;mr.Business_id, &amp;mr.Location_id, &amp;mr.Created_at, &amp;mr.Updated_at, &amp;mr.Revenue_code, &amp;mr.Revenue_id, &amp;mr.Status, &amp;mr.Ipv4_addr)
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
		menuItems := strings.Split(mr.Menu_items.String, &#34;^&#34;)
		items := parseMenuItems(menuItems)
		r = r.AddMenuItems(items)
		b, err := json.Marshal(r)
		if err != nil {
			log.Fatal(err)
		}
		var out bytes.Buffer
		json.Compact(&amp;out, b)
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
		itemParts := strings.Split(v, &#34;|&#34;)
		partsLen := len(itemParts)
		if partsLen &lt; 5 {
			continue
		}
		item.Name = itemParts[0]
		item.Qty, _ = strconv.Atoi(itemParts[1])
		item.Amount, _ = strconv.ParseFloat(itemParts[2], 64)
		item.Item_type = strings.ToUpper(itemParts[3])
		item.Id = itemParts[4]
		if partsLen &gt; 5 {
			item.Family = itemParts[5]
		}
		if partsLen &gt; 6 {
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
