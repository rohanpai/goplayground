package main

import (
	&#34;fmt&#34;
	&#34;runtime&#34;
	&#34;unsafe&#34;
)

func PtrOffset(offset int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(offset))
}

func main() {
	var offset = 10

	ptr := PtrOffset(offset)

	runtime.GC()
	fmt.Println(&#34;Hello, playground&#34;, ptr)
}
