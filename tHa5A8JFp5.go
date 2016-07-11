package main

import (
	&#34;encoding/xml&#34;
	&#34;io&#34;
	&#34;strconv&#34;
	&#34;strings&#34;
	&#34;time&#34;
	&#34;fmt&#34;
)

type Message struct {
	Id        string
	Type      string
	Timestamp time.Time
	Notify    []string
}

func main() {
	x := `&lt;message from=&#34;01234567890@s.whatsapp.net&#34;
         id=&#34;1339831077-7&#34;
         type=&#34;chat&#34;
         timestamp=&#34;1339848755&#34;&gt;
    &lt;notify xmlns=&#34;urn:xmpp:whatsapp&#34;
            name=&#34;Koen&#34; /&gt;
    &lt;request xmlns=&#34;urn:xmpp:receipts&#34; /&gt;
    &lt;body&gt;Hello&lt;/body&gt;
&lt;/message&gt;`

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
			case &#34;message&#34;:
				current = &amp;Message{}
				for _, attr := range t.Attr {
					switch attr.Name.Local {
					case &#34;id&#34;:
						current.Id = attr.Value
					case &#34;type&#34;:
						current.Type = attr.Value
					case &#34;timestamp&#34;:
						i, _ := strconv.Atoi(attr.Value)
						current.Timestamp = time.Unix(int64(i), 0)
					}
				}
				messages = append(messages, current)
			case &#34;notify&#34;:
				for _, attr := range t.Attr {
					switch attr.Name.Local {
					case &#34;name&#34;:
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