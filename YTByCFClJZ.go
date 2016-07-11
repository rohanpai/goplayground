package main

import (
	&#34;bytes&#34;
	&#34;crypto/aes&#34;
	&#34;crypto/cipher&#34;
	&#34;crypto/hmac&#34;
	&#34;crypto/sha1&#34;
	&#34;fmt&#34;
	&#34;io&#34;
	&#34;os&#34;

	&#34;golang.org/x/crypto/pbkdf2&#34;
)

var (
	// AES-256 encrypted file data payload extracted hello-aes.zip
	raw = []byte{
		// Salt (16 bytes)
		0x09, 0x89, 0xB4, 0x63, 0x06, 0xBD, 0x8F, 0x82,
		0x93, 0xA3, 0x89, 0x61, 0x3D, 0xB8, 0x26, 0xD1,

		// Password Verification (2 bytes)
		0xA3, 0xE6,

		// Encrypted File Data (13 bytes)
		0xE0, 0xBA, 0x87, 0x6C, 0xD1, 0x16, 0xA6, 0xDF,
		0x91, 0xCF, 0x7F, 0x8A, 0x14,

		// Authentication code (10 bytes)
		0xB8, 0x9F, 0x23, 0xE3, 0x99, 0xCB, 0x17, 0xD5,
		0x65, 0x1A,
	} // total 41 bytes

	aes128KeyLen = 16
	aes192KeyLen = 24
	aes256KeyLen = 32

	iterationCount = 1000
)

func main() {
	salt := raw[:16]
	pwcode := raw[16:18]
	data := raw[18:31]
	authcode := raw[31:]

	// Print raw contents
	fmt.Printf(&#34;Testing decryption of file data payload.\n&#34;)
	fmt.Printf(&#34;Salt: %x\n&#34;, salt)
	fmt.Printf(&#34;PW Code: %x\n&#34;, pwcode)
	fmt.Printf(&#34;Encrypted contents: %x\n&#34;, data)
	fmt.Printf(&#34;Auth Code: %x\n&#34;, authcode)

	// Generate decryption/auth keys and password verifier from password and salt
	decKey, authKey, pwv := generateKeys([]byte(&#34;golang&#34;), salt, aes256KeyLen)

	fmt.Printf(&#34;Decryption key: %x\n&#34;, decKey)
	fmt.Printf(&#34;Auth key: %x\n&#34;, authKey)
	fmt.Printf(&#34;Password verifier: %x\n&#34;, pwv)

	// Check password verification.
	if !checkPasswordVerification(pwv, pwcode) {
		fmt.Printf(&#34;Password verification failed.\n&#34;)
	}

	// Check MAC authentication code from the payload matches
	// the MAC code generated from auth key and encrypted content.
	if !checkAuthentication(data, authcode, authKey) {
		fmt.Printf(&#34;Authentication failed.\n&#34;)
	}

	// Generate the IV (or counter?)
	var iv [aes.BlockSize]byte
	iv[0] = 1 // Why is this 1 instead of 0?!?!?!?

	// Get the decryption stream
	decStream := decryptStream(data, decKey, iv[:])

	fmt.Printf(&#34;Decrypted Contents: \n&#34;)
	io.Copy(os.Stdout, decStream)
}

// checks if the code from the payload matches the password verification
// pwv generated from the password and salt
func checkPasswordVerification(pwv, code []byte) bool {
	return bytes.Equal(pwv, code)
}

func checkAuthentication(message, authcode, key []byte) bool {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedAuthCode := mac.Sum(nil)
	// Truncate at the first 10 bytes
	expectedAuthCode = expectedAuthCode[:10]
	return bytes.Equal(expectedAuthCode, authcode)
}

func generateKeys(password, salt []byte, keySize int) (encKey, authKey, pwv []byte) {
	totalSize := (keySize * 2) &#43; 2 // enc &#43; auth &#43; pv sizes

	key := pbkdf2.Key(password, salt, iterationCount, totalSize, sha1.New)
	fmt.Printf(&#34;Master key: %x\n&#34;, key)
	encKey = key[:keySize]
	authKey = key[keySize : keySize*2]
	pwv = key[keySize*2:]
	return
}

func decryptStream(ciphertext, key, iv []byte) io.Reader {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	stream := cipher.NewCTR(block, iv)
	reader := cipher.StreamReader{S: stream, R: bytes.NewReader(ciphertext)}
	return reader
}
