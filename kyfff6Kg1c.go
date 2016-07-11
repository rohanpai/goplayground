// Example SVG parser using a combination of xml.Unmarshal and the
// xml.Unmarshaler interface to handle an unknown combination of group
// elements where order is important.

package main

import (
	&#34;encoding/xml&#34;
	&#34;fmt&#34;
	&#34;strconv&#34;
)

type Path struct {
	Id string `xml:&#34;id,attr&#34;`
	D  string `xml:&#34;d, attr&#34;`
}

type Rect struct {
	Id string `xml:&#34;id,attr&#34;`
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
	Title  string  `xml:&#34;title&#34;`
	Groups []Group `xml:&#34;g&#34;`
}

// Implements encoding.xml.Unmarshaler interface
func (g *Group) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case &#34;id&#34;:
			g.Id = attr.Value
		case &#34;stroke&#34;:
			g.Stroke = attr.Value
		case &#34;stroke-width&#34;:
			if intValue, err := strconv.ParseInt(attr.Value, 10, 32); err != nil {
				return err
			} else {
				g.StrokeWidth = int32(intValue)
			}
		case &#34;fill&#34;:
			g.Fill = attr.Value
		case &#34;fill-rule&#34;:
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
			case &#34;rect&#34;:
				elementStruct = &amp;Rect{}
			case &#34;path&#34;:
				elementStruct = &amp;Path{}
			}

			if err = decoder.DecodeElement(elementStruct, &amp;tok); err != nil {
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
	svg := &amp;Svg{}

	err := xml.Unmarshal([]byte(str), &amp;svg)
	if err != nil {
		fmt.Printf(&#34;ParseSvg Error: %v\n&#34;, err)
		return nil
	}
	return svg
}

func main() {
	fmt.Println(&#34;SVG Parsing Test&#34;)
	
	svgStr := `
&lt;?xml version=&#34;1.0&#34; encoding=&#34;UTF-8&#34; standalone=&#34;no&#34;?&gt;
&lt;svg width=&#34;79px&#34; height=&#34;114px&#34; viewBox=&#34;0 0 79 114&#34; version=&#34;1.1&#34; xmlns=&#34;http://www.w3.org/2000/svg&#34; xmlns:xlink=&#34;http://www.w3.org/1999/xlink&#34; xmlns:sketch=&#34;http://www.bohemiancoding.com/sketch/ns&#34;&gt;
    &lt;!-- Generator: Sketch 3.0.4 (8053) - http://www.bohemiancoding.com/sketch --&gt;
    &lt;title&gt;ship&lt;/title&gt;
    &lt;desc&gt;Created with Sketch.&lt;/desc&gt;
    &lt;defs&gt;&lt;/defs&gt;
    &lt;g id=&#34;Page-1&#34; stroke=&#34;none&#34; stroke-width=&#34;1&#34; fill=&#34;none&#34; fill-rule=&#34;evenodd&#34; sketch:type=&#34;MSPage&#34;&gt;
        &lt;path d=&#34;M70.7470703,54.9351921 C59.4438539,23.2101932 39.4404297,-0.0302734375 39.4404297,-0.0302734375 C39.4404297,-0.0302734375 19.8288957,22.9468825 8.09220641,54.9351916 C3.08063764,68.5942062 -0.495117188,83.8962169 -0.495117188,99.7539062 L19.5214844,108.566406 L59.1821013,108.566406 L79.4462891,100.046875 C79.4462891,100.046875 75.1865234,67.3955078 70.7470703,54.9351921 Z&#34; id=&#34;Path-1&#34; fill=&#34;#D8D8D8&#34; sketch:type=&#34;MSShapeGroup&#34;&gt;&lt;/path&gt;
        &lt;rect id=&#34;Rectangle-1&#34; fill=&#34;#F6A623&#34; sketch:type=&#34;MSShapeGroup&#34; x=&#34;22&#34; y=&#34;107&#34; width=&#34;34&#34; height=&#34;7&#34;&gt;&lt;/rect&gt;
    &lt;/g&gt;
&lt;/svg&gt;
`

	svg := ParseSvg(svgStr)
	for _, elem := range svg.Groups[0].Elements {
		fmt.Println(&#34;Group element: &#34;, elem)
	}
}
