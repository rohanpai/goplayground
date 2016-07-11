package main

import (
	&#34;bytes&#34;
	&#34;encoding/binary&#34;
	&#34;fmt&#34;
	&#34;io&#34;
	&#34;testing/quick&#34;
	&#34;unicode/utf16&#34;
	//&#34;unicode/utf8&#34;
)

type String string

func (m *String) ReadFrom(r io.Reader) (n int64, err error) {
	var length uint16

	// Read length of string, 2 bytes
	err = binary.Read(r, binary.BigEndian, &amp;length)
	if err != nil {
		return 0, err
	}

	// Read string, 2 bytes * length
	var contents = make([]uint16, length)
	err = binary.Read(r, binary.BigEndian, &amp;contents)
	if err != nil {
		return 2, err
	}

	// Update String contents
	*m = String(readString(contents))
	//*m = String(utf16.Decode(contents))

	return int64(2 &#43; len(contents)*2), nil
}

func (m *String) WriteTo(w io.Writer) (n int64, err error) {
	var contents = writeString(string(*m))
	var length = uint16(len(contents))

	// Read size of string inside a short
	err = binary.Write(w, binary.BigEndian, length)
	if err != nil {
		return 0, err
	}

	err = binary.Write(w, binary.BigEndian, contents)
	if err != nil {
		return 2, err
	}

	return int64(2 &#43; length*2), nil
}

func (m String) String() string {
	return string(m)
}

// writeString encodes a Go string to UCS-2 (UTF-16) encoded string ([]uint16).
func writeString(s string) []uint16 {
	runes := []rune(s)
	fmt.Printf(&#34;w: %U\n&#34;, runes)
	return utf16.Encode(runes)
}

// readString decodes a UCS-2 (UTF-16) encoded string ([]uint16) into a Go string.
func readString(u16 []uint16) string {
	runes := utf16.Decode(u16)
	fmt.Printf(&#34;r: %U\n&#34;, runes)
	return string(runes)
}

func main() {
	var buf bytes.Buffer

	f := func(v string) bool {
		var rs []rune
		for _, r := range v {
			switch {
			case r &gt;= 0xD800 &amp;&amp; r &lt;= 0xDBFF:
				continue
			case r &gt;= 0xDC00 &amp;&amp; r &lt;= 0xDFFF:
				continue
			}
			rs = append(rs, r)
		}
		v = string(rs)

		var v1 = String(rs)
		v1.WriteTo(&amp;buf)
		var v2 String
		v2.ReadFrom(&amp;buf)
		return v == string(v2)

		//u16 := utf16.Encode([]rune(v))
		//v = string(utf16.Decode(u16))
		//return utf8.ValidString(v)
	}

	if err := quick.Check(f, nil); err != nil {
		fmt.Println(&#34;Error:&#34;, err)
	}
}
