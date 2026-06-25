package mcp

import (
	"testing"

	"github.com/samchon/ttsc/packages/ttsc/internal/graph"
)

// TestFlowResultUsesCompactPathEdges verifies flow output keeps checker-resolved
// value-use evidence as compact handle edges.
//
// This locks the typed index contract: flow answers the selected path, not every
// adjacent branch from a path node, and edge endpoints refer to the returned
// nodes by handle instead of repeating node coordinates.
//
//  1. Build a small graph with one selected path edge and one query-relevant
//     off-path edge.
//  2. Ask for a flow result seeded by the path nodes.
//  3. Assert only the selected edge is present as fromHandle/toHandle.
func TestFlowResultUsesCompactPathEdges(t *testing.T) {
	route := &graph.Node{ID: "route", Name: "Router.route", Kind: graph.NodeFunction, File: "src/main.ts"}
	normalize := &graph.Node{ID: "normalize", Name: "Pipeline.normalize", Kind: graph.NodeFunction, File: "src/main.ts"}
	alias := &graph.Node{ID: "alias", Name: "AliasFactory.aliasFactory", Kind: graph.NodeFunction, File: "src/main.ts"}
	server := &Server{
		graph: &graph.Graph{
			Nodes: map[string]*graph.Node{
				route.ID:     route,
				normalize.ID: normalize,
				alias.ID:     alias,
			},
			Edges: []*graph.Edge{
				{From: route.ID, To: normalize.ID, Kind: graph.EdgeValueCall},
				{From: route.ID, To: alias.ID, Kind: graph.EdgeValueCall},
			},
		},
	}

	result := server.flowResult([]*graph.Node{route, normalize}, "alias")
	if !hasFlowEdge(result, route.ID, normalize.ID) {
		t.Fatalf("flow result omitted selected path edge: %#v", result.Evidence)
	}
	if hasFlowEdge(result, route.ID, alias.ID) {
		t.Fatalf("flow result included off-path edge: %#v", result.Evidence)
	}
}

func hasFlowEdge(result *flowResult, fromID string, toID string) bool {
	for _, edge := range result.Evidence {
		if edge.FromHandle == nodeHandle(fromID) && edge.ToHandle == nodeHandle(toID) {
			return true
		}
	}
	return false
}
