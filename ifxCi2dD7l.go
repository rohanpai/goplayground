package main

import "net/http"
import "fmt"
import "io"
import "strings"
import "time"
import "os"
import "net"
import "bufio"

var _=os.Stdout
var _=time.Second
var _=bufio.ErrBufferFull
var _ net.Conn

var indexHtml=`<!doctype html>
<html>
	<head>
		<script>
		var source=new EventSource("/events");
		source.addEventListener('message', function(e) {
			document.body.innerHTML+=e.data+"<br>";
		}, false);
		source.addEventListener('open', function(e) {
			console.log("Otwarto !");
		}, false);
		source.addEventListener('error', function(e) {
			console.log(e.target.readyState);
			if (e.target.readyState==EventSource.CLOSED) {
				console.log("Zamknieto !");
			}
		}, false);
		</script>
	</head>
	<body>
	</body>
</html>`

func h1(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/event-stream")
	for i:=0; ; i++ {
		if _, err:=fmt.Fprintf(resp, "data: %d\n\n", i); err != nil {
			fmt.Println("write: ", err)
		}
		flush, _:=resp.(http.Flusher)
		flush.Flush()
		fmt.Println(i)
		time.Sleep(1000*time.Millisecond)
	}
}

func h2(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/event-stream")
	flush, _:=resp.(http.Flusher)
	flush.Flush()
	hj, _:=resp.(http.Hijacker)
	conn, _, _:=hj.Hijack()
	for i:=0; ; i++ {
		if _, err:=fmt.Fprintf(conn, "data: %d\n\n", i); err != nil {
			fmt.Println("write: ", err)
			conn.Close()
			return
		}
		fmt.Println(i)
		time.Sleep(1000*time.Millisecond)
	}
}

func main() {
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		io.Copy(resp, strings.NewReader(indexHtml));
	})
	http.HandleFunc("/events", h2)
	if err:=http.ListenAndServe("localhost:3999", nil); err != nil {
		fmt.Println("server: ", err);
	}


}