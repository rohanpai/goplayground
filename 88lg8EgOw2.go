package main

import (
	&#34;fmt&#34;
	&#34;database/sql&#34;
	&#34;encoding/json&#34;
)

type Customer struct {
  Id    int64
  Name  sql.NullString
}

func main() {
	cc := &amp;Customer{5,sql.NullString{&#34;&#34;,true}}
	buf,err := json.Marshal(cc)
	if err != nil{
		panic(err)
	}
	fmt.Println(string(buf))
	
	cc = &amp;Customer{5,sql.NullString{&#34;pizza&#34;,false}}
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