// Go offers excellent support for string formatting in
// the `printf` tradition. Here are some examples of
// common string formatting tasks.

package main

import &#34;fmt&#34;
import &#34;os&#34;

type point struct {
    x, y int
}

func main() {

    // Go offers several printing &#34;verbs&#34; designed to
    // format general Go values. For example, this prints
    // an instance of our `point` struct.
    p := point{1, 2}
    fmt.Printf(&#34;%v\n&#34;, p)

    // If the value is a struct, the `%&#43;v` variant will
    // include the struct&#39;s field names.
    fmt.Printf(&#34;%&#43;v\n&#34;, p)

    // The `%#v` variant prints a Go syntax representation
    // of the value, i.e. the source code snippet that
    // would produce that value.
    fmt.Printf(&#34;%#v\n&#34;, p)

    // To print the type of a value, use `%T`.
    fmt.Printf(&#34;%T\n&#34;, p)

    // Formatting booleans is straight-forward.
    fmt.Printf(&#34;%t\n&#34;, true)

    // There are many options for formatting integers.
    // Use `%d` for standard, base-10 formatting.
    fmt.Printf(&#34;%d\n&#34;, 123)

    // This prints a binary representation.
    fmt.Printf(&#34;%b\n&#34;, 14)

    // This prints the character corresponding to the
    // given integer.
    fmt.Printf(&#34;%c\n&#34;, 33)

    // `%x` provides hex encoding.
    fmt.Printf(&#34;%x\n&#34;, 456)

    // There are also several formatting options for
    // floats. For basic decimal formatting use `%f`.
    fmt.Printf(&#34;%f\n&#34;, 78.9)

    // `%e` and `%E` format the float in (slightly
    // different versions of) scientific notation.
    fmt.Printf(&#34;%e\n&#34;, 123400000.0)
    fmt.Printf(&#34;%E\n&#34;, 123400000.0)

    // For basic string printing use `%s`.
    fmt.Printf(&#34;%s\n&#34;, &#34;\&#34;string\&#34;&#34;)

    // To double-quote strings as in Go source, use `%q`.
    fmt.Printf(&#34;%q\n&#34;, &#34;\&#34;string\&#34;&#34;)

    // As with integers as seen earlier, `%x` renders
    // the string in base-16, with two output characters
    // per byte of input.
    fmt.Printf(&#34;%x\n&#34;, &#34;hex this&#34;)

    // To print a representation of a pointer, use `%p`.
    fmt.Printf(&#34;%p\n&#34;, &amp;p)

    // When formatting numbers you will often want to
    // control the width and precision of the resulting
    // figure. To specify the width of an integer, use a
    // number after the `%` in the verb. By default the
    // result will be right-justified and padded with
    // spaces.
    fmt.Printf(&#34;|%6d|%6d|\n&#34;, 12, 345)

    // You can also specify the width of printed floats,
    // though usually you&#39;ll also want to restrict the
    // decimal precision at the same time with the
    // width.precision syntax.
    fmt.Printf(&#34;|%6.2f|%6.2f|\n&#34;, 1.2, 3.45)

    // To left-justify, use the `-` flag.
    fmt.Printf(&#34;|%-6.2f|%-6.2f|\n&#34;, 1.2, 3.45)

    // You may also want to control width when formatting
    // strings, especially to ensure that they align in
    // table-like output. For basic right-justified width.
    fmt.Printf(&#34;|%6s|%6s|\n&#34;, &#34;foo&#34;, &#34;b&#34;)

    // To left-justify use the `-` flag as with numbers.
    fmt.Printf(&#34;|%-6s|%-6s|\n&#34;, &#34;foo&#34;, &#34;b&#34;)

    // So far we&#39;ve seen `Printf`, which prints the
    // formatted string to `os.Stdout`. `Sprintf` formats
    // and returns a string without printing it anywhere.
    s := fmt.Sprintf(&#34;a %s&#34;, &#34;string&#34;)
    fmt.Println(s)

    // You can format&#43;print to `io.Writers` other than
    // `os.Stdout` using `Fprintf`.
    fmt.Fprintf(os.Stderr, &#34;an %s\n&#34;, &#34;error&#34;)
}
