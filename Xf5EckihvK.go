package main

// char **a;
// char *b[] = {&#34;0&#34;, &#34;1&#34;, &#34;2&#34;};
// void init() {
//   a = b;
// }
import &#34;C&#34;

import (
	&#34;fmt&#34;
	&#34;reflect&#34;
	&#34;unsafe&#34;
)

func main() {
	C.init()

	fmt.Printf(&#34;%T\n&#34;, C.a)
	var A []*C.char
	slice := reflect.SliceHeader{uintptr(unsafe.Pointer(C.a)), 3, 3}
	a := reflect.NewAt(reflect.TypeOf(A), unsafe.Pointer(&amp;slice)).Elem().Interface()
	fmt.Printf(&#34;%T\n&#34;, a)
	for _, s := range a.([]*C.char) {
		fmt.Println(C.GoString(s))
	}
}
