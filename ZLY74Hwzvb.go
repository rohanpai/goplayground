package main

import (
	&#34;fmt&#34;
	&#34;reflect&#34;
	&#34;strings&#34;
)

// Centralised Composite Tree
//
// NOTE: Since interface{} allows any type need to check if the parent and children nodes
// provided are pointers. The implementation allows this to be turned off, but the client of
// the collection needs to ensure that the key of the objects is not changes.
// Also structure may not have a key, eg in the case of one that contains a slice
type Tree interface {
	Root() Node
	Add(parent Node, child Node) (bool, Node)
	Parent(child Node) Node
	Walk(func(depth int, node Node))
}

// Node, an interface{} type store in the tree.
// Tider then writing interface{} everywhere
type Node interface{}

// Fluent API for building a tree
type Builder interface {
	Add(n Node) Builder
	Down() Builder
	Up() Builder
	Build() Tree
}

// Create a new tree with a hidden root and return the root and the tree
func NewTree() (Tree, Node) {
	t := &amp;ctree{
		make(map[Node][]Node),
		make(map[Node]Node),
		new(interface{}),
		true,
	}
	t.c2[t.root] = nil
	return t, t.root
}

// Creates a tree that allow non pointer node.
// Client beware, do not mutate object or pass non keyable objects
func NewTree_ImmutableNodes() (Tree, Node) {
	t := &amp;ctree{
		make(map[Node][]Node),
		make(map[Node]Node),
		new(interface{}),
		false,
	}
	t.c2[(Node)(t.root)] = nil
	return (Tree)(t), t.root
}

// START TEST
// Testing structure 
type mys struct {
	s string
}

func main() {
	ct, r := NewTree()
	a := mys{&#34;a&#34;}
	ct.Add(r, &amp;a)
	b := mys{&#34;b&#34;}
	c := mys{&#34;c&#34;}
	ct.Add(&amp;a, &amp;b)
	ct.Add(&amp;a, &amp;c)
	ct.Add(&amp;a, &amp;c)	// ignored as c already in the tree
	x := 1
	ct.Add(r, &amp;x)
	ct.Walk(p)

	fmt.Println()

	d := &#34;d&#34;
	d2 := &#34;d&#34;
	slice := &amp;[]int{1, 2, 3, 4}
	BuildTree_ImmutableNodes().
		Add(mys{&#34;a&#34;}).Down().
		    Add(&amp;mys{&#34;b&#34;}).
		    Add(mys{&#34;c&#34;}).
		    Add(&amp;d2).
		    Add(d).
		    Add(&#34;d&#34;).Up().
		Add(slice).
		Build().Walk(p)
}

// Function given to the walker for printing nodes
func p(d int, n Node) {
	s := strings.Repeat(&#34; &#34;, d)
	//	fmt.Println( reflect.TypeOf( n ) )
	switch v := n.(type) {
	case *mys:
		fmt.Printf(&#34;1)%s%v\n&#34;, s, v.s)
	case mys:
		fmt.Printf(&#34;2)%s%v\n&#34;, s, v.s)
	case *string:
		fmt.Printf(&#34;3)%s%v\n&#34;, s, v)
	case string:
		fmt.Printf(&#34;4)%s%v\n&#34;, s, v)
	case *int:
		fmt.Printf(&#34;5)%s%v\n&#34;, s, *v)
	case *[]int:
		fmt.Printf(&#34;6)%s%v\n&#34;, s, *v)	
	default:
		fmt.Printf(&#34;*)%s%v - %T\n&#34;, s, n, n)
	}
}
// END TEST

// The internal structure of the centralised tree
type ctree struct {
	// Map of parent to slice of children
	p2	map[Node][]Node
	// Map of child to parent
	c2	map[Node]Node
	// Hidden root node of the tree
	root	Node
	ptrOnly	bool
}

// Add a child to a parent and create the p2c and c2p pointers
// If the child already exists anywhere in the tree the operations fails and the
// existing parent is returned with the boolean status
func (t *ctree) Add(p Node, c Node) (bool, Node) {
	if t.ptrOnly {
		if reflect.TypeOf(c).Kind() != reflect.Ptr {
			panic(&#34;Child node is not a pointer&#34;)
		}
		if reflect.TypeOf(p).Kind() != reflect.Ptr {
			panic(&#34;Parent node is not a pointer&#34;)
		}
	}
	if orginal, exist := t.c2[c]; exist {
		return false, orginal
	}
	t.p2[p] = append(t.p2[p], c)
	t.c2[c] = p
	return true, nil
}

func (t *ctree) Parent(c Node) Node {
	return t.c2[c]
}

// A depth first walk of the tree calling the provided function at each node 
func (t *ctree) Walk(f func(int, Node)) {
	for _, o := range t.p2[t.root] {
		walk(t, o, 0, f)
	}
}

func walk(t *ctree, node Node, depth int, f func(int, Node)) {
	f(depth, node)
	for _, o := range t.p2[node] {
		walk(t, o, depth&#43;1, f)
	}
}

// Return the root of the tree
func (t *ctree) Root() Node {
	return t.root
}

// Builder&#39;s state
type builder struct {
	tree	Tree
	curr	Node
	last	Node
}

func BuildTree() Builder {
	t, r := NewTree()
	b := &amp;builder{t, r, nil}
	return b
}
func BuildTree_ImmutableNodes() Builder {
	t, r := NewTree_ImmutableNodes()
	b := &amp;builder{t, r, nil}
	return b
}
func (b *builder) Add(n Node) Builder {
	ok, _ := b.tree.Add(b.curr, n)
	if ok {
		b.last = n
	}
	return b
}
func (b *builder) Down() Builder {
	b.curr = b.last
	return b
}
func (b *builder) Up() Builder {
	b.curr = b.tree.Parent(b.curr)
	b.last = b.curr
	return b
}
func (b *builder) Build() Tree {
	return b.tree
}
