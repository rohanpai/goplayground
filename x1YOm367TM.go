package main

import (
	&#34;fmt&#34;

	&#34;github.com/soniakeys/graph&#34;
)

// node represents a node in a directed graph.  It represents directed edges
// from the node with the handy DijkstraNeighbor type from the graph package.
type node struct {
	nbs  []graph.DijkstraNeighbor // directed edges as DijkstraNeighbors
	name string                   // example application specific data
}

// edge is a simple number representing an edge length/distance/weight.
type edge float64

// node implements graph.DijkstraNode, also fmt.Stringer
func (n *node) Neighbors([]graph.DijkstraNeighbor) []graph.DijkstraNeighbor {
	return n.nbs
}
func (n *node) String() string { return n.name }

// edge implements graph.DijkstraEdge
func (e edge) Distance() float64 { return float64(e) }

// edgeData struct for simple specification of example data
type edgeData struct {
	v1, v2 string
	l      float64
}

// example data
var (
	exampleEdges = []edgeData{
		{&#34;a&#34;, &#34;b&#34;, 7},
		{&#34;a&#34;, &#34;c&#34;, 9},
		{&#34;a&#34;, &#34;f&#34;, 14},
		{&#34;b&#34;, &#34;c&#34;, 10},
		{&#34;b&#34;, &#34;d&#34;, 15},
		{&#34;c&#34;, &#34;d&#34;, 11},
		{&#34;c&#34;, &#34;f&#34;, 2},
		{&#34;d&#34;, &#34;e&#34;, 6},
		{&#34;e&#34;, &#34;f&#34;, 9},
	}
	exampleStart = &#34;a&#34;
	exampleEnd   = &#34;e&#34;
)

// linkGraph constructs a linked representation of example data.
func linkGraph(g []edgeData, start, end string) (allNodes int, startNode, endNode *node) {
	all := map[string]*node{}
	// one pass over data to collect nodes
	for _, e := range g {
		if all[e.v1] == nil {
			all[e.v1] = &amp;node{name: e.v1}
		}
		if all[e.v2] == nil {
			all[e.v2] = &amp;node{name: e.v2}
		}
	}
	// second pass to link neighbors
	for _, ge := range g {
		n1 := all[ge.v1]
		n1.nbs = append(n1.nbs, graph.DijkstraNeighbor{edge(ge.l), all[ge.v2]})
	}
	return len(all), all[start], all[end]
}

func main() {
	// construct linked representation of example data
	allNodes, startNode, endNode :=
		linkGraph(exampleEdges, exampleStart, exampleEnd)
	// echo initial conditions
	fmt.Printf(&#34;Directed graph with %d nodes, %d edges\n&#34;,
		allNodes, len(exampleEdges))
	// run Dijkstra&#39;s shortest path algorithm
	p, l := graph.DijkstraShortestPath(startNode, endNode)
	if p == nil {
		fmt.Println(&#34;No path from start node to end node&#34;)
		return
	}
	fmt.Println(&#34;Shortest path:&#34;, p)
	fmt.Println(&#34;Path length:&#34;, l)
}
