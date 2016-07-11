package main

import (
	&#34;crypto/rand&#34;
	&#34;crypto/rsa&#34;
	&#34;crypto/sha1&#34;
	&#34;crypto/tls&#34;
	&#34;crypto/x509&#34;
	&#34;crypto/x509/pkix&#34;
	&#34;encoding/pem&#34;
	&#34;errors&#34;
	&#34;fmt&#34;
	&#34;log&#34;
	&#34;math/big&#34;
	&#34;net&#34;
	&#34;net/http&#34;
	&#34;time&#34;
)

var (
	CACert, CAKey = mustNewCA()

	ServerCert, ServerKey = mustNewServerCert()
	KeyBits               = 512
)

const addr = &#34;127.0.0.1:32643&#34;

func main() {
	err := NewServer(addr, ServerCert, ServerKey)
	if err != nil {
		log.Fatal(err)
	}

	client, err := oneCAClient(CACert)
	if err != nil {
		log.Fatalf(&#34;cannot get http client: %v&#34;, err)
	}
	req, err := http.NewRequest(&#34;GET&#34;, &#34;https://&#34;&#43;addr&#43;&#34;/&#34;, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf(&#34;get: %v&#34;, err)
	}
	log.Printf(&#34;got resp %#v&#34;, resp)
}

// Serve serves the given state by accepting requests on the given
// listener, using the given certificate and key (in PEM format) for
// authentication.
func NewServer(addr string, cert, key []byte) error {
	lis, err := net.Listen(&#34;tcp&#34;, addr)
	if err != nil {
		return err
	}
	log.Printf(&#34;listening on %q&#34;, lis.Addr())
	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return err
	}
	// TODO(rog) check that *srvRoot is a valid type for using
	// as an RPC server.
	lis = tls.NewListener(lis, &amp;tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	})
	mux := http.NewServeMux()
	mux.HandleFunc(&#34;/&#34;, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, &#34;ok\n&#34;)
	})
	go func() {
		err := http.Serve(lis, mux)
		log.Fatalf(&#34;server error: %v&#34;, err)
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
	secureConfig := &amp;tls.Config{
		RootCAs:    pool,
		ServerName: &#34;anything&#34;,
	}
	return &amp;http.Client{
		Transport: &amp;http.Transport{
			TLSClientConfig:   secureConfig,
			DisableKeepAlives: true,
		},
	}, nil
}

// ParseCert parses the given PEM-formatted X509 certificate.
func ParseCert(certPEM []byte) (*x509.Certificate, error) {
	for len(certPEM) &gt; 0 {
		var certBlock *pem.Block
		certBlock, certPEM = pem.Decode(certPEM)
		if certBlock == nil {
			break
		}
		if certBlock.Type == &#34;CERTIFICATE&#34; {
			cert, err := x509.ParseCertificate(certBlock.Bytes)
			return cert, err
		}
	}
	return nil, errors.New(&#34;no certificates found&#34;)
}

// NewCA generates a CA certificate/key pair suitable for signing server
// keys for an environment with the given name.
func NewCA(envName string, expiry time.Time) (certPEM, keyPEM []byte, err error) {
	key, err := rsa.GenerateKey(rand.Reader, KeyBits)
	if err != nil {
		return nil, nil, err
	}
	now := time.Now()
	template := &amp;x509.Certificate{
		SerialNumber: new(big.Int),
		Subject: pkix.Name{
			CommonName:   fmt.Sprintf(&#34;juju-generated CA for environment %q&#34;, envName),
			Organization: []string{&#34;juju&#34;},
		},
		NotBefore:             now.UTC().Add(-5 * time.Minute),
		NotAfter:              expiry.UTC(),
		SubjectKeyId:          bigIntHash(key.N),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		IsCA:                  true,
		MaxPathLen:            0, // Disallow delegation for now.
		BasicConstraintsValid: true,
	}
	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &amp;key.PublicKey, key)
	if err != nil {
		return nil, nil, fmt.Errorf(&#34;canot create certificate: %v&#34;, err)
	}
	certPEM = pem.EncodeToMemory(&amp;pem.Block{
		Type:  &#34;CERTIFICATE&#34;,
		Bytes: certDER,
	})
	keyPEM = pem.EncodeToMemory(&amp;pem.Block{
		Type:  &#34;RSA PRIVATE KEY&#34;,
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
		return nil, nil, fmt.Errorf(&#34;more than one certificate for CA&#34;)
	}
	caCert, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return nil, nil, err
	}
	if !caCert.BasicConstraintsValid || !caCert.IsCA {
		return nil, nil, fmt.Errorf(&#34;CA certificate is not a valid CA&#34;)
	}
	caKey, ok := tlsCert.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, fmt.Errorf(&#34;CA private key has unexpected type %T&#34;, tlsCert.PrivateKey)
	}
	key, err := rsa.GenerateKey(rand.Reader, KeyBits)
	if err != nil {
		return nil, nil, fmt.Errorf(&#34;cannot generate key: %v&#34;, err)
	}
	now := time.Now()
	template := &amp;x509.Certificate{
		SerialNumber: new(big.Int),
		Subject: pkix.Name{
			// This won&#39;t match host names with dots. The hostname
			// is hardcoded when connecting to avoid the issue.
			CommonName:   &#34;*&#34;,
			Organization: []string{&#34;juju&#34;},
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
	certDER, err := x509.CreateCertificate(rand.Reader, template, caCert, &amp;key.PublicKey, caKey)
	if err != nil {
		return nil, nil, err
	}
	certPEM = pem.EncodeToMemory(&amp;pem.Block{
		Type:  &#34;CERTIFICATE&#34;,
		Bytes: certDER,
	})
	keyPEM = pem.EncodeToMemory(&amp;pem.Block{
		Type:  &#34;RSA PRIVATE KEY&#34;,
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
	caCert, caKey, err := NewCA(&#34;juju testing&#34;, time.Now().AddDate(10, 0, 0))
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
		return nil, nil, fmt.Errorf(&#34;private key with unexpected type %T&#34;, key)
	}
	return cert, key, nil
}
