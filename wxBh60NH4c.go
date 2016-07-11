// Go 提供内置的 JSON 编解码支持，包括内置或者自定义类
// 型与 JSON 数据之间的转化。

package main

import &#34;encoding/json&#34;
import &#34;fmt&#34;
import &#34;os&#34;

// 下面我们将使用这两个结构体来演示自定义类型的编码和解
// 码。
type Response1 struct {
	Page   int
	Fruits []string
}
type Response2 struct {
	Page   int      `json:&#34;page&#34;`
	Fruits []string `json:&#34;fruits&#34;`
}

func main() {

	// 首先我们来看一下基本数据类型到 JSON 字符串的编码
	// 过程。这里是一些原子值的例子。
	bolB, _ := json.Marshal(true)
	fmt.Println(string(bolB))

	intB, _ := json.Marshal(1)
	fmt.Println(string(intB))

	fltB, _ := json.Marshal(2.34)
	fmt.Println(string(fltB))

	strB, _ := json.Marshal(&#34;gopher&#34;)
	fmt.Println(string(strB))

	// 这里是一些切片和 map 编码成 JSON 数组和对象的例子。
	slcD := []string{&#34;apple&#34;, &#34;peach&#34;, &#34;pear&#34;}
	slcB, _ := json.Marshal(slcD)
	fmt.Println(string(slcB))

	mapD := map[string]int{&#34;apple&#34;: 5, &#34;lettuce&#34;: 7}
	mapB, _ := json.Marshal(mapD)
	fmt.Println(string(mapB))

	// JSON 包可以自动的编码你的自定义类型。编码仅输出可
	// 导出的字段，并且默认使用他们的名字作为 JSON 数据的
	// 键。
	res1D := &amp;Response1{
		Page:   1,
		Fruits: []string{&#34;apple&#34;, &#34;peach&#34;, &#34;pear&#34;}}
	res1B, _ := json.Marshal(res1D)
	fmt.Println(string(res1B))

	// 你可以给结构字段声明标签来自定义编码的 JSON 数据键
	// 名称。在上面 `Response2` 的定义可以作为这个标签这个
	// 的一个例子。
	res2D := &amp;Response2{
		Page:   1,
		Fruits: []string{&#34;apple&#34;, &#34;peach&#34;, &#34;pear&#34;}}
	res2B, _ := json.Marshal(res2D)
	fmt.Println(string(res2B))

	// 现在来看看解码 JSON 数据为 Go 值的过程。这里
	// 是一个普通数据结构的解码例子。
	byt := []byte(`{&#34;num&#34;:6.13,&#34;strs&#34;:[&#34;a&#34;,&#34;b&#34;]}`)

	// 我们需要提供一个 JSON 包可以存放解码数据的变量。这里
	// 的 `map[string]interface{}` 将保存一个 string 为键，
	// 值为任意值的map。
	var dat map[string]interface{}

	// 这里就是实际的解码和相关的错误检查。
	if err := json.Unmarshal(byt, &amp;dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)

	// 为了使用解码 map 中的值，我们需要将他们进行适当的类
	// 型转换。例如这里我们将 `num` 的值转换成 `float64`
	// 类型。
	num := dat[&#34;num&#34;].(float64)
	fmt.Println(num)

	// 访问嵌套的值需要一系列的转化。
	strs := dat[&#34;strs&#34;].([]interface{})
	str1 := strs[0].(string)
	fmt.Println(str1)

	// 我们也可以解码 JSON 值到自定义类型。这个功能的好处就
	// 是可以为我们的程序带来额外的类型安全加强，并且消除在
	// 访问数据时的类型断言。
	str := `{&#34;page&#34;: 1, &#34;fruits&#34;: [&#34;apple&#34;, &#34;peach&#34;]}`
	res := &amp;Response2{}
	json.Unmarshal([]byte(str), &amp;res)
	fmt.Println(res)
	fmt.Println(res.Fruits[0])

	// 在上面的例子中，我们经常使用 byte 和 string 作为使用
	// 标准输出时数据和 JSON 表示之间的中间值。我们也可以和
	// `os.Stdout` 一样，直接将 JSON 编码直接输出至 `os.Writer`
	// 流中，或者作为 HTTP 响应体。
	enc := json.NewEncoder(os.Stdout)
	d := map[string]int{&#34;apple&#34;: 5, &#34;lettuce&#34;: 7}
	enc.Encode(d)
}
