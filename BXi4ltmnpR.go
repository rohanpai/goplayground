// You can edit this code!
// Click here and start typing.
package main

import (
	&#34;bytes&#34;
	&#34;compress/gzip&#34;
	&#34;encoding/hex&#34;
	&#34;fmt&#34;
)

const (
	original = `14f88cbb028f050272780c3974c855650228fcfff23f4c8a24c92581a5845d40c585ca88c4fac29d43c24090a0c028cd5060524af003428a40c0c0e05060`
)

func main() {
	// ヒントの even -&gt; odd はそれぞれ偶数、奇数の意味。
	// 偶数番目と奇数番目の文字のみを取り出す。
	var even, odd []rune
	for i, c := range original {
		if i%2 == 0 { // 偶数
			even = append(even, c)
		} else { // 奇数
			odd = append(odd, c)
		}
	}

	fmt.Printf(&#34;偶数 = %s\n&#34;, string(even))
	fmt.Printf(&#34;奇数 = %s\n&#34;, string(odd))

	s := string(even) &#43; string(odd)
	fmt.Printf(&#34;偶数 -&gt; 奇数 = %s\n&#34;, s)

	// 1f8b から始まるのは、前もやったけど gzip 形式です。
	// 16進数をバイナリ表現に直して、gzip を解凍します。

	binary, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	fmt.Printf(&#34;バイナリ表現 = %v\n&#34;, binary)

	r, err := gzip.NewReader(bytes.NewReader(binary))
	if err != nil {
		panic(err)
	}
	defer r.Close()

	result := make([]byte, 1000)
	_, err = r.Read(result)
	if err != nil {
		panic(err)
	}

	fmt.Printf(&#34;答 = %s\n&#34;, string(result))
}
