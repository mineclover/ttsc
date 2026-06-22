package mcp

import (
  "encoding/json"
  "fmt"
  "sort"
  "strings"

  "github.com/samchon/ttsc/packages/ttsc/internal/graph"
)

// toolsListResult advertises the tool surface. Following codegraph's hard-won
// lesson, graph_explore is the fat, agent-facing default that answers a
// structural question in one round-trip; graph_diagnostics is the focused
// "what's wrong with this file" tool.
func toolsListResult() any {
  return map[string]any{
    "tools": []any{
      map[string]any{
        "name":        "graph_explore",
        "description": "Explore the checker-resolved code graph around a symbol or file: returns the matching nodes and their relationships (what they extend/implement and what derives from them). Start here for structural questions before reading files.",
        "inputSchema": map[string]any{
          "type": "object",
          "properties": map[string]any{
            "query": map[string]any{
              "type":        "string",
              "description": "A symbol name (e.g. \"MyClass\") or a file path fragment (e.g. \"src/service\").",
            },
          },
          "required": []any{"query"},
        },
      },
      map[string]any{
        "name":        "graph_diagnostics",
        "description": "Return the tsc semantic diagnostics for one file, in the same code and location tsgo reports.",
        "inputSchema": map[string]any{
          "type": "object",
          "properties": map[string]any{
            "file": map[string]any{
              "type":        "string",
              "description": "An absolute path or a trailing fragment of a project source file (e.g. \"src/main.ts\").",
            },
          },
          "required": []any{"file"},
        },
      },
    },
  }
}

// callTool routes a tools/call request to the named tool.
func (s *Server) callTool(params json.RawMessage) (any, *rpcError) {
  var call struct {
    Name      string          `json:"name"`
    Arguments json.RawMessage `json:"arguments"`
  }
  if err := json.Unmarshal(params, &call); err != nil {
    return nil, &rpcError{Code: codeInvalidParams, Message: "invalid tools/call params"}
  }
  switch call.Name {
  case "graph_explore":
    return s.explore(call.Arguments)
  case "graph_diagnostics":
    return s.diagnostics(call.Arguments)
  default:
    return nil, &rpcError{Code: codeInvalidParams, Message: "unknown tool: " + call.Name}
  }
}

// textResult wraps plain text in the MCP tools/call content envelope.
func textResult(text string) any {
  return map[string]any{
    "content": []any{map[string]any{"type": "text", "text": text}},
  }
}

// explore returns the nodes matching a query and their checker-resolved
// relationships. v1 renders the relationship map; verbatim source and blast
// radius are layered on later.
func (s *Server) explore(args json.RawMessage) (any, *rpcError) {
  var in struct {
    Query string `json:"query"`
  }
  if err := json.Unmarshal(args, &in); err != nil || strings.TrimSpace(in.Query) == "" {
    return nil, &rpcError{Code: codeInvalidParams, Message: "graph_explore requires a non-empty 'query'"}
  }
  matches := s.matchNodes(in.Query)
  if len(matches) == 0 {
    return textResult(fmt.Sprintf("No graph nodes match %q.", in.Query)), nil
  }
  var b strings.Builder
  for _, node := range matches {
    s.writeNodeRelations(&b, node)
  }
  return textResult(strings.TrimRight(b.String(), "\n")), nil
}

// matchNodes returns the nodes whose name equals the query or whose file path
// contains it, sorted by id for a stable response.
func (s *Server) matchNodes(query string) []*graph.Node {
  matches := make([]*graph.Node, 0)
  for _, node := range s.graph.Nodes {
    if node.Name == query || strings.Contains(node.File, query) {
      matches = append(matches, node)
    }
  }
  sort.Slice(matches, func(i, j int) bool { return matches[i].ID < matches[j].ID })
  return matches
}

// writeNodeRelations renders one node and its outgoing/incoming edges.
func (s *Server) writeNodeRelations(b *strings.Builder, node *graph.Node) {
  external := ""
  if node.External {
    external = " (external)"
  }
  fmt.Fprintf(b, "%s %s%s\n  %s\n", node.Kind, node.Name, external, node.File)
  for _, edge := range s.graph.Edges {
    if edge.From == node.ID {
      if to := s.graph.Nodes[edge.To]; to != nil {
        fmt.Fprintf(b, "  -> %s %s (%s)\n", to.Kind, to.Name, edge.Kind)
      }
    }
    if edge.To == node.ID {
      if from := s.graph.Nodes[edge.From]; from != nil {
        fmt.Fprintf(b, "  <- %s %s (%s)\n", from.Kind, from.Name, edge.Kind)
      }
    }
  }
  b.WriteString("\n")
}

// diagnostics returns a file's tsc semantic diagnostics as text.
func (s *Server) diagnostics(args json.RawMessage) (any, *rpcError) {
  var in struct {
    File string `json:"file"`
  }
  if err := json.Unmarshal(args, &in); err != nil || strings.TrimSpace(in.File) == "" {
    return nil, &rpcError{Code: codeInvalidParams, Message: "graph_diagnostics requires a non-empty 'file'"}
  }
  path, ok := s.resolveFile(in.File)
  if !ok {
    return textResult(fmt.Sprintf("No project source file matches %q.", in.File)), nil
  }
  found := graph.FileDiagnostics(s.prog, path)
  if len(found) == 0 {
    return textResult(fmt.Sprintf("No diagnostics for %s.", path)), nil
  }
  var b strings.Builder
  for _, d := range found {
    fmt.Fprintf(&b, "%s:%d:%d TS%d %s\n", d.File, d.Line, d.Column, d.Code, d.Message)
  }
  return textResult(strings.TrimRight(b.String(), "\n")), nil
}

// resolveFile maps a tool's file argument to a program source-file path: an exact
// match if present, otherwise the unique source file whose path ends with the
// argument. Returns ("", false) when nothing matches.
func (s *Server) resolveFile(file string) (string, bool) {
  for _, source := range s.prog.SourceFiles() {
    if source.FileName() == file {
      return file, true
    }
  }
  for _, source := range s.prog.SourceFiles() {
    if strings.HasSuffix(source.FileName(), file) {
      return source.FileName(), true
    }
  }
  return "", false
}
