package main

import (
	"fmt"
	"github.com/daswasser/validate"
	"github.com/daswasser/validate/web"
)

func main() {
	// Setup a new validator
	v := validate.NewValidator()

	// Create a new Domain object and return the message on failure
	goodDomain :=
		web.NewDomain("www.golang.org").
			MaxSubdomains(2).
			SetMessage("Invalid domain specified!")

	// Validate the good domain
	err := v.Validate(goodDomain)
	if err != nil {
		fmt.Printf("%s error:\n", goodDomain)
		fmt.Println(err)
		fmt.Println(goodDomain.Message())
	} else {
		fmt.Printf("%s is a valid domain\n", goodDomain)
	}

}
