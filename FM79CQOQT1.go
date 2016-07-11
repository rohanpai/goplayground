//
// direct test of secure vs non-secure websockets, wss and ws uri schemes respectively
//
// to run this test:
//
// 1. "go run ws.go"
// 2. in a browser, visit unsecure url "http://localhost:8080" ---
//	this is insecure version, you should see echo messages once per second
// 3. in a browser, visit secure url "https://localhost:8090" ---
//	this is secure version, depending on browser type, you should see same output as insecure version, or error; note
//      that you will have to accept the untrusted certificate below and/or security exception manually when prompted.
//
// in particular, wss does not work in my chrome 28.0.1500.95 (it silently fails on both browser and server sides),
//      but works in "Opera 12.16 Build 1860 for Linux x86_64" and also works in mozilla firefox 23.0.
//
// unsecure version works fine in all three browsers i tested.
//
package main

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		switch r.URL.Path {
		case "/":
			u := func() string {
				if r.TLS == nil {
					return "ws://localhost:8080/sock"
				} else {
					return "wss://localhost:8090/sock"
				}
			}()
			http.ServeContent(w, r, "index.html", time.Now(), bytes.NewReader([]byte(Content(u))))
		case "/sock":
			websocket.Handler(EchoHandler).ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	done := make(chan bool)

	go func() {
		fmt.Println(http.ListenAndServe(":8080", nil))
		done <- true
	}()

	go func() {
		ioutil.WriteFile(CERT_NAME, []byte(CERT), os.ModePerm)
		ioutil.WriteFile(KEY_NAME, []byte(KEY), os.ModePerm)
		fmt.Println(http.ListenAndServeTLS(":8090", CERT_NAME, KEY_NAME, nil))
		done <- true
	}()

	<-done

}

func EchoHandler(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

func Content(ws string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang='en'>
  <head>
    <meta charset='utf-8'/>
    <title>websocket test</title>
    <script src="//ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script>  
  </head>
  <body>
    <p> go-lang / chrome websocket test at %s: </p>
    <ol/>
    <script>

var log = function(e) {
    $('ol').append("<li>"+JSON.stringify(e)+"</li>");
}

var sock = new WebSocket("%s")

sock.onerror = function(e) { log({ERROR:e}) };    
sock.onclose = function(e) { log({CLOSE:e}) };    
sock.onmessage = function(e) { log({MESSAGE:e.data}) };    

sock.onopen = function(e) { 
    log({OPEN:e}) 
    setInterval(function() {
	sock.send("howdy at " + new Date());
    },1000);
};    

    </script>
  </body>
</html>
`, ws, ws)
}

const (
	CERT_NAME = "/tmp/cert.prm"
	CERT      = `-----BEGIN CERTIFICATE-----
MIIC4zCCAc2gAwIBAgIBADALBgkqhkiG9w0BAQUwEjEQMA4GA1UEChMHQWNtZSBD
bzAeFw0xMzA4MTIxOTUzNTdaFw0xNDA4MTIxOTUzNTdaMBIxEDAOBgNVBAoTB0Fj
bWUgQ28wggEgMAsGCSqGSIb3DQEBAQOCAQ8AMIIBCgKCAQEAw917qyvQTFtr65LQ
MimR7obQzOc05YvqH/ytdtqSyPqQsvnOLi0Jk2/t33RpqIGb36dIMGDV3X5SpiLW
9TReidxtz22YbZ1CTRUgdcF5XJiq3CO5a+1vGvi/fIBIqzPb6CQwgG/eFm6xdvgf
IRodZgM6ym6cglm6ndhoB4/TpIni/i+bo757LUxfu78/mkUBfsQK0VDaqq59ZkZW
An8E/4DBeDdg8jq9i+4VuOTSulfiLu06eIWoAhfBeWX6uGaEK44xW5NJpp1y4CKX
BOvVKuz+Dv4jz/vp1e8p1zFqkUzbrtdVv+Z3olfTHQWe1qMpt9l01tedxuPTEHas
HTQiRwIDAQABo0owSDAOBgNVHQ8BAf8EBAMCAKAwEwYDVR0lBAwwCgYIKwYBBQUH
AwEwDAYDVR0TAQH/BAIwADATBgNVHREEDDAKggh4dXZhLmNvbTALBgkqhkiG9w0B
AQUDggEBAFppIRxB25BRCK1w//9U73sddEnw2Q/MUR+V4twB5qVGnnY6VJW+u9W3
9dNs4teemAxUJh7ZBOR2xEj1N6Q9D0H3BFupUNRa9jIxWam+E+mNCF1jQRlTxCur
WHk56ZmsqWyb8yXIB3ymHyAXnJziAUO/US6e8xeMKcvIZCtCzCsaOLI5G35h9xoJ
+mKedqOW7tGmT3suqj/bx2hEq1nPRW/H9XmyLkDMvIiUaU8iopKzShuyG/kosRXY
WsoDWrkv3a0O1crnlQLZXg1IH0/bQg7OwJz5P8qDSNR1uiBzGkmbOFpp0IVCdVIk
UwKAnERUK1MkYu6jh5hXj6mb0fW01/o=
-----END CERTIFICATE-----`
	KEY_NAME = "/tmp/key.pem"
	KEY      = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAw917qyvQTFtr65LQMimR7obQzOc05YvqH/ytdtqSyPqQsvnO
Li0Jk2/t33RpqIGb36dIMGDV3X5SpiLW9TReidxtz22YbZ1CTRUgdcF5XJiq3CO5
a+1vGvi/fIBIqzPb6CQwgG/eFm6xdvgfIRodZgM6ym6cglm6ndhoB4/TpIni/i+b
o757LUxfu78/mkUBfsQK0VDaqq59ZkZWAn8E/4DBeDdg8jq9i+4VuOTSulfiLu06
eIWoAhfBeWX6uGaEK44xW5NJpp1y4CKXBOvVKuz+Dv4jz/vp1e8p1zFqkUzbrtdV
v+Z3olfTHQWe1qMpt9l01tedxuPTEHasHTQiRwIDAQABAoIBAA22+4rf1YUTPbpQ
HG32xTYrkIFYizarlmhI/Ch/Y5nZGbq+jTZkhvAg/UoRT7ix4qVFhGOG1FLfHpBt
jhm7YgdLPREyPmMmiNb27L/yHTpjoksp4Tjydj4wPtBL90qtpe9aYV8M9kMh2yFW
fG+H8ZkMDtjP5/ukptGYrqgg5RP3SGacLGdmdCakesCP/+bjSjKVyHsBThybr5oy
zSsLaGocU6JL7k6x1lbnbToSwQo5ORRgWzsVIc68qm4Nc4+Nn9lrhMqlucYMtrxf
EyNpcCUFuO1BMz1jJwP5zttY0ZRKWitaDRpVyZhKMonUSsYIfRV+kM9+qKatZGbW
Tuh9dLkCgYEA3pSRWsyCHxmVDXpTAw0EPSZISaAqq5RSZoBtJkKBGpqxs245RI4E
FQN2d3oPpI3qBf3J6+lLRy4W0Hju+UxSzqdCHONmngdxbzOuMdoianjwviwr+ZjM
e8kmFtQ72KK6rEYsaywYxH7NVutXmmK3eXkH3MD6/0fwMSFVBwJ1piMCgYEA4UYM
/KOHzosscvGNPFW8U8b+5yNTh/Vm2CVBCe+ZwUA6WEE+nwGh4FVE5tmxzlXOr91h
I6LJ3WeV/byUl+5wAwAe4cBGwSrjOeR5VOgI+zXCH7LzKZG0EMTxhzRKs/R89gWi
zVL8noy1VomNLo2C+O9YBTLYF6S3Uu3PCmwda40CgYEAvwW4XbnILtKwxjlmRucD
7UsOnQl1tX183nV3t286B9AdlAWT5o8PV815/X3nMO2OnAe8JNg6f+NBNzeiuJfV
NX/8UHilGBkBNFOhOy2ffcs/qaaVMwf87nuqUcthdUHrfXBYLL5Sn0jIB8HAlEIG
fpztr3p7r11Y+YFGzNZCjAsCgYAo5qoW+K4Arz4rxHWrPbnK0DeZyc0xwzmgBuuP
HUSiVMIDIh13izlT3Md8zou89dFoFt67NKRIIbWW8zVbfHwz30K8JEf0bJADA9uP
se1nhvQvAzOpGX5DCS79KF5j3AEQPie39dhOBSgrhR/wEttzzSkDEJ8xc8OhN/I+
ZzDURQKBgAS8tKXfammEZBhLc1NGphBYf5HLXQEMm3Uqy7xW82YLjz6KN27x1fhz
YcZbFvY6KAsisGqZFUwlTecbZC8teFkFfoMWtu99fpPuLnlSr1MgLTchqnm1P70f
wWHpDR5CmggfYMrt/Hw1UQkRsDxuYnp6lnD8BOVZPo8jiMpUJmI8
-----END RSA PRIVATE KEY-----`
)
