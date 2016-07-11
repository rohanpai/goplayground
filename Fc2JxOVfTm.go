package main

import (
	&#34;bufio&#34;
	&#34;fmt&#34;
	&#34;io&#34;
	&#34;strings&#34;
)

func main() {
	fmap := Frequency(strings.NewReader(&#34;Francisco, Francisco&#34;))
	fmt.Println(Suggest(&#34;Fransisco&#34;, fmap))
	// Francisco
}

// Frequency counts the frequency of each word.
func Frequency(reader io.Reader) map[string]int {
	scanner := bufio.NewScanner(reader)
	//
	// This must be called before Scan.
	// The default split function is bufio.ScanLines.
	scanner.Split(bufio.ScanWords)
	//
	fmap := make(map[string]int)
	//
	for scanner.Scan() {
		// Remove all leading and trailing Unicode code points.
		word := strings.Trim(scanner.Text(), &#34;,-!;:\&#34;?.&#34;)
		if _, exist := fmap[word]; exist {
			fmap[word]&#43;&#43;
		} else {
			fmap[word] = 1
		}
	}
	return fmap
}

// distanceOne sends all possible corrections
// with edit distance 1 to the channel one.
// This is much more probable than the one with 2 edit distance.
func distanceOne(txt string, one chan string) {
	const alphabet = &#34;abcdefghijklmnopqrstuvwxyz&#34;
	type pair struct {
		front, back string
	}
	pairs := []pair{}
	for i := 0; i &lt;= len(txt); i&#43;&#43; {
		pairs = append(pairs, pair{txt[:i], txt[i:]})
	}
	for _, pair := range pairs {
		// deletion of pair.back[0]
		if len(pair.back) &gt; 0 {
			one &lt;- pair.front &#43; pair.back[1:]
		}
		// transpose of pair.back[0] and pair.back[1]
		if len(pair.back) &gt; 1 {
			one &lt;- pair.front &#43; string(pair.back[1]) &#43; string(pair.back[0]) &#43; pair.back[2:]
		}
		// replace of pair.back[0]
		for _, elem := range alphabet {
			if len(pair.back) &gt; 0 {
				one &lt;- pair.front &#43; string(elem) &#43; pair.back[1:]
			}
		}
		// insertion
		for _, elem := range alphabet {
			one &lt;- pair.front &#43; string(elem) &#43; pair.back
		}
	}
}

// distanceMore sends other possible corrections
// based on the results from distanceOne.
func distanceMore(word string, other chan string) {
	one := make(chan string, 1024*1024)
	go func() {
		distanceOne(word, one)
		close(one)
	}()
	// retrieve from distanceOne results and break when it&#39;s done
	for v := range one {
		// run distanceOne in addition to the results from the first distanceOne
		distanceOne(v, other)
	}
}

// known returns the word with maximum frequencies.
func known(txt string, distFunc func(string, chan string), fmap map[string]int) string {
	words := make(chan string, 1024*1024)
	go func() {
		distFunc(txt, words)
		close(words)
	}()
	maxFq := 0
	suggest := &#34;&#34;
	for wd := range words {
		if freq, exist := fmap[wd]; exist &amp;&amp; freq &gt; maxFq {
			maxFq, suggest = freq, wd
		}
	}
	return suggest
}

// Suggest suggests the correct spelling based on the sample data.
func Suggest(txt string, fmap map[string]int) string {
	// edit distance 0
	if _, exist := fmap[txt]; exist {
		return txt
	}
	if v := known(txt, distanceOne, fmap); v != &#34;&#34; {
		return v
	}
	if v := known(txt, distanceMore, fmap); v != &#34;&#34; {
		return v
	}
	// edit distance 3, 4, and more ...
	return txt
}
