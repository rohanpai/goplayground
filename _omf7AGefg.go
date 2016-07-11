package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func consumeWord(data []byte) (int, []byte, error) {
	var accum []byte
	for i, b := range data {
		if b == ' ' || b == '\n' || b == '\t' || b == '\r' {
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
		if b == ' ' || b == '\n' || b == '\t' || b == '\r' {
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
		if '0' <= b && b <= '9' {
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
		if b == delim && !skip {
			return i + 2, accum, nil
		}
		skip = false
		if b == '\\' {
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
	rdr := &lineTrackingReader{
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
		case '(', ')':
			advance, token, err = 1, data[:1], nil
		case '"', '\'': // TODO(jwall): Rune data?
			advance, token, err = consumeString(data)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			advance, token, err = consumeNum(data)
		case ' ', '\n', '\r', '\t':
			advance, token, err = consumeWhitespace(data)
		default:
			advance, token, err = consumeWord(data)
		}
		if advance > 0 {
			rdr.lastCol = rdr.col
			rdr.lastL = rdr.l
			for _, b := range data[:advance] {
				if b == '\n' || atEOF {
					rdr.l++
					rdr.col = 1
				}
				rdr.col++
			}
		}

		return
	}
	s.Split(split)
	return rdr
}

func main() {
	s := newTrackingReader(strings.NewReader("(foo 1 'bar')\n(baz 2 'quux')"))

	for s.Scan() {
		tok := s.Bytes()
		fmt.Printf("Found token %q at line #%d col #%d\n", string(tok), s.lastL, s.lastCol)
	}
}
