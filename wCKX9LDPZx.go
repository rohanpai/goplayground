package bot

import (
	&#34;code.google.com/p/go.net/websocket&#34;
	&#34;fmt&#34;
	&#34;io&#34;
	&#34;log&#34;
	&#34;math/rand&#34;
	&#34;net&#34;
	&#34;sync&#34;
	&#34;syscall&#34;
	&#34;time&#34;
)

func logRecover() {
	if x := recover(); x != nil {
		log.Printf(&#34;runtime panic: %v&#34;, x)
	}
}

type Bot struct {
	nick     string
	city     int
	server   string
	port     int
	addr     *net.IPAddr
	conn     *websocket.Conn
	sleep    time.Duration
	MsgCount int
	Debug    bool
}

func NewBot(nick string, city int, server string, port int, addr *net.IPAddr, sleep time.Duration) *Bot {
	return &amp;Bot{
		nick:   nick,
		city:   city,
		server: server,
		port:   port,
		addr:   addr,
		sleep:  sleep,
	}
}

func durationJitter(d time.Duration, r *rand.Rand) time.Duration {
	if d == 0 {
		return 0
	}
	return d &#43; time.Duration(r.Int63n(2*int64(d)))
}

func (b *Bot) send() {
	defer logRecover()
	random := rand.New(rand.NewSource(0))
	s := []byte(`{&#34;msg&#34;:&#34;o hai everybody!&#34;}`)
	for sent := 0; b.MsgCount &lt; 1 || sent &lt; b.MsgCount; sent&#43;&#43; {
		// SetWriteDeadline does not cause the hanging Write() calls to timeout properly
		//b.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		n, err := b.conn.Write(s)
		if err != nil {
			log.Printf(&#34;send error: %v&#34;, err)
			break
		}
		if n != len(s) {
			log.Printf(&#34;short write!&#34;)
			break
		}
		if b.Debug {
			log.Println(&#34;Sent message&#34;)
		}
		time.Sleep(durationJitter(b.sleep, random))
	}
}

func (b *Bot) recv() {
	buf := make([]byte, 4096)
	for {
		n, err := b.conn.Read(buf)
		if err != nil {
			if err == io.EOF { // this happens occasionally
				break
			}
			if e, ok := err.(*net.OpError); ok {
				if e.Temporary() || e.Timeout() {
					// I don&#39;t think these actually happen, but we would want to continue if they did...
					continue
				} else if e.Err.Error() == &#34;use of closed network connection&#34; { // happens very frequently
					break
				}
			}
			if err.Error() != &#34;use of closed network connection&#34; {
				// not sure why this is different from the above similar case, but this one happens, also.
				log.Printf(&#34;read error: %v&#34;, err)
			}
			break
		}
		if b.Debug {
			log.Printf(&#34;Read %d bytes: %s\n&#34;, n, buf)
		}
	}
}

func (b *Bot) Run(laddr *net.TCPAddr) {
	defer logRecover()

	// can&#39;t use websocket.Dial(...) because it exhausts the OS limit on threads while resolving
	// the hostname using cgo. Instead, we look it up ahead of time and store in b.addr so that
	// we can connect with net.DialTCP on the pre-resolved addr
	config, err := websocket.NewConfig(fmt.Sprintf(&#34;ws://%s:%d/chat?nick=%s&amp;city=%d&#34;, b.server, b.port, b.nick, b.city),
		fmt.Sprintf(&#34;http://%s:%d/&#34;, b.server, b.port))
	if err != nil {
		panic(fmt.Sprintf(&#34;%s NewConfig failed: %v&#34;, b.nick, err))
	}

	tcpAddr := net.TCPAddr{
		IP:   b.addr.IP,
		Port: b.port,
	}

	var sock *net.TCPConn

	for i := 0; i &lt; 10; i&#43;&#43; {
		sock, err = net.DialTCP(&#34;tcp4&#34;, laddr, &amp;tcpAddr)
		if err == nil {
			break
		}
		switch e := err.(type) {
		case *net.OpError:
			if e.Err == syscall.EADDRNOTAVAIL {
				// I think this happens because of a race between bind() and connect()
				// Try again and hope we don&#39;t see the same problem
				continue
			}
		}
		panic(fmt.Sprintf(&#34;%s dial error: %v&#34;, b.nick, err))
	}
	if err != nil {
		panic(fmt.Sprintf(&#34;%s dial error: %v&#34;, b.nick, err))
	}

	b.conn, err = websocket.NewClient(config, sock)
	if err != nil {
		sock.Close()
		panic(fmt.Sprintf(&#34;%s NewClient error: %v&#34;, b.nick, err))
	}

	// Probably overkill to use a WaitGroup for just one goroutine.
	// Maybe it would be better to let the caller pass a WaitGroup.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer b.conn.Close()
		b.send()
	}()
	b.recv()
	wg.Wait()
}
