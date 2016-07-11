/*
Write a function that takes two equal-length buffers
and produce their `XOR` combination.

Pass this string

1c0111001f010100061a024b53535009181c

and decode the string. And `XOR` against this string

686974207468652062756c6c277320657965

and it should return

746865206b696420646f6e277420706c6179
*/
package main

import (
	&#34;bytes&#34;
	&#34;encoding/hex&#34;
	&#34;fmt&#34;
	&#34;log&#34;
	&#34;strconv&#34;
)

const (
	input1 string = &#34;1c0111001f010100061a024b53535009181c&#34;
	input2 string = &#34;686974207468652062756c6c277320657965&#34;
	output string = &#34;746865206b696420646f6e277420706c6179&#34;
)

func main() {
	func() {
		fmt.Println(&#34;XOR:  x ^ y&#34;)
		x := toUint64(&#34;0101&#34;)
		y := toUint64(&#34;0011&#34;)
		z := x ^ y
		fmt.Printf(&#34;%10b (decimal %d)\n&#34;, z, z)
		/*
		   XOR:  x ^ y
		         0101 (decimal 5)
		         0011 (decimal 3)
		          110 (decimal 6)
		*/
		fmt.Println()
	}()

	decodedHexInput1, err := hex.DecodeString(input1)
	if err != nil {
		panic(err)
	}
	fmt.Println(&#34;decodedHexInput1:&#34;, string(decodedHexInput1), len(decodedHexInput1))
	// decodedHexInput1: KSSP	 18

	decodedHexInput2, err := hex.DecodeString(input2)
	if err != nil {
		panic(err)
	}
	fmt.Println(&#34;decodedHexInput2:&#34;, string(decodedHexInput2), len(decodedHexInput2))
	// decodedHexInput2: hit the bull&#39;s eye 18

	decodedHexOutput, err := hex.DecodeString(output)
	if err != nil {
		panic(err)
	}
	fmt.Println(&#34;decodedHexOutput:&#34;, string(decodedHexOutput), len(decodedHexOutput))
	// decodedHexOutput: the kid don&#39;t play 18

	resultBytes := make([]byte, len(decodedHexOutput))
	for i := 0; i &lt; len(decodedHexOutput); i&#43;&#43; {
		resultBytes[i] = decodedHexInput1[i] ^ decodedHexInput2[i]
	}
	if !bytes.Equal(resultBytes, decodedHexOutput) {
		log.Fatalf(&#34;%s %s&#34;, resultBytes, decodedHexOutput)
	}
	fmt.Println(&#34;resultBytes:&#34;, string(resultBytes))
	// resultBytes: the kid don&#39;t play
}

func toUint64(bstr string) uint64 {
	var num uint64
	if i, err := strconv.ParseUint(bstr, 2, 64); err != nil {
		panic(err)
	} else {
		num = i
	}
	fmt.Printf(&#34;%10s (decimal %d)\n&#34;, bstr, num)
	return num
}
