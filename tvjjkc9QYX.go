package main

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

type Struct struct {
	p unsafe.Pointer // some pointer
}

func main() {
	data := 1

	info := Struct{p: unsafe.Pointer(&data)}

	fmt.Printf("info is %d\n", *(*int)(info.p))

	otherData := 2

	atomic.StorePointer(&info.p, unsafe.Pointer(&otherData))

	fmt.Printf("info is %d\n", *(*int)(info.p))

}
