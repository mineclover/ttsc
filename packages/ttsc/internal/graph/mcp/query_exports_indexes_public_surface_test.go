package mcp_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/samchon/ttsc/packages/ttsc/driver"
	"github.com/samchon/ttsc/packages/ttsc/internal/graph/mcp"
)

// TestQueryExportsIndexesCompilerExports verifies query_exports uses the
// compiler's module export surface, including re-exports, while returning only
// compact coordinates and handles.
//
//  1. Compile a small module with local exports and a barrel re-export.
//  2. Query the exported surface.
//  3. Assert the result lists public declarations with file, line, handle, and
//     compact page metadata.
func TestQueryExportsIndexesCompilerExports(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "tsconfig.json"), `{
  "compilerOptions": {
    "target": "ES2022",
    "module": "commonjs",
    "strict": true,
    "rootDir": "src",
    "outDir": "dist"
  },
  "files": ["src/index.ts", "src/model.ts", "src/hidden.ts"]
}
`)
	writeFile(t, filepath.Join(root, "src", "model.ts"), `export interface PublicModel {
  value: string
}

export class PublicService {
  status = "ready"

  run(model: PublicModel): string {
    return model.value
  }

  private secret(): string {
    return "hidden"
  }
}

class InternalOnly {}
`)
	writeFile(t, filepath.Join(root, "src", "hidden.ts"), `class PublicService {
  shadow(): string {
    return "not exported"
  }
}
`)
	writeFile(t, filepath.Join(root, "src", "index.ts"), `export { PublicService as Service } from "./model"
export type { PublicModel } from "./model"
export const localValue = 1
`)

	prog, diags, err := driver.LoadProgram(root, "tsconfig.json", driver.LoadProgramOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	defer func() { _ = prog.Close() }()

	server := mcp.NewServer(prog)
	result := toolStructured(t, server, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"query_exports","arguments":{}}}`)
	for _, want := range []string{
		"PublicService",
		"Service",
		"PublicService.run",
		"PublicService.status",
		"PublicModel",
		"localValue",
		"n:",
	} {
		if !structuredContains(result, want) {
			t.Fatalf("query_exports missing %q:\n%v", want, result)
		}
	}
	page, ok := result["page"].(map[string]any)
	if !ok || page["totalRecords"] != float64(6) || page["totalPages"] != float64(1) {
		t.Fatalf("query_exports returned wrong page metadata: %v", result["page"])
	}
	exports, ok := result["exports"].([]any)
	if !ok || len(exports) != 6 {
		t.Fatalf("query_exports returned wrong exports array: %v", result["exports"])
	}
	for _, raw := range exports {
		export, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("export record is not an object: %v", raw)
		}
		assertExportShape(t, export)
	}
	if structuredContains(result, "InternalOnly") {
		t.Fatalf("query_exports included a non-exported declaration:\n%v", result)
	}
	if structuredContains(result, "PublicService.secret") {
		t.Fatalf("query_exports included a private member:\n%v", result)
	}
	if structuredContains(result, "PublicService.shadow") {
		t.Fatalf("query_exports included a member from a non-exported same-name owner:\n%v", result)
	}

	filtered := toolStructured(t, server, `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"query_exports","arguments":{"query":"service","limit":1}}}`)
	exports, ok = filtered["exports"].([]any)
	if !ok || len(exports) != 1 || !structuredContains(filtered, "PublicService") || structuredContains(filtered, "PublicModel") {
		t.Fatalf("query_exports filter/limit did not narrow the result:\n%v", filtered)
	}
}

// TestQueryExportsOmitsGitIgnoredGeneratedCode pins the orientation rule that
// gitignored generated TypeScript is not part of the first project map.
func TestQueryExportsOmitsGitIgnoredGeneratedCode(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	writeFile(t, filepath.Join(root, ".gitignore"), "src/generated/\n")
	writeFile(t, filepath.Join(root, "tsconfig.json"), `{
  "compilerOptions": {
    "target": "ES2022",
    "module": "commonjs",
    "strict": true,
    "rootDir": "src",
    "outDir": "dist"
  },
  "files": ["src/app.ts", "src/generated/client.ts"]
}
`)
	writeFile(t, filepath.Join(root, "src", "generated", "client.ts"), `export class GeneratedClient {}
`)
	writeFile(t, filepath.Join(root, "src", "app.ts"), `export class AppService {}
`)

	server := mcp.NewLazyServer(root, "tsconfig.json", driver.LoadProgramOptions{})
	result := toolStructured(t, server, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"query_exports","arguments":{}}}`)
	if structuredContains(result, "GeneratedClient") {
		t.Fatalf("query_exports included gitignored generated code:\n%v", result)
	}
	if !structuredContains(result, "AppService") {
		t.Fatalf("query_exports omitted authored export:\n%v", result)
	}
}

func assertExportShape(t *testing.T, export map[string]any) {
	t.Helper()
	required := map[string]bool{
		"name":   true,
		"kind":   true,
		"file":   true,
		"line":   true,
		"handle": true,
	}
	allowed := map[string]bool{
		"name":       true,
		"exportedAs": true,
		"kind":       true,
		"file":       true,
		"line":       true,
		"handle":     true,
	}
	for key := range required {
		if _, ok := export[key]; !ok {
			t.Fatalf("export record missing %q: %v", key, export)
		}
	}
	for key := range export {
		if !allowed[key] {
			t.Fatalf("export record contains non-index field %q: %v", key, export)
		}
	}
	if handle, _ := export["handle"].(string); !strings.HasPrefix(handle, "n:") {
		t.Fatalf("export record handle is not stable graph handle: %v", export)
	}
}
