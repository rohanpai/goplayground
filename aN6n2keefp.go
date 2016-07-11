package main

import (
	&#34;encoding/json&#34;
	&#34;fmt&#34;
)

var str = `
{
    &#34;glossary&#34;: {
        &#34;title&#34;: &#34;example glossary&#34;,
		&#34;GlossDiv&#34;: {
            &#34;title&#34;: &#34;S&#34;,
			&#34;GlossList&#34;: {
                &#34;GlossEntry&#34;: {
                    &#34;ID&#34;: &#34;SGML&#34;,
					&#34;SortAs&#34;: &#34;SGML&#34;,
					&#34;GlossTerm&#34;: &#34;Standard Generalized Markup Language&#34;,
					&#34;Acronym&#34;: &#34;SGML&#34;,
					&#34;Abbrev&#34;: &#34;ISO 8879:1986&#34;,
					&#34;GlossDef&#34;: {
                        &#34;para&#34;: &#34;A meta-markup language, used to create markup languages such as DocBook.&#34;,
						&#34;GlossSeeAlso&#34;: [&#34;GML&#34;, &#34;XML&#34;]
                    },
					&#34;GlossSee&#34;: &#34;markup&#34;
                }
            }
        }
    }
}
`

func main() {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(str), &amp;m)
	if err != nil {
		panic(err)
	}
	fmt.Printf(&#34;%#v\n&#34;, m)
}
