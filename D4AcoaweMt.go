package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func RunTraceroute(host string) {
	errch := make(chan error, 1)
	cmd := exec.Command("/usr/bin/traceroute", host)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	go func() {
		errch <- cmd.Wait()
	}()
	select {
	case <-time.After(time.Second * 2):
		log.Println("Timeout hit..")
		return
	case err := <-errch:
		if err != nil {
			log.Println("traceroute failed:", err)
		}
	default:
		for _, char := range "|/-\\" {
			fmt.Printf("\r%s...%c", "Running traceroute", char)
			time.Sleep(100 * time.Millisecond)
		}
		scanner := bufio.NewScanner(stdout)
		fmt.Println("")
		for scanner.Scan() {
			line := scanner.Text()
			log.Println(line)
		}
	}
}

func main() {
	RunTraceroute(os.Args[1])
}
