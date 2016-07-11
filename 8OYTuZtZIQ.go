package main

import (
	&#34;net&#34;
	&#34;errors&#34;
	&#34;strings&#34;
  &#34;io/ioutil&#34;
	&#34;encoding/pem&#34;
	&#34;crypto&#34;
	&#34;crypto/tls&#34;
	&#34;crypto/x509&#34;
	&#34;crypto/rsa&#34;
	&#34;crypto/ecdsa&#34;
	&#34;flag&#34;
	&#34;log&#34;
	&#34;net/http&#34;
	&#34;fmt&#34;
)

var url = flag.String(&#34;url&#34;, &#34;https://127.0.0.1:8443&#34;, &#34;the url to get&#34;)
var certFile = flag.String(&#34;certFile&#34;, &#34;cert.pem&#34;, &#34;the cert for client auth&#34;)
var keyFile = flag.String(&#34;keyFile&#34;, &#34;key.pem&#34;, &#34;the key for client auth&#34;)

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set(&#34;Content-Type&#34;, &#34;text/plain&#34;)

	fmt.Fprintln(w, &#34;Your conn state is: &#34;, req.TLS);
	fmt.Fprintln(w, &#34;Your client cert is: &#34;, req.TLS.PeerCertificates);
}

func loadX509KeyPair(certFile, keyFile string) (cert tls.Certificate, err error) {
  certPEMBlock, err := ioutil.ReadFile(certFile)
  if err != nil {
    return
  }
  keyPEMBlock, err := ioutil.ReadFile(keyFile)
  if err != nil {
    return
  }
  return X509KeyPair(certPEMBlock, keyPEMBlock, []byte(&#34;password&#34;))
}

func X509KeyPair(certPEMBlock, keyPEMBlock, pw []byte) (cert tls.Certificate, err error) {
  var certDERBlock *pem.Block
  for {
    certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)
    if certDERBlock == nil {
      break
    }
    if certDERBlock.Type == &#34;CERTIFICATE&#34; {
      cert.Certificate = append(cert.Certificate, certDERBlock.Bytes)
    }
  }

  if len(cert.Certificate) == 0 {
    err = errors.New(&#34;crypto/tls: failed to parse certificate PEM data&#34;)
    return
  }
  var keyDERBlock *pem.Block
  for {
    keyDERBlock, keyPEMBlock = pem.Decode(keyPEMBlock)
    if keyDERBlock == nil {
      err = errors.New(&#34;crypto/tls: failed to parse key PEM data&#34;)
      return
    }
		if x509.IsEncryptedPEMBlock(keyDERBlock) {
      out, err2 := x509.DecryptPEMBlock(keyDERBlock, pw)
			if err2 != nil {
				err = err2
				return
			}
      keyDERBlock.Bytes = out
      break
    }
    if keyDERBlock.Type == &#34;PRIVATE KEY&#34; || strings.HasSuffix(keyDERBlock.Type, &#34; PRIVATE KEY&#34;) {
      break
    }
  }

  cert.PrivateKey, err = parsePrivateKey(keyDERBlock.Bytes)
  if err != nil {
    return
  }
  // We don&#39;t need to parse the public key for TLS, but we so do anyway
  // to check that it looks sane and matches the private key.
  x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
  if err != nil {
    return
  }

  switch pub := x509Cert.PublicKey.(type) {
  case *rsa.PublicKey:
    priv, ok := cert.PrivateKey.(*rsa.PrivateKey)
    if !ok {
      err = errors.New(&#34;crypto/tls: private key type does not match public key type&#34;)
      return
    }
    if pub.N.Cmp(priv.N) != 0 {
      err = errors.New(&#34;crypto/tls: private key does not match public key&#34;)
      return
    }
  case *ecdsa.PublicKey:
    priv, ok := cert.PrivateKey.(*ecdsa.PrivateKey)
    if !ok {
      err = errors.New(&#34;crypto/tls: private key type does not match public key type&#34;)
      return

    }
    if pub.X.Cmp(priv.X) != 0 || pub.Y.Cmp(priv.Y) != 0 {
      err = errors.New(&#34;crypto/tls: private key does not match public key&#34;)
      return
    }
  default:
    err = errors.New(&#34;crypto/tls: unknown public key algorithm&#34;)
    return
  }
return
}

// Attempt to parse the given private key DER block. OpenSSL 0.9.8 generates
// PKCS#1 private keys by default, while OpenSSL 1.0.0 generates PKCS#8 keys.
// OpenSSL ecparam generates SEC1 EC private keys for ECDSA. We try all three.
func parsePrivateKey(der []byte) (crypto.PrivateKey, error) {
  if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
    return key, nil
  }
  if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
    switch key := key.(type) {
    case *rsa.PrivateKey, *ecdsa.PrivateKey:
      return key, nil
    default:
      return nil, errors.New(&#34;crypto/tls: found unknown private key type in PKCS#8 wrapping&#34;)
    }
  }
  if key, err := x509.ParseECPrivateKey(der); err == nil {
    return key, nil
  }

  return nil, errors.New(&#34;crypto/tls: failed to parse private key&#34;)
}

func main() {
	flag.Parse()

	cert, err := loadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatal(err)
	}

	tlscfg := &amp;tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	srv := &amp;http.Server{Addr: &#34;:8443&#34;, Handler: nil}

	l, err := net.Listen(&#34;tcp&#34;, srv.Addr)
	if err != nil {
		log.Fatal(err)
	}
	tl := tls.NewListener(l, tlscfg)

	http.HandleFunc(&#34;/&#34;, handler)
	srv.Serve(tl)
}
