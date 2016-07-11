package main

import (
	&#34;fmt&#34;
	&#34;math&#34;
	&#34;strconv&#34;
)

func main() {
	n := 12345.6789
	formats := []string{
		&#34;&#34;,
		&#34;#&#34;,
		&#34;#.&#34;,
		&#34;#,###.&#34;,
		&#34;#,###&#34;,
		&#34;#,###.##&#34;,
		&#34;#.###,######&#34;,
		&#34;#\u202F###,##&#34;,
	}
	for _, format := range formats {
		fmt.Printf(&#34;%16q =&gt; %s\n&#34;, format, RenderFloat(format, n))
	}
}

var renderFloatPrecisionMultipliers = [10]float64{
	1,
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	1000000000,
}

var renderFloatPrecisionRounders = [10]float64{
	0.5,
	0.05,
	0.005,
	0.0005,
	0.00005,
	0.000005,
	0.0000005,
	0.00000005,
	0.000000005,
	0.0000000005,
}

func RenderFloat(format string, n float64) string {
	// Special cases:
	// NaN = &#34;NaN&#34;
	// &#43;Inf = &#34;&#43;Infinity&#34;
	// -Inf = &#34;-Infinity&#34;
	if math.IsNaN(n) {
		return &#34;NaN&#34;
	}
	if n &gt; math.MaxFloat64 {
		return &#34;Infinity&#34;
	}
	if n &lt; -math.MaxFloat64 {
		return &#34;-Infinity&#34;
	}

	// default format
	precision := 2
	decimalStr := &#34;.&#34;
	thousandStr := &#34;,&#34;
	positiveStr := &#34;&#34;
	negativeStr := &#34;-&#34;

	if len(format) &gt; 0 {
		// If there is an explicit format directive,
		// then default values are these:
		precision = 9
		thousandStr = &#34;&#34;

		// collect indices of meaningful formatting directives
		formatDirectiveChars := []rune(format)
		formatDirectiveIndices := make([]int, 0)
		for i, char := range formatDirectiveChars {
			if char != &#39;#&#39; &amp;&amp; char != &#39;0&#39; {
				formatDirectiveIndices = append(formatDirectiveIndices, i)
			}
		}

		if len(formatDirectiveIndices) &gt; 0 {
			// Directive at index 0:
			// Must be a &#39;&#43;&#39;
			// Raise an error if not the case
			// index: 0123456789
			// &#43;0.000,000
			// &#43;000,000.0
			// &#43;0000.00
			// &#43;0000
			if formatDirectiveIndices[0] == 0 {
				if formatDirectiveChars[formatDirectiveIndices[0]] != &#39;&#43;&#39; {
					panic(&#34;RenderFloat(): invalid positive sign directive&#34;)
				}
				positiveStr = &#34;&#43;&#34;
				formatDirectiveIndices = formatDirectiveIndices[1:]
			}

			// Two directives:
			// First is thousands separator
			// Raise an error if not followed by 3-digit
			// 0123456789
			// 0.000,000
			// 000,000.00
			if len(formatDirectiveIndices) == 2 {
				if (formatDirectiveIndices[1] - formatDirectiveIndices[0]) != 4 {
					panic(&#34;RenderFloat(): thousands separator directive must be followed by 3 digit-specifiers&#34;)
				}
				thousandStr = string(formatDirectiveChars[formatDirectiveIndices[0]])
				formatDirectiveIndices = formatDirectiveIndices[1:]
			}

			// One directive:
			// Directive is decimal separator
			// The number of digit-specifier following the separator indicates wanted precision
			// 0123456789
			// 0.00
			// 000,0000
			if len(formatDirectiveIndices) == 1 {
				decimalStr = string(formatDirectiveChars[formatDirectiveIndices[0]])
				precision = len(formatDirectiveChars) - formatDirectiveIndices[0] - 1
			}
		}
	}

	// generate sign part
	var signStr string
	if n &gt;= 0.000000001 {
		signStr = positiveStr
	} else if n &lt;= -0.000000001 {
		signStr = negativeStr
		n = -n
	} else {
		signStr = &#34;&#34;
		n = 0.0
	}

	// split number into integer and fractional parts
	intf, fracf := math.Modf(n &#43; renderFloatPrecisionRounders[precision])

	// generate integer part string
	intStr := strconv.Itoa(int(intf))

	// add thousand separator if required
	if len(thousandStr) &gt; 0 {
		for i := len(intStr); i &gt; 3; {
			i -= 3
			intStr = intStr[:i] &#43; thousandStr &#43; intStr[i:]
		}
	}

	// no fractional part, we can leave now
	if precision == 0 {
		return signStr &#43; intStr
	}

	// generate fractional part
	fracStr := strconv.Itoa(int(fracf * renderFloatPrecisionMultipliers[precision]))
	// may need padding
	if len(fracStr) &lt; precision {
		fracStr = &#34;000000000000000&#34;[:precision-len(fracStr)] &#43; fracStr
	}

	return signStr &#43; intStr &#43; decimalStr &#43; fracStr
}

func RenderInteger(format string, n int) string {
	return RenderFloat(format, float64(n))
}
