// 常规的通过通道发送和接收数据是阻塞的。然而，我们可以
// 使用带一个 `default` 子句的 `select` 来实现_非阻塞_ 的
// 发送、接收，甚至是非阻塞的多路 `select`。

package main

import &#34;fmt&#34;

func main() {
    messages := make(chan string)
    signals := make(chan bool)

    // 这里是一个非阻塞接收的例子。如果在 `messages` 中
    // 存在，然后 `select` 将这个值带入 `&lt;-messages` `case`
    // 中。如果不是，就直接到 `default` 分支中。
    select {
    case msg := &lt;-messages:
        fmt.Println(&#34;received message&#34;, msg)
    default:
        fmt.Println(&#34;no message received&#34;)
    }

    // 一个非阻塞发送的实现方法和上面一样。
    msg := &#34;hi&#34;
    select {
    case messages &lt;- msg:
        fmt.Println(&#34;sent message&#34;, msg)
    default:
        fmt.Println(&#34;no message sent&#34;)
    }

    // 我们可以在 `default` 前使用多个 `case` 子句来实现
    // 一个多路的非阻塞的选择器。这里我们试图在 `messages`
    // 和 `signals` 上同时使用非阻塞的接受操作。
    select {
    case msg := &lt;-messages:
        fmt.Println(&#34;received message&#34;, msg)
    case sig := &lt;-signals:
        fmt.Println(&#34;received signal&#34;, sig)
    default:
        fmt.Println(&#34;no activity&#34;)
    }
}
