/*
// A Duration represents the elapsed time between two instants as
// an int64 nanosecond count. The representation limits the largest
// representable duration to approximately 290 years.

type Duration int64

// Common durations. There is no definition for units of Day or larger
// to avoid confusion across daylight savings time zone transitions.

const (
        Nanosecond  Duration = 1
        Microsecond          = 1000 * Nanosecond
        Millisecond          = 1000 * Microsecond
        Second               = 1000 * Millisecond
        Minute               = 60 * Second
        Hour                 = 60 * Minute
)

// Add returns the time t&#43;d.
func (t Time) Add(d Duration) Time
*/

// Sample program to show a idiomatic use of named types from the
// standard library and how they work in concert with other Go concepts.
package main

import (
	&#34;fmt&#34;
	&#34;time&#34;
)

// fiveSeconds is a typed constant of type int64.
const fiveSeconds int64 = 5 * time.Second // Duration(5) * Duration(1000000000)

// ./example2.go:34: cannot use time.Duration(5) * time.Second
// (type time.Duration) as type int64 in const initializer

// main is the entry point for the application.
func main() {
	// Use the time package to get the current date/time.
	now := time.Now()

	// Subtract 5 nanoseconds from now time using a literal constant.
	lessFiveNanoseconds := now.Add(-5)

	// Attempt to use the constant of type int64.
	lessFiveSeconds := now.Add(-fiveSeconds)

	// ./example2.go:48: cannot use -fiveSeconds (type int64) as
	// type time.Duration in argument to now.Add

	// Display the values.
	fmt.Printf(&#34;Now     : %v\n&#34;, now)
	fmt.Printf(&#34;Nano    : %v\n&#34;, lessFiveNanoseconds)
	fmt.Printf(&#34;Seconds : %v\n&#34;, lessFiveSeconds)
}