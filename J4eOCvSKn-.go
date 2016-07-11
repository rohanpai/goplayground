package main

type A struct {
	Value int
}

func (a A) getVal() int {
	return a.Value
}

func (a *A) getVal2() int {
	return a.Value
}

func main() {
	a := map[int]A{ 1: A{10} }
	
	//println(&amp;a[1]) // wrong, cannot take the address of a[1]
	
	println(a[1].Value) // ok
	//a[1].Value = 20 // wrong, cannot assign to a[1].Value
	
	println(a[1].getVal()) // ok
	//println(a[1].getVal2()) // wrong, cannot call pointer method on a[1]; cannot take the address of a[1]
}
