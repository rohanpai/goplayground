package main

import &#34;fmt&#34;

func main() {
	b := NewBag()
	b.Set(&#34;a&#34;, 1)
	b.Set(&#34;b&#34;, 2)
	b.Set(&#34;c&#34;, 3)

	it := b.Iter()
	for it.Next() {
		fmt.Println(it.Key, &#34;=&#34;, it.Value)
	}
}

type Bag struct {
	m map[string]interface{}
}

func NewBag() *Bag {
	return &amp;Bag{make(map[string]interface{})}
}

func (b *Bag) Set(key string, value interface{}) {
	b.m[key] = value
}

func (b *Bag) Iter() *Iter {
	it := &amp;Iter{b: b, l: make([]string, len(b.m))}
	n := 0
	for k := range b.m {
		it.l[n] = k
		n&#43;&#43;
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
	it.n&#43;&#43;
	return true
}