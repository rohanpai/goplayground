package main

import (
	"bytes"
	"fmt"

	"code.google.com/p/go.crypto/ssh"
)

func getPass() (string, error) {
	return "super secret", nil
}

func main() {
	config := &ssh.ClientConfig{
		User: "me",
		Auth: []ssh.AuthMethod{
			ssh.PasswordCallback(getPass),
		},
	}
	client, err := ssh.Dial("tcp", "me.example.com:22", config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("uptime"); err != nil {
		panic("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())
}
