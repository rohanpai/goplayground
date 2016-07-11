package main

type T struct{ x int }

func (T) mT() {}

type S struct {
	T
}

func (S) mS() {}

type P *S

type I interface {
	mT()
}

func main() {
	var s S
	s.T.mT()
	s.mT() // == s.T.mT()

	var i I
	_ = i
	i = s.T
	i = s

	var ps = &s
	ps.mS()
	ps.T.mT()
	ps.mT() // == ps.T.mT()

	i = ps.T
	i = ps

	var p P = ps
	(*p).mS()
	// p.mS() // not ok - no automatic dereferencing for method selectors

	i = *p
	// i = p // not ok - like above, p.mS is not ok

	p.T.mT()
	p.mT() // == p.T.mT() (ok!) - this should not be permitted since line 48 is not permitted

	i = p.T
	// i = p // not ok! (but p.mT() is ok)
}
