// I'm trying to write one function that encodes/decodes various types of Messages.
// In OO languages, I would use type inheritance, but Go doesn't have this concept,
// ref: http://golang.org/doc/faq#inheritance, so instead, here I'm trying yo use
// the "marker interface" style to leverage interface inheritance for this purpose instead.

package main

import (
    "fmt"
	"bytes"
	"log"
	"encoding/gob"
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
	enc := gob.NewEncoder(&bb)
	err := enc.Encode(m)
	if err != nil {
		log.Fatal("Cannot encode! err=", err)
	}
	return
}

func decode(m Msger, bb bytes.Buffer) {     // for A
//func decode(m *Msger, bb bytes.Buffer) {  // for B, C
	dec := gob.NewDecoder(&bb)
	err := dec.Decode(&m)
	if err != nil {
		log.Fatal("Cannot decode Msg! err=", err)
	}
	return
}

func main() {
    m1 := ClientMsg{"id_1"}
	b1 := encode(m1)
	p_bb := bytes.NewBuffer(b1.Bytes())

	var mDecoded ClientMsg // for A, B, C
	//var mDecoded interface{} // D: : cannot use mDecoded (type interface {}) as type Msger in argument to decode: interface {} does not implement Msger (missing IsMsg method)
	decode(mDecoded, *p_bb)     // A: gives: Cannot decode Msg! err=gob: local interface type *main.Msger can only be decoded from remote interface type; received concrete type ClientMsg = struct { Id string; }

	//decode(mDecoded, *p_bb)  // B: gives: cannot use mDecoded (type ClientMsg) as type *Msger in argument to decode: *Msger is pointer to interface, not interface
	//decode(&mDecoded, *p_bb) // C: gives: cannot use &mDecoded (type *ClientMsg) as type *Msger in argument to decode: *Msger is pointer to interface, not interface

	fmt.Printf("m1 Decoded='%v'\n", mDecoded)
}
