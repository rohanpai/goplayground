package main

import (
   // use the github.com/clbanning/x2j.go imports
   &#34;bytes&#34;
   &#34;encoding/json&#34;
   &#34;encoding/xml&#34;
   &#34;errors&#34;
   &#34;fmt&#34;
   &#34;regexp&#34;
   &#34;strconv&#34;
   &#34;strings&#34;
)

var doc = `
&lt;doc&gt;
   &lt;image src=&#34;something.png&#34;&gt;&lt;/image&gt;
   &lt;something&gt;
      &lt;image src=&#34;something else.jpg&#34;&gt;&lt;/image&gt;
      &lt;title&gt;something else&lt;/title&gt;
   &lt;/something&gt;
   &lt;more_stuff&gt;
      &lt;some_images&gt;
         &lt;image src=&#34;first.gif&#34;&gt;&lt;/image&gt;
         &lt;image src=&#34;second.gif&#34;&gt;&lt;/image&gt;
      &lt;/some_images&gt;
   &lt;/more_stuff&gt;
&lt;/doc&gt;`

func main() {
   // get all image tag values
   images,err := ValuesForTag(doc,&#34;image&#34;)
   if err != nil {
      fmt.Println(&#34;error parsing doc:&#34;,err.Error())
      return
   }

   sources := make([]string,0)
   for _, v := range images {
       m := make(map[string]interface{},1)
       m[&#34;dummy&#34;] = v
       ss := ValuesForKey(m,&#34;-src&#34;)
       for _,s := range ss {
          sources = append(sources,s.(string))
       }
    }

   for _,src := range sources {
      fmt.Println(src)
   }
}

//==================== github/clbanning/x2j.go ================

type Node struct {
	dup   bool   // is member of a list
	attr  bool   // is an attribute
	key   string // XML tag
	val   string // element value
	nodes []*Node
}

// DocToJson - return an XML doc as a JSON string.
//	If the optional argument &#39;recast&#39; is &#39;true&#39;, then values will be converted to boolean or float64 if possible.
func DocToJson(doc string, recast ...bool) (string, error) {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}
	m, merr := DocToMap(doc, r)
	if m == nil || merr != nil {
		return &#34;&#34;, merr
	}

	b, berr := json.Marshal(m)
	if berr != nil {
		return &#34;&#34;, berr
	}

	// NOTE: don&#39;t have to worry about safe JSON marshaling with json.Marshal, since &#39;&lt;&#39; and &#39;&gt;&#34; are reservedin XML.
	return string(b), nil
}

// DocToJsonIndent - return an XML doc as a prettified JSON string.
//	If the optional argument &#39;recast&#39; is &#39;true&#39;, then values will be converted to boolean or float64 if possible.
//	Note: recasting is only applied to element values, not attribute values.
func DocToJsonIndent(doc string, recast ...bool) (string, error) {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}
	m, merr := DocToMap(doc, r)
	if m == nil || merr != nil {
		return &#34;&#34;, merr
	}

	b, berr := json.MarshalIndent(m, &#34;&#34;, &#34;  &#34;)
	if berr != nil {
		return &#34;&#34;, berr
	}

	// NOTE: don&#39;t have to worry about safe JSON marshaling with json.Marshal, since &#39;&lt;&#39; and &#39;&gt;&#34; are reservedin XML.
	return string(b), nil
}

// DocToMap - convert an XML doc into a map[string]interface{}.
// (This is analogous to unmarshalling a JSON string to map[string]interface{} using json.Unmarshal().)
//	If the optional argument &#39;recast&#39; is &#39;true&#39;, then values will be converted to boolean or float64 if possible.
//	Note: recasting is only applied to element values, not attribute values.
func DocToMap(doc string, recast ...bool) (map[string]interface{}, error) {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}
	n, err := DocToTree(doc)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	m[n.key] = n.treeToMap(r)

	return m, nil
}

// DocToTree - convert an XML doc into a tree of nodes.
func DocToTree(doc string) (*Node, error) {
	// xml.Decoder doesn&#39;t properly handle whitespace in some doc
	// see songTextString.xml test case ...
	reg, _ := regexp.Compile(&#34;[ \t\n\r]*&lt;&#34;)
	doc = reg.ReplaceAllString(doc, &#34;&lt;&#34;)

	b := bytes.NewBufferString(doc)
	p := xml.NewDecoder(b)
	n, berr := xmlToTree(&#34;&#34;, nil, p)
	if berr != nil {
		return nil, berr
	}

	return n, nil
}

// (*Node)WriteTree - convert a tree of nodes into a printable string.
//	&#39;padding&#39; is the starting indentation; typically: n.WriteTree().
func (n *Node) WriteTree(padding ...int) string {
	var indent int
	if len(padding) == 1 {
		indent = padding[0]
	}

	var s string
	if n.val != &#34;&#34; {
		for i := 0; i &lt; indent; i&#43;&#43; {
			s &#43;= &#34;  &#34;
		}
		s &#43;= n.key &#43; &#34; : &#34; &#43; n.val &#43; &#34;\n&#34;
	} else {
		for i := 0; i &lt; indent; i&#43;&#43; {
			s &#43;= &#34;  &#34;
		}
		s &#43;= n.key &#43; &#34; :&#34; &#43; &#34;\n&#34;
		for _, nn := range n.nodes {
			s &#43;= nn.WriteTree(indent &#43; 1)
		}
	}
	return s
}

// xmlToTree - load a &#39;clean&#39; XML doc into a tree of *Node.
func xmlToTree(skey string, a []xml.Attr, p *xml.Decoder) (*Node, error) {
	n := new(Node)
	n.nodes = make([]*Node, 0)

	if skey != &#34;&#34; {
		n.key = skey
		if len(a) &gt; 0 {
			for _, v := range a {
				na := new(Node)
				na.attr = true
				na.key = `-` &#43; v.Name.Local
				na.val = v.Value
				n.nodes = append(n.nodes, na)
			}
		}
	}
	for {
		t, err := p.Token()
		if err != nil {
			return nil, err
		}
		switch t.(type) {
		case xml.StartElement:
			tt := t.(xml.StartElement)
			// handle root
			if n.key == &#34;&#34; {
				n.key = tt.Name.Local
				if len(tt.Attr) &gt; 0 {
					for _, v := range tt.Attr {
						na := new(Node)
						na.attr = true
						na.key = `-` &#43; v.Name.Local
						na.val = v.Value
						n.nodes = append(n.nodes, na)
					}
				}
			} else {
				nn, nnerr := xmlToTree(tt.Name.Local, tt.Attr, p)
				if nnerr != nil {
					return nil, nnerr
				}
				n.nodes = append(n.nodes, nn)
			}
		case xml.EndElement:
			// scan n.nodes for duplicate n.key values
			n.markDuplicateKeys()
			return n, nil
		case xml.CharData:
			tt := string(t.(xml.CharData))
			if len(n.nodes) &gt; 0 {
				nn := new(Node)
				nn.key = &#34;#text&#34;
				nn.val = tt
				n.nodes = append(n.nodes, nn)
			} else {
				n.val = tt
			}
		default:
			// noop
		}
	}
	// Logically we can&#39;t get here, but provide an error message anyway.
	return nil, errors.New(&#34;Unknown parse error in xmlToTree() for: &#34; &#43; n.key)
}

// (*Node)markDuplicateKeys - set node.dup flag for loading map[string]interface{}.
func (n *Node) markDuplicateKeys() {
	l := len(n.nodes)
	for i := 0; i &lt; l; i&#43;&#43; {
		if n.nodes[i].dup {
			continue
		}
		for j := i &#43; 1; j &lt; l; j&#43;&#43; {
			if n.nodes[i].key == n.nodes[j].key {
				n.nodes[i].dup = true
				n.nodes[j].dup = true
			}
		}
	}
}

// (*Node)treeToMap - convert a tree of nodes into a map[string]interface{}.
//	(Parses to map that is structurally the same as from json.Unmarshal().)
// Note: root is not instantiated; call with: &#34;m[n.key] = n.treeToMap(recast)&#34;.
func (n *Node) treeToMap(r bool) interface{} {
	if len(n.nodes) == 0 {
		return recast(n.val, r)
	}

	m := make(map[string]interface{}, 0)
	for _, v := range n.nodes {
		// just a value
		if !v.dup &amp;&amp; len(v.nodes) == 0 {
			m[v.key] = recast(v.val, r)
			continue
		}

		// a list of values
		if v.dup {
			var a []interface{}
			if vv, ok := m[v.key]; ok {
				a = vv.([]interface{})
			} else {
				a = make([]interface{}, 0)
			}
			a = append(a, v.treeToMap(r))
			m[v.key] = interface{}(a)
			continue
		}

		// it&#39;s a unique key
		m[v.key] = v.treeToMap(r)
	}

	return interface{}(m)
}

// recast - try to cast string values to bool or float64
func recast(s string, r bool) interface{} {
	if r {
		// handle numeric strings ahead of boolean
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return interface{}(f)
		}
		// ParseBool treats &#34;1&#34;==true &amp; &#34;0&#34;==false
		if b, err := strconv.ParseBool(s); err == nil {
			return interface{}(b)
		}
	}
	return interface{}(s)
}

// WriteMap - dumps the map[string]interface{} for examination.
//	&#39;offset&#39; is initial indentation count; typically: WriteMap(m).
//	NOTE: with XML all element types are &#39;string&#39;.
//	But code written as generic for use with maps[string]interface{} values from json.Unmarshal().
//	Or it can handle a DocToMap(doc,true) result where values have be recast&#39;d.
func WriteMap(m interface{}, offset ...int) string {
	var indent int
	if len(offset) == 1 {
		indent = offset[0]
	}

	var s string
	switch m.(type) {
	case nil:
		return &#34;[nil] nil&#34;
	case string:
		return &#34;[string] &#34; &#43; m.(string)
	case float64:
		return &#34;[float64] &#34; &#43; strconv.FormatFloat(m.(float64), &#39;e&#39;, 2, 64)
	case bool:
		return &#34;[bool] &#34; &#43; strconv.FormatBool(m.(bool))
	case []interface{}:
		s &#43;= &#34;[[]interface{}]&#34;
		for i, v := range m.([]interface{}) {
			s &#43;= &#34;\n&#34;
			for i := 0; i &lt; indent; i&#43;&#43; {
				s &#43;= &#34;  &#34;
			}
			s &#43;= &#34;[item: &#34; &#43; strconv.FormatInt(int64(i), 10) &#43; &#34;]&#34;
			switch v.(type) {
			case string, float64, bool:
				s &#43;= &#34;\n&#34;
			default:
				// noop
			}
			for i := 0; i &lt; indent; i&#43;&#43; {
				s &#43;= &#34;  &#34;
			}
			s &#43;= WriteMap(v, indent&#43;1)
		}
	case map[string]interface{}:
		for k, v := range m.(map[string]interface{}) {
			s &#43;= &#34;\n&#34;
			for i := 0; i &lt; indent; i&#43;&#43; {
				s &#43;= &#34;  &#34;
			}
			// s &#43;= &#34;[map[string]interface{}] &#34;&#43;k&#43;&#34; :&#34;&#43;WriteMap(v,indent&#43;1)
			s &#43;= k &#43; &#34; :&#34; &#43; WriteMap(v, indent&#43;1)
		}
	default:
		// shouldn&#39;t ever be here ...
		s &#43;= fmt.Sprintf(&#34;unknown type for: %v&#34;, m)
	}
	return s
}

// ------------------------  value extraction from XML doc --------------------------

// DocValue - return a value for a specific tag
//	&#39;doc&#39; is a valid XML message.
//	&#39;path&#39; is a hierarchy of XML tags, e.g., &#34;doc.name&#34;.
//	&#39;attrs&#39; is an optional list of &#34;name:value&#34; pairs for attributes.
//	Note: &#39;recast&#39; is not enabled here. Use DocToMap(), NewAttributeMap(), and MapValue() calls for that.
func DocValue(doc, path string, attrs ...string) (interface{}, error) {
	n, err := DocToTree(doc)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	m[n.key] = n.treeToMap(false)

	a, aerr := NewAttributeMap(attrs...)
	if aerr != nil {
		return nil, aerr
	}
	v, verr := MapValue(m, path, a)
	if verr != nil {
		return nil, verr
	}
	return v, nil
}

// MapValue - retrieves value based on walking the map, &#39;m&#39;.
//	&#39;m&#39; is the map value of interest.
//	&#39;path&#39; is a period-separated hierarchy of keys in the map.
//	&#39;attr&#39; is a map of attribute &#34;name:value&#34; pairs from NewAttributeMap().
//	If the path can&#39;t be traversed, an error is returned.
//	Note: the optional argument &#39;r&#39; can be used to coerce attribute values, &#39;attr&#39;, if done so for &#39;m&#39;.
func MapValue(m map[string]interface{}, path string, attr map[string]interface{}, r ...bool) (interface{}, error) {
	// attribute values may have been recasted during map construction; default is &#39;false&#39;.
	if len(r) == 1 &amp;&amp; r[0] == true {
		for k, v := range attr {
			attr[k] = recast(v.(string), true)
		}
	}

	// parse the path
	keys := strings.Split(path, &#34;.&#34;)

	// initialize return value to &#39;m&#39; so a path of &#34;&#34; will work correctly
	var v interface{} = m
	var ok bool
	var okey string
	var isMap bool = true
	if keys[0] == &#34;&#34; &amp;&amp; len(attr) == 0 {
		return v, nil
	}
	for _, key := range keys {
		if !isMap {
			return nil, errors.New(&#34;no keys beyond: &#34; &#43; okey)
		}
		if v, ok = m[key]; !ok {
			return nil, errors.New(&#34;no key in map: &#34; &#43; key)
		} else {
			switch v.(type) {
			case map[string]interface{}:
				m = v.(map[string]interface{})
				isMap = true
			default:
				isMap = false
			}
		}
		// save &#39;key&#39; for error reporting
		okey = key
	}

	// match attributes; value is &#34;#text&#34; or nil
	return hasAttributes(v, attr)
}

// hasAttributes() - interface{} equality works for string, float64, bool
func hasAttributes(v interface{}, a map[string]interface{}) (interface{}, error) {
	switch v.(type) {
	case []interface{}:
		// run through all entries looking one with matching attributes
		for _, vv := range v.([]interface{}) {
			if vvv, vvverr := hasAttributes(vv, a); vvverr == nil {
				return vvv, nil
			}
		}
		return nil, errors.New(&#34;no list member with matching attributes&#34;)
	case map[string]interface{}:
		// do all attribute name:value pairs match?
		nv := v.(map[string]interface{})
		for key, val := range a {
			if vv, ok := nv[key]; !ok {
				return nil, errors.New(&#34;no attribute with name: &#34; &#43; key[1:])
			} else if val != vv {
				return nil, errors.New(&#34;no attribute key:value pair: &#34; &#43; fmt.Sprintf(&#34;%s:%v&#34;, key[1:], val))
			}
		}
		// they all match; so return value associated with &#34;#text&#34; key.
		if vv, ok := nv[&#34;#text&#34;]; ok {
			return vv, nil
		} else {
			// this happens when another element is value of tag rather than just a string value
			return nv, nil
		}
	}
	return nil, errors.New(&#34;no match for attributes&#34;)
}

// NewAttributeMap() - generate map of attributes=value entries as map[&#34;-&#34;&#43;string]string.
//	&#39;kv&#39; arguments are &#34;name:value&#34; pairs that appear as attributes, name=&#34;value&#34;.
func NewAttributeMap(kv ...string) (map[string]interface{}, error) {
	m := make(map[string]interface{}, 0)
	for _, v := range kv {
		vv := strings.Split(v, &#34;:&#34;)
		if len(vv) != 2 {
			return nil, errors.New(&#34;attribute not \&#34;name:value\&#34; pair: &#34; &#43; v)
		}
		// attributes are stored as keys prepended with hyphen
		m[&#34;-&#34;&#43;vv[0]] = interface{}(vv[1])
	}
	return m, nil
}

//------------------------- get values for key ----------------------------

// ValuesForTag - return all values in doc associated with &#39;tag&#39;.
//	Returns nil if the &#39;tag&#39; does not occur in the doc.
//	If there is an error encounted while parsing doc, that is returned.
//	If you want values &#39;recast&#39; use DocToMap() and ValuesForKey().
func ValuesForTag(doc, tag string) ([]interface{}, error) {
	n, err := DocToTree(doc)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	m[n.key] = n.treeToMap(false)

	return ValuesForKey(m, tag), nil
}

// ValuesForKey - return all values in map associated with &#39;key&#39;
//	Returns nil if the &#39;key&#39; does not occur in the map
func ValuesForKey(m map[string]interface{}, key string) []interface{} {
	ret := make([]interface{}, 0)

	hasKey(m, key, &amp;ret)
	if len(ret) &gt; 0 {
		return ret
	}
	return nil
}

// hasKey - if the map &#39;key&#39; exists append it to array
//          if it doesn&#39;t do nothing except scan array and map values
func hasKey(iv interface{}, key string, ret *[]interface{}) {
	switch iv.(type) {
	case map[string]interface{}:
		vv := iv.(map[string]interface{})
		if v, ok := vv[key]; ok {
			*ret = append(*ret, v)
		}
		for _, v := range iv.(map[string]interface{}) {
			hasKey(v, key, ret)
		}
	case []interface{}:
		for _, v := range iv.([]interface{}) {
			hasKey(v, key, ret)
		}
	}
}

