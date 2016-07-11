package main

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	_ "crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"sync"
	"time"
)

var wg sync.WaitGroup
var host = "127.0.0.1"

func main() {
	// Create a root CA.
	rootCertPem, _ := generate_cert(true, nil)
	rootCert, err := x509.ParseCertificate(rootCertPem.Bytes)
	if err != nil {
		log.Fatalf("failt to make parent: %s", err)
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
		log.Fatalf("making client TLS cert: %s", err)
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{c},
		RootCAs:      pool,
		ServerName:   host,
	}
	config.BuildNameToCertificate()

	client, err := tls.Dial("tcp", host+":6976", config)
	if err != nil {
		log.Fatalf("error connecting: %v\n", err)
	}
	n, err := client.Write([]byte("hello, server"))
	log.Println(n, err)
	client.Close()
}

func server(rootCert *x509.Certificate, pool *x509.CertPool) {
	// Create a server certificate.
	serverCertPem, serverKeyPem := generate_cert(false, rootCert)
	c, err := tls.X509KeyPair(pem.EncodeToMemory(serverCertPem),
		pem.EncodeToMemory(serverKeyPem))
	if err != nil {
		log.Fatalf("making server TLS cert: %s", err)
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{c},
		ServerName:   host,
		ClientCAs:    pool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	config.BuildNameToCertificate()
	server, err := tls.Listen("tcp", host+":6976", config)
	if err != nil {
		log.Fatalf("getting connection: %v\n", err)
	}
	wg.Done()
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println("error accepting:", err)
		} else {
			go handler(conn)
		}
	}
}

func handler(conn net.Conn) {
	defer conn.Close()
	br := bufio.NewReader(conn)
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			log.Println("reading line", conn.LocalAddr(), "-", conn.RemoteAddr(), ":", err)
			return
		}
		log.Println("got message:", line)
	}
}

func generate_cert(ca bool, parent *x509.Certificate) (*pem.Block, *pem.Block) {
	// Generate a key.
	key, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	if err != nil {
		log.Fatalf("failed to generate private key: %s", err)
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
		parent = &template
	}
	// Generate the certificate.
	cert, err := x509.CreateCertificate(rand.Reader, &template, parent, &key.PublicKey, key)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}
	// Marshal the key.
	b, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		log.Fatalf("Failed to marshal ecdsa: %s", err)
	}
	return &pem.Block{Type: "CERTIFICATE", Bytes: cert},
		&pem.Block{Type: "ECDSA PRIVATE KEY", Bytes: b}
}
