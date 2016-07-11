package main

import (
	&#34;log&#34;
	&#34;net/http&#34;
	&#34;net/http/httptest&#34;
	&#34;strings&#34;

	&#34;github.com/gorilla/handlers&#34;
	&#34;github.com/gorilla/websocket&#34;
)

func main() {
	mux := http.NewServeMux()
	mux.Handle(&#34;/&#34;, handlers.CompressHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()
	})))

	server := httptest.NewServer(mux)
	defer server.Close()

	dialer := websocket.Dialer{
		Subprotocols:    []string{},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	url := strings.Replace(server.URL, &#34;http://&#34;, &#34;ws://&#34;, 1)
	header := http.Header{&#34;Accept-Encoding&#34;: []string{&#34;gzip&#34;}}

	conn, _, err := dialer.Dial(url, header)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
}
