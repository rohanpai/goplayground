package main

import (
	&#34;crypto/des&#34;
	&#34;fmt&#34;
)

func main() {
	// 鍵の長さは 8 バイト (64 ビット) にしないとエラー
	key := []byte(&#34;priv-key&#34;)
	// cipher.Block を実装している DES 暗号化オブジェクトを生成する
	c, err := des.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// 平文も暗号化されるのは 8 バイト (64 ビット)
	plainText := []byte(&#34;plaintxt&#34;)
	// 暗号化されたバイト列を格納するスライスを用意する
	encrypted := make([]byte, des.BlockSize)
	// DES で暗号化をおこなう
	c.Encrypt(encrypted, plainText)
	// 結果は暗号化されている
	fmt.Println(string(encrypted)) //=&gt; ����A�

	// 復号する
	decrypted := make([]byte, des.BlockSize)
	c.Decrypt(decrypted, encrypted)
	// 結果は元の平文が得られる
	fmt.Println(string(decrypted)) //=&gt; plaintxt
}
