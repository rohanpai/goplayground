package main

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"reflect"
	"strconv"

	"github.com/ugorji/go/codec"
)

type Item struct {
	Ref         Ref    `json:"ref"`
	Size        uint32 `json:"size"`
	Compression uint8  `json:"compr,omitempty"`
	OSize       uint32 `json:"osize,omitempty"`
}

type Ref struct {
	Type      uint8
	HashBytes []byte
}

func main() {
	var (
		bh codec.BincHandle
		mh codec.MsgpackHandle
	)
	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))

	v := Item{
		Ref:         Ref{Type: 0, HashBytes: []byte("12345678901234567890")},
		Size:        3555,
		Compression: 1,
		OSize:       7643,
	}

	b := make([]byte, 0, 64)
	enc := codec.NewEncoderBytes(&b, &mh)
	err := enc.Encode(v)
	if err != nil {
		log.Printf("error encoding %v to MessagePack: %v", v, err)
	}
	log.Printf("length of %#v as MessagePack: %d\n%v", v, len(b), b)

	b = b[:0]
	enc = codec.NewEncoderBytes(&b, &bh)
	if err = enc.Encode(v); err != nil {
		log.Printf("error encoding %v to Binc: %v", v, err)
	}
	log.Printf("length of %#v as Binc: %d\n%v", v, len(b), b)

	if b, err = json.Marshal(v); err != nil {
		log.Printf("error encoding to json: %v", err)
	}
	log.Printf("length of %#v as JSON: %d\n%s", v, len(b), b)

	b = b[:40]
	hex.Encode(b, v.Ref.HashBytes)
	b = append(append([]byte("[sha1-"), b...),
		[]byte(strconv.FormatUint(uint64(v.Size), 10)+" 1 "+
			strconv.FormatUint(uint64(v.OSize), 10)+"]")...)
	log.Printf("length of %#v as now: %d\n%s", v, len(b), b)
}
