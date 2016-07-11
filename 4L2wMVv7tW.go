// Go offers built-in support for JSON encoding and
// decoding, including to and from built-in and custom
// data types.

package main

import &#34;encoding/json&#34;
import &#34;fmt&#34;
import &#34;os&#34;

// We&#39;ll use these two structs to demonstrate encoding and
// decoding of custom types below.
type Response1 struct {
    Page   int
    Fruits []string
}
type Response2 struct {
    Page   int      `json:&#34;page&#34;`
    Fruits []string `json:&#34;fruits&#34;`
}

func main() {

    // First we&#39;ll look at encoding basic data types to
    // JSON strings. Here are some examples for atomic
    // values.
    bolB, _ := json.Marshal(true)
    fmt.Println(string(bolB))

    intB, _ := json.Marshal(1)
    fmt.Println(string(intB))

    fltB, _ := json.Marshal(2.34)
    fmt.Println(string(fltB))

    strB, _ := json.Marshal(&#34;gopher&#34;)
    fmt.Println(string(strB))

    // And here are some for slices and maps, which encode
    // to JSON arrays and objects as you&#39;d expect.
    slcD := []string{&#34;apple&#34;, &#34;peach&#34;, &#34;pear&#34;}
    slcB, _ := json.Marshal(slcD)
    fmt.Println(string(slcB))

    mapD := map[string]int{&#34;apple&#34;: 5, &#34;lettuce&#34;: 7}
    mapB, _ := json.Marshal(mapD)
    fmt.Println(string(mapB))

    // The JSON package can automatically encode your
    // custom data types. It will only include exported
    // fields in the encoded output and will by default
    // use those names as the JSON keys.
    res1D := &amp;Response1{
        Page:   1,
        Fruits: []string{&#34;apple&#34;, &#34;peach&#34;, &#34;pear&#34;}}
    res1B, _ := json.Marshal(res1D)
    fmt.Println(string(res1B))

    // You can use tags on struct field declarations
    // to customize the encoded JSON key names. Check the
    // definition of `Response2` above to see an example
    // of such tags.
    res2D := &amp;Response2{
        Page:   1,
        Fruits: []string{&#34;apple&#34;, &#34;peach&#34;, &#34;pear&#34;}}
    res2B, _ := json.Marshal(res2D)
    fmt.Println(string(res2B))

    // Now let&#39;s look at decoding JSON data into Go
    // values. Here&#39;s an example for a generic data
    // structure.
    byt := []byte(`{&#34;num&#34;:6.0,&#34;strs&#34;:[&#34;a&#34;,&#34;b&#34;]}`)

    // We need to provide a variable where the JSON
    // package can put the decoded data. This
    // `map[string]interface{}` will hold a map of strings
    // to arbitrary data types.
    var dat map[string]interface{}

    // Here&#39;s the actual decoding, and a check for
    // associated errors.
    if err := json.Unmarshal(byt, &amp;dat); err != nil {
        panic(err)
    }
    fmt.Println(dat)

    // In order to use the values in the decoded map,
    // we&#39;ll need to cast them to their appropriate type.
    // For example here we cast the value in `num` to
    // the expected `float64` type.
    num := dat[&#34;num&#34;].(float64)
    fmt.Println(num)

    // Accessing nested data requires a series of
    // casts.
    strs := dat[&#34;strs&#34;].([]interface{})
    str1 := strs[0].(string)
    fmt.Println(str1)

    // We can also decode JSON into custom data types.
    // This has the advantages of adding additional
    // type-safety to our programs and eliminating the
    // need for type assertions when accessing the decoded
    // data.
    str := `{&#34;page&#34;: 1, &#34;fruits&#34;: [&#34;apple&#34;, &#34;peach&#34;]}`
    res := &amp;Response2{}
    json.Unmarshal([]byte(str), &amp;res)
    fmt.Println(res)
    fmt.Println(res.Fruits[0])

    // In the examples above we always used bytes and
    // strings as intermediates between the data and
    // JSON representation on standard out. We can also
    // stream JSON encodings directly to `os.Writer`s like
    // `os.Stdout` or even HTTP response bodies.
    enc := json.NewEncoder(os.Stdout)
    d := map[string]int{&#34;apple&#34;: 5, &#34;lettuce&#34;: 7}
    enc.Encode(d)
}
