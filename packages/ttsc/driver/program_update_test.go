package driver

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	shimtspath "github.com/microsoft/typescript-go/shim/tspath"
)

// TestUpdateProgramReusesProgramOnBodyEdit verifies that typescript-go's
// incremental program reuse (compiler.Program.UpdateProgram) is reachable
// through ttsc's shim and actually reuses the loaded program when only one
// file's body changed.
//
// This is the foundation primitive for a resident / incremental ttsc host
// (samchon/ttsc#255): @ttsc/metro and @ttsc/unplugin must update a loaded
// program per file instead of recompiling the whole project. UpdateProgram
// re-parses only the changed file and reuses every unchanged file's AST,
// returning reused=true when the changed file's import/reference graph is
// unchanged. ttsc pins a single checker (forceSingleChecker); the reused
// program options preserve that, so the transform-walk invariant survives.
//
// 1. Load a two-file project (a imports b).
// 2. Rewrite b's function body on disk, leaving its signature and imports intact.
// 3. Call Program.UpdateProgram for b; assert it reused the program and the new
//    body is visible in the updated program.
func TestUpdateProgramReusesProgramOnBodyEdit(t *testing.T) {
	dir := t.TempDir()
	write := func(name, content string) string {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
		return p
	}
	write("tsconfig.json", `{"compilerOptions":{"strict":true,"noEmit":true},"files":["a.ts","b.ts"]}`)
	write("a.ts", "import { b } from \"./b\";\nexport const a: number = b();\n")
	bPath := write("b.ts", "export function b(): number {\n\treturn 1;\n}\n")

	prog, diags, err := LoadProgram(dir, "tsconfig.json", LoadProgramOptions{ForceNoEmit: true})
	if err != nil {
		t.Fatalf("LoadProgram: %v", err)
	}
	if prog == nil {
		t.Fatalf("LoadProgram returned nil program (diagnostics: %v)", diags)
	}
	defer prog.Close()

	bFile := prog.SourceFile(bPath)
	if bFile == nil {
		t.Fatalf("b.ts is not in the loaded program source files")
	}

	// Graph-preserving edit: same signature and exports, different body.
	if err := os.WriteFile(bPath, []byte("export function b(): number {\n\treturn 2;\n}\n"), 0o644); err != nil {
		t.Fatalf("rewrite b.ts: %v", err)
	}

	changed := shimtspath.ToPath(bFile.FileName(), dir, prog.FS.UseCaseSensitiveFileNames())
	newHost := DefaultHost(dir, DefaultFS())
	newProg, reused := prog.TSProgram.UpdateProgram(changed, newHost, nil)
	if newProg == nil {
		t.Fatalf("UpdateProgram returned a nil program")
	}
	if !reused {
		t.Fatalf("expected UpdateProgram to reuse the program for a body-only edit, got reused=false")
	}

	var updated string
	for _, file := range newProg.SourceFiles() {
		if file.FileName() == bFile.FileName() {
			updated = file.Text()
			break
		}
	}
	if !strings.Contains(updated, "return 2") {
		t.Fatalf("updated program does not reflect the new body of b.ts; got %q", updated)
	}
}
