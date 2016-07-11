package main

import (
	"fmt"
	"github.com/karalabe/cookiejar/graph"
	"github.com/karalabe/cookiejar/graph/dfs"
)

func main() {
	// Create the graph
	g := graph.New(7)
	g.Connect(0, 1)
	g.Connect(1, 2)
	g.Connect(2, 3)
	g.Connect(3, 4)
	g.Connect(3, 5)

	// Create the depth first search algo structure for g and source node #2
	d := dfs.New(g, 0)

	// Get the path between #0 (source) and #2
	fmt.Println("Path 0->5:", d.Path(5))
	fmt.Println("Order:", d.Order())
	fmt.Println("Reachable #4 #6:", d.Reachable(4), d.Reachable(6))

}
