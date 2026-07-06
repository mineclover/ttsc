package graph

import (
  "encoding/json"
  "path/filepath"
  "reflect"
  "testing"

  "github.com/samchon/ttsc/packages/ttsc/driver"
)

// TestGraphDeclarationAnnotationsEmitOnNodes verifies that declaration-level
// `@graphSemantic` JSDoc tags become neutral node annotations, while the legacy
// semanticTags field stays as a compatibility projection. It also pins the
// non-inheritance rule: class-level tags do not bleed into untagged methods.
func TestGraphDeclarationAnnotationsEmitOnNodes(t *testing.T) {
  root := t.TempDir()
  writeFile(t, filepath.Join(root, "tsconfig.json"), fixtureTSConfig)
  writeFile(t, filepath.Join(root, "src", "main.ts"), `/**
 * @graphSemantic shared-contract route-contract
 * @graphSemantic shared-contract invalid.value
 */
export interface SharedErrorResponse {}

/** @graphSemantic poster-route-contract */
export type PosterRouteSpec = { id: string };

/** @graphSemantic invalid.value */
export type DottedSemantic = string;

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

  shared := graph.Nodes[nodeID(path, "SharedErrorResponse", NodeInterface)]
  assertNodeSemanticTags(
    t,
    shared,
    []string{"shared-contract", "route-contract"},
  )
  assertNodeAnnotations(t, shared, [][]string{
    {"shared-contract", "route-contract"},
    {"shared-contract"},
  })
  assertNodeSemanticTags(
    t,
    graph.Nodes[nodeID(path, "PosterRouteSpec", NodeTypeAlias)],
    []string{"poster-route-contract"},
  )
  assertNodeSemanticTags(t, graph.Nodes[nodeID(path, "DottedSemantic", NodeTypeAlias)], nil)
  assertNodeAnnotations(t, graph.Nodes[nodeID(path, "DottedSemantic", NodeTypeAlias)], nil)
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
  annotations := dumped["src/main.ts#SharedErrorResponse:interface"].Annotations
  if len(annotations) != 2 {
    t.Fatalf("dumped interface annotations = %+v", annotations)
  }
  if got := annotations[0]; got.Source != annotationSourceJSDoc || got.Name != graphSemanticTagName || got.Namespace != graphSemanticNamespace || !reflect.DeepEqual(got.Values, []string{"shared-contract", "route-contract"}) || got.Evidence == nil {
    t.Fatalf("first dumped annotation = %+v", got)
  }
  assertAnnotationEvidence(t, annotations[0], 2, 4)
  if got := annotations[1]; got.Source != annotationSourceJSDoc || got.Name != graphSemanticTagName || got.Namespace != graphSemanticNamespace || !reflect.DeepEqual(got.Values, []string{"shared-contract"}) || got.Evidence == nil {
    t.Fatalf("second dumped annotation = %+v", got)
  }
  assertAnnotationEvidence(t, annotations[1], 3, 4)
  if got := dumped["src/main.ts#DottedSemantic:type"].SemanticTags; got != nil {
    t.Fatalf("dotted semantic value was emitted as semanticTags: %v", got)
  }
  if got := dumped["src/main.ts#DottedSemantic:type"].Annotations; got != nil {
    t.Fatalf("dotted semantic value was emitted as annotations: %+v", got)
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

func assertNodeAnnotations(t *testing.T, node *Node, want [][]string) {
  t.Helper()
  if node == nil {
    t.Fatalf("missing graph node")
  }
  if len(want) == 0 {
    if node.Annotations != nil {
      t.Fatalf("node %s annotations = %+v, want nil", node.ID, node.Annotations)
    }
    return
  }
  if len(node.Annotations) != len(want) {
    t.Fatalf("node %s annotations = %+v, want %v", node.ID, node.Annotations, want)
  }
  for i, values := range want {
    annotation := node.Annotations[i]
    if annotation.Source != annotationSourceJSDoc ||
      annotation.Name != graphSemanticTagName ||
      annotation.Namespace != graphSemanticNamespace ||
      !reflect.DeepEqual(annotation.Values, values) ||
      annotation.Pos <= 0 ||
      annotation.End <= annotation.Pos {
      t.Fatalf("node %s annotation[%d] = %+v, want values %v", node.ID, i, annotation, values)
    }
  }
}

func assertAnnotationEvidence(t *testing.T, annotation DumpAnnotation, wantLine, wantCol int) {
  t.Helper()
  if annotation.Evidence == nil {
    t.Fatalf("annotation evidence is nil: %+v", annotation)
  }
  if annotation.Evidence.File != "src/main.ts" ||
    annotation.Evidence.StartLine != wantLine ||
    annotation.Evidence.StartCol != wantCol {
    t.Fatalf(
      "annotation evidence = %+v, want file src/main.ts line %d col %d",
      annotation.Evidence,
      wantLine,
      wantCol,
    )
  }
}
