package main

import (
	"fmt"
	"database/sql"
	"encoding/json"
)

type Customer struct {
  Id    int64
  Name  sql.NullString
}

func main() {
	cc := &Customer{5,sql.NullString{"",true}}
	buf,err := json.Marshal(cc)
	if err != nil{
		panic(err)
	}
	fmt.Println(string(buf))
	
	cc = &Customer{5,sql.NullString{"pizza",false}}
	buf,err = json.Marshal(cc)
	if err != nil{
		panic(err)
	}
	fmt.Println(string(buf))
}
func (c *Customer) MarshalJSON() ([]byte, error){
    a := struct {Id int64;Name string}{c.Id,c.Name.String}
    return json.Marshal(a)
}