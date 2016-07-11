package main

import (
	"encoding/xml"
	"io"
	"strconv"
	"strings"
	"time"
	"fmt"
)

type Message struct {
	Id        string
	Type      string
	Timestamp time.Time
	Notify    []string
}

func main() {
	x := `<message from="01234567890@s.whatsapp.net"
         id="1339831077-7"
         type="chat"
         timestamp="1339848755">
    <notify xmlns="urn:xmpp:whatsapp"
            name="Koen" />
    <request xmlns="urn:xmpp:receipts" />
    <body>Hello</body>
</message>`

	d := xml.NewDecoder(strings.NewReader(x))
	
	messages := []*Message{}
	var current *Message

	for {
		token, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "message":
				current = &Message{}
				for _, attr := range t.Attr {
					switch attr.Name.Local {
					case "id":
						current.Id = attr.Value
					case "type":
						current.Type = attr.Value
					case "timestamp":
						i, _ := strconv.Atoi(attr.Value)
						current.Timestamp = time.Unix(int64(i), 0)
					}
				}
				messages = append(messages, current)
			case "notify":
				for _, attr := range t.Attr {
					switch attr.Name.Local {
					case "name":
						current.Notify = append(current.Notify, attr.Value)
					}
				}
			}
		case xml.EndElement:
		case xml.CharData:
		case xml.Comment:
		case xml.ProcInst:
		case xml.Directive:
		}
	}
	
	for _, msg := range messages {
		fmt.Println(msg)
	}
}