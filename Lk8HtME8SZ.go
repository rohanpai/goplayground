// download data from ftp site
package main

import (
	&#34;flag&#34;
	&#34;fmt&#34;
	//ftp &#34;github.com/jlaffaye/goftp&#34;
	//ftp &#34;github.com/jawr/ftp.go&#34;
	//&#34;bitbucket.org/zombiezen/ftp&#34;
	ftp &#34;github.com/jum/tinyftp&#34;
	&#34;log&#34;
	&#34;net&#34;
	&#34;os&#34;
	&#34;strings&#34;
)

var verbose = flag.Bool(&#34;v&#34;, false, &#34;verbose&#34;)

func main() {
	flag.Usage = func() {
		fmt.Println(&#34;ftputl &lt;host[:port]&gt; &lt;user&gt; &lt;pass&gt; &lt;dir&gt; &lt;file&gt; [dst]&#34;)
	}
	flag.Parse()
	if flag.NArg() != 6 &amp;&amp; flag.NArg() != 5 {
		flag.Usage()
		flag.PrintDefaults()
		os.Exit(1)
	}
	host, user, pass, dir, file := flag.Arg(0), flag.Arg(1),
		flag.Arg(2), flag.Arg(3), flag.Arg(4)

	if !strings.Contains(host, &#34;:&#34;) {
		host &#43;= &#34;:21&#34;
	}
	dst := file
	if flag.NArg() == 6 {
		dst = flag.Arg(5)
	}

	if *verbose {
		log.Println(&#34;ftp connect&#34;, host)
	}
	c, code, msg, err := ftp.Dial(&#34;tcp&#34;, host)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(code, msg)

	if *verbose {
		log.Println(&#34;ftp login&#34;, host)
	}
	code, msg, err = c.Login(user, pass)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(code, msg)

	if *verbose {
		log.Println(&#34;ftp cd&#34;, dir)
	}
	code, msg, err = c.Cwd(dir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(code, msg)

	addr, code, msg, err := c.Passive()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(addr, code, msg)

	dconn, err := net.Dial(&#34;tcp&#34;, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer dconn.Close()

	if *verbose {
		log.Println(&#34;ftp type I&#34;)
	}
	code, msg, err = c.Type(&#34;I&#34;)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(code, msg)

	if *verbose {
		log.Println(&#34;creating&#34;, dst)
	}
	w, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}

	n, code, msg, err := c.Size(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(code, msg, n)

	if *verbose {
		log.Println(&#34;get&#34;, file)
	}
	n, code, msg, err = c.RetrieveTo(file, dconn, w)
	if err != nil {
		log.Fatal(err)
	}
	w.Close()
	fmt.Println(&#34;done, copy bytes:&#34;, n)
}
