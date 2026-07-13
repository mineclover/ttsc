package linthost

import (
  "path/filepath"
  "testing"
)

// TestEditorFormatOverridesDuplicateLanguageSectionUsesLastValue verifies a
// repeated JSON property replaces its earlier object instead of merging it.
//
// VS Code first parses settings.json into an object, where the final value for
// a duplicate property wins without moving that property's insertion order.
// Treating both occurrences as separate override scopes would preserve stale
// keys that no longer exist in the parsed settings object.
//
// 1. Set tab size in the first occurrence of a combined language property.
// 2. Replace that property with an object that sets only insertSpaces.
// 3. Assert tab size falls back to the top level while the final object applies.
func TestEditorFormatOverridesDuplicateLanguageSectionUsesLastValue(t *testing.T) {
  root := t.TempDir()
  settings := `{
  "editor.tabSize": 8,
  "editor.insertSpaces": false,
  "[javascript][typescript]": { "editor.tabSize": 4 },
  "[javascript][typescript]": { "editor.insertSpaces": true }
}`
  writeFile(t, filepath.Join(root, ".vscode", "settings.json"), settings)

  got := editorFormatOverrides(root, "typescript")
  if got["tabWidth"] != float64(8) {
    t.Fatalf("replaced section must not retain stale tabWidth; want 8, got %v", got["tabWidth"])
  }
  if got["useTabs"] != false {
    t.Fatalf("last section value should set useTabs=false, got %v", got["useTabs"])
  }
}
