package main

import (
    &#34;fmt&#34;
    &#34;encoding/base64&#34;
    &#34;encoding/hex&#34;
    &#34;log&#34;
    )

const (
    cipherHex = &#34;1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736&#34;
    )

func StringToHex(s string) []byte {
    my_hex, err := hex.DecodeString(s)
    if err != nil {
        log.Fatal(&#34;string not converted to hex&#34;)
    }
    return []byte(my_hex)
}

func HexTo64(b []byte) string {
    return base64.StdEncoding.EncodeToString(b)
}

func HexToString(b []byte) string {
    return hex.EncodeToString(b)
}

func FixedXOR(a,b []byte) []byte {
    if len(a) == len(b) {
        c := a
        for i:=0; i&lt;len(a); i&#43;&#43; {
            c[i] = a[i] ^ b[i]
        }
        return c
    }
    log.Fatal(&#34;buffers are not equal length&#34;)
    return nil
}

func MakeXORString(s string, XORlength int) string {
    var key string
    for ;len(key)&lt;XORlength; {
        if len(s) &#43; len(key) &gt; XORlength {
            key &#43;= s[0:(XORlength-len(key))]
        } else {
            key &#43;= s
        }
    }
    return key
}

func MakeXORHex (s string, XORlength int) []byte {
    t := []byte(MakeXORString(s,XORlength))
    return t
}

func ScoreKey(b []byte) int {
    score := int(0)
    for _, char := range b {
        if char&gt;64 &amp;&amp; char&lt;123 &amp;&amp; (char&lt;97 || char&gt;90) {
            score&#43;&#43;
        }
    }
//    fmt.Println(score)
    return score
}

func FindSingleXORKey(cipherhex []byte) []byte {
    t := []byte(&#34;G&#34;)
    for i:=0;i&lt;=255;i&#43;&#43; {
    //loop through ASCII
        a := ScoreKey(t)
        fmt.Printf(&#34;cached plaintext for loop %d is %s with score %d&#34;, i, t, a)
        fmt.Println()
        
        tempkey := MakeXORHex(string(i), len(cipherhex))
        //Make a repeating key at each iteration

        attempt := FixedXOR(cipherhex, tempkey)
        //attempt to decode using the generated temp key
        
        b := ScoreKey(attempt)
        fmt.Printf(&#34;Score of new plaintext %s is %d&#34;, attempt, b)
        fmt.Println()


        if a &lt; b {
        //is the score for the decryption attempt higher than the score for the cached attempt?

            t = attempt
            fmt.Printf(&#34;New string %s score is higher&#34;, attempt)
            fmt.Println()

        } else {
            fmt.Printf(&#34;keeping %s&#34;, t)
            fmt.Println()
        }
    }
    return t
}

func main() {
    text := FindSingleXORKey(StringToHex(cipherHex))
    fmt.Println(string(text))
}
