package bot

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"sync"
	"syscall"
	"time"
)

func logRecover() {
	if x := recover(); x != nil {
		log.Printf("runtime panic: %v", x)
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
	return &Bot{
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
	return d + time.Duration(r.Int63n(2*int64(d)))
}

func (b *Bot) send() {
	defer logRecover()
	random := rand.New(rand.NewSource(0))
	s := []byte(`{"msg":"o hai everybody!"}`)
	for sent := 0; b.MsgCount < 1 || sent < b.MsgCount; sent++ {
		//b.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		n, err := b.conn.Write(s)
		if err != nil {
			log.Printf("send error: %v", err)
			break
		}
		if n != len(s) {
			log.Printf("short write!")
			break
		}
		if b.Debug {
			log.Println("Sent message")
		}
		time.Sleep(durationJitter(b.sleep, random))
	}
}

func (b *Bot) recv() {
	buf := make([]byte, 4096)
	for {
		n, err := b.conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			if e, ok := err.(*net.OpError); ok {
				if e.Temporary() || e.Timeout() {
					continue
				} else if e.Err.Error() == "use of closed network connection" {
					break
				}
			}
			if err.Error() != "use of closed network connection" {
				log.Printf("read error: %v", err)
			}
			break
		}
		if b.Debug {
			log.Printf("Read %d bytes: %s\n", n, buf)
		}
	}
}

func (b *Bot) Run(laddr *net.TCPAddr) {
	defer logRecover()

	config, err := websocket.NewConfig(fmt.Sprintf("ws://%s:%d/chat?nick=%s&city=%d", b.server, b.port, b.nick, b.city),
		fmt.Sprintf("http://%s:%d/", b.server, b.port))
	if err != nil {
		panic(fmt.Sprintf("%s NewConfig failed: %v", b.nick, err))
	}

	tcpAddr := net.TCPAddr{
		IP:   b.addr.IP,
		Port: b.port,
	}

	var sock *net.TCPConn

	for i := 0; i < 10; i++ {
		sock, err = net.DialTCP("tcp4", laddr, &tcpAddr)
		if err == nil {
			break
		}
		switch e := err.(type) {
		case *net.OpError:
			if e.Err == syscall.EADDRNOTAVAIL {
				continue
			}
		}
		panic(fmt.Sprintf("%s dial error: %v", b.nick, err))
	}
	if err != nil {
		panic(fmt.Sprintf("%s dial error: %v", b.nick, err))
	}

	b.conn, err = websocket.NewClient(config, sock)
	if err != nil {
		sock.Close()
		panic(fmt.Sprintf("%s NewClient error: %v", b.nick, err))
	}

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
