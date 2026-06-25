package mcp_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/samchon/ttsc/packages/ttsc/driver"
	"github.com/samchon/ttsc/packages/ttsc/internal/graph/mcp"
)

// TestExploreIncludesLateCallExcerpts verifies query_nodes keeps checker-known
// value-use coordinates visible without returning the long declaration body.
//
// Code-flow questions often need the late call sites that explain why returned
// edges matter. The index exposes those as edge use locations; source stays
// behind exact expand_nodes handles.
//
//  1. Compile a function whose call to helper appears after the source-line cap.
//  2. Explore the caller.
//  3. Assert the body is omitted and the late helper call edge is indexed.
func TestExploreIncludesLateCallExcerpts(t *testing.T) {
	root := t.TempDir()
	var src strings.Builder
	src.WriteString("export function lateFlow(): number {\n")
	for i := 0; i < 40; i++ {
		fmt.Fprintf(&src, "  const value%d = %d;\n", i, i)
	}
	src.WriteString("  return helper()\n")
	src.WriteString("}\n")
	src.WriteString("export function helper(): number {\n")
	src.WriteString("  return 1\n")
	src.WriteString("}\n")

	writeFile(t, filepath.Join(root, "tsconfig.json"), `{
  "compilerOptions": {
    "target": "ES2022",
    "module": "commonjs",
    "strict": true
  },
  "files": ["src/main.ts"]
}
`)
	writeFile(t, filepath.Join(root, "src", "main.ts"), src.String())

	prog, diags, err := driver.LoadProgram(root, "tsconfig.json", driver.LoadProgramOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(diags) != 0 {
		t.Fatalf("unexpected parse diagnostics: %v", diags)
	}
	defer func() { _ = prog.Close() }()

	server := mcp.NewServer(prog)
	text := toolText(t, server, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"query_nodes","arguments":{"query":"lateFlow"}}}`)
	if strings.Contains(text, "const value39") || strings.Contains(text, "return helper()") {
		t.Fatalf("query_nodes included source body instead of index coordinates:\n%s", text)
	}
	if !strings.Contains(text, `"kind": "value-call"`) ||
		!strings.Contains(text, `"name": "helper"`) ||
		!strings.Contains(text, `"line": 42`) {
		t.Fatalf("query_nodes did not include the late value-call edge location:\n%s", text)
	}
}

// TestExploreUsesEdgeSpansForLateAccessExcerpts verifies late value-access
// evidence is sourced from graph edge spans, not from call-looking text search.
// A property read like `this.value` has no trailing `(`, but it is still the
// resolved use an agent needs when a long method body is truncated.
func TestExploreUsesEdgeSpansForLateAccessExcerpts(t *testing.T) {
	root := t.TempDir()
	var src strings.Builder
	src.WriteString("export class Store {\n")
	src.WriteString("  value = 1;\n")
	src.WriteString("  read(): number {\n")
	for i := 0; i < 40; i++ {
		fmt.Fprintf(&src, "    const filler%d = %d;\n", i, i)
	}
	src.WriteString("    return this.value;\n")
	src.WriteString("  }\n")
	src.WriteString("}\n")

	writeFile(t, filepath.Join(root, "tsconfig.json"), `{
  "compilerOptions": {
    "target": "ES2022",
    "module": "commonjs",
    "strict": true
  },
  "files": ["src/main.ts"]
}
`)
	writeFile(t, filepath.Join(root, "src", "main.ts"), src.String())

	prog, diags, err := driver.LoadProgram(root, "tsconfig.json", driver.LoadProgramOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(diags) != 0 {
		t.Fatalf("unexpected parse diagnostics: %v", diags)
	}
	defer func() { _ = prog.Close() }()

	server := mcp.NewServer(prog)
	text := toolText(t, server, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"query_nodes","arguments":{"query":"Store.read"}}}`)
	if strings.Contains(text, "const filler39") || strings.Contains(text, "return this.value;") {
		t.Fatalf("query_nodes included source body instead of index coordinates:\n%s", text)
	}
	if !strings.Contains(text, `"kind": "value-access"`) ||
		!strings.Contains(text, `"name": "Store.value"`) ||
		!strings.Contains(text, `"line": 44`) {
		t.Fatalf("query_nodes did not include the late value-access edge location:\n%s", text)
	}
}
