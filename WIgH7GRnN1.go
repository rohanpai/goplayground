package main

import (
	&#34;math/rand&#34;
	&#34;testing&#34;
	&#34;time&#34;
)

// Implementations

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune(&#34;abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ&#34;)

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

const letterBytes = &#34;abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ&#34;
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1&lt;&lt;letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func RandStringBytesMask(n int) string {
	b := make([]byte, n)
	for i := 0; i &lt; n; {
		if idx := int(rand.Int63() &amp; letterIdxMask); idx &lt; len(letterBytes) {
			b[i] = letterBytes[idx]
			i&#43;&#43;
		}
	}
	return string(b)
}

func RandStringBytesMaskImpr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i &gt;= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache &amp; letterIdxMask); idx &lt; len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache &gt;&gt;= letterIdxBits
		remain--
	}

	return string(b)
}

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i &gt;= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache &amp; letterIdxMask); idx &lt; len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache &gt;&gt;= letterIdxBits
		remain--
	}

	return string(b)
}

// Benchmark functions

const n = 16

func BenchmarkRunes(b *testing.B) {
	for i := 0; i &lt; b.N; i&#43;&#43; {
		RandStringRunes(n)
	}
}

func BenchmarkBytes(b *testing.B) {
	for i := 0; i &lt; b.N; i&#43;&#43; {
		RandStringBytes(n)
	}
}

func BenchmarkBytesRmndr(b *testing.B) {
	for i := 0; i &lt; b.N; i&#43;&#43; {
		RandStringBytesRmndr(n)
	}
}

func BenchmarkBytesMask(b *testing.B) {
	for i := 0; i &lt; b.N; i&#43;&#43; {
		RandStringBytesMask(n)
	}
}

func BenchmarkBytesMaskImpr(b *testing.B) {
	for i := 0; i &lt; b.N; i&#43;&#43; {
		RandStringBytesMaskImpr(n)
	}
}

func BenchmarkBytesMaskImprSrc(b *testing.B) {
	for i := 0; i &lt; b.N; i&#43;&#43; {
		RandStringBytesMaskImprSrc(n)
	}
}
