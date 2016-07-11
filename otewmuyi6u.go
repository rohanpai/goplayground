package main

import "fmt"

const ( // iota is reset to 0
	c0 = iota // c0 == 0
	c1 = iota // c1 == 1
	c2 = iota // c2 == 2
)

const (
	a = 1 << iota // a == 1 (iota has been reset)
	b = 1 << iota // b == 2
	c = 1 << iota // c == 4
)

const (
	u         = iota * 42 // u == 0     (untyped integer constant)
	v float64 = iota * 42 // v == 42.0  (float64 constant)
	w         = iota * 42 // w == 84    (untyped integer constant)
)

const x = iota // x == 0 (iota has been reset)
const y = iota // y == 0 (iota has been reset)
/**
Within an ExpressionList, the value of each iota is the same
because it is only incremented after each ConstSpec:
**/

const (
	bit0, mask0 = 1 << iota, 1<<iota - 1 // bit0 == 1, mask0 == 0
	bit1, mask1                          // bit1 == 2, mask1 == 1
	_, _                                 // skips iota == 2
	bit3, mask3                          // bit3 == 8, mask3 == 7
)

func main() {

	if x == y && y == u && u == c0 {
		fmt.Println("In each Const Declaration iota reset to 0")
		fmt.Println("x, y, u, c0 == 0")
	}

	fmt.Println("For each ConstSpec the IOTA gets incremented only by one")
	fmt.Println("Within an ExpressionList,\nthe value of each iota is the same")

	fmt.Println("\t\tbit0, mask0 = 1 << iota, 1 << iota - 1")
	fmt.Println("Notice the first ConstSpec of last Const expression")
	fmt.Println("The first ConstSpec is defined as \n\t< identifier-1 > < identifier2 > = 1 << iota, 1 << iota - 1")
	fmt.Println("Therefore, each successive ConstSpec will\n\t1. Increase the iota only once. \n\t2. Use the same expression to evaluate the values of Identifiers")
	fmt.Println("\n\tLets see the same behavior in the last consts  ")
	fmt.Println("bit0, mask0 = ", bit1, mask1, " && IOTA is 0")
	fmt.Println("bit1, mask1 = ", bit1, mask1, " && IOTA is 1")
	fmt.Println("bit2 and mask2 are defined as _, _,\n\t but since it is a next ConstSpec, IOTA gets incremented to 2")
	fmt.Println("bit2, mask2 = UNDEfined in Const Spec as _,_", "&& IOTA is 2")
	fmt.Println("bit3, mask3 = ", bit3, mask3, " && IOTA is 3")

}
