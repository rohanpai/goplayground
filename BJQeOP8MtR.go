/*
golang.org/pkg/container/heap
https://code.google.com/p/go/source/browse/#hg%2Fsrc%2Fpkg%2Fcontainer%2Fheap

Package heap provides heap operations for any type that implements heap.Interface. A heap is a tree with the property that each node is the minimum-valued node in its subtree.

The minimum element in the tree is the root, at index 0.

A heap is a common way to implement a priority queue. To build a priority queue, implement the Heap interface with the (negative) priority as the ordering for the Less method, so Push adds items while Pop removes the highest-priority item from the queue.
*/
package main

import (
	"fmt"
	"sort"
)

// From here, it is Go source code.
// http://golang.org/src/pkg/container/heap/heap.go

// Any type that implements heap.Interface may be used as a
// min-heap with the following invariants (established after
// Init has been called or if the data is empty or sorted):
//
//	!h.Less(j, i) for 0 <= i < h.Len() and j = 2*i+1 or 2*i+2 and j < h.Len()
//
// Note that Push and Pop in this interface are for package heap's
// implementation to call.  To add and remove things from the heap,
// use heap.Push and heap.Pop.
type Interface interface {
	sort.Interface
	Push(x interface{}) // add x as element Len()
	Pop() interface{}   // remove and return element Len() - 1.
}

// A heap must be initialized before any of the heap operations
// can be used. Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
// Its complexity is O(n) where n = h.Len().
func Init(h Interface) {
	// heapify
	n := h.Len()
	for i := n/2 - 1; i >= 0; i-- {
		down(h, i, n)
	}
}

// Push pushes the element x onto the heap. The complexity is
// O(log(n)) where n = h.Len().
func Push(h Interface, x interface{}) {
	h.Push(x)
	// Min-Heapify
	up(h, h.Len()-1)
}

// Min-Heapify from the bottom.
// This actually starts from the mid index.
// We need down to heapify all.
func up(h Interface, j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		j = i
	}
}

// Min-Heapify from index
func down(h Interface, i, n int) {
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && !h.Less(j1, j2) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		i = j
	}
}

// Pop removes the minimum element (according to Less) from the heap
// and returns it. The complexity is O(log(n)) where n = h.Len().
// It is equivalent to Remove(h, 0).
func Pop(h Interface) interface{} {
	n := h.Len() - 1
	h.Swap(0, n)
	down(h, 0, n)
	return h.Pop()
}

// Remove removes the element at index i from the heap.
// The complexity is O(log(n)) where n = h.Len().
func Remove(h Interface, i int) interface{} {
	n := h.Len() - 1
	if n != i {
		h.Swap(i, n)
		down(h, i, n)
		up(h, i)
	}
	return h.Pop()
}

// Note that Go source code heap implements "MIN" heap
// Go heap embeds sort.Interface
// Therefore, we need to define our custom Interface

// An IntHeap is a min-heap of ints
// Go's Init uses Build-Min-Heap, not Max Heap
type IntHeap []int

func (s IntHeap) Len() int           { return len(s) }
func (s IntHeap) Less(i, j int) bool { return s[i] < s[j] }
func (s IntHeap) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Push and Pop use pointer receivers
// because they modify the slice's length,
func (s *IntHeap) Push(val interface{}) {
	*s = append(*s, val.(int))
}

// Pop removes the last element and return the new array.
// This is only for sort.Interface.
// Nothing to do with sorting because it just returns the last element.
// Pop in heap does the sorting job.
func (s *IntHeap) Pop() interface{} {
	osl := *s
	n := len(osl)
	ns := osl[n-1]
	*s = osl[0 : n-1]
	return ns
}

func main() {
	slice := &IntHeap{12, 100, -15, 200, -5, 3, -12, 7}
	fmt.Println("Before Push:", slice)
	// Before Push: &[12 100 -15 200 -5 3 -12 7]

	// Heap Push implements "up" that does Min-Heapify.
	// So as we Push, the Heap(tree) is automatically
	// Half-Heapified(Sorted) since only 'up' is implemented
	Push(slice, -10)
	fmt.Println("Before Init:", slice)
	// Before Init: &[-10 12 -15 100 -5 3 -12 7 200]

	Init(slice)
	fmt.Println("After Init:", slice)
	// After Init: &[-15 -5 -12 7 12 3 -10 100 200]

	// Don't use slice.Pop()
	// slice.Pop() is only sort.Interface
	// not doing Heapifying
	fmt.Println(Pop(slice))
	// -15
	// returns a minimum elements

	Push(slice, -100)
	fmt.Println("After Push(slice, -100):", slice)
	// After Push(slice, -100): &[-100 -12 -10 -5 12 3 200 100 7]

	Push(slice, 5)
	fmt.Println("After Push(slice, 5):", slice)
	// After Push(slice, 5): &[-100 -12 -10 -5 5 3 200 100 7 12]

	Push(slice, -300)
	fmt.Println("After Push(slice, -300):", slice)
	// After Push(slice, -300): &[-300 -100 -10 -5 -12 3 200 100 7 12 5]

	Init(slice)
	fmt.Println("After Init:", slice)
	// After Init: &[-300 -100 -10 -5 -12 3 200 100 7 12 5]

	// Heap-Sort is simple
	// We just keep popping off
	// and it will keep returning the minimum elements
	// doing Min-Heapifying at the same time
	print("Heap Sort: ")
	for slice.Len() != 0 {
		fmt.Print(Pop(slice), ",")
	}
	// Heap Sort: -300,-100,-12,-10,-5,3,5,7,12,100,200,
}
