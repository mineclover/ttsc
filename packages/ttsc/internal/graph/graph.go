package graph

// NodeKind classifies a graph node by what its symbol declares.
type NodeKind string

const (
  NodeFunction  NodeKind = "function"
  NodeClass     NodeKind = "class"
  NodeInterface NodeKind = "interface"
  NodeTypeAlias NodeKind = "type"
  NodeEnum      NodeKind = "enum"
  NodeVariable  NodeKind = "variable"
)

// Provenance marks how a node or edge was derived. Every relationship in this
// graph is resolved by the in-process type checker, so the single value is a
// trust signal: the inverse of a tree-sitter tool tagging an uncertain edge
// "heuristic".
const Provenance = "checker-resolved"

// Node is one declared symbol. Its ID is position-invariant, built from the file
// realpath, the declared name, and the kind, so inserting a line above a
// declaration does not re-key it. That keeps a future incremental layer from
// churning the whole graph on every edit, which a byte-offset key would force.
type Node struct {
  ID       string
  Name     string
  Kind     NodeKind
  File     string
  External bool
}

// Graph is the in-memory adjacency the MCP tools query. Edges are added by the
// resolution pass on top of the declaration nodes Build records.
type Graph struct {
  Nodes map[string]*Node
}

// nodeID builds the position-invariant identity for a symbol named name,
// declared as kind in the source file at path.
func nodeID(path string, name string, kind NodeKind) string {
  return path + "#" + name + ":" + string(kind)
}
