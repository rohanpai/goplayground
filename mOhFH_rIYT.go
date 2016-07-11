package main

import (
	&#34;errors&#34;
	&#34;fmt&#34;
	&#34;io&#34;
	&#34;time&#34;
)

func main() {
	Pipe()
}

func Pipe() {
	pipeReader, pipeWriter := io.Pipe()
	go PipeWrite(pipeWriter)
	go PipeRead(pipeReader)
	time.Sleep(1e7)
}

func PipeWrite(pipeWriter *io.PipeWriter) {
	// close the pipe
	pipeWriter.CloseWithError(errors.New(&#34;close...&#34;))
	n, err := pipeWriter.Write([]byte(&#34;http://studygolang.com&#34;))
	fmt.Println(&#34;The number of bytes after writing: &#34;, n, &#34;. The error：&#34;,  err)
}

func PipeRead(pipeReader *io.PipeReader) {
	var (
		err error
		n   int
	)
	data := make([]byte, 1024)
	for n, err = pipeReader.Read(data); err == nil; n, err = pipeReader.Read(data) {
		fmt.Printf(&#34;%s\n&#34;, data[:n])
	}
	fmt.Printf(&#34;The data：%#v\n&#34;, data[:n])
	fmt.Println(&#34;The error：&#34;, err)
}