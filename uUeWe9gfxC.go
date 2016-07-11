package main

import (
	"os"
)

var hint = map[int]string{
	1: "semicolons http://golang.org/doc/go_spec.html#Semicolons",
	2: "multi line string use backquote" +
		"\n  `line 1" +
		"\n  line 2`" +
		"\n  for string include backquote use strings.Replace(`.X..`,\"X\",\"'\",-1)" +
		"\n  or use \"line 1\"+" +
		"\n         \"line 2\"",
}
var dic = map[string]string{
	"auto":   "X",
	"const":  "only for define constant variable in global/function",
	"double": "float64",
}

func main() {
	if len(os.Args) == 1 {
		println(`Usage: c2go [options] <keyword_in_c/c++>
This utility is to let you quick map concept of C/C++ to Go.
eg.
  c2go sprintf explicit
  c2go -t # vaildate all usage
  c2go -n # list syntax/concept differences
  c2go 12 # list No.12 hint
Reference:
  Go syntax          http://golang.org/doc/go_spec.html
  Go packages        http://golang.org/pkg/
  3rd party packages
		http://godashboard.appspot.com/package
		http://godashboard.appspot.com/project
`)
		os.Exit(1)
	}
	for _, s := range os.Args[1:] {
		println(s, "->", dic[s])
	}
}
