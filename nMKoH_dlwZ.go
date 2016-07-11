// Go 语言使用一个独立的·明确的返回值来传递错误信息的。这与
// 使用异常的 Java 和 Ruby 以及在 C 语言中经常见到的超重的
// 单返回值/错误值相比，Go 语言的处理方式能清楚的知道哪个函数
// 返回了错误，并能像调用那些没有出错的函数一样调用。

package main

import &#34;errors&#34;
import &#34;fmt&#34;

// 按照惯例，错误通常是最后一个返回值并且是 `error` 类
// 型，一个内建的接口。
func f1(arg int) (int, error) {
    if arg == 42 {

        // `errors.New` 构造一个使用给定的错误信息的基本
        // `error` 值。
        return -1, errors.New(&#34;can&#39;t work with 42&#34;)

    }

    // 返回错误值为 nil 代表没有错误。
    return arg &#43; 3, nil
}

// 通过实现 `Error` 方法来自定义 `error` 类型是可以的。
// 这里使用自定义错误类型来表示上面的参数错误。
type argError struct {
    arg  int
    prob string
}

func (e *argError) Error() string {
    return fmt.Sprintf(&#34;%d - %s&#34;, e.arg, e.prob)
}

func f2(arg int) (int, error) {
    if arg == 42 {

        // 在这个例子中，我们使用 `&amp;argError` 语法来建立一个
        // 新的结构体，并提供了 `arg` 和 `prob` 这个两个字段
        // 的值。
        return -1, &amp;argError{arg, &#34;can&#39;t work with it&#34;}
    }
    return arg &#43; 3, nil
}

func main() {

    // 下面的两个循环测试了各个返回错误的函数。注意在 `if`
    // 行内的错误检查代码，在 Go 中是一个普遍的用法。
    for _, i := range []int{7, 42} {
        if r, e := f1(i); e != nil {
            fmt.Println(&#34;f1 failed:&#34;, e)
        } else {
            fmt.Println(&#34;f1 worked:&#34;, r)
        }
    }
    for _, i := range []int{7, 42} {
        if r, e := f2(i); e != nil {
            fmt.Println(&#34;f2 failed:&#34;, e)
        } else {
            fmt.Println(&#34;f2 worked:&#34;, r)
        }
    }

    // 你如果想在程序中使用一个自定义错误类型中的数据，你
    // 需要通过类型断言来得到这个错误类型的实例。
    _, e := f2(42)
    if ae, ok := e.(*argError); ok {
        fmt.Println(ae.arg)
        fmt.Println(ae.prob)
    }
}
