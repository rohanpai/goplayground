// sysctl -w fs.file-max=999999
// ulimit -n `cat /proc/sys/fs/file-max`

package main

import (
	&#34;log&#34;
	&#34;net&#34;
	&#34;os/exec&#34;
	&#34;runtime&#34;
	&#34;strings&#34;
	&#34;time&#34;
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
			iprange &lt;- ip.String()
		}
	}

	go status()
	todo = make(chan string, 150000)
	
	for s := 0; s &lt;= 60000; s&#43;&#43; {
		go worker()
	}

	var ipnext int64
	for {
		if ipnext &gt; 184549375 || ipnext &lt; 167772160 {
			ipnext = 167772160
		}

		ipnext&#43;&#43;
		next := inet_ntoa(ipnext)
		todo &lt;- next.String()
	}
}

func status() {
	for {
		//netstat -n | awk &#39;/^tcp/ {t[$NF]&#43;&#43;}END{for(state in t){print state, t[state]} }&#39;
		out, _ := exec.Command(&#34;/tmp/count.sh&#34;).Output()
		log.Println(&#34;GO:&#34;, runtime.NumGoroutine(), &#34;|| TCP:&#34;, strings.Replace(string(out), &#34;\n&#34;, &#34; &#34;, -1), &#34; || ERROR:&#34;, lasterror)
		time.Sleep(5 * time.Second)
	}
}

func worker() {
	for {
		connect(&lt;-todo)
	}
}

func connect(ipaddr string) {
	var err error
	var ipNext string

	RETRY:
	
	ipNext = &lt;-iprange
	iprange &lt;- ipNext

	d := net.Dialer{Timeout: 1 * time.Second, Deadline: time.Now().Add(2 * time.Second)}
	d.LocalAddr, err = net.ResolveTCPAddr(&#34;tcp4&#34;, ipNext&#43;&#34;:0&#34;)
	if err != nil {
		lasterror = err.Error()
		return
	}
	dconn, err := d.Dial(&#34;tcp&#34;, ipaddr&#43;&#34;:80&#34;)

	if err != nil {
		if strings.Contains(err.Error(), &#34;address already in use&#34;) {
			lasterror = err.Error()
			goto RETRY
		}
		return
	}
	defer dconn.Close()
}

func inet_ntoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr &amp; 0xFF)
	bytes[1] = byte((ipnr &gt;&gt; 8) &amp; 0xFF)
	bytes[2] = byte((ipnr &gt;&gt; 16) &amp; 0xFF)
	bytes[3] = byte((ipnr &gt;&gt; 24) &amp; 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}
