package main

import (
	"fmt"
	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
	"github.com/gonum/graph/search"
)

func main() {
	g := concrete.NewGonumGraph(true)
	var n0, n1, n2, n3 concrete.GonumNode = 0, 1, 2, 3
	g.AddNode(n0, []graph.Node{n1, n2})
	g.AddEdge(concrete.GonumEdge{n2, n3})
	path, v := search.BreadthFirstSearch(n0, n3, g)
	fmt.Println("path:", path)
	fmt.Println("nodes visited:", v)
}
