package main

import (
	&#34;bytes&#34;
	&#34;fmt&#34;
	&#34;regexp&#34;
	&#34;strings&#34;
	&#34;unicode&#34;
)

func main() {
	fmt.Println(SwapCase(&#34;Hello Hi&#34;))
	// hELLO hI

	fmt.Println(SwapCaseII(&#34;Hello Hi&#34;))
	// hELLO hI

	fmt.Println(ReverseString(&#34;Hello Hi&#34;))
	// iH olleH

	fmt.Println(CheckPalindrome(&#34;Anne, I vote more cars race Rome-to-Vienna&#34;))
	// true
}

// Convert character cases.
func SwapCase(str string) string {
	b := new(bytes.Buffer)

	// traverse character values, without index
	for _, elem := range str {
		if unicode.IsUpper(elem) {
			b.WriteRune(unicode.ToLower(elem))
		} else {
			b.WriteRune(unicode.ToUpper(elem))
		}
	}

	return b.String()
}

func SwapCaseII(str string) string {
	return strings.Map(SwapRune, str)
}

// rune is variable-length and can be made up of one or more bytes.
// rune literals are mapped to their unicode codepoint.
// For example, a rune literal &#39;a&#39; is a number 97.
// 32 is the offset of the uppercase and lowercase characters.
// So if you add 32 to &#39;A&#39;, you get &#39;a&#39; and vice versa.
func SwapRune(r rune) rune {
	switch {
	case &#39;a&#39; &lt;= r &amp;&amp; r &lt;= &#39;z&#39;:
		return r - &#39;a&#39; &#43; &#39;A&#39;
	case &#39;A&#39; &lt;= r &amp;&amp; r &lt;= &#39;Z&#39;:
		return r - &#39;A&#39; &#43; &#39;a&#39;
	default:
		return r
	}
}

// ReverseString changes the order of the input string.
// It returns the new version of input string.
// We need to use rune because string is immutable.
func ReverseString(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i &lt; j; i, j = i&#43;1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// CheckPalindrome returns true the input string is palindrome.
func CheckPalindrome(str string) bool {
	str = CleanUp(str)
	str = ReplaceNonAlnumWithSpace(str)
	str = strings.ToLower(str)

	var validID = regexp.MustCompile(`\s{1,}`)
	str = validID.ReplaceAllString(str, &#34;&#34;)

	rvs := ReverseString(str)

	if str == rvs {
		return true
	} else {
		return false
	}
}

// CleanUp cleans up unnecessary characters in string.
// It cleans up the blank characters that carry no meaning in context
// , converts all whitespaces into single whitespace.
// String is immutable, which means the original string would not change.
func CleanUp(str string) string {

	// validID := regexp.MustCompile(`\s{2,}`)
	// func TrimSpace(s string) string
	// slicing off all &#34;leading&#34; and
	// &#34;trailing&#34; white space, as defined by Unicode.
	str = strings.TrimSpace(str)

	// func Fields(s string) []string
	// Fields splits the slice s around each instance
	// of &#34;one or more consecutive white space&#34;
	slice := strings.Fields(str)

	// now join them with a single white space character
	return strings.Join(slice, &#34; &#34;)
}

// ReplaceNonAlnumWithSpace removes all alphanumeric characters.
// It replaces them with a single whitespace character.
// It returns the new version of input string, in lower case.
func ReplaceNonAlnumWithSpace(str string) string {
	str = ExpandApostrophe(str)
	// alphanumeric (== [0-9A-Za-z])
	// \s is a white space character
	validID := regexp.MustCompile(`[^[:alnum:]\s]`)
	return validID.ReplaceAllString(str, &#34; &#34;)
}

// ExpandApostrophe expands the apostrophe phrases.
// And convert them to lower case letters.
func ExpandApostrophe(str string) string {
	// assignment between string is not &#34;copy&#34;
	// even if str1 is longer than str2
	// like str1 := &#34;Hello&#34;, str2 = &#34;&#34;
	// str1 = str2 makes str1 &#34;&#34;
	str = strings.Replace(strings.ToLower(str), &#34;&#39;d&#34;, &#34; would&#34;, -1)

	// If n &lt; 0, there is no limit on the number of replacements.
	str = strings.Replace(str, &#34;&#39;ve&#34;, &#34; have&#34;, -1)
	str = strings.Replace(str, &#34;&#39;re&#34;, &#34; are&#34;, -1)
	str = strings.Replace(str, &#34;&#39;m&#34;, &#34; am&#34;, -1)
	str = strings.Replace(str, &#34;t&#39;s&#34;, &#34;t is&#34;, -1)
	str = strings.Replace(str, &#34;&#39;ll&#34;, &#34; will&#34;, -1)

	str = strings.Replace(str, &#34;won&#39;t&#34;, &#34;will not&#34;, -1)
	str = strings.Replace(str, &#34;can&#39;t&#34;, &#34;can not&#34;, -1)

	str = strings.Replace(str, &#34;haven&#39;t&#34;, &#34;have not&#34;, -1)
	str = strings.Replace(str, &#34;hasn&#39;t&#34;, &#34;has not&#34;, -1)

	str = strings.Replace(str, &#34;dn&#39;t&#34;, &#34;d not&#34;, -1)
	str = strings.Replace(str, &#34;don&#39;t&#34;, &#34;do not&#34;, -1)
	str = strings.Replace(str, &#34;doesn&#39;t&#34;, &#34;does not&#34;, -1)
	str = strings.Replace(str, &#34;didn&#39;t&#34;, &#34;did not&#34;, -1)

	return str
}
