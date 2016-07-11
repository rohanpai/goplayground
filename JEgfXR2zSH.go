// Port of some Python&#39;s itertools
// ================================
// A.K.A. my first piece of Go code
//
// Written by Nuno Antunes, 2012-08-08
// GitHub: https://github.com/ntns
//
// Python docs: http://docs.python.org/library/itertools.html
// Python source: http://svn.python.org/view/python/tags/r271/Modules/itertoolsmodule.c?view=markup

package main

import &#34;fmt&#34;

func combinations(iterable []int, r int) {
	pool := iterable
	n := len(pool)

	if r &gt; n {
		return
	}

	indices := make([]int, r)
	for i := range indices {
		indices[i] = i
	}

	result := make([]int, r)
	for i, el := range indices {
		result[i] = pool[el]
	}

	fmt.Println(result)

	for {
		i := r - 1
		for ; i &gt;= 0 &amp;&amp; indices[i] == i&#43;n-r; i -= 1 {
		}

		if i &lt; 0 {
			return
		}

		indices[i] &#43;= 1
		for j := i &#43; 1; j &lt; r; j &#43;= 1 {
			indices[j] = indices[j-1] &#43; 1
		}

		for ; i &lt; len(indices); i &#43;= 1 {
			result[i] = pool[indices[i]]
		}
		fmt.Println(result)

	}

}

func permutations(iterable []int, r int) {
	pool := iterable
	n := len(pool)

	if r &gt; n {
		return
	}

	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}

	cycles := make([]int, r)
	for i := range cycles {
		cycles[i] = n - i
	}

	result := make([]int, r)
	for i, el := range indices[:r] {
		result[i] = pool[el]
	}

	fmt.Println(result)

	for n &gt; 0 {
		i := r - 1
		for ; i &gt;= 0; i -= 1 {
			cycles[i] -= 1
			if cycles[i] == 0 {
				index := indices[i]
				for j := i; j &lt; n-1; j &#43;= 1 {
					indices[j] = indices[j&#43;1]
				}
				indices[n-1] = index
				cycles[i] = n - i
			} else {
				j := cycles[i]
				indices[i], indices[n-j] = indices[n-j], indices[i]

				for k := i; k &lt; r; k &#43;= 1 {
					result[k] = pool[indices[k]]
				}

				fmt.Println(result)

				break
			}
		}

		if i &lt; 0 {
			return
		}

	}

}

func product(argsA, argsB []int) {

	pools := [][]int{argsA, argsB}
	npools := len(pools)
	indices := make([]int, npools)

	result := make([]int, npools)
	for i := range result {
		result[i] = pools[i][0]
	}

	fmt.Println(result)

	for {
		i := npools - 1
		for ; i &gt;= 0; i -= 1 {
			pool := pools[i]
			indices[i] &#43;= 1

			if indices[i] == len(pool) {
				indices[i] = 0
				result[i] = pool[0]
			} else {
				result[i] = pool[indices[i]]
				break
			}

		}

		if i &lt; 0 {
			return
		}

		fmt.Println(result)
	}
}

func main() {
	fmt.Println(&#34;Itertools combinations in Go:&#34;)
	// combinations(&#39;ABCD&#39;, 2) --&gt; AB AC AD BC BD CD
	// combinations(range(4), 3) --&gt; 012 013 023 123
	fmt.Printf(&#34;iterable = %s, r = %d&#34;, &#34;[]int{1, 2, 3, 4, 5, 6}&#34;, 3)
	fmt.Println()
	combinations([]int{1, 2, 3, 4, 5, 6}, 3)

	fmt.Println(&#34;Itertools permutations in Go:&#34;)
	// permutations(&#39;ABCD&#39;, 2) --&gt; AB AC AD BA BC BD CA CB CD DA DB DC
	// permutations(range(3)) --&gt; 012 021 102 120 201 210
	fmt.Printf(&#34;iterable = %s, r = %d&#34;, &#34;[]int{1, 2, 3, 4}&#34;, 3)
	fmt.Println()
	permutations([]int{1, 2, 3, 4}, 3)

	fmt.Println(&#34;Itertools product in Go:&#34;)
	// product(&#39;ABCD&#39;, &#39;xy&#39;) --&gt; Ax Ay Bx By Cx Cy Dx Dy
	// product(range(2), repeat=3) --&gt; 000 001 010 011 100 101 110 111
	fmt.Printf(&#34;iterables = %s, %s&#34;, &#34;[]int{1, 2, 3}&#34;, &#34;[]int{10, 20, 30}&#34;)
	fmt.Println()
	product([]int{1, 2, 3}, []int{10, 20, 30})

}
