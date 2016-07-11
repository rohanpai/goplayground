package master

import (
        &#34;golem/facter&#34;
        &#34;log&#34;
        &#34;net/url&#34;
        &#34;net/http&#34;
        &#34;crypto/tls&#34;
        &#34;crypto/x509&#34;
        &#34;encoding/pem&#34;
        &#34;encoding/json&#34;
        &#34;io/ioutil&#34;
)

const (
        catalogURL = &#34;https://localhost:8140/production/catalog/agent&#34;
        cert       = &#34;/Users/daniel/.puppet/etc-agent/ssl/certs/agent.pem&#34;
        key        = &#34;/Users/daniel/.puppet/etc-agent/ssl/private_keys/agent.pem&#34;
        cacert     = &#34;/Users/daniel/.puppet/etc-agent/ssl/certs/ca.pem&#34;
)

type Catalog struct {
}

func getCACertPool(logger *log.Logger) *x509.CertPool {
        caPEM, err := ioutil.ReadFile(cacert)
        if err != nil {
                logger.Fatalf(&#34;Can&#39;t read CA cert: %v&#34;, err)
        }

        // REVISIT: This totally ignores the possibility of a second
        // certificate in that PEM code, for now.
        caDERBlock, caPEM := pem.Decode(caPEM)
        if caDERBlock.Type != &#34;CERTIFICATE&#34; {
                logger.Fatalf(&#34;CA cert is not a certificate: %v&#34;, caDERBlock)
        }

        caCert, err := x509.ParseCertificate(caDERBlock.Bytes)
        if err != nil {
                logger.Fatalf(&#34;Can&#39;t parse CA cert: %v&#34;, err)
        }

        caPool := x509.NewCertPool()
        caPool.AddCert(caCert)
        return caPool
}

func FetchCatalog(logger *log.Logger, facts facter.Facts) *Catalog {
        clientCert, err := tls.LoadX509KeyPair(cert, key)
        if err != nil {
                logger.Fatalf(&#34;Can&#39;t load clientCert: %v&#34;, err)
        }

        CAPool := getCACertPool(logger)

        tlsConfig := &amp;tls.Config{Certificates: []tls.Certificate{clientCert}, RootCAs: CAPool, InsecureSkipVerify: false}

        transport := &amp;http.Transport{TLSClientConfig: tlsConfig, DisableKeepAlives: true}
        client := &amp;http.Client{Transport: transport}

        factJSON, err := json.Marshal(facts)
        if err != nil {
                logger.Fatalf(&#34;Can&#39;t encode facts to JSON: %v&#34;, err)
        }
        factBody := url.Values{&#34;facts_format&#34;: []string{&#34;pson&#34;}, &#34;facts&#34;: []string{string(factJSON)}}
        response, err := client.PostForm(catalogURL, factBody)
        if err != nil {
                logger.Fatalf(&#34;Can&#39;t fetch catalog: %v&#34;, err)
        }

        log.Printf(&#34;Got response of %d bytes, CTE %v, headers %v&#34;, response.ContentLength, response.TransferEncoding, response.Header)

        return &amp;Catalog{}
}
