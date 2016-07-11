package main

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha1"
    "crypto/x509"
    "encoding/base32"
    //"encoding/pem"
    "flag"
    "fmt"
    "os"
    "runtime/pprof"
    "strings"
)

const encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
var prefix *string = flag.String("prefix", "foo", "Prefix string to search for (e.g. 'foo')" )
var cpuprofile *string = flag.String("cpuprofile", "", "write cpu profile to file")
var ncpus *int = flag.Int("ncpu", 1, "number of cpu's to use")

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
        ch <- kp
    }   
}

func HashCert( prefix string, ch chan KeyPair, done chan bool ) { 

    var kp KeyPair

    fmt.Println("HashCert(): starting with prefix: ", prefix )

    i := 0
    for {

        i++ 

        // get the key pair
        kp = <-ch

        // create the sha1 hasher
        h := sha1.New()

        // hash the DER representation of the public key starting at 23'rd byte
        pubBytes, _ := x509.MarshalPKIXPublicKey( &kp.key.PublicKey ) 
        h.Write( pubBytes[22:] )

        // get the hash
        kp.hash = h.Sum(nil)

        //fmt.Println( "sha1: ", hex.EncodeToString( kp.hash ) )

        // create the base32 encoder
        e := base32.NewEncoding( encodeStd )

        // now base32 encode the first 20 bytes of the hash
        kp.onion = strings.ToLower( e.EncodeToString( kp.hash[:20] ) )[:16]
        //fmt.Println(kp.onion, ".onion")

        if ( strings.HasPrefix( kp.onion, prefix ) ) {

            fmt.Printf( "\nfound match: %s.onion\n", kp.onion )

            // save the key pair
            //keyOut, _ := os.OpenFile( kp.onion + ".private_key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600 )
            //pem.Encode( keyOut, &pem.Block{ Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey( kp.key ) } )
            //keyOut.Close()

            // signal that we're done and return to stop cert processing
            done <- true
            return

        } else {

            if (i % 10 == 0) {
                fmt.Print(".")
            }
            if (i % 800 == 0) {
                i = 0
                fmt.Print("\n")
            }
        }
    }
}


func main() {

    // parse command line flags
    flag.Parse()

    // start cpu profiling
    if *cpuprofile != "" {
        f, _ := os.Create(*cpuprofile)
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }

    // create the channel
    ch := make(chan KeyPair)
    done := make(chan bool)

    // start the create cert goroutine
    for i := 0; i < *ncpus; i++ {
        fmt.Println("CreateCert goroutine")
        go CreateCert( ch )
    }

    // start the hash cert goroutine
    go HashCert( *prefix, ch, done )

    // wait for match
    <- done
}
