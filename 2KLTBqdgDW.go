package main

import (
	&#34;code.google.com/p/go.crypto/bcrypt&#34;
	&#34;fmt&#34;
)

const pw = &#34;p455w0rd&#34;

func main() {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	err = bcrypt.CompareHashAndPassword(hashed, []byte(pw))
	fmt.Println(err)
}
