package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const chanbuf = 1000

func main() {
	var urls = make(chan string, chanbuf)
	var done = make(chan struct{}, chanbuf)
	var proc = 0
	if len(os.Args) != 2 {
		perror("Usage: " + os.Args[0] + " <file name>")
		os.Exit(1)
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		perror(err.Error())
		os.Exit(1)
	}
	defer f.Close()

	go read(urls, f)
	for url := range urls {
		proc++
		go check(done, url)
	}
	for i := 0; i < proc; i++ {
		<-done
	}
}

func perror(msg string) {
	fmt.Fprintln(os.Stderr)
}

func read(out chan string, r io.Reader) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		out <- scanner.Text()
	}
	close(out)
	if scanner.Err() != nil {
		perror(scanner.Err().Error())
	}
}

func check(done chan struct{}, url string) {
	resp, err := http.Head(url)
	if err != nil {
		perror(err.Error())
	}
	fmt.Print(fmt.Sprintf("%v:%s:%d\r\n", time.Now(), url, resp.StatusCode))
	done <- struct{}{}
}