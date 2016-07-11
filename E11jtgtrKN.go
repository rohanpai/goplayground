package main

import (
	&#34;fmt&#34;
	&#34;log&#34;
	&#34;math/big&#34;
	&#34;net/http&#34;

	&#34;github.com/btcsuite/btcec&#34;
	&#34;github.com/btcsuite/btcnet&#34;
	&#34;github.com/btcsuite/btcutil&#34;
)

const ResultsPerPage = 128

const PageTemplateHeader = `&lt;html&gt;
&lt;head&gt;
	&lt;title&gt;All bitcoin private keys&lt;/title&gt;
	&lt;meta charset=&#34;utf8&#34; /&gt;
	&lt;link href=&#34;http://fonts.googleapis.com/css?family=Open&#43;Sans&#34; rel=&#34;stylesheet&#34; type=&#34;text/css&#34;&gt;
	&lt;style&gt;
		body{font-size: 9pt; font-family: &#39;Open Sans&#39;, sans-serif;}
		a{text-decoration: none}
		a:hover {text-decoration: underline}
		.keys &gt; span:hover { background: #f0f0f0; }
		span:target { background: #ccffcc; }
	&lt;/style&gt;
&lt;/head&gt;
&lt;body&gt;
&lt;h1&gt;Bitcoin private key database&lt;/h1&gt;
&lt;h2&gt;Page %s out of %s&lt;/h2&gt;
&lt;a href=&#34;/%s&#34;&gt;previous&lt;/a&gt; | &lt;a href=&#34;/%s&#34;&gt;next&lt;/a&gt;
&lt;pre class=&#34;keys&#34;&gt;
&lt;strong&gt;Private Key&lt;/strong&gt;                                            &lt;strong&gt;Address&lt;/strong&gt;                            &lt;strong&gt;Compressed Address&lt;/strong&gt;
`

const PageTemplateFooter = `&lt;/pre&gt;
&lt;pre style=&#34;margin-top: 1em; font-size: 8pt&#34;&gt;
It took a lot of computing power to generate this database. Donations welcome: 1Bv8dN7pemC5N3urfMDdAFReibefrBqCaK
&lt;/pre&gt;
&lt;a href=&#34;/%s&#34;&gt;previous&lt;/a&gt; | &lt;a href=&#34;/%s&#34;&gt;next&lt;/a&gt;
&lt;/body&gt;
&lt;/html&gt;`

const KeyTemplate = `&lt;span id=&#34;%s&#34;&gt;&lt;a href=&#34;/warning:understand-how-this-works!/%s&#34;&gt;&#43;&lt;/a&gt; &lt;span title=&#34;%s&#34;&gt;%s &lt;/span&gt; &lt;a href=&#34;https://blockchain.info/address/%s&#34;&gt;%34s&lt;/a&gt; &lt;a href=&#34;https://blockchain.info/address/%s&#34;&gt;%34s&lt;/a&gt;&lt;/span&gt;
`

var (
	// Total bitcoins
	total = new(big.Int).SetBytes([]byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE,
		0xBA, 0xAE, 0xDC, 0xE6, 0xAF, 0x48, 0xA0, 0x3B, 0xBF, 0xD2, 0x5E, 0x8C, 0xD0, 0x36, 0x41, 0x40,
	})

	// One
	one = big.NewInt(1)

	// Total pages
	_pages = new(big.Int).Div(total, big.NewInt(ResultsPerPage))
	pages  = _pages.Add(_pages, one)
)

type Key struct {
	private      string
	number       string
	compressed   string
	uncompressed string
}

func compute(count *big.Int) (keys [ResultsPerPage]Key, length int) {
	var padded [32]byte

	var i int
	for i = 0; i &lt; ResultsPerPage; i&#43;&#43; {
		// Increment our counter
		count.Add(count, one)

		// Check to make sure we&#39;re not out of range
		if count.Cmp(total) &gt; 0 {
			break
		}

		// Copy count value&#39;s bytes to padded slice
		copy(padded[32-len(count.Bytes()):], count.Bytes())

		// Get private and public keys
		privKey, public := btcec.PrivKeyFromBytes(btcec.S256(), padded[:])

		// Get compressed and uncompressed addresses for public key
		caddr, _ := btcutil.NewAddressPubKey(public.SerializeCompressed(), &amp;btcnet.MainNetParams)
		uaddr, _ := btcutil.NewAddressPubKey(public.SerializeUncompressed(), &amp;btcnet.MainNetParams)

		// Encode addresses
		wif, _ := btcutil.NewWIF(privKey, &amp;btcnet.MainNetParams, false)
		keys[i].private = wif.String()
		keys[i].number = count.String()
		keys[i].compressed = caddr.EncodeAddress()
		keys[i].uncompressed = uaddr.EncodeAddress()
	}
	return keys, i
}

func PageRequest(w http.ResponseWriter, r *http.Request) {
	// Default page is page 1
	if len(r.URL.Path) &lt;= 1 {
		r.URL.Path = &#34;/1&#34;
	}

	// Convert page number to bignum
	page, success := new(big.Int).SetString(r.URL.Path[1:], 0)
	if !success {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Make sure page number cannot be negative or 0
	page.Abs(page)
	if page.Cmp(one) == -1 {
		page.SetInt64(1)
	}

	// Make sure we&#39;re not above page count
	if page.Cmp(pages) &gt; 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get next and previous page numbers
	previous := new(big.Int).Sub(page, one)
	next := new(big.Int).Add(page, one)

	// Calculate our starting key from page number
	start := new(big.Int).Mul(previous, big.NewInt(ResultsPerPage))

	// Send page header
	fmt.Fprintf(w, PageTemplateHeader, page, pages, previous, next)

	// Send keys
	keys, length := compute(start)
	for i := 0; i &lt; length; i&#43;&#43; {
		key := keys[i]
		fmt.Fprintf(w, KeyTemplate, key.private, key.private, key.number, key.private, key.uncompressed, key.uncompressed, key.compressed, key.compressed)
	}

	// Send page footer
	fmt.Fprintf(w, PageTemplateFooter, previous, next)
}

func RedirectRequest(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[36:]

	wif, err := btcutil.DecodeWIF(key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	page, _ := new(big.Int).DivMod(new(big.Int).SetBytes(wif.PrivKey.D.Bytes()), big.NewInt(ResultsPerPage), big.NewInt(ResultsPerPage))
	page.Add(page, one)

	fragment, _ := btcutil.NewWIF(wif.PrivKey, &amp;btcnet.MainNetParams, false)

	http.Redirect(w, r, &#34;/&#34;&#43;page.String()&#43;&#34;#&#34;&#43;fragment.String(), http.StatusTemporaryRedirect)
}

func main() {
	http.HandleFunc(&#34;/&#34;, PageRequest)
	http.HandleFunc(&#34;/warning:understand-how-this-works!/&#34;, RedirectRequest)

	log.Println(&#34;Listening&#34;)
	log.Fatal(http.ListenAndServe(&#34;:8085&#34;, nil))
}