package linthost

import (
  "strings"
  "testing"

  shimast "github.com/microsoft/typescript-go/shim/ast"
)

// TestEngineIgnoresJsxTextInlineDisableMarkers verifies JSX text cannot create
// an inline lint directive while an expression-container comment still can.
//
// JSX text has a parser-owned lexical goal in which slash-shaped bytes are
// ordinary text. Treating those bytes as a block comment silently suppressed a
// real diagnostic, whereas `{/* ... */}` is genuine JavaScript comment trivia.
//
//  1. Put a disable-next-line-shaped string directly in JSX text before `debugger`.
//  2. Put the same marker in a JSX expression comment before another `debugger`.
//  3. Assert only the first statement is reported at its exact source range.
func TestEngineIgnoresJsxTextInlineDisableMarkers(t *testing.T) {
  const ruleName = "no-debugger"
  source := "const visible = <div>/* eslint-disable-next-line no-debugger */</div>;\ndebugger;\nconst active = <div>{/* eslint-disable-next-line no-debugger */}</div>;\ndebugger;\nJSON.stringify([visible, active]);\n"
  file := parseTSXFile(t, "/virtual/test.tsx", source)
  findings := NewEngine(RuleConfig{ruleName: SeverityError}).Run([]*shimast.SourceFile{file}, nil)
  if len(findings) != 1 {
    t.Fatalf("want 1 finding, got %d (%+v)", len(findings), findings)
  }
  start := strings.Index(source, "debugger;")
  if findings[0].Pos != start || findings[0].End != start+len("debugger") {
    t.Fatalf("want first debugger range [%d,%d), got [%d,%d)", start, start+len("debugger"), findings[0].Pos, findings[0].End)
  }
}
