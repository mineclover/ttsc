package graph

import (
  "strings"
  "unicode"

  shimast "github.com/microsoft/typescript-go/shim/ast"
)

const (
  annotationSourceJSDoc  = "jsdoc"
  graphSemanticTagName   = "graphSemantic"
  graphSemanticNamespace = "graph"
)

type declarationAnnotationSpec struct {
  tagName   string
  namespace string
  values    func(string) []string
}

var graphSemanticAnnotationSpec = declarationAnnotationSpec{
  tagName:   graphSemanticTagName,
  namespace: graphSemanticNamespace,
  values:    graphSemanticValues,
}

// declarationAnnotations extracts declaration-level annotations for graph
// nodes. The graph stores annotations as producer facts; consumers decide which
// namespaces or tag names carry domain meaning.
func declarationAnnotations(declaration *shimast.Node) []Annotation {
  return collectDeclarationAnnotations(declaration, graphSemanticAnnotationSpec)
}

func collectDeclarationAnnotations(declaration *shimast.Node, specs ...declarationAnnotationSpec) []Annotation {
  if declaration == nil || len(specs) == 0 {
    return nil
  }
  byName := make(map[string]declarationAnnotationSpec, len(specs))
  for _, spec := range specs {
    if spec.tagName == "" {
      continue
    }
    byName[spec.tagName] = spec
  }
  if len(byName) == 0 {
    return nil
  }

  annotations := make([]Annotation, 0)
  for _, owner := range declarationAnnotationOwners(declaration) {
    for _, jsdoc := range owner.JSDoc(nil) {
      data := jsdoc.AsJSDoc()
      if data == nil || data.Tags == nil {
        continue
      }
      for _, tag := range data.Tags.Nodes {
        annotation := jsdocDeclarationAnnotation(tag, byName)
        if annotation == nil {
          continue
        }
        annotations = append(annotations, *annotation)
      }
    }
  }
  if len(annotations) == 0 {
    return nil
  }
  return annotations
}

// declarationAnnotationOwners returns the AST nodes whose JSDoc belongs to this
// graph declaration. Variable JSDoc is attached to the enclosing
// VariableStatement in TypeScript, so variable declarations get that narrow
// parent fallback. Other declarations use only their own JSDoc to avoid applying
// class-level annotations to member nodes.
func declarationAnnotationOwners(declaration *shimast.Node) []*shimast.Node {
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

func jsdocDeclarationAnnotation(tag *shimast.Node, specs map[string]declarationAnnotationSpec) *Annotation {
  name, rest, ok := splitJSDocTagText(shimast.NodeText(tag))
  if !ok {
    return nil
  }
  spec, ok := specs[name]
  if !ok {
    return nil
  }
  values := []string{}
  if spec.values != nil {
    values = spec.values(rest)
  }
  if len(values) == 0 {
    return nil
  }
  return &Annotation{
    Source:    annotationSourceJSDoc,
    Name:      name,
    Namespace: spec.namespace,
    Values:    values,
    Pos:       tag.Pos(),
    End:       tag.End(),
  }
}

func splitJSDocTagText(raw string) (name, rest string, ok bool) {
  text := strings.TrimSpace(raw)
  if text == "" || text[0] != '@' {
    return "", "", false
  }
  body := text[1:]
  end := len(body)
  for i, r := range body {
    if unicode.IsSpace(r) {
      end = i
      break
    }
  }
  if end == 0 {
    return "", "", false
  }
  return body[:end], strings.TrimSpace(body[end:]), true
}

func graphSemanticValues(raw string) []string {
  seen := map[string]struct{}{}
  values := make([]string, 0)
  text := strings.TrimSpace(strings.NewReplacer(",", " ", ";", " ").Replace(raw))
  for _, value := range strings.Fields(text) {
    if value == "" || !isGraphSemanticValue(value) {
      continue
    }
    if _, ok := seen[value]; ok {
      continue
    }
    seen[value] = struct{}{}
    values = append(values, value)
  }
  return values
}

func semanticTagsFromAnnotations(annotations []Annotation) []string {
  seen := map[string]struct{}{}
  tags := make([]string, 0)
  for _, annotation := range annotations {
    if annotation.Source != annotationSourceJSDoc || annotation.Name != graphSemanticTagName {
      continue
    }
    for _, tag := range annotation.Values {
      if _, ok := seen[tag]; ok {
        continue
      }
      seen[tag] = struct{}{}
      tags = append(tags, tag)
    }
  }
  if len(tags) == 0 {
    return nil
  }
  return tags
}

func isGraphSemanticValue(value string) bool {
  for _, r := range value {
    if unicode.IsLetter(r) || unicode.IsDigit(r) {
      continue
    }
    switch r {
    case '-', '_', ':', '/':
      continue
    default:
      return false
    }
  }
  return true
}
