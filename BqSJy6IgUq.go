package main

import (
	&#34;bytes&#34;
	&#34;encoding/xml&#34;
	&#34;os&#34;
)

type Rdf_RDF struct {
	Rdf_Description *Rdf_Description `xml:&#34;http://www.w3.org/1999/02/22-rdf-syntax-ns# Description,omitempty&#34;`
	Attr_cd         string           `xml:&#34;xmlns cd,attr&#34;`
	Attr_rdf        string           `xml:&#34;xmlns rdf,attr&#34;`
	XMLName         xml.Name         `xml:&#34;http://www.w3.org/1999/02/22-rdf-syntax-ns# RDF,omitempty&#34;`
}

type Rdf_Description struct {
	Cd_artist      string   `xml:&#34;http://www.recshop.fake/cd# artist,omitempty&#34;`
	Cd_company     string   `xml:&#34;http://www.recshop.fake/cd# company,omitempty&#34;`
	Cd_country     string   `xml:&#34;http://www.recshop.fake/cd# country,omitempty&#34;`
	Cd_price       string   `xml:&#34;http://www.recshop.fake/cd# price,omitempty&#34;`
	Cd_year        string   `xml:&#34;http://www.recshop.fake/cd# year,omitempty&#34;`
	Attr_rdf_about string   `xml:&#34;http://www.w3.org/1999/02/22-rdf-syntax-ns# about,attr&#34;`
	XMLName        xml.Name `xml:&#34;http://www.w3.org/1999/02/22-rdf-syntax-ns# Description,omitempty&#34;`
}

// Derived from http://www.w3schools.com/webservices/ws_rdf_example.asp
var xmlString string = `
&lt;?xml version=&#34;1.0&#34;?&gt;
&lt;rdf:RDF
 xmlns:rdf=&#34;http://www.w3.org/1999/02/22-rdf-syntax-ns#&#34;
 xmlns:cd=&#34;http://www.recshop.fake/cd#&#34;&gt;

&lt;rdf:Description
 rdf:about=&#34;http://www.recshop.fake/cd/Empire Burlesque&#34;&gt;
  &lt;cd:artist&gt;Bob Dylan&lt;/cd:artist&gt;
  &lt;cd:country&gt;USA&lt;/cd:country&gt;
  &lt;cd:company&gt;Columbia&lt;/cd:company&gt;
  &lt;cd:price&gt;10.90&lt;/cd:price&gt;
  &lt;cd:year&gt;1985&lt;/cd:year&gt;
&lt;/rdf:Description&gt;
&lt;/rdf:RDF&gt; `

func main() {
	{
		reader := bytes.NewBufferString(xmlString)
		enc := xml.NewEncoder(os.Stdout)
		enc.Indent(&#34;  &#34;, &#34;    &#34;)
		decoder := xml.NewDecoder(reader)
		for {
			token, _ := decoder.Token()
			if token == nil {
				break
			}
			switch se := token.(type) {
			case xml.StartElement:
				if se.Name.Local == &#34;RDF&#34; &amp;&amp; se.Name.Space == &#34;http://www.w3.org/1999/02/22-rdf-syntax-ns#&#34; {
					var item Rdf_RDF
					decoder.DecodeElement(&amp;item, &amp;se)
					enc.Encode(item)
				}
			}
		}
	}
}
