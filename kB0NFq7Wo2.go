// https://github.com/google/btree/blob/master/btree.go
package main

import (
	&#34;fmt&#34;
	&#34;sort&#34;
)

func main() {
	tr := New(3)
	for i := Int(0); i &lt; 10; i&#43;&#43; {
		tr.ReplaceOrInsert(i)
	}
	fmt.Println(&#34;len:       &#34;, tr.Len())
	fmt.Println(&#34;get3:      &#34;, tr.Get(Int(3)))
	fmt.Println(&#34;get100:    &#34;, tr.Get(Int(100)))
	fmt.Println(&#34;replace5:  &#34;, tr.ReplaceOrInsert(Int(5)))
	fmt.Println(&#34;replace100:&#34;, tr.ReplaceOrInsert(Int(100)))
	/*
	   len:        10
	   get3:       3
	   get100:     &lt;nil&gt;
	   replace5:   5
	   replace100: &lt;nil&gt;
	*/
}

// Item represents a single object in the tree.
type Item interface {
	// Less tests whether the current item is less than the given argument.
	//
	// This must provide a strict weak ordering.
	// If !a.Less(b) &amp;&amp; !b.Less(a), we treat this to mean a == b (i.e. we can only
	// hold one of either a or b in the tree).
	Less(than Item) bool
}

// Int implements the Item interface for integers.
type Int int

// Less returns true if int(a) &lt; int(b).
func (a Int) Less(b Item) bool {
	return a &lt; b.(Int)
}

//------------------------------------
// node
//------------------------------------

// node is an internal node in a tree.
//
// It must at all times maintain the invariant that either
//   * len(children) == 0, len(items) unconstrained
//   * len(children) == len(items) &#43; 1
type node struct {
	items    items
	children children
	t        *BTree
}

// split splits the given node at the given index.  The current node shrinks,
// and this function returns the item that existed at that index and a new node
// containing all items/children after it.
func (n *node) split(i int) (Item, *node) {
	item := n.items[i]
	next := n.t.newNode()
	next.items = append(next.items, n.items[i&#43;1:]...)
	n.items = n.items[:i]
	if len(n.children) &gt; 0 {
		next.children = append(next.children, n.children[i&#43;1:]...)
		n.children = n.children[:i&#43;1]
	}
	return item, next
}

// maybeSplitChild checks if a child should be split, and if so splits it.
// Returns whether or not a split occurred.
func (n *node) maybeSplitChild(i, maxItems int) bool {
	if len(n.children[i].items) &lt; maxItems {
		return false
	}
	first := n.children[i]
	item, second := first.split(maxItems / 2)
	n.items.insertAt(i, item)
	n.children.insertAt(i&#43;1, second)
	return true
}

// insert inserts an item into the subtree rooted at this node, making sure
// no nodes in the subtree exceed maxItems items.  Should an equivalent item be
// be found/replaced by insert, it will be returned.
func (n *node) insert(item Item, maxItems int) Item {
	i, found := n.items.find(item)
	if found {
		out := n.items[i]
		n.items[i] = item
		return out
	}
	if len(n.children) == 0 {
		n.items.insertAt(i, item)
		return nil
	}
	if n.maybeSplitChild(i, maxItems) {
		inTree := n.items[i]
		switch {
		case item.Less(inTree):
			// no change, we want first split node
		case inTree.Less(item):
			i&#43;&#43; // we want second split node
		default:
			out := n.items[i]
			n.items[i] = item
			return out
		}
	}
	return n.children[i].insert(item, maxItems)
}

// get finds the given key in the subtree and returns it.
func (n *node) get(key Item) Item {
	i, found := n.items.find(key)
	if found {
		return n.items[i]
	} else if len(n.children) &gt; 0 {
		return n.children[i].get(key)
	}
	return nil
}

////////////////////////////////////////////////////////////////////

//------------------------------------
// items
//------------------------------------

// items stores items in a node.
type items []Item

// insertAt inserts a value into the given index, pushing all subsequent values
// forward.
func (s *items) insertAt(index int, item Item) {
	*s = append(*s, nil)
	if index &lt; len(*s) {
		copy((*s)[index&#43;1:], (*s)[index:])
	}
	(*s)[index] = item
}

// removeAt removes a value at a given index, pulling all subsequent values
// back.
func (s *items) removeAt(index int) Item {
	item := (*s)[index]
	copy((*s)[index:], (*s)[index&#43;1:])
	*s = (*s)[:len(*s)-1]
	return item
}

// pop removes and returns the last element in the list.
func (s *items) pop() (out Item) {
	index := len(*s) - 1
	out, *s = (*s)[index], (*s)[:index]
	return
}

// find returns the index where the given item should be inserted into this
// list.  &#39;found&#39; is true if the item already exists in the list at the given
// index.
func (s items) find(item Item) (index int, found bool) {
	i := sort.Search(len(s), func(i int) bool {
		return item.Less(s[i])
	})
	if i &gt; 0 &amp;&amp; !s[i-1].Less(item) {
		return i - 1, true
	}
	return i, false
}

////////////////////////////////////////////////////////////////////

//------------------------------------
// children
//------------------------------------

// children stores child nodes in a node.
type children []*node

// insertAt inserts a value into the given index, pushing all subsequent values
// forward.
func (s *children) insertAt(index int, n *node) {
	*s = append(*s, nil)
	if index &lt; len(*s) {
		copy((*s)[index&#43;1:], (*s)[index:])
	}
	(*s)[index] = n
}

// removeAt removes a value at a given index, pulling all subsequent values
// back.
func (s *children) removeAt(index int) *node {
	n := (*s)[index]
	copy((*s)[index:], (*s)[index&#43;1:])
	*s = (*s)[:len(*s)-1]
	return n
}

// pop removes and returns the last element in the list.
func (s *children) pop() (out *node) {
	index := len(*s) - 1
	out, *s = (*s)[index], (*s)[:index]
	return
}

////////////////////////////////////////////////////////////////////

//------------------------------------
// BTree
//------------------------------------

// New creates a new B-Tree with the given degree.
//
// New(2), for example, will create a 2-3-4 tree (each node contains 1-3 items
// and 2-4 children).
func New(degree int) *BTree {
	if degree &lt;= 1 {
		panic(&#34;bad degree&#34;)
	}
	return &amp;BTree{
		degree:   degree,
		freelist: make([]*node, 0, 32),
	}
}

// BTree is an implementation of a B-Tree.
//
// BTree stores Item instances in an ordered structure, allowing easy insertion,
// removal, and iteration.
//
// Write operations are not safe for concurrent mutation by multiple
// goroutines, but Read operations are.
type BTree struct {
	degree   int
	length   int
	root     *node
	freelist []*node
}

// maxItems returns the max number of items to allow per node.
func (t *BTree) maxItems() int {
	return t.degree*2 - 1
}

// minItems returns the min number of items to allow per node (ignored for the
// root node).
func (t *BTree) minItems() int {
	return t.degree - 1
}

func (t *BTree) newNode() (n *node) {
	index := len(t.freelist) - 1
	if index &lt; 0 {
		return &amp;node{t: t}
	}
	t.freelist, n = t.freelist[:index], t.freelist[index]
	return
}

func (t *BTree) freeNode(n *node) {
	if len(t.freelist) &lt; cap(t.freelist) {
		for i := range n.items {
			n.items[i] = nil // clear to allow GC
		}
		n.items = n.items[:0]
		for i := range n.children {
			n.children[i] = nil // clear to allow GC
		}
		n.children = n.children[:0]
		t.freelist = append(t.freelist, n)
	}
}

// ReplaceOrInsert adds the given item to the tree.  If an item in the tree
// already equals the given one, it is removed from the tree and returned.
// Otherwise, nil is returned.
//
// nil cannot be added to the tree (will panic).
func (t *BTree) ReplaceOrInsert(item Item) Item {
	if item == nil {
		panic(&#34;nil item being added to BTree&#34;)
	}
	if t.root == nil {
		t.root = t.newNode()
		t.root.items = append(t.root.items, item)
		t.length&#43;&#43;
		return nil
	} else if len(t.root.items) &gt;= t.maxItems() {
		item2, second := t.root.split(t.maxItems() / 2)
		oldroot := t.root
		t.root = t.newNode()
		t.root.items = append(t.root.items, item2)
		t.root.children = append(t.root.children, oldroot, second)
	}
	out := t.root.insert(item, t.maxItems())
	if out == nil {
		t.length&#43;&#43;
	}
	return out
}

// Len returns the number of items currently in the tree.
func (t *BTree) Len() int {
	return t.length
}

// Get looks for the key item in the tree, returning it.  It returns nil if
// unable to find that item.
func (t *BTree) Get(key Item) Item {
	if t.root == nil {
		return nil
	}
	return t.root.get(key)
}

////////////////////////////////////////////////////////////////////
