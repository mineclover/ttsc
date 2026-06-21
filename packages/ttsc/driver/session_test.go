package driver

import (
  "os"
  "path/filepath"
  "strings"
  "testing"
)

// TestSessionAppliesIncrementalOverlayEdit verifies the resident Session feeds
// an in-memory overlay edit to the live program through an incremental
// UpdateProgram, reflecting the new content without a full reload.
//
// This is the resident-host slice of samchon/ttsc#255: @ttsc/metro must update
// one file's content (Metro's src) against a warm program rather than recompile
// the project. The edit keeps the file's signature and imports, so the update
// must reuse the program and the resident source must show the new body.
//
// 1. Open a Session on a two-file project (a imports b).
// 2. Apply an overlay edit to b's body.
// 3. Assert the update reused the program and SourceText shows the new body.
func TestSessionAppliesIncrementalOverlayEdit(t *testing.T) {
  dir := t.TempDir()
  write := func(name, content string) {
    if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
      t.Fatalf("write %s: %v", name, err)
    }
  }
  write("tsconfig.json", `{"compilerOptions":{"strict":true,"noEmit":true},"files":["a.ts","b.ts"]}`)
  write("a.ts", "import { b } from \"./b\";\nexport const a: number = b();\n")
  write("b.ts", "export function b(): number {\n  return 1;\n}\n")

  sess, diags, err := NewSession(dir, "tsconfig.json", LoadProgramOptions{ForceNoEmit: true})
  if err != nil {
    t.Fatalf("NewSession: %v", err)
  }
  if sess == nil {
    t.Fatalf("NewSession returned nil session (diagnostics: %v)", diags)
  }
  defer sess.Close()

  bAbs := filepath.Join(dir, "b.ts")
  if text, ok := sess.SourceText(bAbs); !ok || !strings.Contains(text, "return 1") {
    t.Fatalf("initial source text mismatch: ok=%v text=%q", ok, text)
  }

  reused := sess.Apply(bAbs, "export function b(): number {\n  return 2;\n}\n")
  if !reused {
    t.Fatalf("expected the overlay edit to reuse the program, got reused=false")
  }

  text, ok := sess.SourceText(bAbs)
  if !ok || !strings.Contains(text, "return 2") {
    t.Fatalf("session did not reflect the overlay edit: ok=%v text=%q", ok, text)
  }
}
