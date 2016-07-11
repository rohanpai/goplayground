package main

import (
	&#34;fmt&#34;
	&#34;reflect&#34;
)

type Z struct {
	Id int
}

type V struct {
	Id int
	F  Z
}

type T struct {
	Id int
	F  V
}

func InspectStructV(val reflect.Value) {
	if val.Kind() == reflect.Interface &amp;&amp; !val.IsNil() {
		elm := val.Elem()
		if elm.Kind() == reflect.Ptr &amp;&amp; !elm.IsNil() &amp;&amp; elm.Elem().Kind() == reflect.Ptr {
			val = elm
		}
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i &lt; val.NumField(); i&#43;&#43; {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		address := &#34;not-addressable&#34;

		if valueField.Kind() == reflect.Interface &amp;&amp; !valueField.IsNil() {
			elm := valueField.Elem()
			if elm.Kind() == reflect.Ptr &amp;&amp; !elm.IsNil() &amp;&amp; elm.Elem().Kind() == reflect.Ptr {
				valueField = elm
			}
		}

		if valueField.Kind() == reflect.Ptr {
			valueField = valueField.Elem()

		}
		if valueField.CanAddr() {
			address = fmt.Sprintf(&#34;0x%X&#34;, valueField.Addr().Pointer())
		}

		fmt.Printf(&#34;Field Name: %s,\t Field Value: %v,\t Address: %v\t, Field type: %v\t, Field kind: %v\n&#34;, typeField.Name,
			valueField.Interface(), address, typeField.Type, valueField.Kind())

		if valueField.Kind() == reflect.Struct {
			InspectStructV(valueField)
		}
	}
}

func InspectStruct(v interface{}) {
	InspectStructV(reflect.ValueOf(v))
}
func main() {
	t := new(T)
	t.Id = 1
	t.F = *new(V)
	t.F.Id = 2
	t.F.F = *new(Z)
	t.F.F.Id = 3

	InspectStruct(t)
}