package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"strings"
)

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}

func main() {
	// Set up authentication information.

	smtpServer := "SERVER HERE"
	auth := smtp.PlainAuth(
		"",
		"EMAIL HERE",
		"PWD HERE",
		smtpServer,
	)

	from := mail.Address{"NAME", "EMAIL"}
	to := mail.Address{"NAME", "EMAIL"}
	title := "Hello"

	body := "This is a message"

	header := make(map[string]string)
	header["Return-Path"] = "EMAIL"
	header["From"] = "EMAIL"
	header["To"] = to.String()
	header["Subject"] = encodeRFC2047(title)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		smtpServer+":25",
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
	)
	if err != nil {
		log.Fatal(err)
	}
}
