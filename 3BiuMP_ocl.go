package main

import (
	"fmt"
	"runtime"
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
		fmt.Println(f.Name() + "()")
		file, line := f.FileLine(f.Entry())
		fmt.Printf("\t%s:%d\n", file, line)
		if f.Name() == "main.main" {
			break
		}
	}
}

func main() {
	a()
}
