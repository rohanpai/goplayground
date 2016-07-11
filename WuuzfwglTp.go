package main

import (
	&#34;fmt&#34;
	&#34;github.com/daswasser/validate&#34;
	&#34;github.com/daswasser/validate/web&#34;
)

func main() {
	// Setup a new validator
	v := validate.NewValidator()

	// Create a new Domain object and return the message on failure
	goodDomain :=
		web.NewDomain(&#34;www.golang.org&#34;).
			MaxSubdomains(2).
			SetMessage(&#34;Invalid domain specified!&#34;)

	// Validate the good domain
	err := v.Validate(goodDomain)
	if err != nil {
		fmt.Printf(&#34;%s error:\n&#34;, goodDomain)
		fmt.Println(err)
		fmt.Println(goodDomain.Message())
	} else {
		fmt.Printf(&#34;%s is a valid domain\n&#34;, goodDomain)
	}

}
