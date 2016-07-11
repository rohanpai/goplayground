package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
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
	fmt.Printf("Testing decryption of file data payload.\n")
	fmt.Printf("Salt: %x\n", salt)
	fmt.Printf("PW Code: %x\n", pwcode)
	fmt.Printf("Encrypted contents: %x\n", data)
	fmt.Printf("Auth Code: %x\n", authcode)

	// Generate decryption/auth keys and password verifier from password and salt
	decKey, authKey, pwv := generateKeys([]byte("golang"), salt, aes256KeyLen)

	fmt.Printf("Decryption key: %x\n", decKey)
	fmt.Printf("Auth key: %x\n", authKey)
	fmt.Printf("Password verifier: %x\n", pwv)

	// Check password verification.
	if !checkPasswordVerification(pwv, pwcode) {
		fmt.Printf("Password verification failed.\n")
	}

	// Check MAC authentication code from the payload matches
	// the MAC code generated from auth key and encrypted content.
	if !checkAuthentication(data, authcode, authKey) {
		fmt.Printf("Authentication failed.\n")
	}

	// Generate the IV (or counter?)
	var iv [aes.BlockSize]byte
	iv[0] = 1 // Why is this 1 instead of 0?!?!?!?

	// Get the decryption stream
	decStream := decryptStream(data, decKey, iv[:])

	fmt.Printf("Decrypted Contents: \n")
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
	totalSize := (keySize * 2) + 2 // enc + auth + pv sizes

	key := pbkdf2.Key(password, salt, iterationCount, totalSize, sha1.New)
	fmt.Printf("Master key: %x\n", key)
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
