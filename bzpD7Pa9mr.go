package main

import (
	&#34;crypto&#34;
	&#34;crypto/rand&#34;
	&#34;crypto/rsa&#34;
	&#34;crypto/sha256&#34;
	&#34;crypto/x509&#34;
	&#34;encoding/base64&#34;
	&#34;encoding/pem&#34;
	&#34;errors&#34;
	&#34;fmt&#34;
)

func main() {
	signer, err := loadPrivateKey(&#34;private.pem&#34;)
	if err != nil {
		fmt.Errorf(&#34;signer is damaged: %v&#34;, err)
	}

	toSign := &#34;date: Thu, 05 Jan 2012 21:31:40 GMT&#34;

	signed, err := signer.Sign([]byte(toSign))
	if err != nil {
		fmt.Errorf(&#34;could not sign request: %v&#34;, err)
	}
	sig := base64.StdEncoding.EncodeToString(signed)
	fmt.Printf(&#34;Signature: %v\n&#34;, sig)

	parser, perr := loadPublicKey(&#34;public.pem&#34;)
	if perr != nil {
		fmt.Errorf(&#34;could not sign request: %v&#34;, err)
	}
	
	err = parser.Unsign([]byte(toSign), signed)
	if err != nil {
		fmt.Errorf(&#34;could not sign request: %v&#34;, err)
	}
	
	fmt.Printf(&#34;Unsign error: %v\n&#34;, err)
}

// loadPrivateKey loads an parses a PEM encoded private key file.
func loadPublicKey(path string) (Unsigner, error) {

	return parsePublicKey([]byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCFENGw33yGihy92pDjZQhl0C3
6rPJj&#43;CvfSC8&#43;q28hxA161QFNUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78&#43;6
Z4UMR7EOcpfdUE9Hf3m/hs&#43;FUR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc&#43;2xjJw
oYi&#43;1hqp1fIekaxsyQIDAQAB
-----END PUBLIC KEY-----`))
}

// parsePublicKey parses a PEM encoded private key.
func parsePublicKey(pemBytes []byte) (Unsigner, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New(&#34;ssh: no key found&#34;)
	}

	var rawkey interface{}
	switch block.Type {
	case &#34;PUBLIC KEY&#34;:
		rsa, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = rsa
	default:
		return nil, fmt.Errorf(&#34;ssh: unsupported key type %q&#34;, block.Type)
	}

	return newUnsignerFromKey(rawkey)
}

// loadPrivateKey loads an parses a PEM encoded private key file.
func loadPrivateKey(path string) (Signer, error) {
	return parsePrivateKey([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDCFENGw33yGihy92pDjZQhl0C36rPJj&#43;CvfSC8&#43;q28hxA161QF
NUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78&#43;6Z4UMR7EOcpfdUE9Hf3m/hs&#43;F
UR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc&#43;2xjJwoYi&#43;1hqp1fIekaxsyQIDAQAB
AoGBAJR8ZkCUvx5kzv&#43;utdl7T5MnordT1TvoXXJGXK7ZZ&#43;UuvMNUCdN2QPc4sBiA
QWvLw1cSKt5DsKZ8UETpYPy8pPYnnDEz2dDYiaew9&#43;xEpubyeW2oH4Zx71wqBtOK
kqwrXa/pzdpiucRRjk6vE6YY7EBBs/g7uanVpGibOVAEsqH1AkEA7DkjVH28WDUg
f1nqvfn2Kj6CT7nIcE3jGJsZZ7zlZmBmHFDONMLUrXR/Zm3pR5m0tCmBqa5RK95u
412jt1dPIwJBANJT3v8pnkth48bQo/fKel6uEYyboRtA5/uHuHkZ6FQF7OUkGogc
mSJluOdc5t6hI1VsLn0QZEjQZMEOWr&#43;wKSMCQQCC4kXJEsHAve77oP6HtG/IiEn7
kpyUXRNvFsDE0czpJJBvL/aRFUJxuRK91jhjC68sA7NsKMGg5OXb5I5Jj36xAkEA
gIT7aFOYBFwGgQAQkWNKLvySgKbAZRTeLBacpHMuQdl1DfdntvAyqpAZ0lY0RKmW
G6aFKaqQfOXKCyWoUiVknQJAXrlgySFci/2ueKlIE1QqIiLSZ8V8OlpFLRnb1pzI
7U1yQXnTAEFYM560yJlzUpOb1V4cScGd365tiSMvxLOvTA==
-----END RSA PRIVATE KEY-----`))
}

// parsePublicKey parses a PEM encoded private key.
func parsePrivateKey(pemBytes []byte) (Signer, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New(&#34;ssh: no key found&#34;)
	}

	var rawkey interface{}
	switch block.Type {
	case &#34;RSA PRIVATE KEY&#34;:
		rsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = rsa
	default:
		return nil, fmt.Errorf(&#34;ssh: unsupported key type %q&#34;, block.Type)
	}
	return newSignerFromKey(rawkey)
}

// A Signer is can create signatures that verify against a public key.
type Signer interface {
	// Sign returns raw signature for the given data. This method
	// will apply the hash specified for the keytype to the data.
	Sign(data []byte) ([]byte, error)
}

// A Signer is can create signatures that verify against a public key.
type Unsigner interface {
	// Sign returns raw signature for the given data. This method
	// will apply the hash specified for the keytype to the data.
	Unsign(data[]byte, sig []byte) error
}

func newSignerFromKey(k interface{}) (Signer, error) {
	var sshKey Signer
	switch t := k.(type) {
	case *rsa.PrivateKey:
		sshKey = &amp;rsaPrivateKey{t}
	default:
		return nil, fmt.Errorf(&#34;ssh: unsupported key type %T&#34;, k)
	}
	return sshKey, nil
}

func newUnsignerFromKey(k interface{}) (Unsigner, error) {
	var sshKey Unsigner
	switch t := k.(type) {
	case *rsa.PublicKey:
		sshKey = &amp;rsaPublicKey{t}
	default:
		return nil, fmt.Errorf(&#34;ssh: unsupported key type %T&#34;, k)
	}
	return sshKey, nil
}

type rsaPublicKey struct {
	*rsa.PublicKey
}

type rsaPrivateKey struct {
	*rsa.PrivateKey
}

// Sign signs data with rsa-sha256
func (r *rsaPrivateKey) Sign(data []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA256, d)
}

// Unsign verifies the message using a rsa-sha256 signature
func (r *rsaPublicKey) Unsign(message []byte, sig []byte) error {
	h := sha256.New()
	h.Write(message)
	d := h.Sum(nil)
	return rsa.VerifyPKCS1v15(r.PublicKey, crypto.SHA256, d, sig)
}