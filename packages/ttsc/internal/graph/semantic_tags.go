package graph

import (
  "strings"
  "unicode"

  shimast "github.com/microsoft/typescript-go/shim/ast"
)

const graphSemanticTag = "@graphSemantic"

// declarationSemanticTags extracts graph-specific semantic labels from
// declaration-level JSDoc. The graph deliberately reads only producer-owned
// annotations here; consumers should use the emitted node field instead of
// reparsing source comments.
func declarationSemanticTags(declaration *shimast.Node) []string {
  seen := map[string]struct{}{}
  tags := make([]string, 0)
  for _, owner := range semanticTagOwners(declaration) {
    for _, jsdoc := range owner.JSDoc(nil) {
      data := jsdoc.AsJSDoc()
      if data == nil || data.Tags == nil {
        continue
      }
      for _, tag := range data.Tags.Nodes {
        appendGraphSemanticTags(&tags, seen, shimast.NodeText(tag))
      }
    }
  }
  if len(tags) == 0 {
    return nil
  }
  return tags
}

// semanticTagOwners returns the AST nodes whose JSDoc belongs to this graph
// declaration. Variable JSDoc is attached to the enclosing VariableStatement in
// TypeScript, so variable declarations get that narrow parent fallback. Other
// declarations use only their own JSDoc to avoid applying class-level tags to
// member nodes.
func semanticTagOwners(declaration *shimast.Node) []*shimast.Node {
  if declaration == nil {
    return nil
  }
  owners := []*shimast.Node{declaration}
  if declaration.Kind != shimast.KindVariableDeclaration {
    return owners
  }
  for node := declaration.Parent; node != nil; node = node.Parent {
    if node.Kind == shimast.KindVariableStatement {
      return append(owners, node)
    }
    if node.Kind == shimast.KindSourceFile {
      break
    }
  }
  return owners
}

func appendGraphSemanticTags(tags *[]string, seen map[string]struct{}, raw string) {
  text := strings.TrimSpace(raw)
  if !strings.HasPrefix(text, graphSemanticTag) {
    return
  }
  rest := text[len(graphSemanticTag):]
  if rest != "" && !isGraphSemanticTagSeparator(rest[0]) {
    return
  }
  rest = strings.TrimSpace(strings.NewReplacer(",", " ", ";", " ").Replace(rest))
  for _, tag := range strings.Fields(rest) {
    if tag == "" || !isGraphSemanticValue(tag) {
      continue
    }
    if _, ok := seen[tag]; ok {
      continue
    }
    seen[tag] = struct{}{}
    *tags = append(*tags, tag)
  }
}

func isGraphSemanticTagSeparator(b byte) bool {
  return b == ' ' || b == '\t' || b == '\r' || b == '\n'
}

func isGraphSemanticValue(value string) bool {
  for _, r := range value {
    if unicode.IsLetter(r) || unicode.IsDigit(r) {
      continue
    }
    switch r {
    case '-', '_', ':', '.', '/':
      continue
    default:
      return false
    }
  }
  return true
}
