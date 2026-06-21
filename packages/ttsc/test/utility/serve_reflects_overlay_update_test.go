package ttsc_test

import (
  "bytes"
  "encoding/json"
  "path/filepath"
  "strings"
  "testing"

  "github.com/samchon/ttsc/packages/ttsc/utility"
)

// serveUpdateReply mirrors the JSON reply RunServe writes for an update request.
type serveUpdateReply struct {
  Updated bool `json:"updated"`
}

// serveUpdateLine encodes one resident-host update request line.
func serveUpdateLine(t *testing.T, file, content string) string {
  t.Helper()
  data, err := json.Marshal(map[string]string{"content": content, "update": file})
  if err != nil {
    t.Fatal(err)
  }
  return string(data)
}

// TestUtilityServeReflectsOverlayUpdate verifies the resident serve host applies
// an in-memory edit and re-transforms, so a later transform request returns the
// edited content without restarting the host.
//
// This is the incremental half of the resident host (samchon/ttsc#255): an
// editor or watch consumer feeds an unsaved buffer through an update request and
// the next transform must reflect it. The host keys the overlay so the edit
// shadows the on-disk file, rebuilds the transform over the new content, and
// keeps serving from the warm process.
//
// 1. Transform index.ts and confirm the original value.
// 2. Update index.ts with new content and confirm the rebuild succeeded.
// 3. Transform index.ts again and confirm the edited value is returned.
func TestUtilityServeReflectsOverlayUpdate(t *testing.T) {
  root := t.TempDir()
  writeProjectFile(t, root, "tsconfig.json", `{
  "compilerOptions": { "module": "commonjs", "target": "es2020", "noEmit": true },
  "files": ["index.ts"]
}
`)
  writeProjectFile(t, root, "index.ts", `export const value: number = 1;
`)
  index := filepath.Join(root, "index.ts")

  requests := serveRequestLine(t, index) + "\n" +
    serveUpdateLine(t, index, "export const value: number = 2;\n") + "\n" +
    serveRequestLine(t, index) + "\n"

  var out bytes.Buffer
  code := utility.RunServe(strings.NewReader(requests), &out, []string{"--cwd", root})
  if code != 0 {
    t.Fatalf("RunServe exit %d; output=%q", code, out.String())
  }

  lines := strings.Split(strings.TrimSpace(out.String()), "\n")
  if len(lines) != 3 {
    t.Fatalf("expected one reply per request, got %d: %q", len(lines), out.String())
  }

  var before serveResponse
  if err := json.Unmarshal([]byte(lines[0]), &before); err != nil {
    t.Fatalf("decode reply 0: %v (%q)", err, lines[0])
  }
  if !before.Found || !strings.Contains(before.TypeScript, "1") {
    t.Fatalf("initial transform did not return the original value: %q", lines[0])
  }

  var updated serveUpdateReply
  if err := json.Unmarshal([]byte(lines[1]), &updated); err != nil {
    t.Fatalf("decode reply 1: %v (%q)", err, lines[1])
  }
  if !updated.Updated {
    t.Fatalf("expected the overlay update to rebuild successfully: %q", lines[1])
  }

  var after serveResponse
  if err := json.Unmarshal([]byte(lines[2]), &after); err != nil {
    t.Fatalf("decode reply 2: %v (%q)", err, lines[2])
  }
  if !after.Found || !strings.Contains(after.TypeScript, "2") {
    t.Fatalf("resident host did not reflect the overlay update: %q", lines[2])
  }
  // The edit must replace, not append: the original value must be gone.
  if strings.Contains(after.TypeScript, "= 1") {
    t.Fatalf("update did not replace the original value: %q", lines[2])
  }
}

// TestUtilityServeUpdateFailureKeepsPreviousTransform verifies that an update
// that does not compile leaves the previous transform in effect (reply
// updated:false), and that the failed edit is rolled back so a later valid
// update still succeeds rather than staying wedged on the broken buffer.
//
// This is the load-bearing half of the update contract (samchon/ttsc#255): an
// editor sends a transient broken buffer mid-keystroke, and the resident host
// must neither crash nor corrupt the cache, and must recover on the next good
// edit.
//
// 1. Transform index.ts (value 1).
// 2. Update with a type error; assert updated:false and that a transform still
//    returns the original value 1.
// 3. Update with valid new content; assert updated:true and the new value, which
//    only holds if the broken edit was rolled back out of the overlay.
func TestUtilityServeUpdateFailureKeepsPreviousTransform(t *testing.T) {
  root := t.TempDir()
  writeProjectFile(t, root, "tsconfig.json", `{
  "compilerOptions": { "module": "commonjs", "target": "es2020", "noEmit": true },
  "files": ["index.ts"]
}
`)
  writeProjectFile(t, root, "index.ts", `export const value: number = 1;
`)
  index := filepath.Join(root, "index.ts")

  requests := serveUpdateLine(t, index, "export const value: number = \"oops\";\n") + "\n" +
    serveRequestLine(t, index) + "\n" +
    serveUpdateLine(t, index, "export const value: number = 3;\n") + "\n" +
    serveRequestLine(t, index) + "\n"

  var out bytes.Buffer
  code := utility.RunServe(strings.NewReader(requests), &out, []string{"--cwd", root})
  if code != 0 {
    t.Fatalf("RunServe exit %d; output=%q", code, out.String())
  }

  lines := strings.Split(strings.TrimSpace(out.String()), "\n")
  if len(lines) != 4 {
    t.Fatalf("expected one reply per request, got %d: %q", len(lines), out.String())
  }

  var failed serveUpdateReply
  if err := json.Unmarshal([]byte(lines[0]), &failed); err != nil {
    t.Fatalf("decode reply 0: %v (%q)", err, lines[0])
  }
  if failed.Updated {
    t.Fatalf("expected the type-erroring update to fail: %q", lines[0])
  }

  var stale serveResponse
  if err := json.Unmarshal([]byte(lines[1]), &stale); err != nil {
    t.Fatalf("decode reply 1: %v (%q)", err, lines[1])
  }
  if !stale.Found || !strings.Contains(stale.TypeScript, "= 1") {
    t.Fatalf("failed update did not keep the previous transform: %q", lines[1])
  }

  var recovered serveUpdateReply
  if err := json.Unmarshal([]byte(lines[2]), &recovered); err != nil {
    t.Fatalf("decode reply 2: %v (%q)", err, lines[2])
  }
  if !recovered.Updated {
    t.Fatalf("a valid update after a failed one should succeed (rollback): %q", lines[2])
  }

  var after serveResponse
  if err := json.Unmarshal([]byte(lines[3]), &after); err != nil {
    t.Fatalf("decode reply 3: %v (%q)", err, lines[3])
  }
  if !after.Found || !strings.Contains(after.TypeScript, "3") {
    t.Fatalf("resident host did not recover to the valid update: %q", lines[3])
  }
}
