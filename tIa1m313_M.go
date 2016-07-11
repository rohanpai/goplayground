package main

import (
	&#34;crypto/x509&#34;
	&#34;encoding/pem&#34;
	&#34;log&#34;
	&#34;strings&#34;
)

const gitlabTestCAPEM = `Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number: 1 (0x1)
    Signature Algorithm: sha1WithRSAEncryption
        Issuer: CN=GitLab Test CA
        Validity
            Not Before: Mar 15 21:17:00 2016 GMT
            Not After : Mar 15 21:17:00 2026 GMT
        Subject: CN=GitLab Test CA
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                Public-Key: (1024 bit)
                Modulus:
                    00:ae:bc:cd:f5:87:24:cb:4c:55:cc:31:df:a5:28:
                    ef:cf:b5:70:22:dc:d6:28:3b:93:4f:98:db:42:9e:
                    3d:b0:52:11:c0:51:57:21:30:10:96:d8:d0:30:04:
                    9d:12:7c:57:4d:99:c1:f3:5e:e3:90:fb:d4:99:bf:
                    a5:79:d3:39:87:13:47:a8:c0:0b:00:24:e3:fa:a8:
                    c3:4f:77:56:a1:94:7f:2f:b1:97:df:ff:6e:86:70:
                    67:af:bc:f6:33:e4:66:12:a2:dc:2e:27:89:da:b3:
                    4b:5b:8b:70:2c:1d:97:23:ae:16:d2:d2:41:b6:10:
                    48:05:a3:32:6d:39:ca:3c:fb
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Basic Constraints: critical
                CA:TRUE
            X509v3 Subject Key Identifier:
                BE:2E:29:CB:51:A1:A5:83:9A:C1:E2:87:79:B9:12:37:7E:C3:F0:4F
            X509v3 Key Usage:
                Certificate Sign, CRL Sign
            X509v3 Name Constraints:
                Permitted:
                  DNS:.with.dot.com
                  DNS:without.dot.com

            Netscape Cert Type:
                SSL CA, S/MIME CA, Object Signing CA
            Netscape Comment:
                xca certificate
    Signature Algorithm: sha1WithRSAEncryption
         68:57:43:1b:b3:da:cd:1d:f7:9a:8d:4c:8e:61:60:65:99:cc:
         87:c4:13:51:e8:86:d9:e2:e3:5f:c4:fc:29:38:de:d1:b5:3f:
         4b:45:1e:3a:12:9c:59:88:a8:56:62:d8:66:d9:46:7b:f4:4a:
         cb:3c:76:4b:56:1d:a8:0a:a0:7d:c5:fe:d6:fc:ea:35:1d:de:
         54:a4:9b:4f:b5:95:cf:8e:20:4e:9b:29:ff:9a:77:93:47:14:
         29:6e:bd:7f:ee:5f:1d:29:e9:dc:a1:5d:11:b7:41:90:91:d1:
         ee:7e:d2:a7:e6:ac:f1:14:a7:6b:e5:2e:a5:cd:db:b7:02:3b:
         39:80
-----BEGIN CERTIFICATE-----
MIICTTCCAbagAwIBAgIBATANBgkqhkiG9w0BAQUFADAZMRcwFQYDVQQDEw5HaXRM
YWIgVGVzdCBDQTAeFw0xNjAzMTUyMTE3MDBaFw0yNjAzMTUyMTE3MDBaMBkxFzAV
BgNVBAMTDkdpdExhYiBUZXN0IENBMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKB
gQCuvM31hyTLTFXMMd&#43;lKO/PtXAi3NYoO5NPmNtCnj2wUhHAUVchMBCW2NAwBJ0S
fFdNmcHzXuOQ&#43;9SZv6V50zmHE0eowAsAJOP6qMNPd1ahlH8vsZff/26GcGevvPYz
5GYSotwuJ4nas0tbi3AsHZcjrhbS0kG2EEgFozJtOco8&#43;wIDAQABo4GkMIGhMA8G
A1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFL4uKctRoaWDmsHih3m5Ejd&#43;w/BPMAsG
A1UdDwQEAwIBBjAvBgNVHR4EKDAmoCQwD4INLndpdGguZG90LmNvbTARgg93aXRo
b3V0LmRvdC5jb20wEQYJYIZIAYb4QgEBBAQDAgAHMB4GCWCGSAGG&#43;EIBDQQRFg94
Y2EgY2VydGlmaWNhdGUwDQYJKoZIhvcNAQEFBQADgYEAaFdDG7PazR33mo1MjmFg
ZZnMh8QTUeiG2eLjX8T8KTje0bU/S0UeOhKcWYioVmLYZtlGe/RKyzx2S1YdqAqg
fcX&#43;1vzqNR3eVKSbT7WVz44gTpsp/5p3k0cUKW69f&#43;5fHSnp3KFdEbdBkJHR7n7S
p&#43;as8RSna&#43;Uupc3btwI7OYA=
-----END CERTIFICATE-----`

func decodePEM(data string) []byte {
	block, _ := pem.Decode([]byte(data))
	if block == nil {
		log.Fatalln(&#34;PEM failed to decode&#34;)
	}
	if block.Type != &#34;CERTIFICATE&#34; {
		log.Fatalln(&#34;PEM invalid block type:&#34;, block.Type)
	}
	return block.Bytes
}

func roots() *x509.CertPool {
	caCert, err := x509.ParseCertificate(decodePEM(gitlabTestCAPEM))
	if err != nil {
		log.Fatalln(&#34;Parse CA PEM:&#34;, err)
	}

	roots := x509.NewCertPool()
	roots.AddCert(caCert)
	println(strings.Join(caCert.PermittedDNSDomains, &#34; &#34;))
	return roots
}

func verifyName(name string, cert *x509.Certificate, expected string) {
	var actual string
	_, err := cert.Verify(x509.VerifyOptions{DNSName: name})
	_, isHostErr := err.(*x509.HostnameError)

	// we ignore hostname errors
	if err != nil &amp;&amp; !isHostErr {
		actual = &#34;FAIL&#34;
	} else {
		actual = &#34;OK&#34;
	}

	if expected != actual {
		log.Println(name, &#34;Expected:&#34;, expected, &#34;Actual:&#34;, actual, &#34;Error:&#34;, err)
	} else {
		log.Println(name, &#34;IS OK&#34;)
	}
}

func main() {
	caCert, err := x509.ParseCertificate(decodePEM(gitlabTestCAPEM))
	if err != nil {
		log.Fatalln(&#34;Parse CA PEM:&#34;, err)
	}

	verifyName(&#34;example.with.dot.com&#34;, caCert, &#34;OK&#34;)
	verifyName(&#34;my.example.with.dot.com&#34;, caCert, &#34;OK&#34;)
	verifyName(&#34;with.dot.com&#34;, caCert, &#34;FAIL&#34;)
	verifyName(&#34;example.without.dot.com&#34;, caCert, &#34;FAIL&#34;)
	verifyName(&#34;without.dot.com&#34;, caCert, &#34;FAIL&#34;)
}
