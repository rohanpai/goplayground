package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var testXML = `
<?xml version="1.0" encoding="UTF-8"?>
<data>
    <netid>
        <disable>no</disable>
        <text1>default:text</text1>
        <word1>default:word</word1>
    </netid>
</data>
`

var testXML2 = `
<?xml version="1.0" encoding="UTF-8"?>
<data>
    <idnet>
        <disable>yes</disable>
        <text1>default:text2</text1>
        <word1>default:word2</word1>
    </idnet>
</data>
`
func main() {
	var docs = []string{testXML,testXML2}
	for _,doc := range docs {
		v,_ := ValuesFromTagPath(doc,"data.*")
		m := v[0].(map[string]interface{})
		fmt.Println(m["disable"].(string))
		fmt.Println(m["text1"].(string))
		fmt.Println(m["word1"].(string))
	}
}

//-----------------------------
// https://github.com/clbanning/x2j/blob/master/x2j.go below

//	Unmarshal an arbitrary XML doc to a map[string]interface{} or a JSON string. 
// Copyright 2012-2013 Charles Banning. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file
/*	
	Unmarshal an arbitrary XML doc to a map[string]interface{} or a JSON string. 

	DocToMap() returns an intermediate result with the XML doc unmarshal'd to a map
	of type map[string]interface{}. It is analogous to unmarshal'ng a JSON string to
	a map using json.Unmarshal(). (This was the original purpose of this library.)

	DocToTree()/WriteTree() let you examine the parsed XML doc.

	XML values are all type 'string'. The optional argument 'recast' for DocToJson()
	and DocToMap() will convert element values to JSON data types - 'float64' and 'bool' -
	if possible.  This, however, should be done with caution as it will recast ALL numeric
	and boolean string values, even those that are meant to be of type 'string'.
 */
/* 
package x2j

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"strconv"
)
*/

type Node struct {
	dup bool			// is member of a list
	attr bool		// is an attribute
	key string		// XML tag
	val string		// element value
	nodes []*Node
}

// DocToJson - return an XML doc as a JSON string.
//	If the optional argument 'recast' is 'true', then values will be converted to boolean or float64 if possible.
func DocToJson(doc string,recast ...bool) (string,error) {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}
	m,merr := DocToMap(doc,r)
	if m == nil || merr != nil {
		return "",merr
	}

	b, berr := json.Marshal(m)
	if berr != nil {
		return "",berr
	}

	// NOTE: don't have to worry about safe JSON marshaling with json.Marshal, since '<' and '>" are reservedin XML.
	return string(b),nil
}

// DocToJsonIndent - return an XML doc as a prettified JSON string.
//	If the optional argument 'recast' is 'true', then values will be converted to boolean or float64 if possible.
//	Note: recasting is only applied to element values, not attribute values.
func DocToJsonIndent(doc string,recast ...bool) (string,error) {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}
	m,merr := DocToMap(doc,r)
	if m == nil || merr != nil {
		return "",merr
	}

	b, berr := json.MarshalIndent(m,"","  ")
	if berr != nil {
		return "",berr
	}

	// NOTE: don't have to worry about safe JSON marshaling with json.Marshal, since '<' and '>" are reservedin XML.
	return string(b),nil
}

// DocToMap - convert an XML doc into a map[string]interface{}.
// (This is analogous to unmarshalling a JSON string to map[string]interface{} using json.Unmarshal().)
//	If the optional argument 'recast' is 'true', then values will be converted to boolean or float64 if possible.
//	Note: recasting is only applied to element values, not attribute values.
func DocToMap(doc string,recast ...bool) (map[string]interface{},error) {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}
	n,err := DocToTree(doc)
	if err != nil {
		return nil,err
	}

	m := make(map[string]interface{})
	m[n.key] = n.treeToMap(r)

	return m,nil
}

// DocToTree - convert an XML doc into a tree of nodes.
func DocToTree(doc string) (*Node, error) {
	// xml.Decoder doesn't properly handle whitespace in some doc
	// see songTextString.xml test case ... 
	reg,_ := regexp.Compile("[ \t\n\r]*<")
	doc = reg.ReplaceAllString(doc,"<")

	b := bytes.NewBufferString(doc)
	p := xml.NewDecoder(b)
	n, berr := xmlToTree("",nil,p)
	if berr != nil {
		return nil, berr
	}

	return n,nil
}

// (*Node)WriteTree - convert a tree of nodes into a printable string.
//	'padding' is the starting indentation; typically: n.WriteTree().
func (n *Node)WriteTree(padding ...int) string {
	var indent int
	if len(padding) == 1 {
		indent = padding[0]
	}

	var s string
	if n.val != "" {
		for i := 0 ; i < indent ; i++ {
			s += "  "
		}
		s += n.key+" : "+n.val+"\n"
	} else {
		for i := 0 ; i < indent ; i++ {
			s += "  "
		}
		s += n.key+" :"+"\n"
		for _,nn := range n.nodes {
			s += nn.WriteTree(indent+1)
		}
	}
	return s
}

// xmlToTree - load a 'clean' XML doc into a tree of *Node.
func xmlToTree(skey string,a []xml.Attr,p *xml.Decoder) (*Node, error) {
	n := new(Node)
	n.nodes = make([]*Node,0)

	if skey != "" {
		n.key = skey
		if len(a) > 0 {
			for _,v := range a {
				na := new(Node)
				na.attr = true
				na.key = `-`+v.Name.Local
				na.val = v.Value
				n.nodes = append(n.nodes,na)
			}
		}
	}
	for {
		t,err := p.Token()
		if err != nil {
			return nil, err
		}
		switch t.(type) {
			case xml.StartElement:
				tt := t.(xml.StartElement)
				// handle root
				if n.key == "" {
					n.key = tt.Name.Local
					if len(tt.Attr) > 0 {
						for _,v := range tt.Attr {
							na := new(Node)
							na.attr = true
							na.key = `-`+v.Name.Local
							na.val = v.Value
							n.nodes = append(n.nodes,na)
						}
					}
				} else {
					nn, nnerr := xmlToTree(tt.Name.Local,tt.Attr,p)
					if nnerr != nil {
						return nil, nnerr
					}
					n.nodes = append(n.nodes,nn)
				}
			case xml.EndElement:
				// scan n.nodes for duplicate n.key values
				n.markDuplicateKeys()
				return n, nil
			case xml.CharData:
				tt := string(t.(xml.CharData))
				if len(n.nodes) > 0 {
					nn := new(Node)
					nn.key = "#text"
					nn.val = tt
					n.nodes = append(n.nodes,nn)
				} else {
					n.val = tt
				}
			default:
				// noop
		}
	}
	// Logically we can't get here, but provide an error message anyway.
	return nil, errors.New("Unknown parse error in xmlToTree() for: "+n.key)
}

// (*Node)markDuplicateKeys - set node.dup flag for loading map[string]interface{}.
func (n *Node)markDuplicateKeys() {
	l := len(n.nodes)
	for i := 0 ; i < l ; i++ {
		if n.nodes[i].dup {
			continue
		}
		for j := i+1 ; j < l ; j++ {
			if n.nodes[i].key == n.nodes[j].key {
				n.nodes[i].dup = true
				n.nodes[j].dup = true
			}
		}
	}
}

// (*Node)treeToMap - convert a tree of nodes into a map[string]interface{}.
//	(Parses to map that is structurally the same as from json.Unmarshal().)
// Note: root is not instantiated; call with: "m[n.key] = n.treeToMap(recast)".
func (n *Node)treeToMap(r bool) interface{} {
	if len(n.nodes) == 0 {
		return recast(n.val,r)
	}

	m := make(map[string]interface{},0)
	for _,v := range n.nodes {
		// just a value
		if !v.dup && len(v.nodes) == 0 {
			m[v.key] = recast(v.val,r)
			continue
		}

		// a list of values
		if v.dup {
			var a []interface{}
			if vv,ok := m[v.key]; ok {
				a = vv.([]interface{})
			} else {
				a = make([]interface{},0)
			}
			a = append(a,v.treeToMap(r))
			m[v.key] = interface{}(a)
			continue
		}

		// it's a unique key
		m[v.key] = v.treeToMap(r)
	}

	return interface{}(m)
}

// recast - try to cast string values to bool or float64
func recast(s string,r bool) interface{} {
	if r {
		// handle numeric strings ahead of boolean
		if f, err := strconv.ParseFloat(s,64); err == nil {
			return interface{}(f)
		}
		// ParseBool treats "1"==true & "0"==false
		if b, err := strconv.ParseBool(s); err == nil {
			return interface{}(b)
		}
	}
	return interface{}(s)
}

// WriteMap - dumps the map[string]interface{} for examination.
//	'offset' is initial indentation count; typically: WriteMap(m).
//	NOTE: with XML all element types are 'string'.
//	But code written as generic for use with maps[string]interface{} values from json.Unmarshal().
//	Or it can handle a DocToMap(doc,true) result where values have be recast'd.
func WriteMap(m interface{}, offset ...int) string {
	var indent int
	if len(offset) == 1 {
		indent = offset[0]
	}

	var s string
	switch m.(type) {
		case nil:
			return "[nil] nil"
		case string:
			return "[string] "+m.(string)
		case float64:
			return "[float64] "+strconv.FormatFloat(m.(float64),'e',2,64)
		case bool:
			return "[bool] "+strconv.FormatBool(m.(bool))
		case []interface{}:
			s += "[[]interface{}]"
			for i,v := range m.([]interface{}) {
				s += "\n"
				for i := 0 ; i < indent ; i++ {
					s += "  "
				}
				s += "[item: "+strconv.FormatInt(int64(i),10)+"]"
				switch v.(type) {
					case string,float64,bool:
						s += "\n"
					default:
						// noop
				}
				for i := 0 ; i < indent ; i++ {
					s += "  "
				}
				s += WriteMap(v,indent+1)
			}
		case map[string]interface{}:
			for k,v := range m.(map[string]interface{}) {
				s += "\n"
				for i := 0 ; i < indent ; i++ {
					s += "  "
				}
				// s += "[map[string]interface{}] "+k+" :"+WriteMap(v,indent+1)
				s += k+" :"+WriteMap(v,indent+1)
		}
		default:
			// shouldn't ever be here ...
			s += fmt.Sprintf("unknown type for: %v",m)
	}
	return s
}

// ------------------------  value extraction from XML doc --------------------------

// ValuesFromTagPath - deliver all values for a path node from a XML doc
// If there are no values for the path 'nil' is returned.
// A return value of (nil, nil) means that there were no values and no errors parsing the doc.
//   'doc' is the XML document
//   'path' is a dot-separated path of tag nodes
//          If a node is '*', then everything beyond is scanned for values.
//          E.g., "doc.books' might return a single value 'book' of type []interface{}, but
//                "doc.books.*" could return all the 'book' entries as []map[string]interface{}.
//                "doc.books.*.author" might return all the 'author' tag values as []string - or
//            		"doc.books.*.author.lastname" might be required, depending on he schema.
func ValuesFromTagPath(doc, path string, getAttrs ...bool) ([]interface{}, error) {
	var a bool
	if len(getAttrs) == 1 {
		a = getAttrs[0]
	}
	m, err := DocToMap(doc)
	if err != nil {
		return nil, err
	}

	v := ValuesFromKeyPath(m, path, a)
	return v, nil
}

// ValuesFromKeyPath - deliver all values for a path node from a map[string]interface{}
// If there are no values for the path 'nil' is returned.
//   'm' is the map to be walked
//   'path' is a dot-separated path of key values
//          If a node is '*', then everything beyond is walked.
//          E.g., see ValuesForTagPath documentation.
func ValuesFromKeyPath(m map[string]interface{}, path string, getAttrs ...bool) []interface{} {
	var a bool
	if len(getAttrs) == 1 {
		a = getAttrs[0]
	}
	keys := strings.Split(path, ".")
	ret := make([]interface{}, 0)
	valuesFromKeyPath(&ret, m, keys, a)
	if len(ret) == 0 {
		return nil
	}
	return ret
}

func valuesFromKeyPath(ret *[]interface{}, m interface{}, keys []string, getAttrs bool) {
	lenKeys := len(keys)

	// load 'm' values into 'ret'
	// expand any lists
	if lenKeys == 0 {
		switch m.(type) {
		case map[string]interface{}:
			*ret = append(*ret, m)
		case []interface{}:
			for _, v := range m.([]interface{}) {
				*ret = append(*ret, v)
			}
		default:
			*ret = append(*ret, m)
		}
		return
	}

	// key of interest
	key := keys[0]
	switch key {
	case "*": // wildcard - scan all values
		switch m.(type) {
		case map[string]interface{}:
			for k, v := range m.(map[string]interface{}) {
				if string(k[:1]) == "-" && !getAttrs { // skip attributes?
					continue
				}
				valuesFromKeyPath(ret, v, keys[1:], getAttrs)
			}
		case []interface{}:
			for _, v := range m.([]interface{}) {
				switch v.(type) {
				// flatten out a list of maps - keys are processed
				case map[string]interface{}:
					for kk, vv := range v.(map[string]interface{}) {
						if string(kk[:1]) == "-" && !getAttrs { // skip attributes?
							continue
						}
						valuesFromKeyPath(ret, vv, keys[1:], getAttrs)
					}
				default:
					valuesFromKeyPath(ret, v, keys[1:], getAttrs)
				}
			}
		}
	default: // key - must be map[string]interface{}
		switch m.(type) {
		case map[string]interface{}:
			if v, ok := m.(map[string]interface{})[key]; ok {
//				if lenKeys == 1 {
//					*ret = append(*ret, v)
//					return
//				}
				valuesFromKeyPath(ret, v, keys[1:], getAttrs)
			}
		case []interface{}: // may be buried in list
			for _, v := range m.([]interface{}) {
				switch v.(type) {
				case map[string]interface{}:
					if vv, ok := v.(map[string]interface{})[key]; ok {
						valuesFromKeyPath(ret, vv, keys[1:], getAttrs)
					}
				}
			}
		}
	}
}