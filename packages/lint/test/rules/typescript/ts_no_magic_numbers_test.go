package linthost

import "testing"

// TestRuleCorpusTypeScriptNoMagicNumbers verifies the lint rule corpus
// fixture typescript-no-magic-numbers.ts.
//
// `typescript/no-magic-numbers` is AST-only: it mirrors the core
// `no-magic-numbers` baseline (whitelisting `-1` / `0` / `1`,
// `const x = N`, `arr[N]`, enum member initializers) and additionally
// skips TS-only positions — numeric literal types and readonly class
// property initializers — that would otherwise fire on every TS file
// that uses literal types or readonly constants.
//
// 1. Load the annotated TypeScript fixture source embedded below.
// 2. Enable the rule severities declared by its // expect: comments.
// 3. Assert the native Engine reports exactly the annotated diagnostics.
func TestRuleCorpusTypeScriptNoMagicNumbers(t *testing.T) {
	assertRuleCorpusCase(t, "typescript-no-magic-numbers.ts", "// Negative: enum member values are intentional named numbers.\nenum HttpStatus {\n  Ok = 200,\n  NotFound = 404,\n  ServerError = 500,\n}\n\n// Negative: literal numeric type — type position.\ntype ZeroOrTwo = 0 | 2;\n\n// Negative: unit values carry intrinsic meaning.\nconst counter = 0;\nconst stepSize = 1;\nconst notFound = -1;\n\n// Negative: `const x = N` is the named binding.\nconst SECONDS_PER_DAY = 86400;\n\n// Negative: readonly class property — the field name is the binding.\nclass Timeout {\n  readonly defaultMs = 5000;\n}\n\n// Positive: bare magic number in a comparison.\nfunction isLong(value: number): boolean {\n  // expect: typescript/no-magic-numbers error\n  return value > 86400;\n}\n\n// Positive: magic number in an arithmetic expression.\nfunction feeFor(amount: number): number {\n  // expect: typescript/no-magic-numbers error\n  return amount * 0.035;\n}\n\n// Positive: `let` cannot anchor a named constant.\n// expect: typescript/no-magic-numbers error\nlet timeoutMs = 5000;\n\nJSON.stringify({\n  HttpStatus,\n  counter,\n  stepSize,\n  notFound,\n  SECONDS_PER_DAY,\n  Timeout,\n  isLong,\n  feeFor,\n  timeoutMs,\n  zeroOrTwo: null as ZeroOrTwo | null,\n});\n")
}
