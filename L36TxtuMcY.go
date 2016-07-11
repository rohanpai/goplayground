package main

import (
	&#34;errors&#34;
	&#34;fmt&#34;
)

type Duration int64

const (
	Nanosecond  Duration = 1
	Microsecond          = 1000 * Nanosecond
	Millisecond          = 1000 * Microsecond
	Second               = 1000 * Millisecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
)

var errLeadingInt = errors.New(&#34;time: bad [0-9]*&#34;) // never printed

// leadingInt consumes the leading [0-9]* from s.
func leadingInt(s string) (x int64, rem string, err error) {
	i := 0
	for ; i &lt; len(s); i&#43;&#43; {
		c := s[i]
		if c &lt; &#39;0&#39; || c &gt; &#39;9&#39; {
			break
		}
		if x &gt; (1&lt;&lt;63-1)/10 {
			// overflow
			return 0, &#34;&#34;, errLeadingInt
		}
		x = x*10 &#43; int64(c) - &#39;0&#39;
		if x &lt; 0 {
			// overflow
			return 0, &#34;&#34;, errLeadingInt
		}
	}
	return x, s[i:], nil
}

var unitMap = map[string]int64{
	&#34;&#34;:   int64(1), // Handle case of unitless &#34;0.0&#34;
	&#34;ns&#34;: int64(Nanosecond),
	&#34;us&#34;: int64(Microsecond),
	&#34;µs&#34;: int64(Microsecond), // U&#43;00B5 = micro symbol
	&#34;μs&#34;: int64(Microsecond), // U&#43;03BC = Greek letter mu
	&#34;ms&#34;: int64(Millisecond),
	&#34;s&#34;:  int64(Second),
	&#34;m&#34;:  int64(Minute),
	&#34;h&#34;:  int64(Hour),
}

// ParseDuration parses a duration string.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as &#34;300ms&#34;, &#34;-1.5h&#34; or &#34;2h45m&#34;.
// Valid time units are &#34;ns&#34;, &#34;us&#34; (or &#34;µs&#34;), &#34;ms&#34;, &#34;s&#34;, &#34;m&#34;, &#34;h&#34;.
func ParseDuration(s string) (Duration, error) {
	// [-&#43;]?([0-9]*(\.[0-9]*)?[a-z]&#43;)&#43;
	orig := s
	var d int64
	neg := false

	// Consume [-&#43;]?
	if s != &#34;&#34; {
		c := s[0]
		if c == &#39;-&#39; || c == &#39;&#43;&#39; {
			neg = c == &#39;-&#39;
			s = s[1:]
		}
	}
	// Special case: if all that is left is &#34;0&#34;, this is zero.
	if s == &#34;0&#34; {
		return 0, nil
	}
	if s == &#34;&#34; {
		return 0, errors.New(&#34;time: invalid duration &#34; &#43; orig)
	}
	for s != &#34;&#34; {
		var (
			v, f  int64       // integers before, after decimal point
			scale float64 = 1 // value = v &#43; f/scale
		)

		var err error

		// The next character must be [0-9.]
		if !(s[0] == &#39;.&#39; || &#39;0&#39; &lt;= s[0] &amp;&amp; s[0] &lt;= &#39;9&#39;) {
			return 0, errors.New(&#34;time: invalid duration &#34; &#43; orig)
		}
		// Consume [0-9]*
		pl := len(s)
		v, s, err = leadingInt(s)
		if err != nil {
			return 0, errors.New(&#34;time: invalid duration &#34; &#43; orig)
		}
		pre := pl != len(s) // whether we consumed anything before a period

		// Consume (\.[0-9]*)?
		post := false
		if s != &#34;&#34; &amp;&amp; s[0] == &#39;.&#39; {
			s = s[1:]
			pl := len(s)
			f, s, err = leadingInt(s)
			if err != nil {
				return 0, errors.New(&#34;time: invalid duration &#34; &#43; orig)
			}
			for n := pl - len(s); n &gt; 0; n-- {
				scale *= 10
			}
			post = pl != len(s)
		}
		if !pre &amp;&amp; !post {
			// no digits (e.g. &#34;.s&#34; or &#34;-.s&#34;)
			return 0, errors.New(&#34;time: invalid duration &#34; &#43; orig)
		}

		// Consume unit.
		i := 0
		for ; i &lt; len(s); i&#43;&#43; {
			c := s[i]
			if c == &#39;.&#39; || &#39;0&#39; &lt;= c &amp;&amp; c &lt;= &#39;9&#39; {
				break
			}
		}
		if i == 0 &amp;&amp; (v != 0 || f != 0) {
			return 0, errors.New(&#34;time: missing unit in duration &#34; &#43; orig)
		}
		u := s[:i]
		s = s[i:]
		unit, ok := unitMap[u]
		if !ok {
			return 0, errors.New(&#34;time: unknown unit &#34; &#43; u &#43; &#34; in duration &#34; &#43; orig)
		}
		if v &gt; (1&lt;&lt;63-1)/unit {
			// overflow
			return 0, errors.New(&#34;time: invalid duration &#34; &#43; orig)
		}
		v *= unit
		if f &gt; 0 {
			// float64 is needed to be nanosecond accurate for fractions of hours.
			// v &gt;= 0 &amp;&amp; (f*unit/scale) &lt;= 3.6e&#43;12 (ns/h, h is the largest unit)
			v &#43;= int64(float64(f) * (float64(unit) / scale))
			if v &lt; 0 {
				// overflow
				return 0, errors.New(&#34;time: invalid duration &#34; &#43; orig)
			}
		}
		d &#43;= v
		if d &lt; 0 {
			// overflow
			return 0, errors.New(&#34;time: invalid duration &#34; &#43; orig)
		}
	}

	if neg {
		d = -d
	}
	return Duration(d), nil
}

func main() {
	fmt.Println(ParseDuration(&#34;s&#34;))
	fmt.Println(ParseDuration(&#34;1&#34;))
	fmt.Println(ParseDuration(&#34;0&#34;))
	fmt.Println(ParseDuration(&#34;00.0&#34;))
	fmt.Println(ParseDuration(&#34;0.0ms&#34;))
}
