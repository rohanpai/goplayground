package main

// === package di ===
import &#34;reflect&#34;

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
		// this will panic if it&#39;s wrong type, should print runtime.Caller(1)
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
	println(d.name, &#34;: woof&#34;)
}

func NewDog(name string) Animal {
	return &amp;Dog{name}
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
	Resolve(&#34;Animal&#34;, &amp;NewAnimal)
	Resolve(&#34;Animal&#34;, &amp;System.NewAnimal)
	Resolve(&#34;Factory&#34;, &amp;Factory)

	// in the binding part
	Bind(&#34;Animal&#34;, NewDog)

	factory := &amp;DogFactory{}
	Bind(&#34;Factory&#34;, factory)

	// somewhere in the module
	a := NewAnimal(&#34;a&#34;)
	a.Speak()

	b := Factory.NewAnimal(&#34;b&#34;)
	b.Speak()

	c := System.NewAnimal(&#34;c&#34;)
	c.Speak()
}
