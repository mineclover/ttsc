package graph

import (
  "strings"

  shimast "github.com/microsoft/typescript-go/shim/ast"

  "github.com/samchon/ttsc/packages/ttsc/driver"
)

// Build walks the program's user-authored source files and records a node for
// each top-level declaration. driver.SourceFiles already drops declaration files
// and the program never compiles a dependency's `.ts`, so every node Build emits
// is workspace source: External is false. External boundary leaves enter the
// graph only as the resolved target of an edge (see Resolve).
func Build(prog *driver.Program) *Graph {
  g := &Graph{Nodes: map[string]*Node{}}
  for _, file := range prog.SourceFiles() {
    collectDeclarations(g, file)
  }
  g.addEdges(prog)
  return g
}

// collectDeclarations records a node for each top-level declaration statement in
// file, plus a method node for each callable member of a class or interface, so
// method-to-method calls have both endpoints. This pass establishes the symbol
// nodes that cross-file edges connect.
func collectDeclarations(g *Graph, file *shimast.SourceFile) {
  if file.Statements == nil {
    return
  }
  path := file.FileName()
  for _, statement := range file.Statements.Nodes {
    switch statement.Kind {
    case shimast.KindFunctionDeclaration:
      addNode(g, path, statement, NodeFunction)
    case shimast.KindClassDeclaration:
      addNode(g, path, statement, NodeClass)
      collectMethods(g, path, statement)
    case shimast.KindInterfaceDeclaration:
      addNode(g, path, statement, NodeInterface)
      collectMethods(g, path, statement)
    case shimast.KindTypeAliasDeclaration:
      addNode(g, path, statement, NodeTypeAlias)
    case shimast.KindEnumDeclaration:
      addNode(g, path, statement, NodeEnum)
    case shimast.KindVariableStatement:
      collectVariables(g, path, statement)
    }
  }
}

// collectVariables records a variable node for each binding in a top-level
// variable statement (both bindings of `const a = 1, b = 2`).
func collectVariables(g *Graph, path string, statement *shimast.Node) {
  variables := statement.AsVariableStatement()
  if variables == nil || variables.DeclarationList == nil {
    return
  }
  list := variables.DeclarationList.AsVariableDeclarationList()
  if list == nil || list.Declarations == nil {
    return
  }
  for _, binding := range list.Declarations.Nodes {
    addNode(g, path, binding, NodeVariable)
  }
}

// addNode records a node for the symbol declared by node under its
// position-invariant id. A declaration the checker did not bind to a single
// named symbol (a destructuring pattern) is skipped, and a redeclaration (a
// merged interface, an overload set) keeps the first node.
func addNode(g *Graph, path string, node *shimast.Node, kind NodeKind) {
  symbol := node.Symbol()
  if symbol == nil || symbol.Name == "" {
    return
  }
  id := nodeID(path, symbol.Name, kind)
  if _, exists := g.Nodes[id]; exists {
    return
  }
  g.Nodes[id] = &Node{
    ID:   id,
    Name: symbol.Name,
    Kind: kind,
    File: path,
    Pos:  node.Pos(),
    End:  node.End(),
  }
}

// collectMethods records a method node for each callable member (method,
// constructor, accessor) of a class or interface declaration, keyed by its
// class-qualified name so a resolved method call lands on the same node.
func collectMethods(g *Graph, path string, statement *shimast.Node) {
  for _, member := range classMembers(statement) {
    if !isMethodMember(member.Kind) {
      continue
    }
    name := methodName(member.Symbol())
    if name == "" {
      continue
    }
    id := nodeID(path, name, NodeMethod)
    if _, exists := g.Nodes[id]; exists {
      continue
    }
    g.Nodes[id] = &Node{
      ID:   id,
      Name: name,
      Kind: NodeMethod,
      File: path,
      Pos:  member.Pos(),
      End:  member.End(),
    }
  }
}

// classMembers returns the member nodes of a class or interface declaration, or
// nil for anything else.
func classMembers(statement *shimast.Node) []*shimast.Node {
  switch statement.Kind {
  case shimast.KindClassDeclaration:
    if decl := statement.AsClassDeclaration(); decl != nil && decl.Members != nil {
      return decl.Members.Nodes
    }
  case shimast.KindInterfaceDeclaration:
    if decl := statement.AsInterfaceDeclaration(); decl != nil && decl.Members != nil {
      return decl.Members.Nodes
    }
  }
  return nil
}

// isMethodMember reports whether a class/interface member kind is a callable the
// graph models as a method node.
func isMethodMember(kind shimast.Kind) bool {
  switch kind {
  case shimast.KindMethodDeclaration, shimast.KindMethodSignature,
    shimast.KindConstructor, shimast.KindGetAccessor, shimast.KindSetAccessor:
    return true
  default:
    return false
  }
}

// methodName returns the class-qualified, printable name of a method symbol
// ("Class.method"), or "" when it has no named parent (a synthesized member).
// symbol.Parent is the class/interface symbol, set by the binder for every
// member. The internal-name prefix on a constructor (\xFE) is escaped to "__".
func methodName(symbol *shimast.Symbol) string {
  if symbol == nil || symbol.Name == "" || symbol.Parent == nil || symbol.Parent.Name == "" {
    return ""
  }
  return symbol.Parent.Name + "." + strings.ReplaceAll(symbol.Name, "\xFE", "__")
}
