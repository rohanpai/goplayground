// Go supports time formatting and parsing via
// pattern-based layouts.

package main

import &#34;fmt&#34;
import &#34;time&#34;

func main() {
	p := fmt.Println

	// Here&#39;s a basic example of formatting a time
	// according to RFC3339, using the corresponding layout
	// constant.
	t := time.Now()
	p(t.Format(time.RFC3339))

	// Time parsing uses the same layout values as `Format`.
	t1, e := time.Parse(
		time.RFC3339,
		&#34;2012-11-01T22:08:41&#43;00:00&#34;)
	p(t1)

	// `Format` and `Parse` uses example-based layouts. Usually
	// you&#39;ll use a constant from `time` for these layouts, but
	// you can also supply custom layouts. Layouts must use the
	// reference time `Mon Jan 2 15:04:05 MST 2006` to show the
	// pattern with which to format/parse a given time/string.
	// The example time must be exactly as shown: the year 2006,
	// 15 for the hour, Monday for the day of the week, etc.
	p(t.Format(&#34;3:04PM&#34;))
	p(t.Format(&#34;Mon Jan _2 15:04:05 2006&#34;))
	p(t.Format(&#34;2006-01-02T15:04:05.999999-07:00&#34;))
	form := &#34;3 04 PM&#34;
	t2, e := time.Parse(form, &#34;8 41 PM&#34;)
	p(t2)

	// For purely numeric representations you can also
	// use standard string formatting with the extracted
	// components of the time value.
	fmt.Printf(&#34;%d-%02d-%02dT%02d:%02d:%02d-00:00\n&#34;,
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	// `Parse` will return an error on malformed input
	// explaining the parsing problem.
	ansic := &#34;Mon Jan _2 15:04:05 2006&#34;
	_, e = time.Parse(ansic, &#34;8:41PM&#34;)
	p(e)
}
