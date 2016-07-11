//
// direct test of secure vs non-secure websockets, wss and ws uri schemes respectively
//
// to run this test:
//
// 1. &#34;go run ws.go&#34;
// 2. in a browser, visit unsecure url &#34;http://localhost:8080&#34; ---
//	this is insecure version, you should see echo messages once per second
// 3. in a browser, visit secure url &#34;https://localhost:8090&#34; ---
//	this is secure version, depending on browser type, you should see same output as insecure version, or error; note
//      that you will have to accept the untrusted certificate below and/or security exception manually when prompted.
//
// in particular, wss does not work in my chrome 28.0.1500.95 (it silently fails on both browser and server sides),
//      but works in &#34;Opera 12.16 Build 1860 for Linux x86_64&#34; and also works in mozilla firefox 23.0.
//
// unsecure version works fine in all three browsers i tested.
//
package main

import (
	&#34;bytes&#34;
	&#34;code.google.com/p/go.net/websocket&#34;
	&#34;fmt&#34;
	&#34;io&#34;
	&#34;io/ioutil&#34;
	&#34;net/http&#34;
	&#34;os&#34;
	&#34;time&#34;
)

func main() {

	http.HandleFunc(&#34;/&#34;, func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		switch r.URL.Path {
		case &#34;/&#34;:
			u := func() string {
				if r.TLS == nil {
					return &#34;ws://localhost:8080/sock&#34;
				} else {
					return &#34;wss://localhost:8090/sock&#34;
				}
			}()
			http.ServeContent(w, r, &#34;index.html&#34;, time.Now(), bytes.NewReader([]byte(Content(u))))
		case &#34;/sock&#34;:
			websocket.Handler(EchoHandler).ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	done := make(chan bool)

	go func() {
		fmt.Println(http.ListenAndServe(&#34;:8080&#34;, nil))
		done &lt;- true
	}()

	go func() {
		ioutil.WriteFile(CERT_NAME, []byte(CERT), os.ModePerm)
		ioutil.WriteFile(KEY_NAME, []byte(KEY), os.ModePerm)
		fmt.Println(http.ListenAndServeTLS(&#34;:8090&#34;, CERT_NAME, KEY_NAME, nil))
		done &lt;- true
	}()

	&lt;-done

}

func EchoHandler(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

func Content(ws string) string {
	return fmt.Sprintf(`&lt;!DOCTYPE html&gt;
&lt;html lang=&#39;en&#39;&gt;
  &lt;head&gt;
    &lt;meta charset=&#39;utf-8&#39;/&gt;
    &lt;title&gt;websocket test&lt;/title&gt;
    &lt;script src=&#34;//ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js&#34;&gt;&lt;/script&gt;  
  &lt;/head&gt;
  &lt;body&gt;
    &lt;p&gt; go-lang / chrome websocket test at %s: &lt;/p&gt;
    &lt;ol/&gt;
    &lt;script&gt;

var log = function(e) {
    $(&#39;ol&#39;).append(&#34;&lt;li&gt;&#34;&#43;JSON.stringify(e)&#43;&#34;&lt;/li&gt;&#34;);
}

var sock = new WebSocket(&#34;%s&#34;)

sock.onerror = function(e) { log({ERROR:e}) };    
sock.onclose = function(e) { log({CLOSE:e}) };    
sock.onmessage = function(e) { log({MESSAGE:e.data}) };    

sock.onopen = function(e) { 
    log({OPEN:e}) 
    setInterval(function() {
	sock.send(&#34;howdy at &#34; &#43; new Date());
    },1000);
};    

    &lt;/script&gt;
  &lt;/body&gt;
&lt;/html&gt;
`, ws, ws)
}

const (
	CERT_NAME = &#34;/tmp/cert.prm&#34;
	CERT      = `-----BEGIN CERTIFICATE-----
MIIC4zCCAc2gAwIBAgIBADALBgkqhkiG9w0BAQUwEjEQMA4GA1UEChMHQWNtZSBD
bzAeFw0xMzA4MTIxOTUzNTdaFw0xNDA4MTIxOTUzNTdaMBIxEDAOBgNVBAoTB0Fj
bWUgQ28wggEgMAsGCSqGSIb3DQEBAQOCAQ8AMIIBCgKCAQEAw917qyvQTFtr65LQ
MimR7obQzOc05YvqH/ytdtqSyPqQsvnOLi0Jk2/t33RpqIGb36dIMGDV3X5SpiLW
9TReidxtz22YbZ1CTRUgdcF5XJiq3CO5a&#43;1vGvi/fIBIqzPb6CQwgG/eFm6xdvgf
IRodZgM6ym6cglm6ndhoB4/TpIni/i&#43;bo757LUxfu78/mkUBfsQK0VDaqq59ZkZW
An8E/4DBeDdg8jq9i&#43;4VuOTSulfiLu06eIWoAhfBeWX6uGaEK44xW5NJpp1y4CKX
BOvVKuz&#43;Dv4jz/vp1e8p1zFqkUzbrtdVv&#43;Z3olfTHQWe1qMpt9l01tedxuPTEHas
HTQiRwIDAQABo0owSDAOBgNVHQ8BAf8EBAMCAKAwEwYDVR0lBAwwCgYIKwYBBQUH
AwEwDAYDVR0TAQH/BAIwADATBgNVHREEDDAKggh4dXZhLmNvbTALBgkqhkiG9w0B
AQUDggEBAFppIRxB25BRCK1w//9U73sddEnw2Q/MUR&#43;V4twB5qVGnnY6VJW&#43;u9W3
9dNs4teemAxUJh7ZBOR2xEj1N6Q9D0H3BFupUNRa9jIxWam&#43;E&#43;mNCF1jQRlTxCur
WHk56ZmsqWyb8yXIB3ymHyAXnJziAUO/US6e8xeMKcvIZCtCzCsaOLI5G35h9xoJ
&#43;mKedqOW7tGmT3suqj/bx2hEq1nPRW/H9XmyLkDMvIiUaU8iopKzShuyG/kosRXY
WsoDWrkv3a0O1crnlQLZXg1IH0/bQg7OwJz5P8qDSNR1uiBzGkmbOFpp0IVCdVIk
UwKAnERUK1MkYu6jh5hXj6mb0fW01/o=
-----END CERTIFICATE-----`
	KEY_NAME = &#34;/tmp/key.pem&#34;
	KEY      = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAw917qyvQTFtr65LQMimR7obQzOc05YvqH/ytdtqSyPqQsvnO
Li0Jk2/t33RpqIGb36dIMGDV3X5SpiLW9TReidxtz22YbZ1CTRUgdcF5XJiq3CO5
a&#43;1vGvi/fIBIqzPb6CQwgG/eFm6xdvgfIRodZgM6ym6cglm6ndhoB4/TpIni/i&#43;b
o757LUxfu78/mkUBfsQK0VDaqq59ZkZWAn8E/4DBeDdg8jq9i&#43;4VuOTSulfiLu06
eIWoAhfBeWX6uGaEK44xW5NJpp1y4CKXBOvVKuz&#43;Dv4jz/vp1e8p1zFqkUzbrtdV
v&#43;Z3olfTHQWe1qMpt9l01tedxuPTEHasHTQiRwIDAQABAoIBAA22&#43;4rf1YUTPbpQ
HG32xTYrkIFYizarlmhI/Ch/Y5nZGbq&#43;jTZkhvAg/UoRT7ix4qVFhGOG1FLfHpBt
jhm7YgdLPREyPmMmiNb27L/yHTpjoksp4Tjydj4wPtBL90qtpe9aYV8M9kMh2yFW
fG&#43;H8ZkMDtjP5/ukptGYrqgg5RP3SGacLGdmdCakesCP/&#43;bjSjKVyHsBThybr5oy
zSsLaGocU6JL7k6x1lbnbToSwQo5ORRgWzsVIc68qm4Nc4&#43;Nn9lrhMqlucYMtrxf
EyNpcCUFuO1BMz1jJwP5zttY0ZRKWitaDRpVyZhKMonUSsYIfRV&#43;kM9&#43;qKatZGbW
Tuh9dLkCgYEA3pSRWsyCHxmVDXpTAw0EPSZISaAqq5RSZoBtJkKBGpqxs245RI4E
FQN2d3oPpI3qBf3J6&#43;lLRy4W0Hju&#43;UxSzqdCHONmngdxbzOuMdoianjwviwr&#43;ZjM
e8kmFtQ72KK6rEYsaywYxH7NVutXmmK3eXkH3MD6/0fwMSFVBwJ1piMCgYEA4UYM
/KOHzosscvGNPFW8U8b&#43;5yNTh/Vm2CVBCe&#43;ZwUA6WEE&#43;nwGh4FVE5tmxzlXOr91h
I6LJ3WeV/byUl&#43;5wAwAe4cBGwSrjOeR5VOgI&#43;zXCH7LzKZG0EMTxhzRKs/R89gWi
zVL8noy1VomNLo2C&#43;O9YBTLYF6S3Uu3PCmwda40CgYEAvwW4XbnILtKwxjlmRucD
7UsOnQl1tX183nV3t286B9AdlAWT5o8PV815/X3nMO2OnAe8JNg6f&#43;NBNzeiuJfV
NX/8UHilGBkBNFOhOy2ffcs/qaaVMwf87nuqUcthdUHrfXBYLL5Sn0jIB8HAlEIG
fpztr3p7r11Y&#43;YFGzNZCjAsCgYAo5qoW&#43;K4Arz4rxHWrPbnK0DeZyc0xwzmgBuuP
HUSiVMIDIh13izlT3Md8zou89dFoFt67NKRIIbWW8zVbfHwz30K8JEf0bJADA9uP
se1nhvQvAzOpGX5DCS79KF5j3AEQPie39dhOBSgrhR/wEttzzSkDEJ8xc8OhN/I&#43;
ZzDURQKBgAS8tKXfammEZBhLc1NGphBYf5HLXQEMm3Uqy7xW82YLjz6KN27x1fhz
YcZbFvY6KAsisGqZFUwlTecbZC8teFkFfoMWtu99fpPuLnlSr1MgLTchqnm1P70f
wWHpDR5CmggfYMrt/Hw1UQkRsDxuYnp6lnD8BOVZPo8jiMpUJmI8
-----END RSA PRIVATE KEY-----`
)
