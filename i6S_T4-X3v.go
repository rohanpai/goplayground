package main

import (
	&#34;fmt&#34;
	&#34;strconv&#34;
	&#34;time&#34;
)

const (
	secondsPerMinute       = 60
	secondsPerHour         = 60 * 60
	secondsPerDay          = 24 * secondsPerHour
	unixToInternal   int64 = (1969*365 &#43; 1969/4 - 1969/100 &#43; 1969/400) * secondsPerDay
)

var (
	MinTime = time.Unix(0, 0)
	MaxTime = time.Unix(1&lt;&lt;63-1 - unixToInternal, 999999999)
)

// For more information, check out:
// http://golang.org/src/time/time.go?s=1855:2369#L29 - Time.sec is only int64
// http://golang.org/src/time/time.go?s=28472:28692#L959 - Time.sec = seconds &#43; unixToInternal
// http://golang.org/src/time/time.go?s=8279:8361#L235 - unixToInternal
//
// http://golang.org/src/time/time.go?s=2582:2615#L53 - Comparing times using Time.sec
// https://golang.org/ref/spec#Integer_overflow - Overflow behavior in Go Spec

func main() {
	var result int64
	var max int64 = 1&lt;&lt;63 - 1

	fmt.Printf(&#34;MinTime: %s. MaxTime: %s\n&#34;, MinTime, MaxTime)
	fmt.Printf(&#34;Is %s before %s? %t\n&#34;, MinTime, MaxTime, MinTime.Before(MaxTime))

	runtimeInt64, err := strconv.ParseInt(strconv.FormatInt(unixToInternal, 10), 10, 64)
	if err != nil {
		fmt.Errorf(&#34;ParseInt Failed: %s&#34;, err)
	}
	result = runtimeInt64 &#43; max
	fmt.Printf(&#34;This would fail at compile time. Instead, at runtime, we get: 1&lt;&lt;63-1 (%v) &#43; unixToInternal (%v) = %v\n&#34;, max, unixToInternal, result)
}