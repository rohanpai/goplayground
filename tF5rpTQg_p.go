package main

import (
	&#34;encoding/hex&#34;
	&#34;encoding/json&#34;
	&#34;log&#34;
	&#34;reflect&#34;
	&#34;strconv&#34;

	&#34;github.com/ugorji/go/codec&#34;
)

type Item struct {
	Ref         Ref    `json:&#34;ref&#34;`
	Size        uint32 `json:&#34;size&#34;`
	Compression uint8  `json:&#34;compr,omitempty&#34;`
	OSize       uint32 `json:&#34;osize,omitempty&#34;`
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
		Ref:         Ref{Type: 0, HashBytes: []byte(&#34;12345678901234567890&#34;)},
		Size:        3555,
		Compression: 1,
		OSize:       7643,
	}

	b := make([]byte, 0, 64)
	enc := codec.NewEncoderBytes(&amp;b, &amp;mh)
	err := enc.Encode(v)
	if err != nil {
		log.Printf(&#34;error encoding %v to MessagePack: %v&#34;, v, err)
	}
	log.Printf(&#34;length of %#v as MessagePack: %d\n%v&#34;, v, len(b), b)

	b = b[:0]
	enc = codec.NewEncoderBytes(&amp;b, &amp;bh)
	if err = enc.Encode(v); err != nil {
		log.Printf(&#34;error encoding %v to Binc: %v&#34;, v, err)
	}
	log.Printf(&#34;length of %#v as Binc: %d\n%v&#34;, v, len(b), b)

	if b, err = json.Marshal(v); err != nil {
		log.Printf(&#34;error encoding to json: %v&#34;, err)
	}
	log.Printf(&#34;length of %#v as JSON: %d\n%s&#34;, v, len(b), b)

	b = b[:40]
	hex.Encode(b, v.Ref.HashBytes)
	b = append(append([]byte(&#34;[sha1-&#34;), b...),
		[]byte(strconv.FormatUint(uint64(v.Size), 10)&#43;&#34; 1 &#34;&#43;
			strconv.FormatUint(uint64(v.OSize), 10)&#43;&#34;]&#34;)...)
	log.Printf(&#34;length of %#v as now: %d\n%s&#34;, v, len(b), b)
}
