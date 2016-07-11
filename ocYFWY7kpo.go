package main

import "fmt"
import "time"

//365.2545
// 2456668.43767
//The reference time used in the layouts is:
//Mon Jan 2 15:04:05 MST 2006
//which is Unix time 1136239445.

//calculate the Julian  date, provided it's within 209 years of Jan 2, 2006.
func Julian(t time.Time) float64 {
	// Julian date, in seconds, of the "Format" standard time.
	// (See http://www.onlineconversion.com/julian_date.htm)
	const julian = 2453738.4195
	// Easiest way to get the time.Time of the Unix time.
	// (See comments for the UnixDate in package Time.)
	unix := time.Unix(1136239445, 0)
	const oneDay = float64(86400. * time.Second)
	return julian + float64(t.Sub(unix))/oneDay
}

const (
	dateFmt = "20060102"
)
// round value - convert to int64
func Round(value float64) int64 {
	if value < 0.0 {
		value -= 0.5
	} else {
		value += 0.5
	}
	return int64(value)
}
func main() {

	//tToday := time.Now()
	//tToday p := tToday 0.UTC().Format(dateFmt)
	//t0, _ := time.ParseInLocation(dateFmt, tToday p, time.UTC)
	
	tToday, _ := time.ParseInLocation(dateFmt, "20140101", time.UTC)
	t1, _ := time.ParseInLocation(dateFmt, "20120401", time.UTC)

	fmt.Printf("Today     %s\n", tToday.Format(dateFmt))
	fmt.Printf("Yesterday %s\n", t1.Format(dateFmt))

	// today
	t0j := Julian(tToday)
	fmt.Printf("Julian date:%f\n", t0j)

	//yesterday
	t1j := Julian(t1)
	fmt.Printf("Julian date:%f\n", t1j)

	t := t0j - t1j
	fmt.Printf("days since:%f\n", t) // 366.000000
	
	d := t/365.2545
	fmt.Printf("years precise:%f\n", d)

	y := Round(t/365.2545)		 // years since:1
	fmt.Printf("full years since:%d\n", y)
	d = t/365.2545
	fmt.Printf("+days:%f\n", d)
	if d > 0 {
	
	}
	
	
}