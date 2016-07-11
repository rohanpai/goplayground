//Inspired by 
//http://fendrich.se/blog/2012/10/31/c-plus-plus-11-and-boost-succinct-like-python/
//http://news.ycombinator.com/item?id=4722836 (Jabbles)
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
	if string(mp3[:3]) != "TAG" {
		return nil
	}
	m["title"] = string(mp3[3:33])
	m["artist"] = string(mp3[33:63])
	m["album"] = string(mp3[63:93])
	m["year"], _ = strconv.Atoi(string(mp3[93:97]))
	m["comment"] = string(mp3[97:126])
	m["genre"] = int(mp3[127])
	return m
}

//If a particular Tag had additional fields (personal rating?)
//we could provide a different function to display them
func (mp3 mp3ID3v1) String() string {
	return defaultFormat(mp3.Parse())
}

func keyEqualsValue(m map[string]interface{}, s string) string {
	return fmt.Sprintf("%s=%v\n", s, m[s])
}

//defaultFormat should display the tag information like the example
func defaultFormat(m map[string]interface{}) (s string) {
	s += keyEqualsValue(m, "album")
	s += keyEqualsValue(m, "artist")
	s += keyEqualsValue(m, "title")
	s += keyEqualsValue(m, "genre")
	s += keyEqualsValue(m, "year")
	s += keyEqualsValue(m, "comment")
	return
}

func main() {
	b, err := getTrailingBytes("sample.mp3", 128)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(mp3ID3v1(b))
}
