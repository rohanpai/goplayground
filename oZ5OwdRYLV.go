package main

import (
	&#34;crypto/aes&#34;
	&#34;crypto/cipher&#34;
	&#34;fmt&#34;
	&#34;errors&#34;
	&#34;bytes&#34;
)

var (
	ErrPaddingSize = errors.New(&#34;padding size error&#34;)
)
var (
	PKCS5 = &amp;pkcs5{}
)
var (
	// difference with pkcs5 only block must be 8
	PKCS7 = &amp;pkcs5{}
)

// pkcs5Padding is a pkcs5 padding struct.
type pkcs5 struct{}

// Padding implements the Padding interface Padding method.
func (p *pkcs5) Padding(src []byte, blockSize int) []byte {
	srcLen := len(src)
	padLen := blockSize - (srcLen % blockSize)
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(src, padText...)
}

// Unpadding implements the Padding interface Unpadding method.
func (p *pkcs5) Unpadding(src []byte, blockSize int) ([]byte, error) {
	srcLen := len(src)
	paddingLen := int(src[srcLen-1])
	if paddingLen &gt;= srcLen || paddingLen &gt; blockSize {
		return nil, ErrPaddingSize
	}
	return src[:srcLen-paddingLen], nil
}

func main() {
	var block cipher.Block
	var originalData, encryptedData, decryptedData []byte
	var err error
	var ebm, dbm cipher.BlockMode

	key := []byte{231, 165, 119, 133, 0, 233, 67, 180, 174, 205, 132, 250, 92, 63, 130, 166}
	iv := []byte{233, 211, 143, 12, 117, 249, 61, 68, 19, 180, 109, 110, 33, 104, 244, 179}
	if block, err = aes.NewCipher(key); err != nil {
		fmt.Printf(&#34;aes.NewCipher() error(%v)&#34;, err)
	}
	ebm = cipher.NewCBCEncrypter(block, iv)
	dbm = cipher.NewCBCDecrypter(block, iv)

	originalData = []byte(&#34;just a test string&#34;)
	if encryptedData, err = Encrypt(ebm, originalData); err != nil {
		fmt.Printf(&#34;encrypt error(%v)&#34;, err)
	}

	for i := 0; i &lt; 5; i&#43;&#43; {
		// dbm = cipher.NewCBCDecrypter(block, iv)
		tmp := make([]byte, len(encryptedData))
		copy(tmp, encryptedData)
		decryptedData, err = Decrypt(dbm, tmp)
		fmt.Println(decryptedData)
	}

}

func Encrypt(encryptor cipher.BlockMode, msg []byte) (cipherText []byte, err error) {
	if msg != nil {
		// let caller do pkcs7 padding
		msg =PKCS7.Padding(msg, encryptor.BlockSize())
		if len(msg) &lt; encryptor.BlockSize() || len(msg)%encryptor.BlockSize() != 0 {
			fmt.Println(&#34;length error&#34;)
			return
		}
		cipherText = msg
		encryptor.CryptBlocks(cipherText, msg)
	}
	return
}

func Decrypt(decryptor cipher.BlockMode, cipherText []byte) (msg []byte, err error) {
	if decryptor != nil {
		if len(cipherText) &lt; decryptor.BlockSize() || len(cipherText)%decryptor.BlockSize() != 0 {
			fmt.Println(&#34;length error&#34;)
			return
		}
		msg = cipherText
		decryptor.CryptBlocks(msg, cipherText)
		// let caller do pkcs7 unpadding
		msg, err =PKCS7.Unpadding(msg, decryptor.BlockSize())
	}
	return
}