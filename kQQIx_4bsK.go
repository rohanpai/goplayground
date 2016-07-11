package main

import (
	"fmt"
	"net/url"
)

func main() {

	var Url *url.URL
	Url, err := url.Parse("http://www.example.com")
	if err != nil {
		panic("boom")
	}

	Url.Path += "/some/path/or/other_with_funny_characters?_or_not/"
	parameters := url.Values{}
	parameters.Add("hello", "42")
	parameters.Add("hello", "54")
	parameters.Add("vegetable", "potato")
	Url.RawQuery = parameters.Encode()

	fmt.Printf("Encoded URL is %q\n", Url.String())
}
