// Port of some Python's itertools
// ================================
// A.K.A. my first piece of Go code
//
// Written by Nuno Antunes, 2012-08-08
// GitHub: https://github.com/ntns
//
// Python docs: http://docs.python.org/library/itertools.html
// Python source: http://svn.python.org/view/python/tags/r271/Modules/itertoolsmodule.c?view=markup

package main

import "fmt"

func combinations(iterable []int, r int) {
	pool := iterable
	n := len(pool)

	if r > n {
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
		for ; i >= 0 && indices[i] == i+n-r; i -= 1 {
		}

		if i < 0 {
			return
		}

		indices[i] += 1
		for j := i + 1; j < r; j += 1 {
			indices[j] = indices[j-1] + 1
		}

		for ; i < len(indices); i += 1 {
			result[i] = pool[indices[i]]
		}
		fmt.Println(result)

	}

}

func permutations(iterable []int, r int) {
	pool := iterable
	n := len(pool)

	if r > n {
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

	for n > 0 {
		i := r - 1
		for ; i >= 0; i -= 1 {
			cycles[i] -= 1
			if cycles[i] == 0 {
				index := indices[i]
				for j := i; j < n-1; j += 1 {
					indices[j] = indices[j+1]
				}
				indices[n-1] = index
				cycles[i] = n - i
			} else {
				j := cycles[i]
				indices[i], indices[n-j] = indices[n-j], indices[i]

				for k := i; k < r; k += 1 {
					result[k] = pool[indices[k]]
				}

				fmt.Println(result)

				break
			}
		}

		if i < 0 {
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
		for ; i >= 0; i -= 1 {
			pool := pools[i]
			indices[i] += 1

			if indices[i] == len(pool) {
				indices[i] = 0
				result[i] = pool[0]
			} else {
				result[i] = pool[indices[i]]
				break
			}

		}

		if i < 0 {
			return
		}

		fmt.Println(result)
	}
}

func main() {
	fmt.Println("Itertools combinations in Go:")
	// combinations('ABCD', 2) --> AB AC AD BC BD CD
	// combinations(range(4), 3) --> 012 013 023 123
	fmt.Printf("iterable = %s, r = %d", "[]int{1, 2, 3, 4, 5, 6}", 3)
	fmt.Println()
	combinations([]int{1, 2, 3, 4, 5, 6}, 3)

	fmt.Println("Itertools permutations in Go:")
	// permutations('ABCD', 2) --> AB AC AD BA BC BD CA CB CD DA DB DC
	// permutations(range(3)) --> 012 021 102 120 201 210
	fmt.Printf("iterable = %s, r = %d", "[]int{1, 2, 3, 4}", 3)
	fmt.Println()
	permutations([]int{1, 2, 3, 4}, 3)

	fmt.Println("Itertools product in Go:")
	// product('ABCD', 'xy') --> Ax Ay Bx By Cx Cy Dx Dy
	// product(range(2), repeat=3) --> 000 001 010 011 100 101 110 111
	fmt.Printf("iterables = %s, %s", "[]int{1, 2, 3}", "[]int{10, 20, 30}")
	fmt.Println()
	product([]int{1, 2, 3}, []int{10, 20, 30})

}
