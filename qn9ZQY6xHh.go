package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/websocket"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", handlers.CompressHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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

	url := strings.Replace(server.URL, "http://", "ws://", 1)
	header := http.Header{"Accept-Encoding": []string{"gzip"}}

	conn, _, err := dialer.Dial(url, header)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
}
