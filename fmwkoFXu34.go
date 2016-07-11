/*
Package that adds a few extra features to the strings package. 

These functions are very high level and may not be very efficient. Most of the 
functions take inspiration from the Data.List package found in the Haskell programming 
language
*/
package strings_ext

import (
	&#34;strings&#34;
	&#34;unicode/utf8&#34;
)

// Head returns the first rune of s which must be non-empty.
func Head(s string) rune {
	if s == &#34;&#34; {
		panic(&#34;empty list&#34;)
	}
	r, _ := utf8.DecodeRuneInString(s)
	return r
}

// Tail returns the the remainder of s minus the first rune of s, which must be non-empty.
func Tail(s string) string {
	if s == &#34;&#34; {
		panic(&#34;empty list&#34;)
	}

	_, sz := utf8.DecodeRuneInString(s)
	return s[sz:]
}

// Take returns the n rune prefix of s or s itself if n &gt; len([]rune(s))
func Take(n int, s string) string {
	for i := range s {
		if n &lt;= 0 {
			return s[0:i]
		}
		n--
	}
	return s
}

// Drop returns the suffix of s after the first n runes, or &#34;&#34; if n &gt; len([]rune(s))
func Drop(n int, s string) string {
	for i := range s {
		if n &lt;= 0 {
			return s[i:]
		}
		n--
	}
	return &#34;&#34;
}

// TakeWhile, applied to a predicate p and a string s, returns the longest 
// prefix (possibly empty) of s of elements that satisfy p
func TakeWhile(p func(rune) bool, s string) string {
	for i, r := range s {
		if !p(r) {
			return s[0:i]
		}
	}
	return s
}

// DropWhile returns the suffix remaining after TakeWhile
func DropWhile(p func(rune) bool, s string) string {
	for i, r := range s {
		if !p(r) {
			return s[i:]
		}
	}
	return &#34;&#34;
}

// Reverse returns the string s in reverse order
func Reverse(s string) string {
	t := make([]byte, 0, len(s))
	for len(s) &gt; 0 {
		n := 1
		if s[len(s)-1] &gt; 0x7f {
			_, n = utf8.DecodeLastRuneInString(s)
			t = append(t, s[len(s)-n:]...)
		} else {
			t = append(t, s[len(s)-1])
		}
		s = s[0 : len(s)-n]
	}
	return string(t)
}

// Filter, applied to a predicate and a string, returns a string of characters 
// (runes) that satisfy the predicate
func Filter(f func(r rune) bool, s string) string {
	return strings.Map(func(r rune) rune {
		if f(r) {
			return r
		}
		return -1
	}, s)
}

// Span, applied to a predicate p and a string s, returns two strings where the 
// first string is longest prefix (possibly empty) of s of characters (runes) that 
// satisfy p and the second string is the remainder of the string 
func Span(p func(rune) bool, s string) (string, string) {
	return TakeWhile(p, s), DropWhile(p, s)
}

// Group takes a string and returns a slice of strings such 
// that the concatenation of the result is equal to the argument.
// Moreover, each sublist in the result contains only equal elements.
func Group(s string) []string {
	return GroupBy(func(a, b rune) bool { return a == b }, s)
}

// GroupBy is the non-overloaded version of Group.
func GroupBy(p func(rune, rune) bool, s string) []string {
	ss := []string{}
	for len(s) &gt; 0 {
		r0, n := utf8.DecodeRuneInString(s)
		t := TakeWhile(func(r rune) bool {
			return p(r0, r)
		}, s[n:])
		n &#43;= len(t)
		ss = append(ss, s[0:n])
		s = s[n:]
	}
	return ss
}

// Distinct removes duplicate elements from a string. 
// In particular, it keeps only the first occurrence of each element. 
func Distinct(s string) string {
	var ascii [256]bool
	var nonascii map[rune]bool
	return strings.Map(func(r rune) rune {
		if r &lt; 0x80 {
			b := byte(r)
			if ascii[b] {
				return -1
			}
			ascii[b] = true
		} else {
			if nonascii == nil {
				nonascii = make(map[rune]bool)
			}
			if nonascii[r] {
				return -1
			}
			nonascii[r] = true
		}
		return r
	}, s)
}

// Last returns the last rune in a string s, which must be non-empty.
func Last(s string) rune {
	if s == &#34;&#34; {
		panic(&#34;empty list&#34;)
	}

	r, _ := utf8.DecodeLastRuneInString(s)
	return r
}

// Init returns all the elements of s except the last one. The string must 
// be non-empty.
func Init(s string) string {
	if s == &#34;&#34; {
		panic(&#34;empty list&#34;)
	}

	_, sz := utf8.DecodeRuneInString(s)
	c := utf8.RuneCountInString(s)
	return s[:(sz*c)-sz]
}

// IsEmpty tests whether the string s is empty
func IsEmpty(s string) bool {
	return s == &#34;&#34;
}

// All applied to a predicate p and a string s, determines if all elements of
// s satisfy p
func All(p func(rune) bool, s string) bool {
	for _, r := range s {
		if !p(r) {
			return false
		}
	}
	return true
}
