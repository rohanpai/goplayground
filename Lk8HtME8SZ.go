// download data from ftp site
package main

import (
	"flag"
	"fmt"
	//ftp "github.com/jlaffaye/goftp"
	//ftp "github.com/jawr/ftp.go"
	//"bitbucket.org/zombiezen/ftp"
	ftp "github.com/jum/tinyftp"
	"log"
	"net"
	"os"
	"strings"
)

var verbose = flag.Bool("v", false, "verbose")

func main() {
	flag.Usage = func() {
		fmt.Println("ftputl <host[:port]> <user> <pass> <dir> <file> [dst]")
	}
	flag.Parse()
	if flag.NArg() != 6 && flag.NArg() != 5 {
		flag.Usage()
		flag.PrintDefaults()
		os.Exit(1)
	}
	host, user, pass, dir, file := flag.Arg(0), flag.Arg(1),
		flag.Arg(2), flag.Arg(3), flag.Arg(4)

	if !strings.Contains(host, ":") {
		host += ":21"
	}
	dst := file
	if flag.NArg() == 6 {
		dst = flag.Arg(5)
	}

	if *verbose {
		log.Println("ftp connect", host)
	}
	c, code, msg, err := ftp.Dial("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(code, msg)

	if *verbose {
		log.Println("ftp login", host)
	}
	code, msg, err = c.Login(user, pass)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(code, msg)

	if *verbose {
		log.Println("ftp cd", dir)
	}
	code, msg, err = c.Cwd(dir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(code, msg)

	addr, code, msg, err := c.Passive()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(addr, code, msg)

	dconn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer dconn.Close()

	if *verbose {
		log.Println("ftp type I")
	}
	code, msg, err = c.Type("I")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(code, msg)

	if *verbose {
		log.Println("creating", dst)
	}
	w, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}

	n, code, msg, err := c.Size(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(code, msg, n)

	if *verbose {
		log.Println("get", file)
	}
	n, code, msg, err = c.RetrieveTo(file, dconn, w)
	if err != nil {
		log.Fatal(err)
	}
	w.Close()
	fmt.Println("done, copy bytes:", n)
}
