package main

import (
	&#34;golang.org/x/net/context&#34;
	&#34;time&#34;
)

type Request struct {
	context context.Context
}

func Increment(c context.Context) context.Context {
	x := c.Value(&#34;x&#34;).(int)
	return context.WithValue(c, &#34;x&#34;, x&#43;1)
}

func outer(r *Request, next func(r *Request)) {
	r.context = context.WithValue(r.context, &#34;x&#34;, 0)
	next(r)

	// Here&#39;s a random time-consuming task.
	time.Sleep(100 * time.Millisecond)

	r.context = Increment(r.context)
	if r.context.Value(&#34;x&#34;).(int) != 2 {
		panic(&#34;oh shiii....&#34;)
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
	r := &amp;Request{context.TODO()}
	outer(r, inner) // Imitate middleware chain
}
