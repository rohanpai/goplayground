package main

import (
    &#34;crypto/rand&#34;
    &#34;crypto/rsa&#34;
    &#34;crypto/sha1&#34;
    &#34;crypto/x509&#34;
    &#34;encoding/base32&#34;
    //&#34;encoding/pem&#34;
    &#34;flag&#34;
    &#34;fmt&#34;
    &#34;os&#34;
    &#34;runtime/pprof&#34;
    &#34;strings&#34;
)

const encodeStd = &#34;ABCDEFGHIJKLMNOPQRSTUVWXYZ234567&#34;
var prefix *string = flag.String(&#34;prefix&#34;, &#34;foo&#34;, &#34;Prefix string to search for (e.g. &#39;foo&#39;)&#34; )
var cpuprofile *string = flag.String(&#34;cpuprofile&#34;, &#34;&#34;, &#34;write cpu profile to file&#34;)
var ncpus *int = flag.Int(&#34;ncpu&#34;, 1, &#34;number of cpu&#39;s to use&#34;)

type KeyPair struct {
    key *rsa.PrivateKey
    hash []byte
    onion string
}

func CreateCert( ch chan KeyPair ) { 

    for {
        // create the new key pair
        kp := KeyPair{}

        // generate a new RSA key pair
        kp.key, _ = rsa.GenerateKey( rand.Reader, 1024 )

        // send the data to the hasher
        ch &lt;- kp
    }   
}

func HashCert( prefix string, ch chan KeyPair, done chan bool ) { 

    var kp KeyPair

    fmt.Println(&#34;HashCert(): starting with prefix: &#34;, prefix )

    i := 0
    for {

        i&#43;&#43; 

        // get the key pair
        kp = &lt;-ch

        // create the sha1 hasher
        h := sha1.New()

        // hash the DER representation of the public key starting at 23&#39;rd byte
        pubBytes, _ := x509.MarshalPKIXPublicKey( &amp;kp.key.PublicKey ) 
        h.Write( pubBytes[22:] )

        // get the hash
        kp.hash = h.Sum(nil)

        //fmt.Println( &#34;sha1: &#34;, hex.EncodeToString( kp.hash ) )

        // create the base32 encoder
        e := base32.NewEncoding( encodeStd )

        // now base32 encode the first 20 bytes of the hash
        kp.onion = strings.ToLower( e.EncodeToString( kp.hash[:20] ) )[:16]
        //fmt.Println(kp.onion, &#34;.onion&#34;)

        if ( strings.HasPrefix( kp.onion, prefix ) ) {

            fmt.Printf( &#34;\nfound match: %s.onion\n&#34;, kp.onion )

            // save the key pair
            //keyOut, _ := os.OpenFile( kp.onion &#43; &#34;.private_key&#34;, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600 )
            //pem.Encode( keyOut, &amp;pem.Block{ Type: &#34;RSA PRIVATE KEY&#34;, Bytes: x509.MarshalPKCS1PrivateKey( kp.key ) } )
            //keyOut.Close()

            // signal that we&#39;re done and return to stop cert processing
            done &lt;- true
            return

        } else {

            if (i % 10 == 0) {
                fmt.Print(&#34;.&#34;)
            }
            if (i % 800 == 0) {
                i = 0
                fmt.Print(&#34;\n&#34;)
            }
        }
    }
}


func main() {

    // parse command line flags
    flag.Parse()

    // start cpu profiling
    if *cpuprofile != &#34;&#34; {
        f, _ := os.Create(*cpuprofile)
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }

    // create the channel
    ch := make(chan KeyPair)
    done := make(chan bool)

    // start the create cert goroutine
    for i := 0; i &lt; *ncpus; i&#43;&#43; {
        fmt.Println(&#34;CreateCert goroutine&#34;)
        go CreateCert( ch )
    }

    // start the hash cert goroutine
    go HashCert( *prefix, ch, done )

    // wait for match
    &lt;- done
}
