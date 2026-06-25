package mcp

import (
	"testing"

	"github.com/samchon/ttsc/packages/ttsc/internal/graph"
)

// TestFlowSourceWindowsIncludeQueryRelevantEdges verifies flow output keeps
// checker-resolved value-use evidence as structured graph edges.
//
// This locks the post-text-renderer contract: flow is an index, so it should
// return path edges and query-relevant off-path edges as typed records instead of
// source-window text.
//
//  1. Build a small graph with one selected path edge and one query-relevant
//     off-path edge.
//  2. Ask for a flow result seeded by the path nodes.
//  3. Assert both edges are present with their `onPath` state.
func TestFlowSourceWindowsIncludeQueryRelevantEdges(t *testing.T) {
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
	if !hasFlowEdge(result, "Router.route", "Pipeline.normalize", true) {
		t.Fatalf("flow result omitted selected path edge: %#v", result.Evidence)
	}
	if !hasFlowEdge(result, "Router.route", "AliasFactory.aliasFactory", false) {
		t.Fatalf("flow result omitted query-relevant off-path edge: %#v", result.Evidence)
	}
}

func hasFlowEdge(result *flowResult, from string, to string, onPath bool) bool {
	for _, edge := range result.Evidence {
		if edge.From.Name == from && edge.To.Name == to && edge.OnPath == onPath {
			return true
		}
	}
	return false
}
