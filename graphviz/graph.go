package graphviz

import (
	"fmt"
	"github.com/awalterschulze/gographviz"
	"github.com/awalterschulze/gographviz/parser"
)

const DiGraphTemplate = "digraph %s { }"
const GraphID = "G"

type Graph struct {
	graph *gographviz.Graph
}

func (g Graph) AddNode(name string, attrs map[string]string) {
	g.graph.AddNode(GraphID, name, attrs)
}

func (g Graph) AddEdge(src, srcPort, dst, dstPort string, attrs map[string]string) {
	g.graph.AddPortEdge(src, srcPort, dst, dstPort, true, attrs)
}

func (g Graph) String() string {
	return g.graph.String()
}

func NewGraph() *Graph {
	graphAst, _ := parser.ParseString(fmt.Sprintf(DiGraphTemplate, GraphID))
	graph := gographviz.NewGraph()
	gographviz.Analyse(graphAst, graph)

	return &Graph{graph: graph}
}
