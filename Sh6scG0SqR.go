package main

import (
	"fmt"
	"net"
)

func whois(dom, server string) string {
	conn, err := net.Dial("tcp", server+":43")
	if err != nil {
		fmt.Println("Error")
	}
	conn.Write([]byte(dom + "\r\n"))
	buf := make([]byte, 1024)
	res := []byte{}
	for {
		numbytes, err := conn.Read(buf)
		sbuf := buf[0:numbytes]
		res = append(res, sbuf...)
		if err != nil {
			break
		}
	}
	conn.Close()
	return string(res)
}

func main() {
	fmt.Println(whois("hello.com", "com.whois-servers.net"))
}
