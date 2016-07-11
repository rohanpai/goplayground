// I&#39;m trying to write one function that encodes/decodes various types of Messages.
// In OO languages, I would use type inheritance, but Go doesn&#39;t have this concept,
// ref: http://golang.org/doc/faq#inheritance, so instead, here I&#39;m trying yo use
// the &#34;marker interface&#34; style to leverage interface inheritance for this purpose instead.

package main

import (
    &#34;fmt&#34;
	&#34;bytes&#34;
	&#34;log&#34;
	&#34;encoding/gob&#34;
)

type Msger interface {
	IsMsg()
}

type ClientMsg struct {
	Id      string
}
type ServerMsg struct {
	Id      string
}

func (m ClientMsg) IsMsg() { }
func (m ServerMsg) IsMsg() { }

func encode(m Msger) (bb bytes.Buffer) {
	enc := gob.NewEncoder(&amp;bb)
	err := enc.Encode(m)
	if err != nil {
		log.Fatal(&#34;Cannot encode! err=&#34;, err)
	}
	return
}

func decode(m Msger, bb bytes.Buffer) {     // for A
//func decode(m *Msger, bb bytes.Buffer) {  // for B, C
	dec := gob.NewDecoder(&amp;bb)
	err := dec.Decode(&amp;m)
	if err != nil {
		log.Fatal(&#34;Cannot decode Msg! err=&#34;, err)
	}
	return
}

func main() {
    m1 := ClientMsg{&#34;id_1&#34;}
	b1 := encode(m1)
	p_bb := bytes.NewBuffer(b1.Bytes())

	var mDecoded ClientMsg // for A, B, C
	//var mDecoded interface{} // D: : cannot use mDecoded (type interface {}) as type Msger in argument to decode: interface {} does not implement Msger (missing IsMsg method)
	decode(mDecoded, *p_bb)     // A: gives: Cannot decode Msg! err=gob: local interface type *main.Msger can only be decoded from remote interface type; received concrete type ClientMsg = struct { Id string; }

	//decode(mDecoded, *p_bb)  // B: gives: cannot use mDecoded (type ClientMsg) as type *Msger in argument to decode: *Msger is pointer to interface, not interface
	//decode(&amp;mDecoded, *p_bb) // C: gives: cannot use &amp;mDecoded (type *ClientMsg) as type *Msger in argument to decode: *Msger is pointer to interface, not interface

	fmt.Printf(&#34;m1 Decoded=&#39;%v&#39;\n&#34;, mDecoded)
}
