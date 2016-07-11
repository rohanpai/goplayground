package main

import (
	&#34;fmt&#34;
	&#34;net&#34;
)

func whois(dom, server string) string {
	conn, err := net.Dial(&#34;tcp&#34;, server&#43;&#34;:43&#34;)
	if err != nil {
		fmt.Println(&#34;Error&#34;)
	}
	conn.Write([]byte(dom &#43; &#34;\r\n&#34;))
	buf := make([]byte, 1024)
	res := []byte{}
	for {
		numbytes, err := conn.Read(buf)
		sbuf := buf[0:numbytes]
		res = append(res, sbuf...)
		if err != nil {
			break
		}
	}
	conn.Close()
	return string(res)
}

func main() {
	fmt.Println(whois(&#34;hello.com&#34;, &#34;com.whois-servers.net&#34;))
}
