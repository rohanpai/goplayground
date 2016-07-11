package main

import (
	&#34;fmt&#34;
	&#34;runtime&#34;
)

func a() {
	b()
}

func b() {
	c()
}

func c() {
	buf := make([]uintptr, 100)
	n := runtime.Callers(1, buf)
	printStack(buf[:n])
}

func printStack(stack []uintptr) {
	for _, pc := range stack {
		f := runtime.FuncForPC(pc)
		fmt.Println(f.Name() &#43; &#34;()&#34;)
		file, line := f.FileLine(f.Entry())
		fmt.Printf(&#34;\t%s:%d\n&#34;, file, line)
		if f.Name() == &#34;main.main&#34; {
			break
		}
	}
}

func main() {
	a()
}
