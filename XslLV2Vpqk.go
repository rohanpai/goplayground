package main

import (
	"fmt"
	"runtime"
	"unsafe"
)

func PtrOffset(offset int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(offset))
}

func main() {
	var offset = 10

	ptr := PtrOffset(offset)

	runtime.GC()
	fmt.Println("Hello, playground", ptr)
}
