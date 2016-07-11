package main

import (
	&#34;bufio&#34;
	&#34;fmt&#34;
	&#34;io&#34;
	&#34;net/http&#34;
	&#34;os&#34;
	&#34;time&#34;
)

const chanbuf = 1000

func main() {
	var urls = make(chan string, chanbuf)
	var done = make(chan struct{}, chanbuf)
	var proc = 0
	if len(os.Args) != 2 {
		perror(&#34;Usage: &#34; &#43; os.Args[0] &#43; &#34; &lt;file name&gt;&#34;)
		os.Exit(1)
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		perror(err.Error())
		os.Exit(1)
	}
	defer f.Close()

	go read(urls, f)
	for url := range urls {
		proc&#43;&#43;
		go check(done, url)
	}
	for i := 0; i &lt; proc; i&#43;&#43; {
		&lt;-done
	}
}

func perror(msg string) {
	fmt.Fprintln(os.Stderr)
}

func read(out chan string, r io.Reader) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		out &lt;- scanner.Text()
	}
	close(out)
	if scanner.Err() != nil {
		perror(scanner.Err().Error())
	}
}

func check(done chan struct{}, url string) {
	resp, err := http.Head(url)
	if err != nil {
		perror(err.Error())
	}
	fmt.Print(fmt.Sprintf(&#34;%v:%s:%d\r\n&#34;, time.Now(), url, resp.StatusCode))
	done &lt;- struct{}{}
}