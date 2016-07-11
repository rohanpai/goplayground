package main

import (
	&#34;fmt&#34;
	&#34;sync/atomic&#34;
	&#34;unsafe&#34;
)

type Struct struct {
	p unsafe.Pointer // some pointer
}

func main() {
	data := 1

	info := Struct{p: unsafe.Pointer(&amp;data)}

	fmt.Printf(&#34;info is %d\n&#34;, *(*int)(info.p))

	otherData := 2

	atomic.StorePointer(&amp;info.p, unsafe.Pointer(&amp;otherData))

	fmt.Printf(&#34;info is %d\n&#34;, *(*int)(info.p))

}
