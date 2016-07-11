package main

import (
	&#34;bufio&#34;
	&#34;fmt&#34;
	&#34;log&#34;
	&#34;os&#34;
	&#34;os/exec&#34;
	&#34;time&#34;
)

func RunTraceroute(host string) {
	errch := make(chan error, 1)
	cmd := exec.Command(&#34;/usr/bin/traceroute&#34;, host)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	go func() {
		errch &lt;- cmd.Wait()
	}()
	select {
	case &lt;-time.After(time.Second * 2):
		log.Println(&#34;Timeout hit..&#34;)
		return
	case err := &lt;-errch:
		if err != nil {
			log.Println(&#34;traceroute failed:&#34;, err)
		}
	default:
		for _, char := range &#34;|/-\\&#34; {
			fmt.Printf(&#34;\r%s...%c&#34;, &#34;Running traceroute&#34;, char)
			time.Sleep(100 * time.Millisecond)
		}
		scanner := bufio.NewScanner(stdout)
		fmt.Println(&#34;&#34;)
		for scanner.Scan() {
			line := scanner.Text()
			log.Println(line)
		}
	}
}

func main() {
	RunTraceroute(os.Args[1])
}
