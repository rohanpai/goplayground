package main

import (
	&#34;bufio&#34;
	&#34;fmt&#34;
	&#34;io&#34;
	&#34;strings&#34;
)

func consumeWord(data []byte) (int, []byte, error) {
	var accum []byte
	for i, b := range data {
		if b == &#39; &#39; || b == &#39;\n&#39; || b == &#39;\t&#39; || b == &#39;\r&#39; {
			return i, accum, nil
		} else {
			accum = append(accum, b)
		}
	}
	return 0, nil, nil
}

func consumeWhitespace(data []byte) (int, []byte, error) {
	var accum []byte
	for i, b := range data {
		if b == &#39; &#39; || b == &#39;\n&#39; || b == &#39;\t&#39; || b == &#39;\r&#39; {
			accum = append(accum, b)
		} else {
			return i, accum, nil
		}
	}
	return 0, nil, nil
}

func consumeNum(data []byte) (int, []byte, error) {
	var accum []byte
	for i, b := range data {
		if &#39;0&#39; &lt;= b &amp;&amp; b &lt;= &#39;9&#39; {
			accum = append(accum, b)
		} else {
			return i, accum, nil
		}
	}
	return len(accum), accum, nil
}

func consumeString(data []byte) (int, []byte, error) {
	delim := data[0]
	skip := false
	accum := []byte{data[0]}
	for i, b := range data[1:] {
		if b == delim &amp;&amp; !skip {
			return i &#43; 2, accum, nil
		}
		skip = false
		if b == &#39;\\&#39; {
			skip = true
			continue
		}
		accum = append(accum, b)
	}
	return 0, nil, nil
}

type lineTrackingReader struct {
	*bufio.Scanner
	lastL, lastCol int
	l, col         int
}

func newTrackingReader(r io.Reader) *lineTrackingReader {
	s := bufio.NewScanner(r)
	rdr := &amp;lineTrackingReader{
		Scanner: s,
	}
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if rdr.l == 0 {
			rdr.l = 1
			rdr.col = 1
			rdr.lastL = 1
			rdr.lastCol = 1
		}
		switch data[0] {
		case &#39;(&#39;, &#39;)&#39;:
			advance, token, err = 1, data[:1], nil
		case &#39;&#34;&#39;, &#39;\&#39;&#39;: // TODO(jwall): Rune data?
			advance, token, err = consumeString(data)
		case &#39;0&#39;, &#39;1&#39;, &#39;2&#39;, &#39;3&#39;, &#39;4&#39;, &#39;5&#39;, &#39;6&#39;, &#39;7&#39;, &#39;8&#39;, &#39;9&#39;:
			advance, token, err = consumeNum(data)
		case &#39; &#39;, &#39;\n&#39;, &#39;\r&#39;, &#39;\t&#39;:
			advance, token, err = consumeWhitespace(data)
		default:
			advance, token, err = consumeWord(data)
		}
		if advance &gt; 0 {
			rdr.lastCol = rdr.col
			rdr.lastL = rdr.l
			for _, b := range data[:advance] {
				if b == &#39;\n&#39; || atEOF {
					rdr.l&#43;&#43;
					rdr.col = 1
				}
				rdr.col&#43;&#43;
			}
		}

		return
	}
	s.Split(split)
	return rdr
}

func main() {
	s := newTrackingReader(strings.NewReader(&#34;(foo 1 &#39;bar&#39;)\n(baz 2 &#39;quux&#39;)&#34;))

	for s.Scan() {
		tok := s.Bytes()
		fmt.Printf(&#34;Found token %q at line #%d col #%d\n&#34;, string(tok), s.lastL, s.lastCol)
	}
}
