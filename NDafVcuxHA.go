package main

import (
	&#34;fmt&#34;
	&#34;github.com/karalabe/cookiejar/graph&#34;
	&#34;github.com/karalabe/cookiejar/graph/dfs&#34;
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
	fmt.Println(&#34;Path 0-&gt;5:&#34;, d.Path(5))
	fmt.Println(&#34;Order:&#34;, d.Order())
	fmt.Println(&#34;Reachable #4 #6:&#34;, d.Reachable(4), d.Reachable(6))

}
