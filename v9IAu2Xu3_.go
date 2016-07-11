package main

import (
	&#34;fmt&#34;
	&#34;io&#34;
	&#34;net&#34;
	&#34;net/http&#34;
	&#34;time&#34;
)

func Index(w http.ResponseWriter, r *http.Request) {
	url := &#34;http://upload.wikimedia.org/wikipedia/en/b/bc/Wiki.png&#34;

	timeout := time.Duration(5) * time.Second
	transport := &amp;http.Transport{
		ResponseHeaderTimeout: timeout,
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		},
		DisableKeepAlives: true,
	}
	client := &amp;http.Client{
		Transport: transport,
	}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	//copy the relevant headers. If you want to preserve the downloaded file name, extract it with go&#39;s url parser.
	w.Header().Set(&#34;Content-Disposition&#34;, &#34;attachment; filename=Wiki.png&#34;)
	w.Header().Set(&#34;Content-Type&#34;, r.Header.Get(&#34;Content-Type&#34;))
	w.Header().Set(&#34;Content-Length&#34;, r.Header.Get(&#34;Content-Length&#34;))

	//stream the body to the client without fully loading it into memory
	io.Copy(w, resp.Body)
}

func main() {
	http.HandleFunc(&#34;/&#34;, Index)
	err := http.ListenAndServe(&#34;:8000&#34;, nil)

	if err != nil {
		fmt.Println(err)
	}
}
