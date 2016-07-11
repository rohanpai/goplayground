// [_Command-line flags_](http://en.wikipedia.org/wiki/Command-line_interface#Command-line_option)
// are a common way to specify options for command-line
// programs. For example, in `wc -l` the `-l` is a
// command-line flag.

package main

// Go provides a `flag` package supporting basic
// command-line flag parsing. We&#39;ll use this package to
// implement our example command-line program.
import &#34;flag&#34;
import &#34;fmt&#34;

func main() {

    // Basic flag declarations are available for string,
    // integer, and boolean options. Here we declare a
    // string flag `word` with a default value `&#34;foo&#34;`
    // and a short description. This `flag.String` function
    // returns a string pointer (not a string value);
    // we&#39;ll see how to use this pointer below.
    wordPtr := flag.String(&#34;word&#34;, &#34;foo&#34;, &#34;a string&#34;)

    // This declares `numb` and `fork` flags, using a
    // similar approach to the `word` flag.
    numbPtr := flag.Int(&#34;numb&#34;, 42, &#34;an int&#34;)
    boolPtr := flag.Bool(&#34;fork&#34;, false, &#34;a bool&#34;)

    // It&#39;s also possible to declare an option that uses an
    // existing var declared elsewhere in the program.
    // Note that we need to pass in a pointer to the flag
    // declaration function.
    var svar string
    flag.StringVar(&amp;svar, &#34;svar&#34;, &#34;bar&#34;, &#34;a string var&#34;)

    // Once all flags are declared, call `flag.Parse()`
    // to execute the command-line parsing.
    flag.Parse()

    // Here we&#39;ll just dump out the parsed options and
    // any trailing positional arguments. Note that we
    // need to dereference the points with e.g. `*wordPtr`
    // to get the actual option values.
    fmt.Println(&#34;word:&#34;, *wordPtr)
    fmt.Println(&#34;numb:&#34;, *numbPtr)
    fmt.Println(&#34;fork:&#34;, *boolPtr)
    fmt.Println(&#34;svar:&#34;, svar)
    fmt.Println(&#34;tail:&#34;, flag.Args())
}
