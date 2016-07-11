package main

import "fmt"

func main() {
	b := NewBag()
	b.Set("a", 1)
	b.Set("b", 2)
	b.Set("c", 3)

	it := b.Iter()
	for it.Next() {
		fmt.Println(it.Key, "=", it.Value)
	}
}

type Bag struct {
	m map[string]interface{}
}

func NewBag() *Bag {
	return &Bag{make(map[string]interface{})}
}

func (b *Bag) Set(key string, value interface{}) {
	b.m[key] = value
}

func (b *Bag) Iter() *Iter {
	it := &Iter{b: b, l: make([]string, len(b.m))}
	n := 0
	for k := range b.m {
		it.l[n] = k
		n++
	}
	return it
}

type Iter struct {
	Key   string
	Value interface{}

	b *Bag
	n int
	l []string
}

func (it *Iter) Next() bool {
	if it.n == len(it.l) {
		return false
	}
	it.Key = it.l[it.n]
	it.Value = it.b.m[it.Key]
	it.n++
	return true
}