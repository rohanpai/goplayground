package main

import (
	&#34;encoding/json&#34;
	&#34;fmt&#34;
	&#34;time&#34;
)

type Prefix struct {
	Network string
	Mask    int
}

func main() {

	valueStruct()
	time.Sleep(time.Second * 1)
	literalStruct()
	time.Sleep(time.Second * 1)
	pointerStruct()
}

func valueStruct() {
	// struct as a value
	var nw Prefix
	nw.Network = &#34;10.1.1.0&#34;
	nw.Mask = 24
	fmt.Println(&#34;### struct as a pointer ###&#34;)
	PrettyPrint(&amp;nw)
}

func literalStruct() {
	// literal structs are the shortest LOC
	nw2 := &amp;Prefix{&#34;10.1.2.0&#34;, 30}
	fmt.Println(&#34;### struct as a literal ###&#34;)
	PrettyPrint(nw2)
}

func pointerStruct() {
	// struct as a pointer
	nw3 := new(Prefix)
	// very similar to setters/getters in OOP
	nw3.Network = &#34;10.1.1.0&#34;
	// or even like so
	(*nw3).Mask = 28
	fmt.Println(&#34;### struct as a pointer ###&#34;)
	PrettyPrint(nw3)
}

// print the contents of the network obj
func PrettyPrint(data interface{}) {
	var p []byte
	//    var err := error
	p, err := json.MarshalIndent(data, &#34;&#34;, &#34;\t&#34;)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf(&#34;%s \n&#34;, p)
}
