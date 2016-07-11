//Decrypts obfuscated passwords by Remmina - The GTK&#43; Remote Desktop Client
//written by Michael Cochez
package main

import (
	&#34;crypto/cipher&#34;
	&#34;crypto/des&#34;
	&#34;encoding/base64&#34;
	&#34;fmt&#34;
	&#34;log&#34;
)

//set the variables here

var base64secret = &#34;yoursecret&#34;
var base64password = &#34;theconnectionpassword&#34;

//The secret is used for encrypting the passwords. This can typically be found from ~/.remmina/remmina.pref on the line containing &#39;secret=&#39;.
//&#34;The encrypted password used for the connection. This can typically be found from /.remmina/dddddddddddd.remmina &#34; on the line containing &#39;password=&#39;.
//Copy everything after the &#39;=&#39; sign. Also include final &#39;=&#39; signs if they happen to be there.

//returns a function which can be used for decrypting passwords
func makeRemminaDecrypter(base64secret string) func(string) string {
	//decode the secret
	secret, err := base64.StdEncoding.DecodeString(base64secret)
	if err != nil {
		log.Fatal(&#34;Base 64 decoding failed:&#34;, err)
	}
	if len(secret) != 32 {
		log.Fatal(&#34;the secret is not 32 bytes long&#34;)
	}
	//the key is the 24 first bits of the secret
	key := secret[:24]
	//3DES cipher
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		log.Fatal(&#34;Failed creating the 3Des cipher block&#34;, err)
	}
	//the rest of the secret is the iv
	iv := secret[24:]
	decrypter := cipher.NewCBCDecrypter(block, iv)

	return func(encodedEncryptedPassword string) string {
		encryptedPassword, err := base64.StdEncoding.DecodeString(encodedEncryptedPassword)
		if err != nil {
			log.Fatal(&#34;Base 64 decoding failed:&#34;, err)
		}
		//in place decryption
		decrypter.CryptBlocks(encryptedPassword, encryptedPassword)
		return string(encryptedPassword)
	}
}

func main() {

	if base64secret == &#34;yoursecret&#34; || base64password == &#34;theconnectionpassword&#34; {

		log.Fatal(&#34;both base64secret and base64password variables must be set&#34;)
	}

	decrypter := makeRemminaDecrypter(base64secret)

	fmt.Printf(&#34;Passwd : %v\n&#34;, decrypter(base64password))

}