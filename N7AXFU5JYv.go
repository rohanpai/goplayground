// http://golang.org/src/pkg/encoding/json/decode.go
// Sample program to show how to implement a custom error type
// based on the json package in the standard library.
package main

import (
	&#34;fmt&#34;
	&#34;reflect&#34;
)

// An UnmarshalTypeError describes a JSON value that was
// not appropriate for a value of a specific Go type.
type UnmarshalTypeError struct {
	Value string       // description of JSON value
	Type  reflect.Type // type of Go value it could not be assigned to
}

// Error implements the error interface.
func (e *UnmarshalTypeError) Error() string {
	return &#34;json: cannot unmarshal &#34; &#43; e.Value &#43; &#34; into Go value of type &#34; &#43; e.Type.String()
}

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

// Error implements the error interface.
func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return &#34;json: Unmarshal(nil)&#34;
	}

	if e.Type.Kind() != reflect.Ptr {
		return &#34;json: Unmarshal(non-pointer &#34; &#43; e.Type.String() &#43; &#34;)&#34;
	}
	return &#34;json: Unmarshal(nil &#34; &#43; e.Type.String() &#43; &#34;)&#34;
}

// user is a type for use in the Unmarshal call.
type user struct {
	Name int
}

// main is the entry point for the application.
func main() {
	var u user
	err := Unmarshal([]byte(`{&#34;name&#34;:&#34;bill&#34;}`), u) // Run with a value and pointer.
	if err != nil {
		switch e := err.(type) {
		case *UnmarshalTypeError:
			fmt.Printf(&#34;UnmarshalTypeError: Value[%s] Type[%v]\n&#34;, e.Value, e.Type)
		case *InvalidUnmarshalError:
			fmt.Printf(&#34;InvalidUnmarshalError: Type[%v]\n&#34;, e.Type)
		default:
			fmt.Println(err)
		}
		return
	}

	fmt.Println(&#34;Name:&#34;, u.Name)
}

// Unmarshal simulates an unmarshal call that always fails.
func Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &amp;InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	return &amp;UnmarshalTypeError{&#34;string&#34;, reflect.TypeOf(v)}
}
