package main

import (
	&#34;crypto/aes&#34;
	&#34;crypto/cipher&#34;
	&#34;encoding/base64&#34;
	&#34;fmt&#34;
)

/**
 *	由于不知道怎么用go来使用PKCS7的填充，所以自己写了一个方法
 *	
 */
func PKCS7Padding(data []byte) []byte {

	// 数据的长度
	dataLen := len(data)

	var bit16 int

	if dataLen%16 == 0 {
		bit16 = dataLen
	} else {
		// 计算补足的位数，填补16的位数，例如 10 = 16, 17 = 32, 33 = 48
		bit16 = int(dataLen/16&#43;1) * 16
	}

	// 需要填充的数量
	paddingNum := bit16 - dataLen

	bitCode := byte(paddingNum)

	padding := make([]byte, paddingNum)
	for i := 0; i &lt; paddingNum; i&#43;&#43; {
		padding[i] = bitCode

	}
	return append(data, padding...)
}

/**
 *	去除PKCS7的补码
 */
func UnPKCS7Padding(data []byte) []byte {
	dataLen := len(data)

	// 在使用PKCS7会以16的倍数减去数据的长度=补位的字节数作为填充的补码，所以现在获取最后一位字节数进行切割
	endIndex := int(data[dataLen-1])

	// 验证结尾字节数是否符合标准，PKCS7的补码字节只会是1-15的字节数
	if 16 &gt; endIndex {

		// 判断结尾的补码是否相同 TODO 不相同也先不管了，暂时不知道怎么处理
		if 1 &lt; endIndex {
			for i := dataLen - endIndex; i &lt; dataLen; i&#43;&#43; {
				if data[dataLen-1] != data[i] {
					fmt.Println(&#34;不同的字节码，尾部字节码:&#34;, data[dataLen-1], &#34;  下标：&#34;, i, &#34;  字节码：&#34;, data[i])
				}
			}
		}

		return data[:dataLen-endIndex]
	}

	fmt.Println(endIndex)

	return nil
}

func main() {

	key := &#34;mofeimofeimofeimofeimofeimofeimo&#34;
	ckey, err := aes.NewCipher([]byte(key))
	if nil != err {
		fmt.Println(&#34;钥匙创建错误:&#34;, err)
	}

	str := []byte(&#34;1234567890&#34;)
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
