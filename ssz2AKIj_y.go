// There is a problem that comes up in C&#43;&#43; a lot (in my experience) which has no
// good C&#43;&#43; answer, but which Go has a great answer for (OMG, Tim said something
// nice about Go!).  The problem is: wrappers.
//
// I often find myself wanting to implement a particular interface (abstract base
// class in C&#43;&#43;) by delegating all method calls to another implementation of that
// interface, but intercepting just one or two of them.  I can&#39;t use inheritance
// because I don&#39;t know the concrete type of the wrappee.  So I implement my own
// wrapper type, and I implement every one of the interface&#39;s methods as simple
// calls to the wrappee.  This works, but it can be INCREDIBLY tedious to do,
// especially if the interface has a lot of methods.  It works right up until
// someone adds a new method to the base interface -- then my build breaks and I
// have to implement it too.
//
// Go&#39;s answer is very elegant - embed the interface type to &#34;inherit&#34; methods,
// &#34;override&#34; mthods you care about, and init from an instance of that interface.
//
// thockin@google.com

package main

import &#34;fmt&#34;
import &#34;math/rand&#34;

// Base interface I want to wrap.
type Frobber interface {
	FrobGently()
	FrobAggressively()
	FrobWithPrejudice()
}

// One implementation.
type wetFrobber struct{}

func (wet wetFrobber) FrobGently() {
	fmt.Println(&#34;wet.FrobGently()&#34;)
}

func (wet wetFrobber) FrobAggressively() {
	fmt.Println(&#34;wet.FrobAggressively()&#34;)
}

func (wet wetFrobber) FrobWithPrejudice() {
	fmt.Println(&#34;wet.FrobWithPrejudice()&#34;)
}

// Another implementation.
type dryFrobber struct{}

func (dry dryFrobber) FrobGently() {
	fmt.Println(&#34;dry.FrobGently()&#34;)
}

func (dry dryFrobber) FrobAggressively() {
	fmt.Println(&#34;dry.FrobAggressively()&#34;)
}

func (dry dryFrobber) FrobWithPrejudice() {
	fmt.Println(&#34;dry.FrobWithPrejudice()&#34;)
}

// Choose a random implementation.
func newFrobber() Frobber {
	// Yes, I know playground is not actually random.
	if rand.Intn(2) == 0 {
		return dryFrobber{}
	}
	return wetFrobber{}
}

// My wrapper.
type FrobWrapper struct {
	Frobber
}

// Override just one method.
func (wrap FrobWrapper) FrobWithPrejudice() {
	fmt.Printf(&#34;OMG: &#34;)
	wrap.Frobber.FrobWithPrejudice()
}

func main() {
	var f Frobber = FrobWrapper{newFrobber()}
	f.FrobGently()
	f.FrobAggressively()
	f.FrobWithPrejudice()
}
