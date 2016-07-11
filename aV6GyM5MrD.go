package main

import (
	&#34;fmt&#34;
	&#34;reflect&#34;
)

var dict = map[string]string{
	&#34;Hello!&#34;:                 &#34;Hallo!&#34;,
	&#34;What&#39;s up?&#34;:             &#34;Was geht?&#34;,
	&#34;translate this&#34;:         &#34;übersetze dies&#34;,
	&#34;point here&#34;:             &#34;zeige hier her&#34;,
	&#34;translate this as well&#34;: &#34;übersetze dies auch...&#34;,
	&#34;and one more&#34;:           &#34;und noch eins&#34;,
	&#34;deep&#34;:                   &#34;tief&#34;,
}

type I interface{}

type A struct {
	Greeting string
	Message  string
	Pi       float64
}

type B struct {
	Struct    A
	Ptr       *A
	Answer    int
	Map       map[string]string
	StructMap map[string]interface{}
	Slice     []string
}

func create() I {
	// The type C is actually hidden, but reflection allows us to look inside it
	type C struct {
		String string
	}

	return B{
		Struct: A{
			Greeting: &#34;Hello!&#34;,
			Message:  &#34;translate this&#34;,
			Pi:       3.14,
		},
		Ptr: &amp;A{
			Greeting: &#34;What&#39;s up?&#34;,
			Message:  &#34;point here&#34;,
			Pi:       3.14,
		},
		Map: map[string]string{
			&#34;Test&#34;: &#34;translate this as well&#34;,
		},
		StructMap: map[string]interface{}{
			&#34;C&#34;: C{
				String: &#34;deep&#34;,
			},
		},
		Slice: []string{
			&#34;and one more&#34;,
		},
		Answer: 42,
	}
}

func main() {
	// Some example test cases so you can mess around and see if it&#39;s working
	// To check if it&#39;s correct look at the output, no automated checking here

	// Test the simple cases
	{
		fmt.Println(&#34;Test with nil pointer to struct:&#34;)
		var original *B
		translated := translate(original)
		fmt.Println(&#34;original:  &#34;, original)
		fmt.Println(&#34;translated:&#34;, translated)
		fmt.Println()
	}
	{
		fmt.Println(&#34;Test with nil pointer to interface:&#34;)
		var original *I
		translated := translate(original)
		fmt.Println(&#34;original:  &#34;, original)
		fmt.Println(&#34;translated:&#34;, translated)
		fmt.Println()
	}
	{
		fmt.Println(&#34;Test with struct that has no elements:&#34;)
		type E struct {
		}
		var original E
		translated := translate(original)
		fmt.Println(&#34;original:  &#34;, original)
		fmt.Println(&#34;translated:&#34;, translated)
		fmt.Println()
	}
	{
		fmt.Println(&#34;Test with empty struct:&#34;)
		var original B
		translated := translate(original)
		fmt.Println(&#34;original:  &#34;, original, &#34;-&gt;&#34;, original.Ptr)
		fmt.Println(&#34;translated:&#34;, translated, &#34;-&gt;&#34;, translated.(B).Ptr)
		fmt.Println()
	}

	// Imagine we have no influence on the value returned by create()
	created := create()
	{
		// Assume we know that `created` is of type B
		fmt.Println(&#34;Translating a struct:&#34;)
		original := created.(B)
		translated := translate(original)
		fmt.Println(&#34;original:  &#34;, original, &#34;-&gt;&#34;, original.Ptr)
		fmt.Println(&#34;translated:&#34;, translated, &#34;-&gt;&#34;, translated.(B).Ptr)
		fmt.Println()
	}
	{
		// Assume we don&#39;t know created&#39;s type
		fmt.Println(&#34;Translating a struct wrapped in an interface:&#34;)
		original := created
		translated := translate(original)
		fmt.Println(&#34;original:  &#34;, original, &#34;-&gt;&#34;, original.(B).Ptr)
		fmt.Println(&#34;translated:&#34;, translated, &#34;-&gt;&#34;, translated.(B).Ptr)
		fmt.Println()
	}
	{
		// Assume we don&#39;t know B&#39;s type and want to pass a pointer
		fmt.Println(&#34;Translating a pointer to a struct wrapped in an interface:&#34;)
		original := &amp;created
		translated := translate(original)
		fmt.Println(&#34;original:  &#34;, (*original), &#34;-&gt;&#34;, (*original).(B).Ptr)
		fmt.Println(&#34;translated:&#34;, (*translated.(*I)), &#34;-&gt;&#34;, (*translated.(*I)).(B).Ptr)
		fmt.Println()
	}
	{
		// Assume we have a struct that contains an interface of an unknown type
		fmt.Println(&#34;Translating a struct containing a pointer to a struct wrapped in an interface:&#34;)
		type D struct {
			Payload *I
		}
		original := D{
			Payload: &amp;created,
		}
		translated := translate(original)
		fmt.Println(&#34;original:  &#34;, original, &#34;-&gt;&#34;, (*original.Payload), &#34;-&gt;&#34;, (*original.Payload).(B).Ptr)
		fmt.Println(&#34;translated:&#34;, translated, &#34;-&gt;&#34;, (*translated.(D).Payload), &#34;-&gt;&#34;, (*(translated.(D).Payload)).(B).Ptr)
		fmt.Println()
	}
}

func translate(obj interface{}) interface{} {
	// Wrap the original in a reflect.Value
	original := reflect.ValueOf(obj)

	copy := reflect.New(original.Type()).Elem()
	translateRecursive(copy, original)

	// Remove the reflection wrapper
	return copy.Interface()
}

func translateRecursive(copy, original reflect.Value) {
	switch original.Kind() {
	// The first cases handle nested structures and translate them recursively

	// If it is a pointer we need to unwrap and call once again
	case reflect.Ptr:
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don&#39;t end up in
		// an infinite recursion
		originalValue := original.Elem()
		// Check if the pointer is nil
		if !originalValue.IsValid() {
			return
		}
		// Allocate a new object and set the pointer to it
		copy.Set(reflect.New(originalValue.Type()))
		// Unwrap the newly created pointer
		translateRecursive(copy.Elem(), originalValue)

	// If it is an interface (which is very similar to a pointer), do basically the
	// same as for the pointer. Though a pointer is not the same as an interface so
	// note that we have to call Elem() after creating a new object because otherwise
	// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		translateRecursive(copyValue, originalValue)
		copy.Set(copyValue)

	// If it is a struct we translate each field
	case reflect.Struct:
		for i := 0; i &lt; original.NumField(); i &#43;= 1 {
			translateRecursive(copy.Field(i), original.Field(i))
		}

	// If it is a slice we create a new slice and translate each element
	case reflect.Slice:
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i &lt; original.Len(); i &#43;= 1 {
			translateRecursive(copy.Index(i), original.Index(i))
		}

	// If it is a map we create a new map and translate each value
	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			// New gives us a pointer, but again we want the value
			copyValue := reflect.New(originalValue.Type()).Elem()
			translateRecursive(copyValue, originalValue)
			copy.SetMapIndex(key, copyValue)
		}

	// Otherwise we cannot traverse anywhere so this finishes the the recursion

	// If it is a string translate it (yay finally we&#39;re doing what we came for)
	case reflect.String:
		translatedString := dict[original.Interface().(string)]
		copy.SetString(translatedString)

	// And everything else will simply be taken from the original
	default:
		copy.Set(original)
	}

}