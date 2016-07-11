package main

import (
	&#34;fmt&#34;
	&#34;github.com/soniakeys/meeus/base&#34;
	&#34;github.com/soniakeys/meeus/julian&#34;
	&#34;github.com/soniakeys/meeus/line&#34;
	&#34;math&#34;
	&#34;time&#34;
)

func main() {
	// Example 19.a, p. 121.

	// convert degree data to radians
	r1 := 113.56833 * math.Pi / 180
	d1 := 31.89756 * math.Pi / 180
	r2 := 116.25042 * math.Pi / 180
	d2 := 28.03681 * math.Pi / 180
	r3 := make([]float64, 5)
	for i, ri := range []float64{
		118.98067, 119.59396, 120.20413, 120.81108, 121.41475} {
		r3[i] = ri * math.Pi / 180
	}
	d3 := make([]float64, 5)
	for i, di := range []float64{
		21.68417, 21.58983, 21.49394, 21.39653, 21.29761} {
		d3[i] = di * math.Pi / 180
	}
	// use JD as time to handle month boundary
	jd, err := line.Time(r1, d1, r2, d2, r3, d3,
		julian.CalendarGregorianToJD(1994, 9, 29),
		julian.CalendarGregorianToJD(1994, 10, 3))
	if err != nil {
		fmt.Println(err)
		return
	}
	y, m, d := julian.JDToCalendar(jd)
	dInt, dFrac := math.Modf(d)
	fmt.Printf(&#34;%d %s %.4f\n&#34;, y, time.Month(m), d)
	fmt.Printf(&#34;%d %s %d, at %.64s TD(UT)\n&#34;, y, time.Month(m), int(dInt),
		base.NewFmtTime(dFrac*24*3600))
}
