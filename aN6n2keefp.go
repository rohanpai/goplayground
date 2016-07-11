package main

import (
	"encoding/json"
	"fmt"
)

var str = `
{
    "glossary": {
        "title": "example glossary",
		"GlossDiv": {
            "title": "S",
			"GlossList": {
                "GlossEntry": {
                    "ID": "SGML",
					"SortAs": "SGML",
					"GlossTerm": "Standard Generalized Markup Language",
					"Acronym": "SGML",
					"Abbrev": "ISO 8879:1986",
					"GlossDef": {
                        "para": "A meta-markup language, used to create markup languages such as DocBook.",
						"GlossSeeAlso": ["GML", "XML"]
                    },
					"GlossSee": "markup"
                }
            }
        }
    }
}
`

func main() {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(str), &m)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", m)
}
