package mcp_test

import (
	"path/filepath"
	"testing"

	"github.com/samchon/ttsc/packages/ttsc/driver"
	"github.com/samchon/ttsc/packages/ttsc/internal/graph/mcp"
)

// TestQueryFilesOutlinesEachFile verifies query_files returns one structured
// location record per requested location, in input order, each a cheap roster of
// that file without verbatim bodies.
//
//  1. Compile a two-file fixture, each with a couple of declarations.
//  2. Call query_files with both paths.
//  3. Assert two location records, in order, each naming its file's declarations
//     but not dumping their bodies.
func TestQueryFilesOutlinesEachFile(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "tsconfig.json"), `{
  "compilerOptions": { "target": "ES2022", "module": "commonjs", "strict": true, "rootDir": "src", "outDir": "dist" },
  "files": ["src/a.ts", "src/b.ts"]
}
`)
	writeFile(t, filepath.Join(root, "src", "a.ts"), `export class Alpha {
  ping(): number {
    return 1
  }
}
`)
	writeFile(t, filepath.Join(root, "src", "b.ts"), `export function beta(): string {
  return "b"
}
`)

	prog, _, err := driver.LoadProgram(root, "tsconfig.json", driver.LoadProgramOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = prog.Close() }()
	server := mcp.NewServer(prog)

	result := toolStructured(t, server, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"query_files","arguments":{"locations":["src/b.ts","src/a.ts"]}}}`)
	locations, ok := result["locations"].([]any)
	if !ok || len(locations) != 2 {
		t.Fatalf("expected one location record per input, got %v", result["locations"])
	}
	// Input order is preserved: b.ts first, a.ts second.
	if !structuredContains(locations[0], "src/b.ts") || !structuredContains(locations[0], "beta") {
		t.Fatalf("first location did not index src/b.ts:\n%v", locations[0])
	}
	if !structuredContains(locations[1], "src/a.ts") || !structuredContains(locations[1], "Alpha") {
		t.Fatalf("second location did not index src/a.ts:\n%v", locations[1])
	}
	// The roster lists declarations but not their bodies: it is a cheap index, so
	// the verbatim source (here a `return` statement) must not appear.
	if structuredContains(result, "return") {
		t.Fatalf("query_files dumped a body instead of a roster:\n%v", result)
	}
}
