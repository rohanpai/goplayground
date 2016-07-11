package main

import (
	&#34;fmt&#34;
	&#34;reflect&#34;
)

type A struct{ Name string }
type B struct{ Name string }
type C struct{ Name string }

func main() {

	a := &amp;A{Name: &#34;A&#34;}
	b := &amp;B{Name: &#34;B&#34;}
	c := &amp;C{Name: &#34;C&#34;}

	t := make([]interface{}, 0)
	t = append(t, a)
	t = append(t, b)
	t = append(t, c)

	fmt.Println(&#34;--- without type switching&#34;)
	for _, thing := range t {
		switch thing.(type) {
		case *A:
			fmt.Println(reflect.TypeOf(thing), &#34;is of type A. Name is:&#34;, thing.(*A).Name)
		case *B:
			fmt.Println(reflect.TypeOf(thing), &#34;is of type B. Name is:&#34;, thing.(*B).Name)
		case *C:
			fmt.Println(reflect.TypeOf(thing), &#34;is of type C. Name is:&#34;, thing.(*C).Name)
		}
	}

	fmt.Println(&#34;--- type switching on the item&#34;)
	for _, thing := range t {
		switch o := thing.(type) {
		case *A:
			fmt.Println(reflect.TypeOf(o), &#34;is of type A. Name is:&#34;, o.Name)
		case *B:
			fmt.Println(reflect.TypeOf(o), &#34;is of type B. Name is:&#34;, o.Name)
		case *C:
			fmt.Println(reflect.TypeOf(o), &#34;is of type C. Name is:&#34;, o.Name)
		}
	}

	// *main.A is of type A. Name is: A
	// *main.B is of type B. Name is: B
	// *main.C is of type C. Name is: C

	// for _, thing := range t {
	//  switch o := thing.(type) {
	//    case *A, *B:
	//      fmt.Println(reflect.TypeOf(o), &#34;is of type A or B. Name is:&#34;, o.Name)
	//    case *C:
	//      fmt.Println(reflect.TypeOf(o), &#34;is of type C. Name is:&#34;, o.Name)
	//  }
	// }
	// ERROR:
	// o.Name undefined (type interface {} has no field or method Name)

	// It works, as long as we don&#39;t try to get the Name field from the struct
	fmt.Println(&#34;--- Grouping types&#34;)
	for _, thing := range t {
		switch o := thing.(type) {
		case *A, *B:
			fmt.Println(reflect.TypeOf(o), &#34;is of type A or B.&#34;)
		case *C:
			fmt.Println(reflect.TypeOf(o), &#34;is of type C.&#34;)
		}
	}
	// RESULT:
	// *main.A is of type A or B.
	// *main.B is of type A or B.
	// *main.C is of type C.

	fmt.Println(&#34;--- Double type switch all the way&#34;)
	for _, thing := range t {
		switch o := thing.(type) {
		case *A, *B:
			switch oo := o.(type) {
			case *A:
				fmt.Println(reflect.TypeOf(o), &#34;is of type A or B. Name is:&#34;, oo.Name)
			case *B:
				fmt.Println(reflect.TypeOf(o), &#34;is of type A or B. Name is:&#34;, oo.Name)
			}
		case *C:
			fmt.Println(reflect.TypeOf(o), &#34;is of type C. Name is:&#34;, o.Name)
		}
	}

	fmt.Println(&#34;--- The Go way™&#34;)
	for _, thing := range t {
		switch o := thing.(type) {
		case *A:
			fmt.Println(reflect.TypeOf(o), &#34;is of type A or B. Name is:&#34;, o.Name)
		case *B:
			fmt.Println(reflect.TypeOf(o), &#34;is of type A or B. Name is:&#34;, o.Name)
		case *C:
			fmt.Println(reflect.TypeOf(o), &#34;is of type C. Name is:&#34;, o.Name)
		}
	}

}