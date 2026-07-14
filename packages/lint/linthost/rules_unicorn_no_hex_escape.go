// unicorn/no-hex-escape: `\xHH` escapes inside string and template
// literals are legal but they read worse than the equivalent
// `\u00HH` Unicode escape — there's only one way to write a Unicode
// escape, while a hex escape silently widens to a Unicode escape at
// runtime anyway. The rule nudges authors toward the Unicode form.
//
// AST-only: visit `KindStringLiteral`, `KindNoSubstitutionTemplateLiteral`,
// and the `KindTemplateHead`/`Middle`/`Tail` elements an interpolated
// template (or a template literal type) is built from, read each token's raw
// source text via `nodeText` (the parser already decodes escapes into the
// `.Text` value, so a normal accessor would see `©` instead of `\xA9`), and
// fire when the text contains a `\xHH` occurrence where HH is two hex digits.
// Tagged templates are skipped: their tag observes the raw text, where `\xA9`
// and `©` differ.
//
// https://github.com/sindresorhus/eslint-plugin-unicorn/blob/main/docs/rules/no-hex-escape.md
package linthost

import (
  "regexp"

  shimast "github.com/microsoft/typescript-go/shim/ast"
)

var unicornNoHexEscapePattern = regexp.MustCompile(`\\x[0-9A-Fa-f]{2}`)

type unicornNoHexEscape struct{}

func (unicornNoHexEscape) Name() string { return "unicorn/no-hex-escape" }
func (unicornNoHexEscape) Visits() []shimast.Kind {
  return []shimast.Kind{
    shimast.KindStringLiteral,
    shimast.KindNoSubstitutionTemplateLiteral,
    shimast.KindTemplateHead,
    shimast.KindTemplateMiddle,
    shimast.KindTemplateTail,
  }
}
func (unicornNoHexEscape) Check(ctx *Context, node *shimast.Node) {
  if isTaggedTemplateElement(node) {
    return
  }
  source := nodeText(ctx.File, node)
  if source == "" {
    return
  }
  if unicornNoHexEscapePattern.MatchString(source) {
    ctx.Report(node, "Prefer Unicode escapes (`\\uXXXX`) over hex escapes (`\\xHH`).")
  }
}

func init() {
  Register(unicornNoHexEscape{})
}
