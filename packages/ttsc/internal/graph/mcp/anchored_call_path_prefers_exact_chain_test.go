package mcp

import (
  "reflect"
  "testing"

  "github.com/samchon/ttsc/packages/ttsc/internal/graph"
)

// TestAnchoredCallPathPrefersExactChain verifies exact symbol anchors in a flow
// query are stitched through the graph before broad BFS expansion.
//
// When a user names the endpoints of a chain, the graph should return the actual
// value-flow path between those anchors rather than pulling in every sibling
// call from the first node.
//
//  1. Build a synthetic value-call graph with one on-chain branch and one
//     sibling branch from the public entry.
//  2. Ask for the exact public and terminal symbols.
//  3. Assert the selected flow follows the on-chain branch and omits the sibling.
func TestAnchoredCallPathPrefersExactChain(t *testing.T) {
  nodes := map[string]*graph.Node{
    "gateway":     {ID: "gateway", Name: "Gateway.fetch", Kind: graph.NodeMethod, File: "src/main.ts"},
    "coordinator": {ID: "coordinator", Name: "Coordinator.fetch", Kind: graph.NodeMethod, File: "src/main.ts"},
    "bridge":      {ID: "bridge", Name: "Coordinator.bridge", Kind: graph.NodeMethod, File: "src/main.ts"},
    "pipeline":    {ID: "pipeline", Name: "Pipeline.buildSteps", Kind: graph.NodeMethod, File: "src/main.ts"},
    "logger":      {ID: "logger", Name: "Logger.record", Kind: graph.NodeMethod, File: "src/main.ts"},
  }
  server := &Server{
    graph: &graph.Graph{Nodes: nodes},
    forwardCallAdj: map[string][]string{
      "gateway":     {"logger", "coordinator"},
      "coordinator": {"bridge"},
      "bridge":      {"pipeline"},
    },
  }
  seeds := []*graph.Node{nodes["gateway"], nodes["logger"], nodes["pipeline"]}
  gotNodes := server.withCallPath(seeds, maxPathNodes, "Trace Gateway.fetch to Pipeline.buildSteps")
  got := make([]string, 0, len(gotNodes))
  for _, node := range gotNodes {
    got = append(got, node.Name)
  }
  want := []string{
    "Gateway.fetch",
    "Coordinator.fetch",
    "Coordinator.bridge",
    "Pipeline.buildSteps",
  }
  if !reflect.DeepEqual(got, want) {
    t.Fatalf("withCallPath() = %#v; want %#v", got, want)
  }
}
