package main

import (
	&#34;encoding/base64&#34;
	&#34;fmt&#34;
	&#34;log&#34;
	&#34;net/mail&#34;
	&#34;net/smtp&#34;
	&#34;strings&#34;
)

func encodeRFC2047(String string) string {
	// use mail&#39;s rfc2047 to encode any string
	addr := mail.Address{String, &#34;&#34;}
	return strings.Trim(addr.String(), &#34; &lt;&gt;&#34;)
}

func main() {
	// Set up authentication information.

	smtpServer := &#34;SERVER HERE&#34;
	auth := smtp.PlainAuth(
		&#34;&#34;,
		&#34;EMAIL HERE&#34;,
		&#34;PWD HERE&#34;,
		smtpServer,
	)

	from := mail.Address{&#34;NAME&#34;, &#34;EMAIL&#34;}
	to := mail.Address{&#34;NAME&#34;, &#34;EMAIL&#34;}
	title := &#34;Hello&#34;

	body := &#34;This is a message&#34;

	header := make(map[string]string)
	header[&#34;Return-Path&#34;] = &#34;EMAIL&#34;
	header[&#34;From&#34;] = &#34;EMAIL&#34;
	header[&#34;To&#34;] = to.String()
	header[&#34;Subject&#34;] = encodeRFC2047(title)
	header[&#34;MIME-Version&#34;] = &#34;1.0&#34;
	header[&#34;Content-Type&#34;] = &#34;text/plain; charset=\&#34;utf-8\&#34;&#34;
	header[&#34;Content-Transfer-Encoding&#34;] = &#34;base64&#34;

	message := &#34;&#34;
	for k, v := range header {
		message &#43;= fmt.Sprintf(&#34;%s: %s\r\n&#34;, k, v)
	}
	message &#43;= &#34;\r\n&#34; &#43; base64.StdEncoding.EncodeToString([]byte(body))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		smtpServer&#43;&#34;:25&#34;,
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
	)
	if err != nil {
		log.Fatal(err)
	}
}
