package main

import (
	&#34;encoding/json&#34;
	&#34;fmt&#34;
	&#34;io&#34;
)

type Record struct {
	Author Author `json:&#34;author&#34;`
	Title  string `json:&#34;title&#34;`
	URL    string `json:&#34;url&#34;`
}

type Author struct {
	ID    uint64 `json:&#34;id&#34;`
	Email string `json:&#34;email&#34;`
}

type author Author

func (a *Author) UnmarshalJSON(b []byte) (err error) {
	j, s, n := author{}, &#34;&#34;, uint64(0)
	if err = json.Unmarshal(b, &amp;j); err == nil {
		*a = Author(j)
		return
	}
	if err = json.Unmarshal(b, &amp;s); err == nil {
		a.Email = s
		return
	}
	if err = json.Unmarshal(b, &amp;n); err == nil {
		a.ID = n
	}
	return
}

func Decode(r io.Reader) (x *Record, err error) {
	x = new(Record)
	err = json.NewDecoder(r).Decode(x)
	return
}

func main() {
	var r []Record
	fmt.Println(&#34;ERROR:&#34;, json.Unmarshal([]byte(`[{
	  &#34;author&#34;: &#34;attila@attilaolah.eu&#34;,
	  &#34;title&#34;:  &#34;My Blog&#34;,
	  &#34;url&#34;:    &#34;http://attilaolah.eu&#34;
	}, {
	  &#34;author&#34;: 1234567890,
	  &#34;title&#34;:  &#34;Westartup&#34;,
	  &#34;url&#34;:    &#34;http://www.westartup.eu&#34;
	}, {
	  &#34;author&#34;: {
	    &#34;id&#34;:    1234567890,
	    &#34;email&#34;: &#34;nospam@westartup.eu&#34;
	  },
	  &#34;title&#34;:  &#34;Westartup&#34;,
	  &#34;url&#34;:    &#34;http://www.westartup.eu&#34;
	}]`), &amp;r))
	fmt.Println(&#34;RECORDS:&#34;, r)
}
