package main

import (
	&#34;crypto/rand&#34;
	&#34;encoding/json&#34;
	&#34;fmt&#34;
	&#34;io&#34;
	&#34;strings&#34;
)

type UUID []byte

func (uuid UUID) String() string {
	if uuid == nil || len(uuid) != 16 {
		return &#34;&#34;
	}
	b := []byte(uuid)
	return fmt.Sprintf(&#34;%08x-%04x-%04x-%04x-%012x&#34;,
		b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func NewRandom() UUID {
	uuid := make([]byte, 16)
	randomBits([]byte(uuid))
	uuid[6] = (uuid[6] &amp; 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] &amp; 0x3f) | 0x80 // Variant is 10
	return uuid
}

func Parse(s string) UUID {
	if len(s) == 36&#43;9 {
		if strings.ToLower(s[:9]) != &#34;urn:uuid:&#34; {
			return nil
		}
		s = s[9:]
	} else if len(s) != 36 {
		return nil
	}
	if s[8] != &#39;-&#39; || s[13] != &#39;-&#39; || s[18] != &#39;-&#39; || s[23] != &#39;-&#39; {
		return nil
	}
	uuid := make([]byte, 16)
	for i, x := range []int{
		0, 2, 4, 6,
		9, 11,
		14, 16,
		19, 21,
		24, 26, 28, 30, 32, 34} {
		if v, ok := xtob(s[x:]); !ok {
			return nil
		} else {
			uuid[i] = v
		}
	}
	return uuid
}

func randomBits(b []byte) {
	if _, err := io.ReadFull(rander, b); err != nil {
		panic(err.Error()) // rand should never fail
	}
}

// xvalues returns the value of a byte as a hexadecimal digit or 255.
var xvalues = []byte{
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 255, 255, 255, 255, 255, 255,
	255, 10, 11, 12, 13, 14, 15, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 10, 11, 12, 13, 14, 15, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
}

// xtob converts the the first two hex bytes of x into a byte.
func xtob(x string) (byte, bool) {
	b1 := xvalues[x[0]]
	b2 := xvalues[x[1]]
	return (b1 &lt;&lt; 4) | b2, b1 != 255 &amp;&amp; b2 != 255
}

var rander = rand.Reader // random function

func (uuid UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.String())
}

func (uuid *UUID) UnmarshalJSON(in []byte) error {
	var str string
	err := json.Unmarshal(in, &amp;str)
	if err != nil {
		return err
	}
	*uuid = (*uuid)[:0]
	id := Parse(str)
	if id != nil {
		*uuid = append(*uuid, id...)
	} else {
		// return an error here
	}
	return nil
}

func main() {
	id := NewRandom()
	fmt.Println(&#34;First ID:&#34;, id)
	b, err := json.Marshal(id)
	if err != nil {
		panic(err)
	}
	fmt.Println(&#34;JSON:&#34;, string(b))
	var newId UUID
	err = json.Unmarshal(b, &amp;newId)
	if err != nil {
		panic(err)
	}
	fmt.Println(&#34;Second ID:&#34;, newId)
}
