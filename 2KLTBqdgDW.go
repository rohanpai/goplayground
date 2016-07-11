package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
)

const pw = "p455w0rd"

func main() {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	err = bcrypt.CompareHashAndPassword(hashed, []byte(pw))
	fmt.Println(err)
}
