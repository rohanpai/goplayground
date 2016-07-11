package main

import (
	&#34;bytes&#34;
	&#34;fmt&#34;
)

// Hamming returns the normalized similarity value.
// hamming distance is the number of differing &#34;bits&#34;.
// hamming distance is minimum number of substitutions
// required to change one string into the other
// (https://en.wikipedia.org/wiki/Hamming_distance)
func Hamming(txt1, txt2 []byte) float64 {
	switch bytes.Compare(txt1, txt2) {
	case 0: // txt1 == txt2
	case 1: // txt1 &gt; txt2
		temp := make([]byte, len(txt1))
		copy(temp, txt2)
		txt2 = temp
	case -1: // txt1 &lt; txt2
		temp := make([]byte, len(txt2))
		copy(temp, txt1)
		txt1 = temp
	}
	if len(txt1) != len(txt2) {
		panic(&#34;Undefined for sequences of unequal length&#34;)
	}
	count := 0
	for idx, b1 := range txt1 {
		b2 := txt2[idx]
		xor := b1 ^ b2 // 1 if bits are different
		//
		// bit count (number of 1)
		// http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetNaive
		//
		// repeat shifting from left to right (divide by 2)
		// until all bits are zero
		for x := xor; x &gt; 0; x &gt;&gt;= 1 {
			// check if lowest bit is 1
			if int(x&amp;1) == 1 {
				count&#43;&#43;
			}
		}
	}
	if count == 0 {
		// similarity is 1 for equal texts.
		return 1
	}
	return float64(1) / float64(count)
}

func main() {
	fmt.Println(Hamming([]byte(&#34;A&#34;), []byte(&#34;A&#34;)))             // 1
	fmt.Println(Hamming([]byte(&#34;A&#34;), []byte(&#34;a&#34;)))             // 1
	fmt.Println(Hamming([]byte(&#34;a&#34;), []byte(&#34;A&#34;)))             // 1
	fmt.Println(Hamming([]byte(&#34;aaa&#34;), []byte(&#34;aba&#34;)))         // 0.5
	fmt.Println(Hamming([]byte(&#34;aaa&#34;), []byte(&#34;aBa&#34;)))         // 0.333
	fmt.Println(Hamming([]byte(&#34;aaa&#34;), []byte(&#34;a a&#34;)))         // 0.5
	fmt.Println(Hamming([]byte(&#34;karolin&#34;), []byte(&#34;kathrin&#34;))) // 0.1111111111111111

	fmt.Println(Hamming([]byte(&#34;Hello&#34;), []byte(&#34;Hello&#34;)))
	// 1

	fmt.Println(Hamming([]byte(&#34;Hello&#34;), []byte(&#34;Hel lo&#34;)))
	// 0.2

	fmt.Println(Hamming([]byte(&#34;&#34;), []byte(&#34;Hello&#34;)))
	// 0.05

	fmt.Println(Hamming([]byte(&#34;hello&#34;), []byte(&#34;Hello&#34;)))
	// 1

	fmt.Println(Hamming([]byte(&#34;abc&#34;), []byte(&#34;bcd&#34;)))
	// 0.16666666666666666
}
