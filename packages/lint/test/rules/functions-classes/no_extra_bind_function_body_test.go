package linthost

import (
  "slices"
  "sort"
  "strings"
  "testing"
)

// TestRuleNoExtraBindFunctionBody verifies argument shape and `this` scope
// boundaries are both required before reporting a bind call.
//
// A nested arrow inherits the enclosing regular function's receiver, while a
// nested regular function owns a different receiver. Bind arguments after the
// receiver and spread calls are partial-application shapes, even for arrows.
//
// 1. Exercise direct, zero-argument, partial, spread, and dynamic-member calls.
// 2. Place `this` in a parameter default, nested arrow, and nested regular function.
// 3. Assert only the truly unnecessary regular-function binds are reported.
func TestRuleNoExtraBindFunctionBody(t *testing.T) {
  source := `declare const receiver: { value: number };
declare const bindArguments: [unknown];
declare const dynamicKey: string;

const direct = function () { return 1; }.bind(receiver); // diagnostic
const noReceiver = function () { return 2; }.bind();
const partial = function (value: number) { return value; }.bind(null, 1);
const arrowPartial = ((value: number) => value).bind(null, 1);
const spread = function () { return 3; }.bind(...bindArguments);
const dynamic = function () { return 4; }[dynamicKey](receiver);
const parameterDefault = function (this: { value: number }, value = this.value) { return value; }.bind(receiver);
const nestedArrow = function (this: { value: number }) { return () => this.value; }.bind(receiver);
const nestedRegular = function () { return function (this: { value: number }) { return this.value; }; }.bind(receiver); // diagnostic
`
  expectedLines := make([]int, 0)
  for index, line := range strings.Split(source, "\n") {
    if strings.Contains(line, "// diagnostic") {
      expectedLines = append(expectedLines, index+1)
    }
  }

  _, _, findings := runRuleFindingsSnapshot(t, "no-extra-bind", source, nil)
  actualLines := make([]int, 0, len(findings))
  for _, finding := range findings {
    actualLines = append(actualLines, strings.Count(source[:finding.Pos], "\n")+1)
  }
  sort.Ints(actualLines)
  if !slices.Equal(actualLines, expectedLines) {
    t.Fatalf("diagnostic lines mismatch: want %v, got %v", expectedLines, actualLines)
  }
}
