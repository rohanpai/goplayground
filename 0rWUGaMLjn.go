package main

// Caesar Cipher
// (description more or less taken from Wikipedia)
//
//	In cryptography, a Caesar cipher, also known as Caesar&#39;s cipher,
//	the shift cipher, Caesar&#39;s code or Caesar shift, is one of the
//	simplest and most widely known encryption techniques. It is a type
//	of substitution cipher in which each letter in the plaintext is
//	replaced by a letter some fixed number of positions down the
//	alphabet. For example, with a left shift of 3, D would be replaced
//	by A, E would become B, and so on. The method is named after Julius
//	Caesar, who used it in his private correspondence.

// cipher takes in the text to be ciphered along with the direction that
// is being taken; -1 means encoding, &#43;1 means decoding.
func cipher(text string, direction int) string {
	// shift -&gt; number of letters to move to right or left
	// offset -&gt; size of the alphabet, in this case the plain ASCII
	shift, offset := rune(3), rune(26)

	// string-&gt;rune conversion
	runes := []rune(text)

	for index, char := range runes {
		// Iterate over all runes, and perform substitution
		// wherever possible. If the letter is not in the range
		// [1 .. 25], the offset defined above is added or
		// subtracted.
		switch direction {
		case -1: // encoding
			if char &gt;= &#39;a&#39;&#43;shift &amp;&amp; char &lt;= &#39;z&#39; ||
				char &gt;= &#39;A&#39;&#43;shift &amp;&amp; char &lt;= &#39;Z&#39; {
				char = char - shift
			} else if char &gt;= &#39;a&#39; &amp;&amp; char &lt; &#39;a&#39;&#43;shift ||
				char &gt;= &#39;A&#39; &amp;&amp; char &lt; &#39;A&#39;&#43;shift {
				char = char - shift &#43; offset
			}
		case &#43;1: // decoding
			if char &gt;= &#39;a&#39; &amp;&amp; char &lt;= &#39;z&#39;-shift ||
				char &gt;= &#39;A&#39; &amp;&amp; char &lt;= &#39;Z&#39;-shift {
				char = char &#43; shift
			} else if char &gt; &#39;z&#39;-shift &amp;&amp; char &lt;= &#39;z&#39; ||
				char &gt; &#39;Z&#39;-shift &amp;&amp; char &lt;= &#39;Z&#39; {
				char = char &#43; shift - offset
			}
		}

		// Above `if`s handle both upper and lower case ASCII
		// characters; anything else is returned as is (includes
		// numbers, punctuation and space).
		runes[index] = char
	}

	return string(runes)
}

// encode and decode provide the API for encoding and decoding text using
// the Caesar Cipher algorithm.
func encode(text string) string { return cipher(text, -1) }
func decode(text string) string { return cipher(text, &#43;1) }

// A simple test
func main() {
	println(&#34;the text is `das fuchedes 666`&#34;)
	encoded := encode(&#34;das fuchedes 666&#34;)
	println(&#34;  encoded: &#34; &#43; encoded)
	decoded := decode(encoded)
	println(&#34;  decoded: &#34; &#43; decoded)
}
