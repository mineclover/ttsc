package graph

import (
  "encoding/json"
  "path/filepath"
  "reflect"
  "testing"

  "github.com/samchon/ttsc/packages/ttsc/driver"
)

// TestGraphSemanticTagsEmitOnNodes verifies that declaration-level
// `@graphSemantic` JSDoc tags become node semanticTags in both the internal
// graph and the JSON dump. It also pins the non-inheritance rule: class-level
// tags do not bleed into untagged method nodes.
func TestGraphSemanticTagsEmitOnNodes(t *testing.T) {
  root := t.TempDir()
  writeFile(t, filepath.Join(root, "tsconfig.json"), fixtureTSConfig)
  writeFile(t, filepath.Join(root, "src", "main.ts"), `/**
 * @graphSemantic shared-contract route-contract
 * @graphSemantic shared-contract
 */
export interface SharedErrorResponse {}

/** @graphSemantic poster-route-contract */
export type PosterRouteSpec = { id: string };

/** @graphSemantic route-config */
export const posterRoute = { path: "/poster" };

/** @graphSemantic controller-boundary */
export class PosterController {
  /** @graphSemantic action-handler */
  add(): void {}
  remove(): void {}
}
`)

  prog, diags, err := driver.LoadProgram(root, "tsconfig.json", driver.LoadProgramOptions{})
  if err != nil {
    t.Fatal(err)
  }
  if len(diags) != 0 {
    t.Fatalf("unexpected diagnostics: %v", diags)
  }
  defer func() { _ = prog.Close() }()

  graph := Build(prog)
  path := sourceFile(t, prog, "main.ts").FileName()

  assertNodeSemanticTags(
    t,
    graph.Nodes[nodeID(path, "SharedErrorResponse", NodeInterface)],
    []string{"shared-contract", "route-contract"},
  )
  assertNodeSemanticTags(
    t,
    graph.Nodes[nodeID(path, "PosterRouteSpec", NodeTypeAlias)],
    []string{"poster-route-contract"},
  )
  assertNodeSemanticTags(
    t,
    graph.Nodes[nodeID(path, "posterRoute", NodeVariable)],
    []string{"route-config"},
  )
  assertNodeSemanticTags(
    t,
    graph.Nodes[nodeID(path, "PosterController", NodeClass)],
    []string{"controller-boundary"},
  )
  assertNodeSemanticTags(
    t,
    graph.Nodes[nodeID(path, "PosterController.add", NodeMethod)],
    []string{"action-handler"},
  )
  assertNodeSemanticTags(t, graph.Nodes[nodeID(path, "PosterController.remove", NodeMethod)], nil)

  data, err := MarshalDump(graph, root, "tsconfig.json", nil, SourceTexts(prog), false)
  if err != nil {
    t.Fatalf("MarshalDump: %v", err)
  }
  var dump Dump
  if err := json.Unmarshal(data, &dump); err != nil {
    t.Fatalf("dump is not valid JSON: %v\n%s", err, data)
  }
  dumped := map[string]DumpNode{}
  for _, node := range dump.Nodes {
    dumped[node.ID] = node
  }
  if got := dumped["src/main.ts#SharedErrorResponse:interface"].SemanticTags; !reflect.DeepEqual(got, []string{"shared-contract", "route-contract"}) {
    t.Fatalf("dumped interface semanticTags = %v", got)
  }
  if got := dumped["src/main.ts#PosterController.remove:method"].SemanticTags; got != nil {
    t.Fatalf("untagged method inherited semanticTags: %v", got)
  }
}

func assertNodeSemanticTags(t *testing.T, node *Node, want []string) {
  t.Helper()
  if node == nil {
    t.Fatalf("missing graph node")
  }
  if !reflect.DeepEqual(node.SemanticTags, want) {
    t.Fatalf("node %s semanticTags = %v, want %v", node.ID, node.SemanticTags, want)
  }
}
