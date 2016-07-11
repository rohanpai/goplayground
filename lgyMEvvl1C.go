package main

import &#34;fmt&#34;

type Publisher interface {
    Publish(value interface{})
}

type Observer interface {
    Notify(value interface{})
}

type ObserverFunc func(value interface{})

func (fn ObserverFunc) Notify(value interface{}){
    fn(value)
}

type Observable []Observer

func (observers *Observable) AddObserver(a Observer){
    *observers = append(*observers, a)
}

func (observers Observable) Publish(value interface{}){
    for _, obs := range observers {
        obs.Notify(value)
    }
}

type Field struct {
	Value int64
	Observable
}

func (f *Field) Set(v int64){
	f.Value = v
	f.Publish(v)
}

func Listen(value interface{}){
	fmt.Printf(&#34;new value 1: %v\n&#34;, value)
}

func Listen2(value interface{}){
	fmt.Printf(&#34;new value 2: %v\n&#34;, value)
}

func main() {
	v := &amp;Field{}
	v.AddObserver(ObserverFunc(Listen))
	v.AddObserver(ObserverFunc(Listen2))
	v.Set(105)
	
	fmt.Println(&#34;Hello, playground&#34;)
}