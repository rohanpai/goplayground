package main

import &#34;fmt&#34;
import &#34;strings&#34;

func main() {
	var str = &#34;1234567890&#34;
	fmt.Println(str)
	fmt.Println(&#34;leftPad(str, \&#34;*\&#34;, 3) &#34;, leftPad(str, &#34;*&#34;, 3))
	fmt.Println(&#34;leftPad2Len(str, \&#34;*-\&#34;, 13) &#34;, leftPad2Len(str, &#34;*-&#34;, 13))
	fmt.Println(&#34;leftPad2Len(str, \&#34;*-\&#34;, 14) &#34;, leftPad2Len(str, &#34;*-&#34;, 14))
	fmt.Println(&#34;leftPad2Len(str, \&#34;*\&#34;, 14) &#34;, leftPad2Len(str, &#34;*&#34;, 14))
	fmt.Println(&#34;leftPad2Len(str, \&#34;*-x\&#34;, 14) &#34;, leftPad2Len(str, &#34;*-x&#34;, 14))
	fmt.Println(&#34;leftPad2Len(str, \&#34;ABCDE\&#34;, 14) &#34;, leftPad2Len(str, &#34;ABCDE&#34;, 14))
	fmt.Println(&#34;leftPad2Len(str, \&#34;ABCDE\&#34;, 4) &#34;, leftPad2Len(str, &#34;ABCDE&#34;, 4))
	fmt.Println(&#34;rightPad(str, \&#34;*\&#34;, 3) &#34;, rightPad(str, &#34;*&#34;, 3))
	fmt.Println(&#34;rightPad(str, \&#34;*!\&#34;, 3) &#34;, rightPad(str, &#34;*!&#34;, 3))
	fmt.Println(&#34;rightPad2Len(str, \&#34;*-\&#34;, 13) &#34;, rightPad2Len(str, &#34;*-&#34;, 13))
	fmt.Println(&#34;rightPad2Len(str, \&#34;*-\&#34;, 14) &#34;, rightPad2Len(str, &#34;*-&#34;, 14))
	fmt.Println(&#34;rightPad2Len(str, \&#34;*\&#34;, 14) &#34;, rightPad2Len(str, &#34;*&#34;, 14))
	fmt.Println(&#34;rightPad2Len(str, \&#34;*-x\&#34;, 14) &#34;, rightPad2Len(str, &#34;*-x&#34;, 14))
	fmt.Println(&#34;rightPad2Len(str, \&#34;ABCDE\&#34;, 14) &#34;, rightPad2Len(str, &#34;ABCDE&#34;, 14))
	fmt.Println(&#34;rightPad2Len(str, \&#34;ABCDE\&#34;, 4) &#34;, rightPad2Len(str, &#34;ABCDE&#34;, 4))
}

//TODO convert these into a
/*
* leftPad and rightPad just repoeat the padStr the indicated
* number of times
*
 */
func leftPad(s string, padStr string, pLen int) string {
	return strings.Repeat(padStr, pLen) &#43; s
}
func rightPad(s string, padStr string, pLen int) string {
	return s &#43; strings.Repeat(padStr, pLen)
}

/* the Pad2Len functions are generally assumed to be padded with short sequences of strings
* in many cases with a single character sequence
*
* so we assume we can build the string out as if the char seq is 1 char and then
* just substr the string if it is longer than needed
*
* this means we are wasting some cpu and memory work
* but this always get us to want we want it to be
*
* in short not optimized to for massive string work
*
* If the overallLen is shorter than the original string length
* the string will be shortened to this length (substr)
*
 */
func rightPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 &#43; ((overallLen - len(padStr)) / len(padStr))
	var retStr = s &#43; strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}
func leftPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 &#43; ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) &#43; s
	return retStr[(len(retStr) - overallLen):]
}
