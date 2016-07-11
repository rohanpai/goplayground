package main

import (
	&#34;bytes&#34;
	&#34;fmt&#34;

	&#34;code.google.com/p/go.crypto/ssh&#34;
)

func getPass() (string, error) {
	return &#34;super secret&#34;, nil
}

func main() {
	config := &amp;ssh.ClientConfig{
		User: &#34;me&#34;,
		Auth: []ssh.AuthMethod{
			ssh.PasswordCallback(getPass),
		},
	}
	client, err := ssh.Dial(&#34;tcp&#34;, &#34;me.example.com:22&#34;, config)
	if err != nil {
		panic(&#34;Failed to dial: &#34; &#43; err.Error())
	}
	session, err := client.NewSession()
	if err != nil {
		panic(&#34;Failed to create session: &#34; &#43; err.Error())
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &amp;b
	if err := session.Run(&#34;uptime&#34;); err != nil {
		panic(&#34;Failed to run: &#34; &#43; err.Error())
	}
	fmt.Println(b.String())
}
