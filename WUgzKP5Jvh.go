package master

import (
        "golem/facter"
        "log"
        "net/url"
        "net/http"
        "crypto/tls"
        "crypto/x509"
        "encoding/pem"
        "encoding/json"
        "io/ioutil"
)

const (
        catalogURL = "https://localhost:8140/production/catalog/agent"
        cert       = "/Users/daniel/.puppet/etc-agent/ssl/certs/agent.pem"
        key        = "/Users/daniel/.puppet/etc-agent/ssl/private_keys/agent.pem"
        cacert     = "/Users/daniel/.puppet/etc-agent/ssl/certs/ca.pem"
)

type Catalog struct {
}

func getCACertPool(logger *log.Logger) *x509.CertPool {
        caPEM, err := ioutil.ReadFile(cacert)
        if err != nil {
                logger.Fatalf("Can't read CA cert: %v", err)
        }

        // REVISIT: This totally ignores the possibility of a second
        // certificate in that PEM code, for now.
        caDERBlock, caPEM := pem.Decode(caPEM)
        if caDERBlock.Type != "CERTIFICATE" {
                logger.Fatalf("CA cert is not a certificate: %v", caDERBlock)
        }

        caCert, err := x509.ParseCertificate(caDERBlock.Bytes)
        if err != nil {
                logger.Fatalf("Can't parse CA cert: %v", err)
        }

        caPool := x509.NewCertPool()
        caPool.AddCert(caCert)
        return caPool
}

func FetchCatalog(logger *log.Logger, facts facter.Facts) *Catalog {
        clientCert, err := tls.LoadX509KeyPair(cert, key)
        if err != nil {
                logger.Fatalf("Can't load clientCert: %v", err)
        }

        CAPool := getCACertPool(logger)

        tlsConfig := &tls.Config{Certificates: []tls.Certificate{clientCert}, RootCAs: CAPool, InsecureSkipVerify: false}

        transport := &http.Transport{TLSClientConfig: tlsConfig, DisableKeepAlives: true}
        client := &http.Client{Transport: transport}

        factJSON, err := json.Marshal(facts)
        if err != nil {
                logger.Fatalf("Can't encode facts to JSON: %v", err)
        }
        factBody := url.Values{"facts_format": []string{"pson"}, "facts": []string{string(factJSON)}}
        response, err := client.PostForm(catalogURL, factBody)
        if err != nil {
                logger.Fatalf("Can't fetch catalog: %v", err)
        }

        log.Printf("Got response of %d bytes, CTE %v, headers %v", response.ContentLength, response.TransferEncoding, response.Header)

        return &Catalog{}
}
