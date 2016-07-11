// Sample program that adds a few more features. If -z is passed, we want any
// DestFile&#39;s to be gzipped. If -md5 is passed, we want print the md5sum of the
// data that&#39;s been transfered instead of the data itself.
package main

import (
	&#34;compress/gzip&#34;
	&#34;crypto/md5&#34;
	&#34;flag&#34;
	&#34;fmt&#34;
	&#34;io&#34;
	&#34;net/http&#34;
	&#34;os&#34;
)

// Config contains program configuration options.
var Config struct {
	Silent   bool
	DestFile string
	Gzip     bool
	Md5      bool
}

// init is called before main.
func init() {
	flag.StringVar(&amp;Config.DestFile, &#34;o&#34;, &#34;&#34;, &#34;output file&#34;)
	flag.BoolVar(&amp;Config.Silent, &#34;s&#34;, false, &#34;silent (do not output to stdout)&#34;)
	flag.BoolVar(&amp;Config.Gzip, &#34;z&#34;, false, &#34;gzip file output&#34;)
	flag.BoolVar(&amp;Config.Md5, &#34;md5&#34;, false, &#34;stdout md5sum instead of body&#34;)
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println(&#34;Usage: ./example6 [options] &lt;url&gt;&#34;)
		os.Exit(-1)
	}
}

// main is the entry point for the application.
func main() {
	// Capture the url from the arguments.
	url := flag.Args()[0]

	// r here is a response, and r.Body is an io.Reader.
	r, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Our Md5 hash destination, which is an io.Writer that computes the
	// hash of whatever is written to it.
	hash := md5.New()
	var writers []io.Writer

	// If we aren&#39;t in Silent mode, lets add Stdout to our writers.
	if !Config.Silent {
		// If -md5 was passed, write to the hash instead of os.Stdout.
		if Config.Md5 {
			writers = append(writers, hash)
		} else {
			writers = append(writers, os.Stdout)
		}
	}

	// If DestFile was provided, lets try to create it and add to the writers.
	if len(Config.DestFile) &gt; 0 {
		// By declaring a Writer here as a WriteCloser, we&#39;re saying that we don&#39;t care
		// what the underlying implementation is at all, all we require is something that
		// can Write and Close;  both os.File and the gzip.Writer are WriteClosers.
		var writer io.WriteCloser
		writer, err := os.Create(Config.DestFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		// If we&#39;re in Gzip mode, wrap the writer in gzip
		if Config.Gzip {
			writer = gzip.NewWriter(writer)
		}

		writers = append(writers, writer)
		defer writer.Close()
	}

	// MultiWriter(io.Writer...) returns a single writer which multiplexes its
	// writes across all of the writers we pass in.
	dest := io.MultiWriter(writers...)

	// Write to dest the same way as before, copying from the Body.
	io.Copy(dest, r.Body)
	if err = r.Body.Close(); err != nil {
		fmt.Println(err)
		return
	}

	// If we were in Md5 output mode, lets output the checksum and url.
	if Config.Md5 {
		fmt.Printf(&#34;%x  %s\n&#34;, hash.Sum(nil), url)
	}
}
