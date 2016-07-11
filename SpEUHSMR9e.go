package main

import (
	&#34;crypto/aes&#34;
	&#34;crypto/cipher&#34;
	&#34;encoding/base64&#34;
	&#34;bytes&#34;
	&#34;fmt&#34;
)

/**
 *	PKCS7补码
 *	这里可以参考下http://blog.studygolang.com/167.html
 */
func PKCS7Padding(data []byte) []byte {
	blockSize := 16
	padding := blockSize - len(data)%blockSize
    	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    	return append(data, padtext...)

}

/**
 *	去除PKCS7的补码
 */
func UnPKCS7Padding(data []byte) []byte {
	length := len(data)
     	// 去掉最后一个字节 unpadding 次
     	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

func main() {

	key := &#34;mofeimofeimofeimofeimofeimofeimo&#34;
	ckey, err := aes.NewCipher([]byte(key))
	if nil != err {
		fmt.Println(&#34;钥匙创建错误:&#34;, err)
	}

	str := []byte(`polaris@studygo`)
	str = append(str, 0x02)
	iv := []byte(&#34;1234567890123456&#34;)
	fmt.Println(&#34;加密的字符串&#34;, string(str), &#34;\n加密钥匙&#34;, key, &#34;\n向量IV&#34;, string(iv))

	fmt.Println(&#34;加密前的字节：&#34;, str, &#34;\n&#34;)

	encrypter := cipher.NewCBCEncrypter(ckey, iv)

	// PKCS7补码
	str = PKCS7Padding(str)
	out := make([]byte, len(str))

	encrypter.CryptBlocks(out, str)
	fmt.Println(&#34;加密后字节：&#34;, out)

	base64Out := base64.URLEncoding.EncodeToString(out)
	fmt.Println(&#34;Base64后：&#34;, base64Out)

	fmt.Println(&#34;\n开始解码&#34;)
	decrypter := cipher.NewCBCDecrypter(ckey, iv)
	base64In, _ := base64.URLEncoding.DecodeString(base64Out)
	in := make([]byte, len(base64In))
	decrypter.CryptBlocks(in, base64In)

	fmt.Println(&#34;解密后的字节：&#34;, in)

	// 去除PKCS7补码
	in = UnPKCS7Padding(in)

	fmt.Println(&#34;去PKCS7补码：&#34;, in)
	fmt.Println(&#34;解密：&#34;, string(in))
}