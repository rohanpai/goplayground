package main

import (
	&#34;crypto/aes&#34;
	&#34;crypto/cipher&#34;
	&#34;crypto/md5&#34;
	&#34;encoding/base64&#34;
	&#34;fmt&#34;
)

func main(){
	result, _ := DecryptString(`password`, `U2FsdGVkX1&#43;ywYxveBnekSnx6ZP25nyPsWHS3oqcuTo=`)
	fmt.Printf(&#34;Decrypted string is: %s&#34;, result)
}

var openSSLSaltHeader string = &#34;Salted_&#34; // OpenSSL salt is always this string &#43; 8 bytes of actual salt

type OpenSSLCreds struct {
	key []byte
	iv  []byte
}

// Decrypt string that was encrypted using OpenSSL and AES-256-CBC
func DecryptString(passphrase, encryptedBase64String string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedBase64String)
	if err != nil {
		return nil, err
	}
	saltHeader := data[:aes.BlockSize]
	if string(saltHeader[:7]) != openSSLSaltHeader {
		return nil, fmt.Errorf(&#34;Does not appear to have been encrypted with OpenSSL, salt header missing.&#34;)
	}
	salt := saltHeader[8:]
	creds, err := extractOpenSSLCreds([]byte(passphrase), salt)
	if err != nil {
		return nil, err
	}
	return decrypt(creds.key, creds.iv, data)
}

func decrypt(key, iv, data []byte) ([]byte, error) {
	if len(data) == 0 || len(data)%aes.BlockSize != 0 {
		return nil, fmt.Errorf(&#34;bad blocksize(%v), aes.BlockSize = %v\n&#34;, len(data), aes.BlockSize)
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cbc := cipher.NewCBCDecrypter(c, iv)
	cbc.CryptBlocks(data[aes.BlockSize:], data[aes.BlockSize:])
	out, err := pkcs7Unpad(data[aes.BlockSize:], aes.BlockSize)
	if out == nil {
		return nil, err
	}
	return out, nil
}

// openSSLEvpBytesToKey follows the OpenSSL (undocumented?) convention for extracting the key and IV from passphrase.
// It uses the EVP_BytesToKey() method which is basically:
// D_i = HASH^count(D_(i-1) || password || salt) where || denotes concatentaion, until there are sufficient bytes available
// 48 bytes since we&#39;re expecting to handle AES-256, 32bytes for a key and 16bytes for the IV
func extractOpenSSLCreds(password, salt []byte) (OpenSSLCreds, error) {
	m := make([]byte, 48)
	prev := []byte{}
	for i := 0; i &lt; 3; i&#43;&#43; {
		prev = hash(prev, password, salt)
		copy(m[i*16:], prev)
	}
	return OpenSSLCreds{key: m[:32], iv: m[32:]}, nil
}

func hash(prev, password, salt []byte) []byte {
	a := make([]byte, len(prev)&#43;len(password)&#43;len(salt))
	copy(a, prev)
	copy(a[len(prev):], password)
	copy(a[len(prev)&#43;len(password):], salt)
	return md5sum(a)
}

func md5sum(data []byte) []byte {
	h := md5.New()
	h.Write(data)
	return h.Sum(nil)
}

// pkcs7Unpad returns slice of the original data without padding.
func pkcs7Unpad(data []byte, blocklen int) ([]byte, error) {
	if blocklen &lt;= 0 {
		return nil, fmt.Errorf(&#34;invalid blocklen %d&#34;, blocklen)
	}
	if len(data)%blocklen != 0 || len(data) == 0 {
		return nil, fmt.Errorf(&#34;invalid data len %d&#34;, len(data))
	}
	padlen := int(data[len(data)-1])
	if padlen &gt; blocklen || padlen == 0 {
		return nil, fmt.Errorf(&#34;invalid padding&#34;)
	}
	pad := data[len(data)-padlen:]
	for i := 0; i &lt; padlen; i&#43;&#43; {
		if pad[i] != byte(padlen) {
			return nil, fmt.Errorf(&#34;invalid padding&#34;)
		}
	}
	return data[:len(data)-padlen], nil
}