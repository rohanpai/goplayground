package main

import (
	"encoding/xml"
	"fmt"
	"log"
)

type Person struct {
	Xmlns_xsi string `xml:"xmlns xsi,attr,omitempty"`
	Xsi_type  string `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr,omitempty"`
	ID        string `xml:",omitempty"`
	Name      string `xml:",omitempty"`
}

type People struct {
	Person []Person `xml:"Person,omitempty"`
}

func main() {
	fmt.Println("----------- unmarshal -----------")
	unmarshal()
	fmt.Println("------------ marshal ------------")
	marshal()
}

func unmarshal() {
	xmldoc := `
    <People>
      <Person xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="Student">
        <ID>2222</ID>
        <Name>Alice</Name>
      </Person>
      <Person xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="Student">
        <ID>3333</ID>
        <Name>Bob</Name>
      </Person>
    </People>`

	v := People{}
	if err := xml.Unmarshal([]byte(xmldoc), &v); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%+v\n", v)
	}
}

func marshal() {

	xmlacl := People{
		Person: []Person{
			Person{Xmlns_xsi: "http://www.w3.org/2001/XMLSchema-instance", Xsi_type: "Student", ID: "2222", Name: "Alice"},
			Person{Xmlns_xsi: "http://www.w3.org/2001/XMLSchema-instance", Xsi_type: "Student", ID: "3333", Name: "Bob"},
		},
	}
	if v, err := xml.MarshalIndent(&xmlacl, "", "  "); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(string(v))
	}
}