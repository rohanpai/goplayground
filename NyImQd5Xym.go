package main

import (
	&#34;bufio&#34;
	&#34;crypto/ecdsa&#34;
	&#34;crypto/elliptic&#34;
	&#34;crypto/rand&#34;
	_ &#34;crypto/sha512&#34;
	&#34;crypto/tls&#34;
	&#34;crypto/x509&#34;
	&#34;crypto/x509/pkix&#34;
	&#34;encoding/pem&#34;
	&#34;log&#34;
	&#34;math/big&#34;
	&#34;net&#34;
	&#34;sync&#34;
	&#34;time&#34;
)

var wg sync.WaitGroup
var host = &#34;127.0.0.1&#34;

func main() {
	// Create a root CA.
	rootCertPem, _ := generate_cert(true, nil)
	rootCert, err := x509.ParseCertificate(rootCertPem.Bytes)
	if err != nil {
		log.Fatalf(&#34;failt to make parent: %s&#34;, err)
	}
	// Make the CertPool.
	pool := x509.NewCertPool()
	pool.AddCert(rootCert)

	// Start the server
	wg.Add(1)
	go server(rootCert, pool)
	wg.Wait()

	// Create client certificate.
	clientCertPem, clientKeyPem := generate_cert(false, rootCert)
	c, err := tls.X509KeyPair(pem.EncodeToMemory(clientCertPem),
		pem.EncodeToMemory(clientKeyPem))
	if err != nil {
		log.Fatalf(&#34;making client TLS cert: %s&#34;, err)
	}
	config := &amp;tls.Config{
		Certificates: []tls.Certificate{c},
		RootCAs:      pool,
		ServerName:   host,
	}
	config.BuildNameToCertificate()

	client, err := tls.Dial(&#34;tcp&#34;, host&#43;&#34;:6976&#34;, config)
	if err != nil {
		log.Fatalf(&#34;error connecting: %v\n&#34;, err)
	}
	n, err := client.Write([]byte(&#34;hello, server&#34;))
	log.Println(n, err)
	client.Close()
}

func server(rootCert *x509.Certificate, pool *x509.CertPool) {
	// Create a server certificate.
	serverCertPem, serverKeyPem := generate_cert(false, rootCert)
	c, err := tls.X509KeyPair(pem.EncodeToMemory(serverCertPem),
		pem.EncodeToMemory(serverKeyPem))
	if err != nil {
		log.Fatalf(&#34;making server TLS cert: %s&#34;, err)
	}
	config := &amp;tls.Config{
		Certificates: []tls.Certificate{c},
		ServerName:   host,
		ClientCAs:    pool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	config.BuildNameToCertificate()
	server, err := tls.Listen(&#34;tcp&#34;, host&#43;&#34;:6976&#34;, config)
	if err != nil {
		log.Fatalf(&#34;getting connection: %v\n&#34;, err)
	}
	wg.Done()
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println(&#34;error accepting:&#34;, err)
		} else {
			go handler(conn)
		}
	}
}

func handler(conn net.Conn) {
	defer conn.Close()
	br := bufio.NewReader(conn)
	for {
		line, err := br.ReadString(&#39;\n&#39;)
		if err != nil {
			log.Println(&#34;reading line&#34;, conn.LocalAddr(), &#34;-&#34;, conn.RemoteAddr(), &#34;:&#34;, err)
			return
		}
		log.Println(&#34;got message:&#34;, line)
	}
}

func generate_cert(ca bool, parent *x509.Certificate) (*pem.Block, *pem.Block) {
	// Generate a key.
	key, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	if err != nil {
		log.Fatalf(&#34;failed to generate private key: %s&#34;, err)
	}
	// Fill out the template.
	template := x509.Certificate{
		SerialNumber:          new(big.Int).SetInt64(0),
		Subject:               pkix.Name{Organization: []string{host}},
		NotBefore:             time.Now(),
		NotAfter:              time.Date(2049, 12, 31, 23, 59, 59, 0, time.UTC),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP(host)},
	}
	if ca {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}
	if parent == nil {
		parent = &amp;template
	}
	// Generate the certificate.
	cert, err := x509.CreateCertificate(rand.Reader, &amp;template, parent, &amp;key.PublicKey, key)
	if err != nil {
		log.Fatalf(&#34;Failed to create certificate: %s&#34;, err)
	}
	// Marshal the key.
	b, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		log.Fatalf(&#34;Failed to marshal ecdsa: %s&#34;, err)
	}
	return &amp;pem.Block{Type: &#34;CERTIFICATE&#34;, Bytes: cert},
		&amp;pem.Block{Type: &#34;ECDSA PRIVATE KEY&#34;, Bytes: b}
}
