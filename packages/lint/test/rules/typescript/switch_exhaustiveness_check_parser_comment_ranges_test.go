package linthost

import "testing"

// TestSwitchExhaustivenessCheckParserCommentRanges verifies the real command
// path recognizes only a trailing parser-classified default marker.
//
// The switch's AST-bounded trailing trivia gap is scanned independently. A
// Unicode-separated real comment must suppress the open-switch diagnostic,
// while identical bytes inside a template in the last clause remain outside
// that eligible gap.
//
//  1. Put a template-shaped marker in one open switch and expect a diagnostic.
//  2. Put a U+2028-separated real marker after another switch's last clause.
//  3. Assert the negative twin reports once and the real marker suppresses.
func TestSwitchExhaustivenessCheckParserCommentRanges(t *testing.T) {
  assertSwitchExhaustivenessCheckForTest(t, `
declare const value: string;
switch (value) {
  case "known":
    `${"// No Default"}`;
    break;
}
`, map[string]any{
    "requireDefaultForNonUnion": true,
  }, 1, map[string]int{"Cases not matched: default": 1})

  assertSwitchExhaustivenessCheckForTest(t, "\ndeclare const value: string;\nswitch (value) {\n  case \"known\": break;\u2028// No Default\u2029}\n", map[string]any{
    "requireDefaultForNonUnion": true,
  }, 0, nil)
}
