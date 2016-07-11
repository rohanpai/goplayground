package main

import (
	"golang.org/x/net/context"
	"time"
)

type Request struct {
	context context.Context
}

func Increment(c context.Context) context.Context {
	x := c.Value("x").(int)
	return context.WithValue(c, "x", x+1)
}

func outer(r *Request, next func(r *Request)) {
	r.context = context.WithValue(r.context, "x", 0)
	next(r)

	// Here's a random time-consuming task.
	time.Sleep(100 * time.Millisecond)

	r.context = Increment(r.context)
	if r.context.Value("x").(int) != 2 {
		panic("oh shiii....")
	}
}

func inner(r *Request) {
	// Here we decide to run some time-consuming operation in background, that
	// eventually saves its result into context
	go func() {
		// This time interval may or may not be equal to outer wait interval
		time.Sleep(100 * time.Millisecond)
		r.context = Increment(r.context)
	}()
}

func main() {
	r := &Request{context.TODO()}
	outer(r, inner) // Imitate middleware chain
}
