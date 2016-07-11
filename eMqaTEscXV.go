// Go 拥有各值类型，包括字符串，整形，浮点型，布尔
// 型等。下面是一些基本的例子。

package main

import &#34;fmt&#34;

func main() {

    // 字符串可以通过 `&#43;` 连接。
    fmt.Println(&#34;go&#34; &#43; &#34;lang&#34;)

    // 整数和浮点数
    fmt.Println(&#34;1&#43;1 =&#34;, 1&#43;1)
    fmt.Println(&#34;7.0/3.0 =&#34;, 7.0/3.0)

    // 布尔型，还有你想要的逻辑运算符。
    fmt.Println(true &amp;&amp; false)
    fmt.Println(true || false)
    fmt.Println(!true)
}
