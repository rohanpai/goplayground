package main

// === package di ===
import "reflect"

var (
	// define how what name binds to which value
	bindings = make(map[string]reflect.Value)
	// define where to bind the values by name
	targets = make(map[string][]reflect.Value)
)

// add a new binding target
func Resolve(name string, service interface{}) {
	if _, ok := targets[name]; !ok {
		targets[name] = make([]reflect.Value, 0)
	}
	v := reflect.ValueOf(service)

	// for better binding messages also store runtime.Caller(1)
	//   that info could be also used for more complex bindings
	targets[name] = append(targets[name], v)

	if binding, ok := bindings[name]; ok {
		v.Set(binding)
	}
}

// bind a name to the targets
func Bind(name string, service interface{}) {
	v := reflect.ValueOf(service)
	bindings[name] = v
	for _, target := range targets[name] {
		// this will panic if it's wrong type, should print runtime.Caller(1)
		target.Elem().Set(v)
	}
}

// func Check() { ... check whether everything got bound ... }

// === end package di ===

// generic interface for an animal
type Animal interface {
	Speak()
}

// example animal
type Dog struct{ name string }

func (d *Dog) Speak() {
	println(d.name, ": woof")
}

func NewDog(name string) Animal {
	return &Dog{name}
}

// example factory
type DogFactory struct{}

func (ds *DogFactory) NewAnimal(name string) Animal {
	return NewDog(name)
}

// different binding points

// function binding point in a module
var NewAnimal func(string) Animal

// interface binding point in a module
var Factory interface {
	NewAnimal(string) Animal
}

// binding points can be inside a struct
var System struct {
	NewAnimal func(string) Animal
}

func main() {
	// in the module init()
	Resolve("Animal", &NewAnimal)
	Resolve("Animal", &System.NewAnimal)
	Resolve("Factory", &Factory)

	// in the binding part
	Bind("Animal", NewDog)

	factory := &DogFactory{}
	Bind("Factory", factory)

	// somewhere in the module
	a := NewAnimal("a")
	a.Speak()

	b := Factory.NewAnimal("b")
	b.Speak()

	c := System.NewAnimal("c")
	c.Speak()
}
