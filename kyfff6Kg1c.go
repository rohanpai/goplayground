// Example SVG parser using a combination of xml.Unmarshal and the
// xml.Unmarshaler interface to handle an unknown combination of group
// elements where order is important.

package main

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

type Path struct {
	Id string `xml:"id,attr"`
	D  string `xml:"d, attr"`
}

type Rect struct {
	Id string `xml:"id,attr"`
}

type Group struct {
	Id          string
	Stroke      string
	StrokeWidth int32
	Fill        string
	FillRule    string
	Elements    []interface{}
}

type Svg struct {
	Title  string  `xml:"title"`
	Groups []Group `xml:"g"`
}

// Implements encoding.xml.Unmarshaler interface
func (g *Group) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "id":
			g.Id = attr.Value
		case "stroke":
			g.Stroke = attr.Value
		case "stroke-width":
			if intValue, err := strconv.ParseInt(attr.Value, 10, 32); err != nil {
				return err
			} else {
				g.StrokeWidth = int32(intValue)
			}
		case "fill":
			g.Fill = attr.Value
		case "fill-rule":
			g.FillRule = attr.Value
		}
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch tok := token.(type) {
		case xml.StartElement:
			var elementStruct interface{}

			switch tok.Name.Local {
			case "rect":
				elementStruct = &Rect{}
			case "path":
				elementStruct = &Path{}
			}

			if err = decoder.DecodeElement(elementStruct, &tok); err != nil {
				return err
			} else {
				g.Elements = append(g.Elements, elementStruct)
			}

			fmt.Println(tok.Name)

		case xml.EndElement:
			return nil
		}
	}
}

func ParseSvg(str string) *Svg {
	svg := &Svg{}

	err := xml.Unmarshal([]byte(str), &svg)
	if err != nil {
		fmt.Printf("ParseSvg Error: %v\n", err)
		return nil
	}
	return svg
}

func main() {
	fmt.Println("SVG Parsing Test")
	
	svgStr := `
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<svg width="79px" height="114px" viewBox="0 0 79 114" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xmlns:sketch="http://www.bohemiancoding.com/sketch/ns">
    <!-- Generator: Sketch 3.0.4 (8053) - http://www.bohemiancoding.com/sketch -->
    <title>ship</title>
    <desc>Created with Sketch.</desc>
    <defs></defs>
    <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" sketch:type="MSPage">
        <path d="M70.7470703,54.9351921 C59.4438539,23.2101932 39.4404297,-0.0302734375 39.4404297,-0.0302734375 C39.4404297,-0.0302734375 19.8288957,22.9468825 8.09220641,54.9351916 C3.08063764,68.5942062 -0.495117188,83.8962169 -0.495117188,99.7539062 L19.5214844,108.566406 L59.1821013,108.566406 L79.4462891,100.046875 C79.4462891,100.046875 75.1865234,67.3955078 70.7470703,54.9351921 Z" id="Path-1" fill="#D8D8D8" sketch:type="MSShapeGroup"></path>
        <rect id="Rectangle-1" fill="#F6A623" sketch:type="MSShapeGroup" x="22" y="107" width="34" height="7"></rect>
    </g>
</svg>
`

	svg := ParseSvg(svgStr)
	for _, elem := range svg.Groups[0].Elements {
		fmt.Println("Group element: ", elem)
	}
}
