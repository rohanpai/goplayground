package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var inputFile = flag.String("infile", "freebase-rdf", "Input file path")
var filter, _ = regexp.Compile("^file:.*|^talk:.*|^special:.*|^wikipedia:.*|^wiktionary:.*|^user:.*|^user_talk:.*")

type Redirect struct {
	Title string `xml:"title,attr"`
}

type Page struct {
	Title    string `xml:"title"`
	Abstract string `xml:""`
}

func CanonicaliseTitle(title string) string {
	can := strings.ToLower(title)
	can = strings.Replace(can, " ", "_", -1)
	can = url.QueryEscape(can)
	return can
}

func convertFreebaseId(uri string) string {
	if strings.HasPrefix(uri, "<") && strings.HasSuffix(uri, ">") {
		var id = uri[1 : len(uri)-1]
		id = strings.Replace(id, "http://rdf.freebase.com/ns", "", -1)
		id = strings.Replace(id, ".", "/", -1)
		return id
	}
	return uri
}

func parseTriple(line string) (string, string, string) {
	var parts = strings.Split(line, "\t")
	subject := convertFreebaseId(parts[0])
	predicate := convertFreebaseId(parts[1])
	object := convertFreebaseId(parts[2])
	return subject, predicate, object
}

func processTopic(id string, properties map[string][]string, file io.Writer) {
	fmt.Fprint(file, "<card>\n")
	fmt.Fprintf(file, `<title>"%s"</title>\n`, properties["/type/object/name"])
	fmt.Fprintf(file, `<image>"%s"</image>\n`, "https://usercontent.googleapis.com/freebase/v1/image", id)
	fmt.Fprintf(file, `<text>"%s"</text>\n`, properties["/common/document/text"])
	fmt.Fprintf(file, "<facts>")
	for k, v := range properties {
		for _, value := range v {
			fmt.Fprintf(file, "<fact property=\"%s\">%s</fact>\n", k, value)
		}
	}
	fmt.Fprintf(file, "</facts>\n")
	fmt.Fprintf(file, "</card>\n")
}

func main() {
	var current_mid = ""
	current_topic := make(map[string][]string)
	f, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	r := bufio.NewReader(f)
	xmlFile, _ := os.Create("freebase.xml")
	line, err := r.ReadString('\n')
	for err == nil {
		subject, predicate, object := parseTriple(line)
		if subject == current_mid {
			current_topic[predicate] = append(current_topic[predicate], object)
		} else if len(current_mid) > 0 {
			processTopic(current_mid, current_topic, xmlFile)
			current_topic = make(map[string][]string)
		}
		current_mid = subject
		line, err = r.ReadString('\n')
	}
	processTopic(current_mid, current_topic, xmlFile)
	if err != io.EOF {
		fmt.Println(err)
		return
	}
}
