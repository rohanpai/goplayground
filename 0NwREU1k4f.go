package main

import &#34;fmt&#34;

// 数列を作る
func seq(a int, b int) (ret []int) {
	if a &gt;= b {
		panic(&#34;invalid argument: a must be lower than b&#34;)
	}
	ret = make([]int, b-a&#43;1)
	for i := 0; i &lt; len(ret); i&#43;&#43; {
		ret[i] = a &#43; i
	}
	return
}

// 2つのスライスを連結する
func concat(a []int, b []int) []int {
	for _, elm := range b {
		a = append(a, elm)
	}
	return a
}

func foldl(f func(int, int) int, acc int, list []int) (ret int) {
	x, xs := list[0], list[1:]
	if len(xs) == 0 {
		ret = f(acc, x)
	} else {
		ret = foldl(f, f(acc, x), xs)
	}
	return
}

func foldr(f func(int, int) int, op int, list []int) (ret int) {
	x, xs := list[0], list[1:]
	if len(xs) == 0 {
		ret = f(x, op)
	} else {
		ret = f(foldr(f, op, xs), x)
	}
	return
}

func mapf(f func(int) int, list []int) (ret []int) {
	x, xs := list[0], list[1:]
	if len(xs) == 0 {
		ret = []int{ f(x) }
	} else {
		ret = concat([]int{ f(x) }, mapf(f, xs))
	}
	return ret
}

func filter(pred func(int) bool, list []int) (ret []int) {
	x, xs := list[0], list[1:]
	if len(xs) == 0 {
		ret = []int{}
	} else if pred(x) {
		ret = concat([]int{ x }, filter(pred, xs))
	} else {
		ret = filter(pred, xs)
	}
	return
}

func main() {
	plus := func(a int, b int) int {
		fmt.Println(&#34;---&gt; (&#34;, a, &#34;&#43;&#34;, b, &#34;) =&#34;, a&#43;b)
		return a &#43; b
	}
	// 左からの畳み込み
	fmt.Println(foldl(plus, 0, seq(0, 9)))
	// 右からの畳み込み
	fmt.Println(foldr(plus, 0, seq(0, 9)))
	// マップ (各要素を二倍)
	fmt.Println(mapf(func(a int) int { return a*a }, seq(0, 9)))
	// フィルター (偶数だけ通す)
	fmt.Println(filter(func(a int) bool { return a % 2 == 0 }, seq(0, 9)))
}
