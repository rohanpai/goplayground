// /tmp/5892Fk2.go
package main

// Error handling left out for brevity
//
// This only demonstrates some of the possibilities that the new XML
// interfaces create.

import (
	"encoding/xml"
	"fmt"
)

type PlainVector []int
type CopyVector []int
type ManualVector []int

func (v CopyVector) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	// Encapsulating the "copy to different structure" pattern
	vX := struct{ X, Y, Z int }{v[0], v[1], v[2]}
	e.EncodeElement(vX, start)

	return nil
}

func (v *CopyVector) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	vX := struct{ X, Y, Z int }{}
	d.DecodeElement(&vX, &start)

	*v = CopyVector{vX.X, vX.Y, vX.Z}

	return nil
}

func (v ManualVector) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	// Encoding tokens manually to inject attributes and comments
	start.Attr = []xml.Attr{xml.Attr{Name: xml.Name{Local: "method"}, Value: "unicorns"}}
	e.EncodeToken(start)
	e.EncodeToken(xml.Comment("We can emit comments, too."))

	e.EncodeToken(xml.Comment("Value of X"))
	e.EncodeElement(v[0], xml.StartElement{Name: xml.Name{Local: "X"}})

	e.EncodeToken(xml.Comment("Value of Y"))
	e.EncodeElement(v[1], xml.StartElement{Name: xml.Name{Local: "Y"}})

	e.EncodeToken(xml.Comment("Value of Z"))
	e.EncodeElement(v[2], xml.StartElement{Name: xml.Name{Local: "Z"}})

	e.EncodeToken(xml.EndElement{Name: start.Name})

	return nil
}

func (v *ManualVector) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// We don't really care about comments or unicorns
	vX := struct{ X, Y, Z int }{}
	d.DecodeElement(&vX, &start)

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
	err := xml.Unmarshal(encoded[0], &vP)
	fmt.Println(vP, err)

	vC := CopyVector{}
	err = xml.Unmarshal(encoded[1], &vC)
	fmt.Println(vC, err)

	vM := ManualVector{}
	err = xml.Unmarshal(encoded[2], &vM)
	fmt.Println(vM, err)
}
