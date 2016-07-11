package main

import (
	&#34;fmt&#34;
	&#34;math&#34;
)

type Point struct {
	X, Y float64
}

type BinaryPointFn func(Point, Point) float64

func Dist(p1, p2 Point) float64 {
	fmt.Println(&#34;calculating new value...&#34;)
	return math.Sqrt(math.Pow(p2.X-p1.Y, 2.0) &#43; math.Pow(p2.Y-p1.Y, 2.0))
}

func Memoize(fn BinaryPointFn) BinaryPointFn {
	type T struct {
		p1	Point
		p2	Point
	}
	history := make(map[T]float64)
	return func(p1, p2 Point) float64 {
		if res, ok := history[T{p1, p2}]; ok {
			fmt.Println(&#34;reading from history...&#34;)
			return res
		}
		res := fn(p1, p2)
		history[T{p1, p2}] = res
		return res
	}
}

func main() {
	p1 := Point{3.1, 10.3}
	p2 := Point{-4.7, 28.9}
	fn := Memoize(Dist)
	fmt.Println(fn(p1, p2))
	fmt.Println(fn(p1, p2))
	fmt.Println(fn(p1, p2))
	fmt.Println(fn(p1, p2))
	fmt.Println(fn(p1, p2))
}
