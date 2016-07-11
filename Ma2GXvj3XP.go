package main

import (
	&#34;fmt&#34;
)

func Merge(l, r []int) []int {
	ret := make([]int, 0, len(l)&#43;len(r))
	for len(l) &gt; 0 || len(r) &gt; 0 {
		if len(l) == 0 {
			return append(ret, r...)
		}
		if len(r) == 0 {
			return append(ret, l...)
		}
		if l[0] &lt;= r[0] {
			ret = append(ret, l[0])
			l = l[1:]
		} else {
			ret = append(ret, r[0])
			r = r[1:]
		}
	}
	return ret
}

func MergeSort(s []int) []int {
	if len(s) &lt;= 1 {
		return s
	}
	n := len(s) / 2
	l := MergeSort(s[:n])
	r := MergeSort(s[n:])
	return Merge(l, r)
}

func main() {
	s := []int{9, 4, 3, 6, 1, 2, 10, 5, 7, 8}
	fmt.Printf(&#34;%v\n%v\n&#34;, s, MergeSort(s))
}
