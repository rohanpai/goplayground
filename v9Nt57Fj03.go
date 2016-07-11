// /tmp/5892Fk2.go
package main

// Error handling left out for brevity
//
// This only demonstrates some of the possibilities that the new XML
// interfaces create.

import (
	&#34;encoding/xml&#34;
	&#34;fmt&#34;
)

type PlainVector []int
type CopyVector []int
type ManualVector []int

func (v CopyVector) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	// Encapsulating the &#34;copy to different structure&#34; pattern
	vX := struct{ X, Y, Z int }{v[0], v[1], v[2]}
	e.EncodeElement(vX, start)

	return nil
}

func (v *CopyVector) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	vX := struct{ X, Y, Z int }{}
	d.DecodeElement(&amp;vX, &amp;start)

	*v = CopyVector{vX.X, vX.Y, vX.Z}

	return nil
}

func (v ManualVector) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	// Encoding tokens manually to inject attributes and comments
	start.Attr = []xml.Attr{xml.Attr{Name: xml.Name{Local: &#34;method&#34;}, Value: &#34;unicorns&#34;}}
	e.EncodeToken(start)
	e.EncodeToken(xml.Comment(&#34;We can emit comments, too.&#34;))

	e.EncodeToken(xml.Comment(&#34;Value of X&#34;))
	e.EncodeElement(v[0], xml.StartElement{Name: xml.Name{Local: &#34;X&#34;}})

	e.EncodeToken(xml.Comment(&#34;Value of Y&#34;))
	e.EncodeElement(v[1], xml.StartElement{Name: xml.Name{Local: &#34;Y&#34;}})

	e.EncodeToken(xml.Comment(&#34;Value of Z&#34;))
	e.EncodeElement(v[2], xml.StartElement{Name: xml.Name{Local: &#34;Z&#34;}})

	e.EncodeToken(xml.EndElement{Name: start.Name})

	return nil
}

func (v *ManualVector) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// We don&#39;t really care about comments or unicorns
	vX := struct{ X, Y, Z int }{}
	d.DecodeElement(&amp;vX, &amp;start)

	*v = ManualVector{vX.X, vX.Y, vX.Z}
	return nil
}

func main() {
	var encoded [][]byte
	vectors := []interface{}{
		PlainVector{1, 2, 3},
		CopyVector{4, 5, 6},
		ManualVector{7, 8, 9},
	}

	for _, vector := range vectors {
		b, err := xml.Marshal(vector)
		encoded = append(encoded, b)

		fmt.Println(string(b), err)
	}

	fmt.Println()

	vP := PlainVector{}
	err := xml.Unmarshal(encoded[0], &amp;vP)
	fmt.Println(vP, err)

	vC := CopyVector{}
	err = xml.Unmarshal(encoded[1], &amp;vC)
	fmt.Println(vC, err)

	vM := ManualVector{}
	err = xml.Unmarshal(encoded[2], &amp;vM)
	fmt.Println(vM, err)
}
