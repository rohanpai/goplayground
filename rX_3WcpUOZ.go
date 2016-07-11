package main

import (
	&#34;fmt&#34;
	&#34;math/rand&#34;
	&#34;time&#34;
)

const (
	width    = 76
	height   = 10
	numBalls = 10
)

type screen [height][width]byte

func (s *screen) print() {
	fmt.Print(&#34;\x0c&#34;)
	for _, row := range *s {
		for i, b := range row {
			if b == 0 {
				row[i] = &#39; &#39;
			}
		}
		fmt.Println(string(row[:]))
	}
}

func (s *screen) paint(x, y int, b byte) {
	s[y][x] = b
}

type ball struct {
	x, y   int
	xd, yd int
}

func newBall() *ball {
	return &amp;ball{
		x:  rand.Intn(width),
		y:  rand.Intn(height),
		xd: (rand.Intn(2) * 2) - 1,
		yd: (rand.Intn(2) * 2) - 1,
	}
}

func (b *ball) move() {
	moveBounce(&amp;b.x, &amp;b.xd, width)
	moveBounce(&amp;b.y, &amp;b.yd, height)
}

func moveBounce(v, d *int, max int) {
	*v &#43;= *d
	if *v &lt; 0 {
		*v = -*v
		*d = -*d
	} else if *v &gt;= max {
		*v = max - 2 - (*v - max)
		*d = -*d
	}
}

func main() {
	balls := []*ball{}
	for i := 0; i &lt; numBalls; i&#43;&#43; {
		balls = append(balls, newBall())
	}
	for i := 200; i &gt;= 0; i-- {
		var s screen
		for _, b := range balls {
			b.move()
			s.paint(b.x, b.y, &#39;*&#39;)
		}
		s.print()
		time.Sleep(time.Second / 20)
	}
}
