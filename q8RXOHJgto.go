// 这里说明为什么需要谓的 MiniGo.
// 在这个这个例子中 origin 和 another 两段 Go 代码除了排版(注释)实际上是一样的.
// 目前 (go 1.4), 无法简单分析出他们是一样的, 除非直接分析 AST.
// 如果 &#34;go/printer&#34; 能支持正确剔除回车, 空行, 缩进并补全分号, 即所谓的 MiniGo.
// 那在就可以简单的分析出两段代码是否是一样(不同的只是排版和注释).
// 典型的应用场景是翻译 Go package, 翻译者可能因翻译需求, 或者不小心改动了非注释的代码部分.
// 有了 MiniGo 支持, 就可以编写程序判断代码是否只是排版/注释不同, 从而判定翻译对代码的改动是否可允许.
// 请至 main 函数中改动 io.Writer 参数, 查看不同输出.
package main

import (
	&#34;fmt&#34;
	&#34;go/parser&#34;
	&#34;go/printer&#34;
	&#34;go/token&#34;
	&#34;io&#34;
	&#34;os&#34;
)

const origin = `
package semicolon

import (
	&#34;pkg&#34;
)

const (
	CA = pkg.CONSTA
	CB = pkg.CONSTB
)

func Fn() (n int,
	err error) {
	return
}
`

const another = `
package semicolon

import (
	&#34;pkg&#34;
)

const (
	CA = pkg.CONSTA

	CB = pkg.CONSTB
)

func Fn() (n int, err error) {
	return
}
`

// TrimNewLine 模拟预想的 MiniGo 中的剔除输出中的换行功能.
// 实际的情况更复杂, 不应该用 io.Writer 实现.
// 比如: `strings` 这种字符串包含的换行是应该被保留的.
type TrimNewLine struct {
	io.Writer
}

func (w TrimNewLine) Write(p []byte) (n int, err error) {
	var c byte
	var e, i int
	n = len(p)
	if n == 0 {
		return
	}
	for i, c = range p {
		if c != &#39;\n&#39; {
			continue
		}
		if e != i {
			_, err = w.Writer.Write(p[e:i])
		}
		e = i &#43; 1
		if err != nil {
			return
		}
	}

	if e != n {
		_, err = w.Writer.Write(p[e:n])
	}
	return
}

func fprint(src string, w io.Writer) {
	cfg := printer.Config{Indent: 0}
	fset := token.NewFileSet()

	file, _ := parser.ParseFile(fset, &#34;&#34;, src, 0)

	cfg.Fprint(w, fset, file)
}

var trim = TrimNewLine{os.Stdout}

func main() {
	// w 换成 trim 试试, 当然目前没法正确插入分号.
	// trim 和 os.Stdout 的结果略有不同, 即便一致, 因输出是非法的, 也不能作为判定依据.
	w := os.Stdout
	fmt.Println(&#34;========= origin ========&#34;)
	fprint(origin, w)
	fmt.Println(&#34;\n========= another =======&#34;)
	fprint(another, w)
}

// 从输出中可以看出不支持 MiniGo 功能的输出,
// 无法简单判断一个节点的源码只是版式或者注释的变化.
// 对于这个例子来说,两个版本的 MiniGo 的结果应该是一样的, 可能是这个样子
const mini = `
package semicolon;import (&#34;pkg&#34;);const (CA= pkg.CONSTA;CB= pkg.CONSTB;);func Fn() (n int,err error) {return;};
`
