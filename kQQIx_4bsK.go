package main

import (
	&#34;fmt&#34;
	&#34;net/url&#34;
)

func main() {

	var Url *url.URL
	Url, err := url.Parse(&#34;http://www.example.com&#34;)
	if err != nil {
		panic(&#34;boom&#34;)
	}

	Url.Path &#43;= &#34;/some/path/or/other_with_funny_characters?_or_not/&#34;
	parameters := url.Values{}
	parameters.Add(&#34;hello&#34;, &#34;42&#34;)
	parameters.Add(&#34;hello&#34;, &#34;54&#34;)
	parameters.Add(&#34;vegetable&#34;, &#34;potato&#34;)
	Url.RawQuery = parameters.Encode()

	fmt.Printf(&#34;Encoded URL is %q\n&#34;, Url.String())
}
