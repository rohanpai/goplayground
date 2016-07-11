package main

import &#34;fmt&#34;

func Max(more ...int) int {
	max_num := more[0]
	for _, elem := range more {
		if max_num &lt; elem {
			max_num = elem
		}
	}
	return max_num
}

func LCS(str1, str2 string) (int, string) {
	len1 := len(str1)
	len2 := len(str2)

	table := make([][]int, len1&#43;1)
	for i := range table {
		table[i] = make([]int, len2&#43;1)
	}

	i, j := 0, 0
	for i = 0; i &lt;= len1; i&#43;&#43; {
		for j = 0; j &lt;= len2; j&#43;&#43; {
			if i == 0 || j == 0 {
				table[i][j] = 0
			} else if str1[i-1] == str2[j-1] {
				table[i][j] = table[i-1][j-1] &#43; 1
			} else {
				table[i][j] = Max(table[i-1][j], table[i][j-1])
			}
		}
	}
	return table[len1][len2], Back(table, str1, str2, len1-1, len2-1)
}

//http://en.wikipedia.org/wiki/Longest_common_subsequence_problem
func Back(table [][]int, str1, str2 string, i, j int) string {
	if i == 0 || j == 0 {
		return &#34;&#34;
	} else if str1[i] == str2[j] {
		return Back(table, str1, str2, i-1, j-1) &#43; string(str1[i])
	} else {
		if table[i][j-1] &gt; table[i-1][j] {
			return Back(table, str1, str2, i, j-1)
		} else {
			return Back(table, str1, str2, i-1, j)
		}
	}
}

func main() {
	str1 := &#34;AGGTABTABTABTAB&#34;
	str2 := &#34;GXTXAYBTABTABTAB&#34;
	fmt.Println(LCS(str1, str2))
	//Actual Longest Common Subsequence: GTABTABTABTAB
	//13 ABTABTABTAB

	str3 := &#34;AGGTABGHSRCBYJSVDWFVDVSBCBVDWFDWVV&#34;
	str4 := &#34;GXTXAYBRGDVCBDVCCXVXCWQRVCBDJXCVQSQQ&#34;
	fmt.Println(LCS(str3, str4))
	//Actual Longest Common Subsequence: ?
	//14 ABGCBDV
}
