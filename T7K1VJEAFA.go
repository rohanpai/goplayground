package main

import (
	&#34;fmt&#34;
	&#34;time&#34;
)

const Towers = 3
const Disks = 5

type Hanoi [Towers][]int

func main() {
	var state Hanoi
	state.init(Disks)
	state.move(Disks, 0, 1, 2)
}

func (h *Hanoi) init(n int) {
	h[0] = make([]int, n)
	for i := range h[0] {
		h[0][i] = n - i
	}
	h.print()
}

func (h *Hanoi) move(n, a, b, c int) {
	if n &lt;= 0 {
		return
	}
	h.move(n-1, a, c, b)
	disk := h[a][len(h[a])-1]
	h[a] = h[a][:len(h[a])-1]
	h[c] = append(h[c], disk)
	h.print()
	h.move(n-1, b, a, c)
}

func (h *Hanoi) print() {
	fmt.Print(&#34;\f&#34;)
	for i := Disks; i &gt;= 0; i-- {
		for j := 0; j &lt; Towers; j&#43;&#43; {
			if i == 0 {
				fmt.Print(&#34;_/||\\_&#34;)
			} else if len(h[j]) &gt;= i {
				fmt.Printf(&#34;  %02d  &#34;, h[j][i-1])
			} else {
				fmt.Print(&#34;  ||  &#34;)
			}
		}
		fmt.Println()
	}
	fmt.Println()
	time.Sleep(time.Second / 5)
}
