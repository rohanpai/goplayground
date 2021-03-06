// An example package using Go-as-if-it-had-parametric-polymorphism,
// with an "iter" package in the standard library, following an
// already idiomatic Go iteration pattern.
//
// The feature is entirely imaginary, but I've tried to write
// the code to fit as much within Go's existing idioms as possible.
//
// Generics questions that this code does not attempt to resolve:
//	- What semantics are there for equality of values
//	with parametric types?
//	- what happens if you convert a value with parametric type
//	to an interface?
//
// Points given for pointing out syntax errors, logical inconsistencies
// or language-feature implementation roadtraps.
package main
import (
	"os"
	"log"
	"iter"
)

func main() {
	f, err := os.Open("/etc/passwd")
	if err != nil {
		log.Fatal(err)
	}
	r := strings.NewReader("light blue\nfaded khaki\n")
	iter0 := iter.BufioScanner(bufio.NewScanner(f))
	iter1 := iter.Slice<string>{"one", "two"}
	both := iter.Sequence(iter0, iter1)
	prefixed := iter.Map(both, func(s string) string {
		return "x: " + s
	})
	foo, err := iter.Gather(prefixed)
	if err != nil {
		log.Fatal(err)
	}
	// foo is []string{"x: light blue", "x: faded khaki", "x: one", "x: two"}
	fmt.Printf("%v\n", foo)
}

///////////////////////////////////////////////////////////

// The iter package implements a general iterator interface type, Iter,
// and some functions that operate on values of that type.
package iter

// Iter represents an iterable collection of values.
type Iter<T> interface {
	// Next advances the iterator to the next value, which will then
	// be available through the Value method. It returns
	// false when the iteration stops, either by reaching the end
	// or an error. After Scan returns false, the Err method
	// will return any error that occurred during iteration..
	Next() bool

	// Value returns the most recent value generated by a call to Next.
	// It may be called any number of times between calls to Next.
	// If called after Next has returned false, it returns the zero value.
	Value() T

	// Err returns the first error encountered.
	Err() error

	// Close closes the iterator and frees any associated resources.
	Close() error
}

// Gather iterates through all the items in iter
// and returns them as a slice.
func Gather<T>(iter Iter<T>) ([]T, error) {
	var slice []T
	while iter.Next() {
		slice = append(slice, iter.Value()
	}
	return slice, iter.Err()
}

// Identity returns the identity function for
// a given type.
func Identity<T>() func(T) T {
	return func(t T) T {
		return t
	}
}

// Map returns an iterator that produces 
// a value f(x) for every value in the given iterator.
// Any non-nil error returned from the underlying iterator
// will be transformed by the given err function.
func Map<S, T>(
	iter Iter<S>,
	transformError func(error) error,
	f func(S) T,
) Iter<T> {
	if transformError == nil {
		transformError = Identity<error>()
	}
	return &mapping{
		Iter: iter,
		f: f,
	}
}

// mapping implements the iterator returned by Map.
// Note the embedding of a type name with
// a parametric type parameter.
// When checking for interface type compatibility,
// methods must be compatible. So if S==T
// then mapping will automatically
// implement Iter<S> but not Iter<T>.
// It will implement interface {
//	Next() bool
//	Close() error
//	Err() error
// }
type mapping<S, T> struct {
	Iter<S>
	f func(S) T
	transformError func(error) error
}

func (m *mapping) Value() bool {
	return m.iter.Next()
}

func (m *mapping) Err() error {
	return m.transformError
}

func BufioScanner(r *bufio.Scanner) Iter<string> {
	return bufioScanner(r)
}

type bufioScanner struct {
	r *bufio.Scanner
}

func (b bufioScanner) Next() bool {
	return b.r.Scan()
}

func (b bufioScanner) Err() error {
	return b.r.Err()
}

func (b bufioScanner) Value() string {
	return b.r.Text()
}

func (b bufioScanner) Close() error {
	return nil
}

type slice<T> struct {
	first bool
	values []T
}

// Slice implements Iter on a slice.
// The values are traversed from beginning to end.
func NewSlice<T>(values []T) Iter<T> {
	return &slice{
		first: true,
		values: values,
	}
}

func (s *slice) Next() bool {
	if s.first {
		s.first = false
	} else {
		s.values = s.elems[1:]
	}
	return len(s.values) > 0
}

func (s *slice) Err() error {
	return nil
}

func (s *slice) Close() error {
	return nil
}

func (s *slice) Value() s.T {
	return s.values[0]
}


// Sequence returns an iterator that iterates
// through each of the given iterators in turn.
// The first error will terminate any remaining
// iterators. All are closed in turn when or before
// the returned iterator is closed.
//
func Sequence(iters ...Iter<T>) Iter<T> {
	return &sequence{
		iters: iters,
	}
}

type sequence struct<T> {
	iters []Iter<T>
	err error
}

func (s *sequence) Next() bool {
	for len(s.iters) > 0 {
		iter := s.iters[0]
		if iter.Next() {
			return true
		}
		if err := iter.Err(); err != nil {
			s.err = err
			s.Close()
			s.iters = nil
			return false
		}
		s.iters = s.iters[1:]
	}
	return false
}

func (s *sequence) Err() error {
	if s.err != nil {
		return s.err
	}
	if len(s.iters) > 0 {
		s.err = s.iters[0].Err()
		return s.err
	}
	return nil
}

func (s *sequence) Value() s.T {
	if len(s.iters) == 0 {
		// Note the use of "zero" here, as a more
		// general form of "nil" that is a valid
		// zero value for any type, not just pointers.
		return zero
	}
	return s.iters[0].Value()
}

func (s *sequence) Close() error {
	var closeErr error
	for len(s.iters) > 0 {
		if err := s.iters[0].Close(); err != nil && closeErr == nil {
			closeErr = err
		}
		s.iters = s.iters[1:]
	}
	return closeErr
}

// Concurrent returns a channel reading the
// results of al the given iterators running
// concurrently.
func Concurrent<T>(iters ...Iter<T>) chan T {
	c := make(chan T)
	go func() {
		var wg sync.WaitGroup
		for _, iter := range iters {
			iter := iter
			go func() {
				defer wg.Done()
				defer iter.Close()
				for iter.Next() {
					c <- iter.Value()
				}
			}()
		}
		wg.Wait()
		close(c)
	}()
	return c
}
