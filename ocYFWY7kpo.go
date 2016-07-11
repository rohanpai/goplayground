package main

import &#34;fmt&#34;
import &#34;time&#34;

//365.2545
// 2456668.43767
//The reference time used in the layouts is:
//Mon Jan 2 15:04:05 MST 2006
//which is Unix time 1136239445.

//calculate the Julian  date, provided it&#39;s within 209 years of Jan 2, 2006.
func Julian(t time.Time) float64 {
	// Julian date, in seconds, of the &#34;Format&#34; standard time.
	// (See http://www.onlineconversion.com/julian_date.htm)
	const julian = 2453738.4195
	// Easiest way to get the time.Time of the Unix time.
	// (See comments for the UnixDate in package Time.)
	unix := time.Unix(1136239445, 0)
	const oneDay = float64(86400. * time.Second)
	return julian &#43; float64(t.Sub(unix))/oneDay
}

const (
	dateFmt = &#34;20060102&#34;
)
// round value - convert to int64
func Round(value float64) int64 {
	if value &lt; 0.0 {
		value -= 0.5
	} else {
		value &#43;= 0.5
	}
	return int64(value)
}
func main() {

	//tToday := time.Now()
	//tToday p := tToday 0.UTC().Format(dateFmt)
	//t0, _ := time.ParseInLocation(dateFmt, tToday p, time.UTC)
	
	tToday, _ := time.ParseInLocation(dateFmt, &#34;20140101&#34;, time.UTC)
	t1, _ := time.ParseInLocation(dateFmt, &#34;20120401&#34;, time.UTC)

	fmt.Printf(&#34;Today     %s\n&#34;, tToday.Format(dateFmt))
	fmt.Printf(&#34;Yesterday %s\n&#34;, t1.Format(dateFmt))

	// today
	t0j := Julian(tToday)
	fmt.Printf(&#34;Julian date:%f\n&#34;, t0j)

	//yesterday
	t1j := Julian(t1)
	fmt.Printf(&#34;Julian date:%f\n&#34;, t1j)

	t := t0j - t1j
	fmt.Printf(&#34;days since:%f\n&#34;, t) // 366.000000
	
	d := t/365.2545
	fmt.Printf(&#34;years precise:%f\n&#34;, d)

	y := Round(t/365.2545)		 // years since:1
	fmt.Printf(&#34;full years since:%d\n&#34;, y)
	d = t/365.2545
	fmt.Printf(&#34;&#43;days:%f\n&#34;, d)
	if d &gt; 0 {
	
	}
	
	
}