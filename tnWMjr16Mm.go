// Single line comment
/* Multi-
   line comment */

// A package clause starts every source file.
// Main is a special name declaring an executable rather than a library.
package main

// Import declaration declares library packages referenced in this file.
import (
	&#34;fmt&#34;
	// A package in the Go standard library.
	// Implements some I/O utility functions.
	m &#34;math&#34;
	// Math library with local alias m.
	// Yes, a web server!
	&#34;strconv&#34; // String conversions.
)

// A function definition.  Main is special.  It is the entry point for the
// executable program.  Love it or hate it, Go uses brace brackets.
func main() {
	// Println outputs a line to stdout.
	// Qualify it with the package name, fmt.
	fmt.Println(&#34;Hello world!&#34;)

	// Call another function within this package.
	beyondHello()
}

// Functions have parameters in parentheses.
// If there are no parameters, empty parentheses are still required.
func beyondHello() {
	var x int // Variable declaration. Variables must be declared before use.
	x = 3     // Variable assignment.
	// &#34;Short&#34; declarations use := to infer the type, declare, and assign.
	y := 4
	sum, prod := learnMultiple(x, y)        // Function returns two values.
	fmt.Println(&#34;sum:&#34;, sum, &#34;prod:&#34;, prod) // Simple output.
	learnTypes()                            // &lt; y minutes, learn more!
}

// Functions can have parameters and (multiple!) return values.
func learnMultiple(x, y int) (sum, prod int) {
	return x &#43; y, x * y // Return two values.
}

// Some built-in types and literals.
func learnTypes() {
	// Short declaration usually gives you what you want.
	s := &#34;Learn Go!&#34; // string type.

	s2 := `A &#34;raw&#34; string literal
can include line breaks.` // Same string type.

	// Non-ASCII literal.  Go source is UTF-8.
	g := &#39;Σ&#39; // rune type, an alias for uint32, holds a unicode code point.

	f := 3.14195 // float64, an IEEE-754 64-bit floating point number.
	c := 3 &#43; 4i  // complex128, represented internally with two float64&#39;s.

	// Var syntax with an initializers.
	var u uint = 7 // Unsigned, but implementation dependent size as with int.
	var pi float32 = 22. / 7

	// Conversion syntax with a short declaration.
	n := byte(&#39;\n&#39;) // byte is an alias for uint8.

	// Arrays have size fixed at compile time.
	var a4 [4]int           // An array of 4 ints, initialized to all 0.
	a3 := [...]int{3, 1, 5} // An array of 3 ints, initialized as shown.

	// Slices have dynamic size.  Arrays and slices each have advantages
	// but use cases for slices are much more common.
	s3 := []int{4, 5, 9}    // Compare to a3.  No ellipsis here.
	s4 := make([]int, 4)    // Allocates slice of 4 ints, initialized to all 0.
	var d2 [][]float64      // Declaration only, nothing allocated here.
	bs := []byte(&#34;a slice&#34;) // Type conversion syntax.

	p, q := learnMemory() // Declares p, q to be type pointer to int.
	fmt.Println(*p, *q)   // * follows a pointer.  This prints two ints.

	// Maps are a dynamically growable associative array type, like the
	// hash or dictionary types of some other languages.
	m := map[string]int{&#34;three&#34;: 3, &#34;four&#34;: 4}
	m[&#34;one&#34;] = 1

	// Unused variables are an error in Go.
	// The underbar lets you &#34;use&#34; a variable but discard its value.
	_, _, _, _, _, _, _, _, _ = s2, g, f, u, pi, n, a3, s4, bs
	// Output of course counts as using a variable.
	fmt.Println(s, c, a4, s3, d2, m)

	learnFlowControl() // Back in the flow.
}

// It is possible, unlike in many other languages for functions in go
// to have named return values.
// Assigning a name to the type being returned in the function declaration line
// allows us to easily return from multiple points in a function as well as to
// only use the return keyword, without anything further.
func learnNamedReturns(x, y int) (z int) {
	z = x * y
	return // z is implicit here, because we named it earlier.
}

// Go is fully garbage collected.  It has pointers but no pointer arithmetic.
// You can make a mistake with a nil pointer, but not by incrementing a pointer.
func learnMemory() (p, q *int) {
	// Named return values p and q have type pointer to int.
	p = new(int) // Built-in function new allocates memory.
	// The allocated int is initialized to 0, p is no longer nil.
	s := make([]int, 20) // Allocate 20 ints as a single block of memory.
	s[3] = 7             // Assign one of them.
	r := -2              // Declare another local variable.
	return &amp;s[3], &amp;r     // &amp; takes the address of an object.
}

func expensiveComputation() float64 {
	return m.Exp(10)
}

func learnFlowControl() {
	// If statements require brace brackets, and do not require parens.
	if true {
		fmt.Println(&#34;told ya&#34;)
	}
	// Formatting is standardized by the command line command &#34;go fmt.&#34;
	if false {
		// Pout.
	} else {
		// Gloat.
	}
	// Use switch in preference to chained if statements.
	x := 42.0
	switch x {
	case 0:
	case 1:
	case 42:
		// Cases don&#39;t &#34;fall through&#34;.
	case 43:
		// Unreached.
	}
	// Like if, for doesn&#39;t use parens either.
	// Variables declared in for and if are local to their scope.
	for x := 0; x &lt; 3; x&#43;&#43; { // &#43;&#43; is a statement.
		fmt.Println(&#34;iteration&#34;, x)
	}
	// x == 42 here.

	// For is the only loop statement in Go, but it has alternate forms.
	for { // Infinite loop.
		break    // Just kidding.
		continue // Unreached.
	}
	// As with for, := in an if statement means to declare and assign
	// y first, then test y &gt; x.
	if y := expensiveComputation(); y &gt; x {
		x = y
	}
	// Function literals are closures.
	xBig := func() bool {
		return x &gt; 10000 // References x declared above switch statement.
	}
	fmt.Println(&#34;xBig:&#34;, xBig()) // true (we last assigned e^10 to x).
	x = 1.3e3                    // This makes x == 1300
	fmt.Println(&#34;xBig:&#34;, xBig()) // false now.

	// What&#39;s more is function literals may be defined and called inline,
	// acting as an argument to function, as long as:
	// a) function literal is called immediately (),
	// b) result type matches expected type of argument.
	fmt.Println(&#34;Add &#43; double two numbers: &#34;,
		func(a, b int) int {
			return (a &#43; b) * 2
		}(10, 2)) // Called with args 10 and 2
	// =&gt; Add &#43; double two numbers:  24

	// When you need it, you&#39;ll love it.
	goto love
love:

	learnFunctionFactory() // func returning func is fun(3)(3)
	learnDefer()      // A quick detour to an important keyword.
	learnInterfaces() // Good stuff coming up!
}

func learnFunctionFactory() {
	// Next two are equivalent, with second being more practical
	fmt.Println(sentenceFactory(&#34;summer&#34;)(&#34;A beautiful&#34;, &#34;day!&#34;))

	d := sentenceFactory(&#34;summer&#34;)
	fmt.Println(d(&#34;A beautiful&#34;, &#34;day!&#34;))
	fmt.Println(d(&#34;A lazy&#34;, &#34;afternoon!&#34;))
}

// Decorators are common in other languages. Same can be done in Go
// with function literals that accept arguments.
func sentenceFactory(mystring string) func(before, after string) string {
	return func(before, after string) string {
		return fmt.Sprintf(&#34;%s %s %s&#34;, before, mystring, after) // new string
	}
}

func learnDefer() (ok bool) {
	// Deferred statements are executed just before the function returns.
	defer fmt.Println(&#34;deferred statements execute in reverse (LIFO) order.&#34;)
	defer fmt.Println(&#34;\nThis line is being printed first because&#34;)
	// Defer is commonly used to close a file, so the function closing the
	// file stays close to the function opening the file.
	return true
}

// Define Stringer as an interface type with one method, String.
type Stringer interface {
	String() string
}

// Define pair as a struct with two fields, ints named x and y.
type pair struct {
	x, y int
}

// Define a method on type pair.  Pair now implements Stringer.
func (p pair) String() string { // p is called the &#34;receiver&#34;
	// Sprintf is another public function in package fmt.
	// Dot syntax references fields of p.
	return fmt.Sprintf(&#34;(%d, %d)&#34;, p.x, p.y)
}

func learnInterfaces() {
	// Brace syntax is a &#34;struct literal.&#34;  It evaluates to an initialized
	// struct.  The := syntax declares and initializes p to this struct.
	p := pair{3, 4}
	fmt.Println(p.String()) // Call String method of p, of type pair.
	var i Stringer          // Declare i of interface type Stringer.
	i = p                   // Valid because pair implements Stringer
	// Call String method of i, of type Stringer.  Output same as above.
	fmt.Println(i.String())

	// Functions in the fmt package call the String method to ask an object
	// for a printable representation of itself.
	fmt.Println(p) // Output same as above. Println calls String method.
	fmt.Println(i) // Output same as above.

	learnVariadicParams(&#34;great&#34;, &#34;learning&#34;, &#34;here!&#34;)
}

// Functions can have variadic parameters.
func learnVariadicParams(myStrings ...interface{}) {
	// Iterate each value of the variadic.
	// The underbar here is ignoring the index argument of the array.
	for _, param := range myStrings {
		fmt.Println(&#34;param:&#34;, param)
	}

	// Pass variadic value as a variadic parameter.
	fmt.Println(&#34;params:&#34;, fmt.Sprintln(myStrings...))

	learnErrorHandling()
}

func learnErrorHandling() {
	// &#34;, ok&#34; idiom used to tell if something worked or not.
	m := map[int]string{3: &#34;three&#34;, 4: &#34;four&#34;}
	if x, ok := m[1]; !ok { // ok will be false because 1 is not in the map.
		fmt.Println(&#34;no one there&#34;)
	} else {
		fmt.Print(x) // x would be the value, if it were in the map.
	}
	// An error value communicates not just &#34;ok&#34; but more about the problem.
	if _, err := strconv.Atoi(&#34;non-int&#34;); err != nil { // _ discards value
		// prints &#39;strconv.ParseInt: parsing &#34;non-int&#34;: invalid syntax&#39;
		fmt.Println(err)
	}
	// We&#39;ll revisit interfaces a little later.  Meanwhile,
	learnConcurrency()
}

// c is a channel, a concurrency-safe communication object.
func inc(i int, c chan int) {
	c &lt;- i &#43; 1 // &lt;- is the &#34;send&#34; operator when a channel appears on the left.
}

// We&#39;ll use inc to increment some numbers concurrently.
func learnConcurrency() {
	// Same make function used earlier to make a slice.  Make allocates and
	// initializes slices, maps, and channels.
	c := make(chan int)
	// Start three concurrent goroutines.  Numbers will be incremented
	// concurrently, perhaps in parallel if the machine is capable and
	// properly configured.  All three send to the same channel.
	go inc(0, c) // go is a statement that starts a new goroutine.
	go inc(10, c)
	go inc(-805, c)
	// Read three results from the channel and print them out.
	// There is no telling in what order the results will arrive!
	fmt.Println(&lt;-c, &lt;-c, &lt;-c) // channel on right, &lt;- is &#34;receive&#34; operator.

	cs := make(chan string)       // Another channel, this one handles strings.
	ccs := make(chan chan string) // A channel of string channels.
	go func() { c &lt;- 84 }()       // Start a new goroutine just to send a value.
	go func() { cs &lt;- &#34;wordy&#34; }() // Again, for cs this time.
	// Select has syntax like a switch statement but each case involves
	// a channel operation.  It selects a case at random out of the cases
	// that are ready to communicate.
	select {
	case i := &lt;-c: // The value received can be assigned to a variable,
		fmt.Printf(&#34;it&#39;s a %T&#34;, i)
	case &lt;-cs: // or the value received can be discarded.
		fmt.Println(&#34;it&#39;s a string&#34;)
	case &lt;-ccs: // Empty channel, not ready for communication.
		fmt.Println(&#34;didn&#39;t happen.&#34;)
	}
	// At this point a value was taken from either c or cs.  One of the two
	// goroutines started above has completed, the other will remain blocked.

}
