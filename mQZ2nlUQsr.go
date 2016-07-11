// sysctl -w fs.file-max=999999
// ulimit -n `cat /proc/sys/fs/file-max`

package main

import (
	"log"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	iprange   chan string
	todo      chan string
	lasterror string
)

func main() {
	//runtime.GOMAXPROCS(30)
	ipAddrs, _ := net.InterfaceAddrs()
	iprange = make(chan string, len(ipAddrs))

	for i := range ipAddrs {
		ip, _, _ := net.ParseCIDR(ipAddrs[i].String())
		if ip.IsGlobalUnicast() {
			iprange <- ip.String()
		}
	}

	go status()
	todo = make(chan string, 150000)
	
	for s := 0; s <= 60000; s++ {
		go worker()
	}

	var ipnext int64
	for {
		if ipnext > 184549375 || ipnext < 167772160 {
			ipnext = 167772160
		}

		ipnext++
		next := inet_ntoa(ipnext)
		todo <- next.String()
	}
}

func status() {
	for {
		//netstat -n | awk '/^tcp/ {t[$NF]++}END{for(state in t){print state, t[state]} }'
		out, _ := exec.Command("/tmp/count.sh").Output()
		log.Println("GO:", runtime.NumGoroutine(), "|| TCP:", strings.Replace(string(out), "\n", " ", -1), " || ERROR:", lasterror)
		time.Sleep(5 * time.Second)
	}
}

func worker() {
	for {
		connect(<-todo)
	}
}

func connect(ipaddr string) {
	var err error
	var ipNext string

	RETRY:
	
	ipNext = <-iprange
	iprange <- ipNext

	d := net.Dialer{Timeout: 1 * time.Second, Deadline: time.Now().Add(2 * time.Second)}
	d.LocalAddr, err = net.ResolveTCPAddr("tcp4", ipNext+":0")
	if err != nil {
		lasterror = err.Error()
		return
	}
	dconn, err := d.Dial("tcp", ipaddr+":80")

	if err != nil {
		if strings.Contains(err.Error(), "address already in use") {
			lasterror = err.Error()
			goto RETRY
		}
		return
	}
	defer dconn.Close()
}

func inet_ntoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}
