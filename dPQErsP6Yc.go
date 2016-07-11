// _Slices_ are a key data type in Go, giving a more
// powerful interface to sequences than arrays.

package main

import &#34;fmt&#34;

func main() {

    // Unlike arrays, slices are typed only by the
    // elements they contain (not the number of elements).
    // To create an empty slice with non-zero length, use
    // the builtin `make`. Here we make a slice of
    // `string`s of length `3` (initially zero-valued).
    s := make([]string, 3)
    fmt.Println(&#34;emp:&#34;, s)

    // We can set and get just like with arrays.
    s[0] = &#34;a&#34;
    s[1] = &#34;b&#34;
    s[2] = &#34;c&#34;
    fmt.Println(&#34;set:&#34;, s)
    fmt.Println(&#34;get:&#34;, s[2])

    // `len` returns the length of the slice as expected.
    fmt.Println(&#34;len:&#34;, len(s))

    // In addition to these basic operations, slices
    // support several more that make them richer than
    // arrays. One is the builtin `append`, which
    // returns a slice containing one or more new values.
    // Note that we need to accept a return value from
    // append as we may get a new slice value.
    s = append(s, &#34;d&#34;)
    s = append(s, &#34;e&#34;, &#34;f&#34;)
    fmt.Println(&#34;apd:&#34;, s)

    // Slices can also be `copy`&#39;d. Here we create an
    // empty slice `c` of the same length as `s` and copy
    // into `c` from `s`.
    c := make([]string, len(s))
    copy(c, s)
    fmt.Println(&#34;cpy:&#34;, c)

    // Slices support a &#34;slice&#34; operator with the syntax
    // `slice[low:high]`. For example, this gets a slice
    // of the elements `s[2]`, `s[3]`, and `s[4]`.
    l := s[2:5]
    fmt.Println(&#34;sl1:&#34;, l)

    // This slices up to (but excluding) `s[5]`.
    l = s[:5]
    fmt.Println(&#34;sl2:&#34;, l)

    // And this slices up from (and including) `s[2]`.
    l = s[2:]
    fmt.Println(&#34;sl3:&#34;, l)

    // We can declare and initialize a variable for slice
    // in a single line as well.
    t := []string{&#34;g&#34;, &#34;h&#34;, &#34;i&#34;}
    fmt.Println(&#34;dcl:&#34;, t)

    // Slices can be composed into multi-dimensional data
    // structures. The length of the inner slices can
    // vary, unlike with multi-dimensional arrays.
    twoD := make([][]int, 3)
    for i := 0; i &lt; 3; i&#43;&#43; {
        innerLen := i &#43; 1
        twoD[i] = make([]int, innerLen)
        for j := 0; j &lt; innerLen; j&#43;&#43; {
            twoD[i][j] = i &#43; j
        }
    }
    fmt.Println(&#34;2d: &#34;, twoD)
}
