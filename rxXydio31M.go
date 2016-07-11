package main

import (
	&#34;encoding/xml&#34;
	&#34;fmt&#34;
	&#34;log&#34;
)

type Person struct {
	Xmlns_xsi string `xml:&#34;xmlns xsi,attr,omitempty&#34;`
	Xsi_type  string `xml:&#34;http://www.w3.org/2001/XMLSchema-instance type,attr,omitempty&#34;`
	ID        string `xml:&#34;,omitempty&#34;`
	Name      string `xml:&#34;,omitempty&#34;`
}

type People struct {
	Person []Person `xml:&#34;Person,omitempty&#34;`
}

func main() {
	fmt.Println(&#34;----------- unmarshal -----------&#34;)
	unmarshal()
	fmt.Println(&#34;------------ marshal ------------&#34;)
	marshal()
}

func unmarshal() {
	xmldoc := `
    &lt;People&gt;
      &lt;Person xmlns:xsi=&#34;http://www.w3.org/2001/XMLSchema-instance&#34; xsi:type=&#34;Student&#34;&gt;
        &lt;ID&gt;2222&lt;/ID&gt;
        &lt;Name&gt;Alice&lt;/Name&gt;
      &lt;/Person&gt;
      &lt;Person xmlns:xsi=&#34;http://www.w3.org/2001/XMLSchema-instance&#34; xsi:type=&#34;Student&#34;&gt;
        &lt;ID&gt;3333&lt;/ID&gt;
        &lt;Name&gt;Bob&lt;/Name&gt;
      &lt;/Person&gt;
    &lt;/People&gt;`

	v := People{}
	if err := xml.Unmarshal([]byte(xmldoc), &amp;v); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf(&#34;%&#43;v\n&#34;, v)
	}
}

func marshal() {

	xmlacl := People{
		Person: []Person{
			Person{Xmlns_xsi: &#34;http://www.w3.org/2001/XMLSchema-instance&#34;, Xsi_type: &#34;Student&#34;, ID: &#34;2222&#34;, Name: &#34;Alice&#34;},
			Person{Xmlns_xsi: &#34;http://www.w3.org/2001/XMLSchema-instance&#34;, Xsi_type: &#34;Student&#34;, ID: &#34;3333&#34;, Name: &#34;Bob&#34;},
		},
	}
	if v, err := xml.MarshalIndent(&amp;xmlacl, &#34;&#34;, &#34;  &#34;); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(string(v))
	}
}