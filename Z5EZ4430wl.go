package main

import (
	&#34;container/list&#34;
	&#34;fmt&#34;
	&#34;regexp&#34;
	&#34;strconv&#34;
	&#34;strings&#34;
)

func main() {
	str1 := `
D|B,1|C,1
E|C,1|F,1
B|A,1
A|F,1
`
	G1 := ConstructDirectedGraph(str1)
	vertex_D := G1.GetVertexByID(&#34;D&#34;)
	fmt.Print(&#34;D&#39;s Predecessor Vertices: &#34;)
	for e := vertex_D.Predecessor.Front(); e != nil; e = e.Next() {
		fmt.Printf(&#34;%v &#34;, e.Value.(*Vertex).ID)
	}
	// D&#39;s Predecessor Vertices:
	println()

	outcome1 := DAG(G1)
	for e := outcome1.Front(); e != nil; e = e.Next() {
		fmt.Printf(&#34;%v → &#34;, e.Value.(*Vertex).ID)
	}
}

func DAG(G *Graph) *list.List {
	nolist := list.New()
	for vtx := G.GetVertexList().Front(); vtx != nil; vtx = vtx.Next() {
		if vtx.Value.(*Vertex).Predecessor.Len() == 0 {
			nolist.PushBack(vtx)
		}
	}

	result := list.New()
	for nolist.Len() != 0 {
		n := nolist.Front()
		nolist.Remove(nolist.Front())
		result.PushBack(n)
		for m := n.Value.(*Vertex).GetAdjacentVertices().Front(); m != nil; m = m.Next() {
			G.DeleteEdgeFrom(n.Value.(*Vertex), m.Value.(*Vertex))
			if m.Value.(*Vertex).Predecessor.Len() == 0 {
				nolist.PushBack(m)
			}
		}
	}

	if G.EdgeList.Len() != 0 {
		panic(&#34;The graph has a cycle; NOT a DAG !!!&#34;)
	}

	return result
}

type Graph struct {
	VertexList *list.List
	EdgeList   *list.List
}

func NewGraph() *Graph {
	return &amp;Graph{
		list.New(),
		list.New(),
	}
}

type Vertex struct {
	ID                  string
	Color               string
	EdgesFromThisVertex *list.List
	Predecessor         *list.List
	timestamp_d         int64
	timestamp_f         int64
}

func NewVertex(input_id string) *Vertex {
	return &amp;Vertex{
		ID:                  input_id,
		Color:               &#34;white&#34;,
		EdgesFromThisVertex: list.New(),
		Predecessor:         list.New(),
		timestamp_d:         9999999999,
		timestamp_f:         9999999999,
	}
}

type Edge struct {
	SourceVertex      *Vertex
	DestinationVertex *Vertex
	Weight            int
}

func NewEdge(source, destination *Vertex, weight int) *Edge {
	return &amp;Edge{
		source,
		destination,
		weight,
	}
}

func (A *Vertex) ConnectEdgeWithVertex(edges ...*Edge) {
	for _, edge := range edges {
		A.EdgesFromThisVertex.PushBack(edge)
	}
}

func (A *Vertex) GetEdgesFromThisVertex() chan *Edge {
	edgechan := make(chan *Edge)
	go func() {
		defer close(edgechan)
		for e := A.EdgesFromThisVertex.Front(); e != nil; e = e.Next() {
			edgechan &lt;- e.Value.(*Edge)
		}
	}()
	return edgechan
}

func (A *Vertex) GetAdjacentVertices() *list.List {
	result := list.New()
	for edge := range A.GetEdgesFromThisVertex() {
		result.PushBack(edge.DestinationVertex)
	}
	return result
}

func ConstructDirectedGraph(input_str string) *Graph {
	var validID = regexp.MustCompile(`\t{1,}`)
	newstr := validID.ReplaceAllString(input_str, &#34; &#34;)
	newstr = strings.TrimSpace(newstr)
	lines := strings.Split(newstr, &#34;\n&#34;)

	new_graph := NewGraph()

	for _, line := range lines {
		fields := strings.Split(line, &#34;|&#34;)

		SourceID := fields[0]
		edgepairs := fields[1:]

		new_graph.FindOrConstruct(SourceID)

		for _, pair := range edgepairs {
			if len(strings.Split(pair, &#34;,&#34;)) == 1 {
				continue
			}
			DestinationID := strings.Split(pair, &#34;,&#34;)[0]
			weight := StrToInt(strings.Split(pair, &#34;,&#34;)[1])

			src_vertex := new_graph.FindOrConstruct(SourceID)
			des_vertex := new_graph.FindOrConstruct(DestinationID)

			edge := NewEdge(src_vertex, des_vertex, weight)
			src_vertex.ConnectEdgeWithVertex(edge)
			des_vertex.Predecessor.PushBack(src_vertex)

			new_graph.EdgeList.PushBack(edge)
		}
	}
	return new_graph
}

func StrToInt(input_str string) int {
	result, err := strconv.Atoi(input_str)
	if err != nil {
		panic(&#34;failed to convert string&#34;)
	}
	return result
}

func (G *Graph) GetVertexByID(id string) *Vertex {
	for vtx := G.VertexList.Front(); vtx != nil; vtx = vtx.Next() {
		// NOT  vtx.Value.(Vertex).ID
		if vtx.Value.(*Vertex).ID == id {
			return vtx.Value.(*Vertex)
		}
	}
	return nil
}

func (G *Graph) FindOrConstruct(id string) *Vertex {
	vertex := G.GetVertexByID(id)
	if vertex == nil {
		vertex = NewVertex(id)

		// then add this vertex to the graph
		G.VertexList.PushBack(vertex)
	}
	return vertex
}

func (G *Graph) GetVertexList() *list.List {
	return G.VertexList
}

func (G *Graph) DeleteVertex(A *Vertex) {
	for vtx := G.VertexList.Front(); vtx != nil; vtx = vtx.Next() {
		if vtx.Value.(*Vertex) == A {
			// remove from the graph
			G.VertexList.Remove(vtx)
		}
	}
	for edge := range A.GetEdgesFromThisVertex() {
		G.DeleteEdgeFrom(A, edge.DestinationVertex)
		for vtx := edge.DestinationVertex.Predecessor.Front(); vtx != nil; vtx = vtx.Next() {
			edge.DestinationVertex.Predecessor.Remove(vtx)
		}
	}

}

func (G *Graph) DeleteEdgeFrom(A, B *Vertex) {
	for edge := G.EdgeList.Front(); edge != nil; edge = edge.Next() {
		if edge.Value.(*Edge).SourceVertex == A &amp;&amp; edge.Value.(*Edge).DestinationVertex == B {
			G.EdgeList.Remove(edge)
			for vtx := edge.Value.(*Edge).DestinationVertex.Predecessor.Front(); vtx != nil; vtx = vtx.Next() {
				if vtx.Value.(*Vertex) == A {
					edge.Value.(*Edge).DestinationVertex.Predecessor.Remove(vtx)
				}
			}
		}
	}
}

func (G *Graph) GetVertexSize() int {
	return G.VertexList.Len()
}

func (G *Graph) GetEdgeSize() int {
	return G.EdgeList.Len()
}
