package main

import (
	&#34;fmt&#34;
	&#34;github.com/soniakeys/meeus/base&#34;
	&#34;github.com/soniakeys/meeus/conjunction&#34;
	&#34;github.com/soniakeys/meeus/deltat&#34;
	&#34;github.com/soniakeys/meeus/julian&#34;
	&#34;math&#34;
	&#34;time&#34;
)

func main() {
	// Example 18.a, p. 117.

	// Day of month is sufficient for a time scale.
	day1 := 5.
	day5 := 9.

	// Text asks for Mercury-Venus conjunction, so r1, d1 is Venus ephemeris,
	// r2, d2 is Mercury ephemeris.

	// Venus
	r1 := []float64{
		base.NewRA(10, 27, 27.175).Rad(),
		base.NewRA(10, 26, 32.410).Rad(),
		base.NewRA(10, 25, 29.042).Rad(),
		base.NewRA(10, 24, 17.191).Rad(),
		base.NewRA(10, 22, 57.024).Rad(),
	}
	d1 := []float64{
		base.NewAngle(false, 4, 04, 41.83).Rad(),
		base.NewAngle(false, 3, 55, 54.66).Rad(),
		base.NewAngle(false, 3, 48, 03.51).Rad(),
		base.NewAngle(false, 3, 41, 10.25).Rad(),
		base.NewAngle(false, 3, 35, 16.61).Rad(),
	}
	// Mercury
	r2 := []float64{
		base.NewRA(10, 24, 30.125).Rad(),
		base.NewRA(10, 25, 00.342).Rad(),
		base.NewRA(10, 25, 12.515).Rad(),
		base.NewRA(10, 25, 06.235).Rad(),
		base.NewRA(10, 24, 41.185).Rad(),
	}
	d2 := []float64{
		base.NewAngle(false, 6, 26, 32.05).Rad(),
		base.NewAngle(false, 6, 10, 57.72).Rad(),
		base.NewAngle(false, 5, 57, 33.08).Rad(),
		base.NewAngle(false, 5, 46, 27.07).Rad(),
		base.NewAngle(false, 5, 37, 48.45).Rad(),
	}
	// compute conjunction
	day, dd, err := conjunction.Planetary(day1, day5, r1, d1, r2, d2)
	if err != nil {
		fmt.Println(err)
		return
	}
	// time of conjunction
	fmt.Printf(&#34;1991 August %.5f\n&#34;, day)

	// more useful clock format
	dInt, dFrac := math.Modf(day)
	fmt.Printf(&#34;1991 August %d at %s TD\n&#34;, int(dInt),
		base.NewFmtTime(dFrac*24*3600))

	// deltat func needs jd
	jd := julian.CalendarGregorianToJD(1991, 8, day)
	// compute UT = TD - ΔT, and separate back into calendar components.
	// (we could use our known calendar components, but this illustrates
	// the more general technique that would allow for rollovers.)
	y, m, d := julian.JDToCalendar(jd - deltat.Interp10A(jd)/(3600*24))
	// format as before
	dInt, dFrac = math.Modf(d)
	fmt.Printf(&#34;%d %s %d at %s UT\n&#34;, y, time.Month(m), int(dInt),
		base.NewFmtTime(dFrac*24*3600))

	// Δδ
	fmt.Printf(&#34;Δδ = %s\n&#34;, base.NewFmtAngle(dd))

}
