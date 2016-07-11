package main

import (
	"bytes"
	"fmt"
)

var trivialXml = []byte(`<root>a<foo>b</foo>c<bar>d</bar>e<bar>f</bar>g</root>`)

func main() {
	node, err := xmlpath.Parse(bytes.NewBuffer(trivialXml))
	if err != nil {
		panic(err)
	}
	path1 := xmlpath.MustCompile("root")
	path2 := xmlpath.MustCompile("foo")
	result1, ok1 := path1.String(node)
	result2, ok2 := path2.String(node)
	iter := path1.Iter(node)
	if !iter.Next() {
		panic("must exist")
	}
	result3, ok3 := path2.String(iter.Node())
	fmt.Printf("result1: %v %s\n", ok1, result1)
	fmt.Printf("result2: %v %s\n", ok2, result2)
	fmt.Printf("result3: %v %s\n", ok3, result3)
}
