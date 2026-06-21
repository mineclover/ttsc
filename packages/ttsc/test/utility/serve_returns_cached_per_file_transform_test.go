package ttsc_test

import (
  "bytes"
  "encoding/json"
  "path/filepath"
  "strings"
  "testing"

  "github.com/samchon/ttsc/packages/ttsc/utility"
)

// serveResponse mirrors the JSON reply RunServe writes per request.
type serveResponse struct {
  TypeScript string `json:"typescript"`
  Found      bool   `json:"found"`
}

// serveRequestLine encodes one resident-host request line for the given file.
func serveRequestLine(t *testing.T, file string) string {
  t.Helper()
  data, err := json.Marshal(map[string]string{"file": file})
  if err != nil {
    t.Fatal(err)
  }
  return string(data)
}

// TestUtilityServeReturnsCachedPerFileTransform verifies the resident serve host
// transforms the project once and then answers per-file requests from its cache.
//
// This is the resident transform host of samchon/ttsc#255: one warm process can
// serve every Metro worker, replacing the per-call transform subcommand that
// recompiles the project on each invocation. The host must key its cache exactly
// like the transform envelope (project-relative paths), accept absolute request
// paths, and report not-found for files outside the program.
//
// 1. Build a single-file project (no plugins, so transform yields the source).
// 2. Feed RunServe two requests: the project file, then a non-project file.
// 3. Assert the first reply returns the transformed source and the second is
//    reported not-found.
func TestUtilityServeReturnsCachedPerFileTransform(t *testing.T) {
  root := t.TempDir()
  writeProjectFile(t, root, "tsconfig.json", `{
  "compilerOptions": { "module": "commonjs", "target": "es2020", "noEmit": true },
  "files": ["index.ts"]
}
`)
  writeProjectFile(t, root, "index.ts", `export const value: number = 1;
`)

  requests := serveRequestLine(t, filepath.Join(root, "index.ts")) + "\n" +
    serveRequestLine(t, filepath.Join(root, "missing.ts")) + "\n"

  var out bytes.Buffer
  code := utility.RunServe(strings.NewReader(requests), &out, []string{"--cwd", root})
  if code != 0 {
    t.Fatalf("RunServe exit %d; output=%q", code, out.String())
  }

  lines := strings.Split(strings.TrimSpace(out.String()), "\n")
  if len(lines) != 2 {
    t.Fatalf("expected one reply per request, got %d: %q", len(lines), out.String())
  }

  var present serveResponse
  if err := json.Unmarshal([]byte(lines[0]), &present); err != nil {
    t.Fatalf("decode reply 0: %v (%q)", err, lines[0])
  }
  if !present.Found {
    t.Fatalf("expected the project file to be found: %q", lines[0])
  }
  if !strings.Contains(present.TypeScript, "value") {
    t.Fatalf("resident serve did not return the transformed source: %q", present.TypeScript)
  }

  var missing serveResponse
  if err := json.Unmarshal([]byte(lines[1]), &missing); err != nil {
    t.Fatalf("decode reply 1: %v (%q)", err, lines[1])
  }
  if missing.Found {
    t.Fatalf("expected a non-project file to be reported not-found: %q", lines[1])
  }
}
