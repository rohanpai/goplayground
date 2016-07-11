package main

import (
	&#34;os&#34;
)

var hint = map[int]string{
	1: &#34;semicolons http://golang.org/doc/go_spec.html#Semicolons&#34;,
	2: &#34;multi line string use backquote&#34; &#43;
		&#34;\n  `line 1&#34; &#43;
		&#34;\n  line 2`&#34; &#43;
		&#34;\n  for string include backquote use strings.Replace(`.X..`,\&#34;X\&#34;,\&#34;&#39;\&#34;,-1)&#34; &#43;
		&#34;\n  or use \&#34;line 1\&#34;&#43;&#34; &#43;
		&#34;\n         \&#34;line 2\&#34;&#34;,
}
var dic = map[string]string{
	&#34;auto&#34;:   &#34;X&#34;,
	&#34;const&#34;:  &#34;only for define constant variable in global/function&#34;,
	&#34;double&#34;: &#34;float64&#34;,
}

func main() {
	if len(os.Args) == 1 {
		println(`Usage: c2go [options] &lt;keyword_in_c/c&#43;&#43;&gt;
This utility is to let you quick map concept of C/C&#43;&#43; to Go.
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
		println(s, &#34;-&gt;&#34;, dic[s])
	}
}
