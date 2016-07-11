package main

import (
	"bytes"
	"encoding/xml"
	"os"
)

type Rdf_RDF struct {
	Rdf_Description *Rdf_Description `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# Description,omitempty"`
	Attr_cd         string           `xml:"xmlns cd,attr"`
	Attr_rdf        string           `xml:"xmlns rdf,attr"`
	XMLName         xml.Name         `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# RDF,omitempty"`
}

type Rdf_Description struct {
	Cd_artist      string   `xml:"http://www.recshop.fake/cd# artist,omitempty"`
	Cd_company     string   `xml:"http://www.recshop.fake/cd# company,omitempty"`
	Cd_country     string   `xml:"http://www.recshop.fake/cd# country,omitempty"`
	Cd_price       string   `xml:"http://www.recshop.fake/cd# price,omitempty"`
	Cd_year        string   `xml:"http://www.recshop.fake/cd# year,omitempty"`
	Attr_rdf_about string   `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# about,attr"`
	XMLName        xml.Name `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# Description,omitempty"`
}

// Derived from http://www.w3schools.com/webservices/ws_rdf_example.asp
var xmlString string = `
<?xml version="1.0"?>
<rdf:RDF
 xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
 xmlns:cd="http://www.recshop.fake/cd#">

<rdf:Description
 rdf:about="http://www.recshop.fake/cd/Empire Burlesque">
  <cd:artist>Bob Dylan</cd:artist>
  <cd:country>USA</cd:country>
  <cd:company>Columbia</cd:company>
  <cd:price>10.90</cd:price>
  <cd:year>1985</cd:year>
</rdf:Description>
</rdf:RDF> `

func main() {
	{
		reader := bytes.NewBufferString(xmlString)
		enc := xml.NewEncoder(os.Stdout)
		enc.Indent("  ", "    ")
		decoder := xml.NewDecoder(reader)
		for {
			token, _ := decoder.Token()
			if token == nil {
				break
			}
			switch se := token.(type) {
			case xml.StartElement:
				if se.Name.Local == "RDF" && se.Name.Space == "http://www.w3.org/1999/02/22-rdf-syntax-ns#" {
					var item Rdf_RDF
					decoder.DecodeElement(&item, &se)
					enc.Encode(item)
				}
			}
		}
	}
}
