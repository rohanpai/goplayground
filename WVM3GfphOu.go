package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"time"
)

var (
	CACert, CAKey = mustNewCA()

	ServerCert, ServerKey = mustNewServerCert()
	KeyBits               = 512
)

const addr = "127.0.0.1:32643"

func main() {
	err := NewServer(addr, ServerCert, ServerKey)
	if err != nil {
		log.Fatal(err)
	}

	client, err := oneCAClient(CACert)
	if err != nil {
		log.Fatalf("cannot get http client: %v", err)
	}
	req, err := http.NewRequest("GET", "https://"+addr+"/", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("get: %v", err)
	}
	log.Printf("got resp %#v", resp)
}

// Serve serves the given state by accepting requests on the given
// listener, using the given certificate and key (in PEM format) for
// authentication.
func NewServer(addr string, cert, key []byte) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Printf("listening on %q", lis.Addr())
	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return err
	}
	// TODO(rog) check that *srvRoot is a valid type for using
	// as an RPC server.
	lis = tls.NewListener(lis, &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok\n")
	})
	go func() {
		err := http.Serve(lis, mux)
		log.Fatalf("server error: %v", err)
	}()
	return nil
}

// oneCAClient returns a *http.Client, which includes the
// given CA certificate in its trusted list, so it can be used to
// connect to an API server
func oneCAClient(caCert []byte) (*http.Client, error) {
	pool := x509.NewCertPool()
	xcert, err := ParseCert(caCert)
	if err != nil {
		return nil, err
	}
	pool.AddCert(xcert)
	secureConfig := &tls.Config{
		RootCAs:    pool,
		ServerName: "anything",
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   secureConfig,
			DisableKeepAlives: true,
		},
	}, nil
}

// ParseCert parses the given PEM-formatted X509 certificate.
func ParseCert(certPEM []byte) (*x509.Certificate, error) {
	for len(certPEM) > 0 {
		var certBlock *pem.Block
		certBlock, certPEM = pem.Decode(certPEM)
		if certBlock == nil {
			break
		}
		if certBlock.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(certBlock.Bytes)
			return cert, err
		}
	}
	return nil, errors.New("no certificates found")
}

// NewCA generates a CA certificate/key pair suitable for signing server
// keys for an environment with the given name.
func NewCA(envName string, expiry time.Time) (certPEM, keyPEM []byte, err error) {
	key, err := rsa.GenerateKey(rand.Reader, KeyBits)
	if err != nil {
		return nil, nil, err
	}
	now := time.Now()
	template := &x509.Certificate{
		SerialNumber: new(big.Int),
		Subject: pkix.Name{
			CommonName:   fmt.Sprintf("juju-generated CA for environment %q", envName),
			Organization: []string{"juju"},
		},
		NotBefore:             now.UTC().Add(-5 * time.Minute),
		NotAfter:              expiry.UTC(),
		SubjectKeyId:          bigIntHash(key.N),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		IsCA:                  true,
		MaxPathLen:            0, // Disallow delegation for now.
		BasicConstraintsValid: true,
	}
	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return nil, nil, fmt.Errorf("canot create certificate: %v", err)
	}
	certPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})
	keyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	return certPEM, keyPEM, nil
}

// NewServerCert generates a certificate/key pair suitable for use by a server.
func NewServerCert(caCertPEM, caKeyPEM []byte, expiry time.Time, hostnames []string) (certPEM, keyPEM []byte, err error) {
	return newLeaf(caCertPEM, caKeyPEM, expiry, hostnames, nil)
}

// newLeaf generates a certificate/key pair suitable for use by a leaf node.
func newLeaf(caCertPEM, caKeyPEM []byte, expiry time.Time, hostnames []string, extKeyUsage []x509.ExtKeyUsage) (certPEM, keyPEM []byte, err error) {
	tlsCert, err := tls.X509KeyPair(caCertPEM, caKeyPEM)
	if err != nil {
		return nil, nil, err
	}
	if len(tlsCert.Certificate) != 1 {
		return nil, nil, fmt.Errorf("more than one certificate for CA")
	}
	caCert, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return nil, nil, err
	}
	if !caCert.BasicConstraintsValid || !caCert.IsCA {
		return nil, nil, fmt.Errorf("CA certificate is not a valid CA")
	}
	caKey, ok := tlsCert.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, fmt.Errorf("CA private key has unexpected type %T", tlsCert.PrivateKey)
	}
	key, err := rsa.GenerateKey(rand.Reader, KeyBits)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot generate key: %v", err)
	}
	now := time.Now()
	template := &x509.Certificate{
		SerialNumber: new(big.Int),
		Subject: pkix.Name{
			// This won't match host names with dots. The hostname
			// is hardcoded when connecting to avoid the issue.
			CommonName:   "*",
			Organization: []string{"juju"},
		},
		NotBefore: now.UTC().Add(-5 * time.Minute),
		NotAfter:  expiry.UTC(),

		SubjectKeyId: bigIntHash(key.N),
		KeyUsage:     x509.KeyUsageDataEncipherment,
		ExtKeyUsage:  extKeyUsage,
	}
	for _, hostname := range hostnames {
		if ip := net.ParseIP(hostname); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, hostname)
		}
	}
	certDER, err := x509.CreateCertificate(rand.Reader, template, caCert, &key.PublicKey, caKey)
	if err != nil {
		return nil, nil, err
	}
	certPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})
	keyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	return certPEM, keyPEM, nil
}

func bigIntHash(n *big.Int) []byte {
	h := sha1.New()
	h.Write(n.Bytes())
	return h.Sum(nil)
}

func mustParseCertAndKey(certPEM, keyPEM []byte) (*x509.Certificate, *rsa.PrivateKey) {
	cert, key, err := ParseCertAndKey(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return cert, key
}

func mustNewCA() ([]byte, []byte) {
	caCert, caKey, err := NewCA("juju testing", time.Now().AddDate(10, 0, 0))
	if err != nil {
		panic(err)
	}
	return caCert, caKey
}

func mustNewServerCert() ([]byte, []byte) {
	var hostnames []string
	srvCert, srvKey, err := NewServerCert([]byte(CACert), []byte(CAKey), time.Now().AddDate(10, 0, 0), hostnames)
	if err != nil {
		panic(err)
	}
	return srvCert, srvKey
}

// ParseCert parses the given PEM-formatted X509 certificate
// and RSA private key.
func ParseCertAndKey(certPEM, keyPEM []byte) (*x509.Certificate, *rsa.PrivateKey, error) {
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return nil, nil, err
	}

	key, ok := tlsCert.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, fmt.Errorf("private key with unexpected type %T", key)
	}
	return cert, key, nil
}
