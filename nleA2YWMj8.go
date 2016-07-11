package main

import (
	"fmt"
	"reflect"
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
	if val.Kind() == reflect.Interface && !val.IsNil() {
		elm := val.Elem()
		if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
			val = elm
		}
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		address := "not-addressable"

		if valueField.Kind() == reflect.Interface && !valueField.IsNil() {
			elm := valueField.Elem()
			if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
				valueField = elm
			}
		}

		if valueField.Kind() == reflect.Ptr {
			valueField = valueField.Elem()

		}
		if valueField.CanAddr() {
			address = fmt.Sprintf("0x%X", valueField.Addr().Pointer())
		}

		fmt.Printf("Field Name: %s,\t Field Value: %v,\t Address: %v\t, Field type: %v\t, Field kind: %v\n", typeField.Name,
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