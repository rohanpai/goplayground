package main

import (
	&#34;fmt&#34;
	&#34;time&#34;
)

type Stopwatch struct {
	start, stop time.Time       // no need for lap, see mark
	mark        time.Duration   // mark is the duration from the start that the most recent lap was started
	laps        []time.Duration //
}

// New creates a new stopwatch with starting time offset by
// a user defined value. Negative offsets result in a countdown
// prior to the start of the stopwatch.
func New(offset time.Duration, active bool) *Stopwatch {
	var sw Stopwatch
	sw.Reset(offset, active)
	return &amp;sw
}

// Reset allows the re-use of a Stopwatch instead of creating
// a new one.
func (s *Stopwatch) Reset(offset time.Duration, active bool) {
	now := time.Now()
	s.start = now.Add(-offset)
	if active {
		s.stop = time.Time{}
	} else {
		s.stop = now
	}
	s.mark = 0
	s.laps = nil
}

// Active returns true if the stopwatch is active (counting up)
func (s *Stopwatch) Active() bool {
	return s.stop.IsZero()
}

// Stop makes the stopwatch stop counting up
func (s *Stopwatch) Stop() {
	if s.Active() {
		s.stop = time.Now()
	}
}

// Start intiates, or resumes the counting up process
func (s *Stopwatch) Start() {
	if !s.Active() {
		diff := time.Now().Sub(s.stop)
		s.start = s.start.Add(diff)
		s.stop = time.Time{}
	}
}

// Elapsed time is the time the stopwatch has been active
func (s *Stopwatch) ElapsedTime() time.Duration {
	if s.Active() {
		return time.Since(s.start)
	}
	return s.stop.Sub(s.start)
}

// LapTime is the time since the start of the lap
func (s *Stopwatch) LapTime() time.Duration {
	return s.ElapsedTime() - s.mark
}

// Lap starts a new lap, and returns the length of
// the previous one.
func (s *Stopwatch) Lap() (lap time.Duration) {
	lap = s.ElapsedTime() - s.mark
	s.mark = s.ElapsedTime()
	s.laps = append(s.laps, lap)
	return
}

// Laps returns a slice of completed lap times
func (s *Stopwatch) Laps() (laps []time.Duration) {
	laps = make([]time.Duration, len(s.laps))
	copy(laps, s.laps)
	return
}

func main() {
	s := New(-5*time.Second, false)
	t1 := time.Tick(time.Second)
	t2 := time.After(time.Second * 25 / 10)
	var t3 &lt;-chan time.Time

	for i := 0; i &lt; 20; i&#43;&#43; {
		// at first time instance after start, run lap ticker
		if s.ElapsedTime() &gt;= 0 &amp;&amp; t3 == nil {
			fmt.Println(&#34;Starting lap ticker @&#34;, s.ElapsedTime())
			t3 = time.Tick(time.Second * 175 / 100)
		}
		select {
		case &lt;-t1:
			fmt.Printf(&#34;Elapsed: %s, Lap: %s, Laps: %s\n&#34;, s.ElapsedTime(), s.LapTime(), s.Laps())
		case &lt;-t2:
			fmt.Println(&#34;Starting Stopwatch&#34;)
			s.Start()
		case &lt;-t3:
			fmt.Printf(&#34;Lap complete at %s, Lap was %s\n&#34;, s.ElapsedTime(), s.Lap())
		}
	}
}
