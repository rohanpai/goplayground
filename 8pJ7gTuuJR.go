package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

type IntMap interface {
	Get(k int) (v int, found bool)
	Set(k int, v int)
}

// Single lock around the map
type Locked struct {
	e map[int]int
	m sync.RWMutex
}

func NewLocked() *Locked {
	return &Locked{make(map[int]int), sync.RWMutex{}}
}

func (c *Locked) Get(k int) (v int, ok bool) {
	c.m.RLock()
	v, ok = c.e[k]
	c.m.RUnlock()
	return
}

func (c *Locked) Set(k int, v int) {
	c.m.Lock()
	c.e[k] = v
	c.m.Unlock()
}

// Multiple locks around a map
type stripe struct {
	e map[int]int
	m sync.RWMutex
}

const stripes = 53

type Striped [stripes]stripe

func NewStriped() *Striped {
	m := &Striped{}
	for i := range m {
		m[i].e = make(map[int]int)
	}
	return m
}

func (c *Striped) Get(k int) (v int, ok bool) {
	s := &c[k%stripes]
	s.m.RLock()
	v, ok = s.e[k]
	s.m.RUnlock()
	return
}

func (c *Striped) Set(k int, v int) {
	s := &c[k%stripes]
	s.m.Lock()
	s.e[k] = v
	s.m.Unlock()
}

// atomically swapped map
type Atomic struct {
	items unsafe.Pointer // map[int]int
	m     sync.Mutex
}

func NewAtomic() *Atomic {
	items := make(map[int]int)
	return &Atomic{unsafe.Pointer(&items), sync.Mutex{}}
}

func (c *Atomic) get() map[int]int {
	return *(*map[int]int)(atomic.LoadPointer(&c.items))
}

func (c *Atomic) set(m map[int]int) {
	atomic.StorePointer(&c.items, unsafe.Pointer(&m))
}

func (c *Atomic) Get(k int) (v int, ok bool) {
	v, ok = c.get()[k]
	return
}

func (c *Atomic) Set(k int, v int) {
	c.m.Lock()
	m := c.get()
	cp := make(map[int]int, len(m)+1)
	for k, v := range m {
		cp[k] = v
	}
	cp[k] = v
	c.set(cp)
	c.m.Unlock()
	return
}

func main() {
	c := NewAtomic()
	c.Set(123, 312)
	c.Set(212, 412)
	a, b := c.Get(123)
	fmt.Println(a, b)
	a, b = c.Get(212)
	fmt.Println(a, b)
}
