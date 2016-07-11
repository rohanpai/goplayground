package main

import &#34;net/http&#34;
import &#34;fmt&#34;
import &#34;io&#34;
import &#34;strings&#34;
import &#34;time&#34;
import &#34;os&#34;
import &#34;net&#34;
import &#34;bufio&#34;

var _=os.Stdout
var _=time.Second
var _=bufio.ErrBufferFull
var _ net.Conn

var indexHtml=`&lt;!doctype html&gt;
&lt;html&gt;
	&lt;head&gt;
		&lt;script&gt;
		var source=new EventSource(&#34;/events&#34;);
		source.addEventListener(&#39;message&#39;, function(e) {
			document.body.innerHTML&#43;=e.data&#43;&#34;&lt;br&gt;&#34;;
		}, false);
		source.addEventListener(&#39;open&#39;, function(e) {
			console.log(&#34;Otwarto !&#34;);
		}, false);
		source.addEventListener(&#39;error&#39;, function(e) {
			console.log(e.target.readyState);
			if (e.target.readyState==EventSource.CLOSED) {
				console.log(&#34;Zamknieto !&#34;);
			}
		}, false);
		&lt;/script&gt;
	&lt;/head&gt;
	&lt;body&gt;
	&lt;/body&gt;
&lt;/html&gt;`

func h1(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add(&#34;Content-Type&#34;, &#34;text/event-stream&#34;)
	for i:=0; ; i&#43;&#43; {
		if _, err:=fmt.Fprintf(resp, &#34;data: %d\n\n&#34;, i); err != nil {
			fmt.Println(&#34;write: &#34;, err)
		}
		flush, _:=resp.(http.Flusher)
		flush.Flush()
		fmt.Println(i)
		time.Sleep(1000*time.Millisecond)
	}
}

func h2(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add(&#34;Content-Type&#34;, &#34;text/event-stream&#34;)
	flush, _:=resp.(http.Flusher)
	flush.Flush()
	hj, _:=resp.(http.Hijacker)
	conn, _, _:=hj.Hijack()
	for i:=0; ; i&#43;&#43; {
		if _, err:=fmt.Fprintf(conn, &#34;data: %d\n\n&#34;, i); err != nil {
			fmt.Println(&#34;write: &#34;, err)
			conn.Close()
			return
		}
		fmt.Println(i)
		time.Sleep(1000*time.Millisecond)
	}
}

func main() {
	http.HandleFunc(&#34;/&#34;, func(resp http.ResponseWriter, req *http.Request) {
		io.Copy(resp, strings.NewReader(indexHtml));
	})
	http.HandleFunc(&#34;/events&#34;, h2)
	if err:=http.ListenAndServe(&#34;localhost:3999&#34;, nil); err != nil {
		fmt.Println(&#34;server: &#34;, err);
	}


}