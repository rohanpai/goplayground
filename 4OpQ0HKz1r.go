package main

import (
	&#34;bufio&#34;
	&#34;flag&#34;
	&#34;fmt&#34;
	&#34;io&#34;
	&#34;net/url&#34;
	&#34;os&#34;
	&#34;regexp&#34;
	&#34;strings&#34;
)

var inputFile = flag.String(&#34;infile&#34;, &#34;freebase-rdf&#34;, &#34;Input file path&#34;)
var filter, _ = regexp.Compile(&#34;^file:.*|^talk:.*|^special:.*|^wikipedia:.*|^wiktionary:.*|^user:.*|^user_talk:.*&#34;)

type Redirect struct {
	Title string `xml:&#34;title,attr&#34;`
}

type Page struct {
	Title    string `xml:&#34;title&#34;`
	Abstract string `xml:&#34;&#34;`
}

func CanonicaliseTitle(title string) string {
	can := strings.ToLower(title)
	can = strings.Replace(can, &#34; &#34;, &#34;_&#34;, -1)
	can = url.QueryEscape(can)
	return can
}

func convertFreebaseId(uri string) string {
	if strings.HasPrefix(uri, &#34;&lt;&#34;) &amp;&amp; strings.HasSuffix(uri, &#34;&gt;&#34;) {
		var id = uri[1 : len(uri)-1]
		id = strings.Replace(id, &#34;http://rdf.freebase.com/ns&#34;, &#34;&#34;, -1)
		id = strings.Replace(id, &#34;.&#34;, &#34;/&#34;, -1)
		return id
	}
	return uri
}

func parseTriple(line string) (string, string, string) {
	var parts = strings.Split(line, &#34;\t&#34;)
	subject := convertFreebaseId(parts[0])
	predicate := convertFreebaseId(parts[1])
	object := convertFreebaseId(parts[2])
	return subject, predicate, object
}

func processTopic(id string, properties map[string][]string, file io.Writer) {
	fmt.Fprint(file, &#34;&lt;card&gt;\n&#34;)
	fmt.Fprintf(file, `&lt;title&gt;&#34;%s&#34;&lt;/title&gt;\n`, properties[&#34;/type/object/name&#34;])
	fmt.Fprintf(file, `&lt;image&gt;&#34;%s&#34;&lt;/image&gt;\n`, &#34;https://usercontent.googleapis.com/freebase/v1/image&#34;, id)
	fmt.Fprintf(file, `&lt;text&gt;&#34;%s&#34;&lt;/text&gt;\n`, properties[&#34;/common/document/text&#34;])
	fmt.Fprintf(file, &#34;&lt;facts&gt;&#34;)
	for k, v := range properties {
		for _, value := range v {
			fmt.Fprintf(file, &#34;&lt;fact property=\&#34;%s\&#34;&gt;%s&lt;/fact&gt;\n&#34;, k, value)
		}
	}
	fmt.Fprintf(file, &#34;&lt;/facts&gt;\n&#34;)
	fmt.Fprintf(file, &#34;&lt;/card&gt;\n&#34;)
}

func main() {
	var current_mid = &#34;&#34;
	current_topic := make(map[string][]string)
	f, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	r := bufio.NewReader(f)
	xmlFile, _ := os.Create(&#34;freebase.xml&#34;)
	line, err := r.ReadString(&#39;\n&#39;)
	for err == nil {
		subject, predicate, object := parseTriple(line)
		if subject == current_mid {
			current_topic[predicate] = append(current_topic[predicate], object)
		} else if len(current_mid) &gt; 0 {
			processTopic(current_mid, current_topic, xmlFile)
			current_topic = make(map[string][]string)
		}
		current_mid = subject
		line, err = r.ReadString(&#39;\n&#39;)
	}
	processTopic(current_mid, current_topic, xmlFile)
	if err != io.EOF {
		fmt.Println(err)
		return
	}
}
