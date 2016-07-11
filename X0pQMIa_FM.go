package main

import (
	&#34;code.google.com/p/go.net/websocket&#34;
	&#34;encoding/base64&#34;
	&#34;flag&#34;
	&#34;log&#34;
	&#34;net&#34;
	&#34;net/http&#34;
)

var (
	listen = flag.String(&#34;listen&#34;, &#34;:6080&#34;, &#34;Location to listen for connections&#34;)
)

func main() {
	flag.Parse()
	var wsConfig *websocket.Config
	var err error
	if wsConfig, err = websocket.NewConfig(&#34;ws://127.0.0.1:6080/&#34;, &#34;http://127.0.0.1:6080&#34;); err != nil {
		log.Fatalf(err.Error())
		return
	}

	// wsConfig.Protocol = []string{&#34;base64&#34;}
	http.Handle(&#34;/websockify&#34;, websocket.Server{Handler: wsh,
		Config: *wsConfig,
		Handshake: func(ws *websocket.Config, req *http.Request) error {
			ws.Protocol = []string{&#34;base64&#34;}
			return nil
		}})
	http.Handle(&#34;/novnc/&#34;, http.StripPrefix(&#34;/novnc/&#34;, http.FileServer(http.Dir(&#34;./novnc/&#34;))))
	log.Fatal(http.ListenAndServe(*listen, nil))
}

func wsh(ws *websocket.Conn) {
	loc := &#34;127.0.0.1:5901&#34;
	vc, err := net.Dial(&#34;tcp&#34;, loc)
	defer vc.Close()
	if err != nil {
		log.Print(err)
		return
	}
	go func() {
		sbuf := make([]byte, 32*1024)
		dbuf := make([]byte, 32*1024)
		for {
			n, e := ws.Read(sbuf)
			if e != nil {
				return
			}
			n, e = base64.StdEncoding.Decode(dbuf, sbuf[0:n])
			if e != nil {
				return
			}
			n, e = vc.Write(dbuf[0:n])
			if e != nil {
				return
			}
		}
	}()
	go func() {
		sbuf := make([]byte, 32*1024)
		dbuf := make([]byte, 64*1024)
		for {
			n, e := vc.Read(sbuf)
			if e != nil {
				return
			}
			base64.StdEncoding.Encode(dbuf, sbuf[0:n])
			n = ((n &#43; 2) / 3) * 4
			ws.Write(dbuf[0:n])
			if e != nil {
				return
			}
		}
	}()
	select {}
}
