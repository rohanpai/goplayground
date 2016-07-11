package main

import (
	&#34;fmt&#34;
	&#34;github.com/gonum/graph&#34;
	&#34;github.com/gonum/graph/concrete&#34;
	&#34;github.com/gonum/graph/search&#34;
)

func main() {
	g := concrete.NewGonumGraph(true)
	var n0, n1, n2, n3 concrete.GonumNode = 0, 1, 2, 3
	g.AddNode(n0, []graph.Node{n1, n2})
	g.AddEdge(concrete.GonumEdge{n2, n3})
	path, v := search.BreadthFirstSearch(n0, n3, g)
	fmt.Println(&#34;path:&#34;, path)
	fmt.Println(&#34;nodes visited:&#34;, v)
}
