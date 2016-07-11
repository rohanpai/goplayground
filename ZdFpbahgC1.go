// Go&#39;s `math/rand` package provides
// [pseudorandom number](http://en.wikipedia.org/wiki/Pseudorandom_number_generator)
// generation.

package main

import &#34;time&#34;
import &#34;fmt&#34;
import &#34;math/rand&#34;

func main() {

    // For example, `rand.Intn` returns a random `int` n,
    // `0 &lt;= n &lt; 100`.
    fmt.Print(rand.Intn(100), &#34;,&#34;)
    fmt.Print(rand.Intn(100))
    fmt.Println()

    // `rand.Float64` returns a `float64` `f`,
    // `0.0 &lt;= f &lt; 1.0`.
    fmt.Println(rand.Float64())

    // This can be used to generate random floats in
    // other ranges, for example `5.0 &lt;= f&#39; &lt; 10.0`.
    fmt.Print((rand.Float64()*5)&#43;5, &#34;,&#34;)
    fmt.Print((rand.Float64() * 5) &#43; 5)
    fmt.Println()

    // The default number generator is deterministic, so it&#39;ll
    // produce the same sequence of numbers each time by default.
    // To produce varying sequences, give it a seed that changes.
    // Note that this is not safe to use for random numbers you
    // intend to be secret, use `crypto/rand` for those.
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)

    // Call the resulting `rand.Rand` just like the
    // functions on the `rand` package.
    fmt.Print(r1.Intn(100), &#34;,&#34;)
    fmt.Print(r1.Intn(100))
    fmt.Println()

    // If you seed a source with the same number, it
    // produces the same sequence of random numbers.
    s2 := rand.NewSource(42)
    r2 := rand.New(s2)
    fmt.Print(r2.Intn(100), &#34;,&#34;)
    fmt.Print(r2.Intn(100))
    fmt.Println()
    s3 := rand.NewSource(42)
    r3 := rand.New(s3)
    fmt.Print(r3.Intn(100), &#34;,&#34;)
    fmt.Print(r3.Intn(100))
}