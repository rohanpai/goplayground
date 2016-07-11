package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

func main() {
	fmt.Println(SwapCase("Hello Hi"))
	// hELLO hI

	fmt.Println(SwapCaseII("Hello Hi"))
	// hELLO hI

	fmt.Println(ReverseString("Hello Hi"))
	// iH olleH

	fmt.Println(CheckPalindrome("Anne, I vote more cars race Rome-to-Vienna"))
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
// For example, a rune literal 'a' is a number 97.
// 32 is the offset of the uppercase and lowercase characters.
// So if you add 32 to 'A', you get 'a' and vice versa.
func SwapRune(r rune) rune {
	switch {
	case 'a' <= r && r <= 'z':
		return r - 'a' + 'A'
	case 'A' <= r && r <= 'Z':
		return r - 'A' + 'a'
	default:
		return r
	}
}

// ReverseString changes the order of the input string.
// It returns the new version of input string.
// We need to use rune because string is immutable.
func ReverseString(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
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
	str = validID.ReplaceAllString(str, "")

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
	// slicing off all "leading" and
	// "trailing" white space, as defined by Unicode.
	str = strings.TrimSpace(str)

	// func Fields(s string) []string
	// Fields splits the slice s around each instance
	// of "one or more consecutive white space"
	slice := strings.Fields(str)

	// now join them with a single white space character
	return strings.Join(slice, " ")
}

// ReplaceNonAlnumWithSpace removes all alphanumeric characters.
// It replaces them with a single whitespace character.
// It returns the new version of input string, in lower case.
func ReplaceNonAlnumWithSpace(str string) string {
	str = ExpandApostrophe(str)
	// alphanumeric (== [0-9A-Za-z])
	// \s is a white space character
	validID := regexp.MustCompile(`[^[:alnum:]\s]`)
	return validID.ReplaceAllString(str, " ")
}

// ExpandApostrophe expands the apostrophe phrases.
// And convert them to lower case letters.
func ExpandApostrophe(str string) string {
	// assignment between string is not "copy"
	// even if str1 is longer than str2
	// like str1 := "Hello", str2 = ""
	// str1 = str2 makes str1 ""
	str = strings.Replace(strings.ToLower(str), "'d", " would", -1)

	// If n < 0, there is no limit on the number of replacements.
	str = strings.Replace(str, "'ve", " have", -1)
	str = strings.Replace(str, "'re", " are", -1)
	str = strings.Replace(str, "'m", " am", -1)
	str = strings.Replace(str, "t's", "t is", -1)
	str = strings.Replace(str, "'ll", " will", -1)

	str = strings.Replace(str, "won't", "will not", -1)
	str = strings.Replace(str, "can't", "can not", -1)

	str = strings.Replace(str, "haven't", "have not", -1)
	str = strings.Replace(str, "hasn't", "has not", -1)

	str = strings.Replace(str, "dn't", "d not", -1)
	str = strings.Replace(str, "don't", "do not", -1)
	str = strings.Replace(str, "doesn't", "does not", -1)
	str = strings.Replace(str, "didn't", "did not", -1)

	return str
}
