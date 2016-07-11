package main

// char **a;
// char *b[] = {"0", "1", "2"};
// void init() {
//   a = b;
// }
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	C.init()

	fmt.Printf("%T\n", C.a)
	var A []*C.char
	slice := reflect.SliceHeader{uintptr(unsafe.Pointer(C.a)), 3, 3}
	a := reflect.NewAt(reflect.TypeOf(A), unsafe.Pointer(&slice)).Elem().Interface()
	fmt.Printf("%T\n", a)
	for _, s := range a.([]*C.char) {
		fmt.Println(C.GoString(s))
	}
}
