//Decrypts obfuscated passwords by Remmina - The GTK+ Remote Desktop Client
//written by Michael Cochez
package main

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"fmt"
	"log"
)

//set the variables here

var base64secret = "yoursecret"
var base64password = "theconnectionpassword"

//The secret is used for encrypting the passwords. This can typically be found from ~/.remmina/remmina.pref on the line containing 'secret='.
//"The encrypted password used for the connection. This can typically be found from /.remmina/dddddddddddd.remmina " on the line containing 'password='.
//Copy everything after the '=' sign. Also include final '=' signs if they happen to be there.

//returns a function which can be used for decrypting passwords
func makeRemminaDecrypter(base64secret string) func(string) string {
	//decode the secret
	secret, err := base64.StdEncoding.DecodeString(base64secret)
	if err != nil {
		log.Fatal("Base 64 decoding failed:", err)
	}
	if len(secret) != 32 {
		log.Fatal("the secret is not 32 bytes long")
	}
	//the key is the 24 first bits of the secret
	key := secret[:24]
	//3DES cipher
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		log.Fatal("Failed creating the 3Des cipher block", err)
	}
	//the rest of the secret is the iv
	iv := secret[24:]
	decrypter := cipher.NewCBCDecrypter(block, iv)

	return func(encodedEncryptedPassword string) string {
		encryptedPassword, err := base64.StdEncoding.DecodeString(encodedEncryptedPassword)
		if err != nil {
			log.Fatal("Base 64 decoding failed:", err)
		}
		//in place decryption
		decrypter.CryptBlocks(encryptedPassword, encryptedPassword)
		return string(encryptedPassword)
	}
}

func main() {

	if base64secret == "yoursecret" || base64password == "theconnectionpassword" {

		log.Fatal("both base64secret and base64password variables must be set")
	}

	decrypter := makeRemminaDecrypter(base64secret)

	fmt.Printf("Passwd : %v\n", decrypter(base64password))

}