package main

import (
	&#34;bytes&#34;
	&#34;fmt&#34;
	&#34;sort&#34;
	&#34;strings&#34;
)

func main() {
	func() {
		slice := []string{&#34;D&#34;, &#34;C&#34;, &#34;A&#34;, &#34;E&#34;, &#34;X&#34;}
		rs := permuteStrings(slice)
		// var buf bytes.Buffer
		buf := new(bytes.Buffer)
		for _, elem1 := range rs {
			for _, elem2 := range elem1 {
				buf.WriteString(elem2)
			}
		}
		if len(buf.String()) != 600 {
			fmt.Errorf(&#34;Length should be 600 but %d&#34;, len(buf.String()))
		}
	}()

	func() {
		slice := []string{&#34;E&#34;, &#34;A&#34;, &#34;A&#34;}
		rs := permuteStrings(slice)
		// var buf bytes.Buffer
		buf := new(bytes.Buffer)
		for _, elem1 := range rs {
			for _, elem2 := range elem1 {
				buf.WriteString(elem2)
			}
		}
		if len(buf.String()) != 9 || len(rs) != 3 {
			fmt.Errorf(&#34;Length should be 9 but %d&#34;, len(buf.String()))
		}
	}()

	func() {
		slice := []string{&#34;1&#34;, &#34;2&#34;, &#34;3&#34;}
		rs := permuteStrings(slice)
		// var buf bytes.Buffer
		buf := new(bytes.Buffer)
		for _, elem1 := range rs {
			for _, elem2 := range elem1 {
				buf.WriteString(elem2)
			}
			buf.WriteString(&#34;---&#34;)
		}
		rStr := buf.String()
		if len(rStr) != 36 {
			fmt.Errorf(&#34;Length should be 600 but %s&#34;, rStr)
		}
		if !strings.Contains(rStr, &#34;123&#34;) ||
			!strings.Contains(rStr, &#34;132&#34;) ||
			!strings.Contains(rStr, &#34;213&#34;) ||
			!strings.Contains(rStr, &#34;231&#34;) ||
			!strings.Contains(rStr, &#34;312&#34;) ||
			!strings.Contains(rStr, &#34;321&#34;) {
			fmt.Errorf(&#34;Missing a permutation %s&#34;, rStr)
		}
	}()
}

func first(data sort.Interface) {
	sort.Sort(data)
}

// next returns false when it cannot permute any more
// http://en.wikipedia.org/wiki/Permutation#Generation_in_lexicographic_order
func next(data sort.Interface) bool {
	var k, l int
	for k = data.Len() - 2; ; k-- {
		if k &lt; 0 {
			return false
		}
		if data.Less(k, k&#43;1) {
			break
		}
	}
	for l = data.Len() - 1; !data.Less(k, l); l-- {
	}
	data.Swap(k, l)
	for i, j := k&#43;1, data.Len()-1; i &lt; j; i&#43;&#43; {
		data.Swap(i, j)
		j--
	}
	return true
}

// permuteStrings returns all possible permutations of string slice.
func permuteStrings(slice []string) [][]string {
	first(sort.StringSlice(slice))

	copied1 := make([]string, len(slice)) // we need to make a copy!
	copy(copied1, slice)
	result := [][]string{copied1}

	for {
		isDone := next(sort.StringSlice(slice))
		if !isDone {
			break
		}

		// https://groups.google.com/d/msg/golang-nuts/ApXxTALc4vk/z1-2g1AH9jQJ
		// Lesson from Dave Cheney:
		// A slice is just a pointer to the underlying back array, your storing multiple
		// copies of the slice header, but they all point to the same backing array.

		// NOT
		// result = append(result, slice)

		copied2 := make([]string, len(slice))
		copy(copied2, slice)
		result = append(result, copied2)
	}

	combNum := 1
	for i := 0; i &lt; len(slice); i&#43;&#43; {
		combNum *= i &#43; 1
	}
	if len(result) != combNum {
		fmt.Printf(&#34;Expected %d combinations but %&#43;v because of duplicate elements&#34;, combNum, result)
	}

	return result
}
