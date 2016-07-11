package main

import (
    "fmt"
    "encoding/base64"
    "encoding/hex"
    "log"
    )

const (
    cipherHex = "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
    )

func StringToHex(s string) []byte {
    my_hex, err := hex.DecodeString(s)
    if err != nil {
        log.Fatal("string not converted to hex")
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
        for i:=0; i<len(a); i++ {
            c[i] = a[i] ^ b[i]
        }
        return c
    }
    log.Fatal("buffers are not equal length")
    return nil
}

func MakeXORString(s string, XORlength int) string {
    var key string
    for ;len(key)<XORlength; {
        if len(s) + len(key) > XORlength {
            key += s[0:(XORlength-len(key))]
        } else {
            key += s
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
        if char>64 && char<123 && (char<97 || char>90) {
            score++
        }
    }
//    fmt.Println(score)
    return score
}

func FindSingleXORKey(cipherhex []byte) []byte {
    t := []byte("G")
    for i:=0;i<=255;i++ {
    //loop through ASCII
        a := ScoreKey(t)
        fmt.Printf("cached plaintext for loop %d is %s with score %d", i, t, a)
        fmt.Println()
        
        tempkey := MakeXORHex(string(i), len(cipherhex))
        //Make a repeating key at each iteration

        attempt := FixedXOR(cipherhex, tempkey)
        //attempt to decode using the generated temp key
        
        b := ScoreKey(attempt)
        fmt.Printf("Score of new plaintext %s is %d", attempt, b)
        fmt.Println()


        if a < b {
        //is the score for the decryption attempt higher than the score for the cached attempt?

            t = attempt
            fmt.Printf("New string %s score is higher", attempt)
            fmt.Println()

        } else {
            fmt.Printf("keeping %s", t)
            fmt.Println()
        }
    }
    return t
}

func main() {
    text := FindSingleXORKey(StringToHex(cipherHex))
    fmt.Println(string(text))
}
