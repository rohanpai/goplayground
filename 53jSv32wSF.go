//Inspired by 
//http://fendrich.se/blog/2012/10/31/c-plus-plus-11-and-boost-succinct-like-python/
//http://news.ycombinator.com/item?id=4722836 (Jabbles)
package main

import (
	&#34;fmt&#34;
	&#34;log&#34;
	&#34;os&#34;
	&#34;strconv&#34;
)

//getTrailingBytes opens a file and reads the last n bytes
func getTrailingBytes(filename string, n int) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = f.Seek(-int64(n), os.SEEK_END)
	if err != nil {
		return nil, err
	}
	b := make([]byte, n)
	_, err = f.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//Tag represents anything that can produce a list of details
type Tag interface {
	//Parse returns the complete list of all data found in the tag
	Parse() map[string]interface{}
	//String returns the canonical formatted string
	String() string
}

//mp3ID3v1 is a specific kind of tagging
type mp3ID3v1 []byte

//Parse decodes the ID3v1 tag
//According to wikipedia, track number is in here somewhere too
//http://en.wikipedia.org/wiki/ID3#Layout
func (mp3 mp3ID3v1) Parse() map[string]interface{} {
	m := make(map[string]interface{}, 8)
	if string(mp3[:3]) != &#34;TAG&#34; {
		return nil
	}
	m[&#34;title&#34;] = string(mp3[3:33])
	m[&#34;artist&#34;] = string(mp3[33:63])
	m[&#34;album&#34;] = string(mp3[63:93])
	m[&#34;year&#34;], _ = strconv.Atoi(string(mp3[93:97]))
	m[&#34;comment&#34;] = string(mp3[97:126])
	m[&#34;genre&#34;] = int(mp3[127])
	return m
}

//If a particular Tag had additional fields (personal rating?)
//we could provide a different function to display them
func (mp3 mp3ID3v1) String() string {
	return defaultFormat(mp3.Parse())
}

func keyEqualsValue(m map[string]interface{}, s string) string {
	return fmt.Sprintf(&#34;%s=%v\n&#34;, s, m[s])
}

//defaultFormat should display the tag information like the example
func defaultFormat(m map[string]interface{}) (s string) {
	s &#43;= keyEqualsValue(m, &#34;album&#34;)
	s &#43;= keyEqualsValue(m, &#34;artist&#34;)
	s &#43;= keyEqualsValue(m, &#34;title&#34;)
	s &#43;= keyEqualsValue(m, &#34;genre&#34;)
	s &#43;= keyEqualsValue(m, &#34;year&#34;)
	s &#43;= keyEqualsValue(m, &#34;comment&#34;)
	return
}

func main() {
	b, err := getTrailingBytes(&#34;sample.mp3&#34;, 128)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(mp3ID3v1(b))
}
