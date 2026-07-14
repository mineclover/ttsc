package linthost

import (
  "sort"
  "strconv"
  "strings"
  "unicode"

  shimast "github.com/microsoft/typescript-go/shim/ast"
  shimscanner "github.com/microsoft/typescript-go/shim/scanner"
)

// regexpSourceRule implements the high-confidence, source-text subset of
// eslint-plugin-regexp. The first wave intentionally targets regex literals
// only: TypeScript-Go has already parsed those successfully, so these checks can
// stay AST-only and avoid pretending to evaluate arbitrary `RegExp(...)`
// constructor strings.
type regexpSourceRule struct {
  name  string
  check func(regexpLiteralParts) bool
}

type regexpLiteralParts struct {
  raw     string
  pattern string
  flags   string
}

func (r regexpSourceRule) Name() string { return r.name }
func (regexpSourceRule) Visits() []shimast.Kind {
  return []shimast.Kind{shimast.KindRegularExpressionLiteral}
}
func (r regexpSourceRule) Check(ctx *Context, node *shimast.Node) {
  parts, ok := parseRegexpLiteralParts(ctx, node)
  if !ok || r.check == nil || !r.check(parts) {
    return
  }
  ctx.Report(node, regexpRuleMessage(r.name))
}

type regexpNoUselessEscapeAlias struct{}

func (regexpNoUselessEscapeAlias) Name() string { return "regexp/no-useless-escape" }
func (regexpNoUselessEscapeAlias) Visits() []shimast.Kind {
  return []shimast.Kind{shimast.KindRegularExpressionLiteral}
}
func (regexpNoUselessEscapeAlias) Check(ctx *Context, node *shimast.Node) {
  pos, end := tokenRange(ctx.File, node)
  if pos < 0 {
    return
  }
  reportRegexEscapes(ctx, ctx.File.Text()[pos:end], pos)
}

func parseRegexpLiteralParts(ctx *Context, node *shimast.Node) (regexpLiteralParts, bool) {
  raw := nodeText(ctx.File, node)
  if len(raw) < 2 || raw[0] != '/' {
    return regexpLiteralParts{}, false
  }
  closing := strings.LastIndexByte(raw, '/')
  if closing <= 0 {
    return regexpLiteralParts{}, false
  }
  return regexpLiteralParts{
    raw:     raw,
    pattern: raw[1:closing],
    flags:   raw[closing+1:],
  }, true
}

func regexpRuleMessage(name string) string {
  switch name {
  case "regexp/no-control-character":
    return "Unexpected control character in regular expression."
  case "regexp/no-dupe-characters-character-class":
    return "Unexpected duplicate character in character class."
  case "regexp/no-empty-alternative":
    return "Unexpected empty alternative."
  case "regexp/no-empty-capturing-group":
    return "Unexpected empty capturing group."
  case "regexp/no-empty-character-class":
    return "Unexpected empty character class."
  case "regexp/no-empty-group":
    return "Unexpected empty group."
  case "regexp/no-empty-lookarounds-assertion":
    return "Unexpected empty lookaround assertion."
  case "regexp/no-misleading-unicode-character":
    return "Unexpected misleading Unicode character in character class."
  case "regexp/no-useless-character-class":
    return "Unexpected character class with one character."
  case "regexp/no-useless-flag":
    return "Unexpected useless regular expression flag."
  case "regexp/no-useless-quantifier":
    return "Unexpected useless quantifier."
  case "regexp/no-useless-two-nums-quantifier":
    return "Unexpected quantifier with equal minimum and maximum."
  case "regexp/no-zero-quantifier":
    return "Unexpected quantifier that repeats zero times."
  case "regexp/prefer-d":
    return "Prefer \\d over [0-9]."
  case "regexp/prefer-plus-quantifier":
    return "Prefer + over {1,}."
  case "regexp/prefer-question-quantifier":
    return "Prefer ? over {0,1}."
  case "regexp/prefer-star-quantifier":
    return "Prefer * over {0,}."
  case "regexp/prefer-w":
    return "Prefer \\w over [A-Za-z0-9_]."
  case "regexp/require-unicode-regexp":
    return "Regular expression should use the u or v flag."
  case "regexp/require-unicode-sets-regexp":
    return "Regular expression should use the v flag."
  case "regexp/sort-flags":
    return "Regular expression flags should be sorted."
  default:
    return "Unexpected regular expression pattern."
  }
}

func regexpHasEmptyAlternative(parts regexpLiteralParts) bool {
  return scanRegexpPattern(parts.pattern, func(pattern string, i int) bool {
    if pattern[i] != '|' {
      return false
    }
    return i == 0 || i == len(pattern)-1 || pattern[i-1] == '|' || pattern[i-1] == '(' || pattern[i+1] == '|' || pattern[i+1] == ')'
  })
}

func regexpHasEmptyCapturingGroup(parts regexpLiteralParts) bool {
  return scanRegexpPattern(parts.pattern, func(pattern string, i int) bool {
    return pattern[i] == '(' && i+1 < len(pattern) && pattern[i+1] == ')'
  })
}

func regexpHasEmptyGroup(parts regexpLiteralParts) bool {
  return scanRegexpPattern(parts.pattern, func(pattern string, i int) bool {
    return strings.HasPrefix(pattern[i:], "(?:)")
  })
}

func regexpHasEmptyLookaround(parts regexpLiteralParts) bool {
  return scanRegexpPattern(parts.pattern, func(pattern string, i int) bool {
    return strings.HasPrefix(pattern[i:], "(?=)") ||
      strings.HasPrefix(pattern[i:], "(?!)") ||
      strings.HasPrefix(pattern[i:], "(?<=)") ||
      strings.HasPrefix(pattern[i:], "(?<!)")
  })
}

// regexpHasEmptyCharacterClass reports non-negated character classes whose
// parsed element list is empty. The compiler's regexp parser is the syntax
// authority; the structural walk runs only after that parser accepts the full
// literal, so malformed escapes, flags, and class-set expressions cannot
// create lint findings.
func regexpHasEmptyCharacterClass(parts regexpLiteralParts) bool {
  if !shimscanner.IsValidRegularExpressionLiteral(parts.raw) {
    return false
  }
  unicodeSets := strings.Contains(parts.flags, "v")
  depth := 0
  for i := 0; i < len(parts.pattern); i++ {
    switch parts.pattern[i] {
    case '\\':
      i++
    case '[':
      if depth != 0 && !unicodeSets {
        continue
      }
      depth++
      content := i + 1
      negated := content < len(parts.pattern) && parts.pattern[content] == '^'
      if negated {
        content++
      }
      if !negated && content < len(parts.pattern) && parts.pattern[content] == ']' {
        return true
      }
    case ']':
      if depth > 0 {
        depth--
      }
    }
  }
  return false
}

func regexpHasZeroQuantifier(parts regexpLiteralParts) bool {
  return scanRegexpQuantifiers(parts.pattern, func(min, max int, hasComma bool) bool {
    return min == 0 && (!hasComma || max == 0)
  })
}

func regexpHasUselessTwoNumsQuantifier(parts regexpLiteralParts) bool {
  return scanRegexpQuantifiers(parts.pattern, func(min, max int, hasComma bool) bool {
    return hasComma && min == max && min >= 0
  })
}

func regexpHasUselessQuantifier(parts regexpLiteralParts) bool {
  return scanRegexpQuantifiers(parts.pattern, func(min, max int, hasComma bool) bool {
    return !hasComma && min == 1 && max == -1
  })
}

func regexpHasPlusQuantifierCandidate(parts regexpLiteralParts) bool {
  return scanRegexpQuantifiers(parts.pattern, func(min, max int, hasComma bool) bool {
    return hasComma && min == 1 && max == -1
  })
}

func regexpHasStarQuantifierCandidate(parts regexpLiteralParts) bool {
  return scanRegexpQuantifiers(parts.pattern, func(min, max int, hasComma bool) bool {
    return hasComma && min == 0 && max == -1
  })
}

func regexpHasQuestionQuantifierCandidate(parts regexpLiteralParts) bool {
  return scanRegexpQuantifiers(parts.pattern, func(min, max int, hasComma bool) bool {
    return hasComma && min == 0 && max == 1
  })
}

func regexpNeedsUnicodeFlag(parts regexpLiteralParts) bool {
  return !strings.ContainsAny(parts.flags, "uv")
}

func regexpNeedsUnicodeSetsFlag(parts regexpLiteralParts) bool {
  return !strings.Contains(parts.flags, "v")
}

func regexpFlagsUnsorted(parts regexpLiteralParts) bool {
  sorted := []byte(parts.flags)
  sort.SliceStable(sorted, func(i, j int) bool {
    return regexpFlagOrder(sorted[i]) < regexpFlagOrder(sorted[j])
  })
  return string(sorted) != parts.flags
}

func regexpFlagOrder(flag byte) int {
  const order = "dgimsuvy"
  if i := strings.IndexByte(order, flag); i >= 0 {
    return i
  }
  return len(order) + int(flag)
}

// regexpHasUselessFlag decides `regexp/no-useless-flag` for the two flags whose
// effect is visible in the pattern alone: `i` (nothing it could re-case) and `m`
// (no `^`/`$` for it to re-anchor).
//
// Both questions are answered on the regexp AST from regex_tree.go rather than
// on a byte scan. `scanRegexpPattern` never enters a character class -- correct
// for the rules that hunt `|`, `(`, `{` outside classes, and correct for `m`
// (`^`/`$` are literals inside a class), but fatal for `i`, because `[a-z]` is
// exactly where the flag earns its keep.
//
// A literal the parser rejects yields no finding at all: the rule tells people
// to delete a flag, so it stays silent whenever it cannot see the whole pattern.
func regexpHasUselessFlag(parts regexpLiteralParts) bool {
  ignoreCase := strings.Contains(parts.flags, "i")
  multiline := strings.Contains(parts.flags, "m")
  if !ignoreCase && !multiline {
    return false
  }
  parsed, err := regexParseLiteral(parts.raw)
  if err != nil {
    return false
  }
  if ignoreCase && !regexpNodeIsCaseVariant(parsed.Body, strings.ContainsAny(parts.flags, "uv")) {
    return true
  }
  if multiline && !regexpNodeHasLineAnchor(parsed.Body) {
    return true
  }
  return false
}

func regexpHasPreferD(parts regexpLiteralParts) bool {
  return strings.Contains(parts.pattern, "[0-9]")
}

func regexpHasPreferW(parts regexpLiteralParts) bool {
  return strings.Contains(parts.pattern, "[A-Za-z0-9_]") || strings.Contains(parts.pattern, "[a-zA-Z0-9_]")
}

func regexpHasDuplicateClassCharacter(parts regexpLiteralParts) bool {
  return walkRegexpCharacterClasses(parts.pattern, func(content string) bool {
    if classHasRange(content) {
      return false
    }
    seen := map[byte]struct{}{}
    for i := 0; i < len(content); i++ {
      ch := content[i]
      if ch == '\\' {
        i++
        continue
      }
      if ch == '^' && i == 0 {
        continue
      }
      if ch == '-' {
        continue
      }
      if _, ok := seen[ch]; ok {
        return true
      }
      seen[ch] = struct{}{}
    }
    return false
  })
}

func regexpHasUselessCharacterClass(parts regexpLiteralParts) bool {
  return walkRegexpCharacterClasses(parts.pattern, func(content string) bool {
    if len(content) != 1 {
      return false
    }
    ch := content[0]
    return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')
  })
}

func scanRegexpPattern(pattern string, visit func(pattern string, i int) bool) bool {
  inClass := false
  for i := 0; i < len(pattern); i++ {
    switch pattern[i] {
    case '\\':
      i++
    case '[':
      if !inClass {
        inClass = true
      }
    case ']':
      inClass = false
    default:
      if !inClass && visit(pattern, i) {
        return true
      }
    }
  }
  return false
}

func scanRegexpQuantifiers(pattern string, visit func(min, max int, hasComma bool) bool) bool {
  return scanRegexpPattern(pattern, func(pattern string, i int) bool {
    if pattern[i] != '{' {
      return false
    }
    end := i + 1
    for end < len(pattern) && pattern[end] != '}' {
      end++
    }
    if end >= len(pattern) {
      return false
    }
    body := pattern[i+1 : end]
    hasComma := strings.Contains(body, ",")
    min, max := -1, -1
    if hasComma {
      pair := strings.SplitN(body, ",", 2)
      min = parseRegexpQuantifierNumber(pair[0])
      if pair[1] != "" {
        max = parseRegexpQuantifierNumber(pair[1])
      }
    } else {
      min = parseRegexpQuantifierNumber(body)
    }
    return min >= 0 && visit(min, max, hasComma)
  })
}

func parseRegexpQuantifierNumber(text string) int {
  if text == "" {
    return -1
  }
  for i := 0; i < len(text); i++ {
    if text[i] < '0' || text[i] > '9' {
      return -1
    }
  }
  value, err := strconv.Atoi(text)
  if err != nil {
    return -1
  }
  return value
}

func walkRegexpCharacterClasses(pattern string, visit func(content string) bool) bool {
  for i := 0; i < len(pattern); i++ {
    if pattern[i] == '\\' {
      i++
      continue
    }
    if pattern[i] != '[' {
      continue
    }
    start := i + 1
    for j := start; j < len(pattern); j++ {
      if pattern[j] == '\\' {
        j++
        continue
      }
      if pattern[j] == ']' {
        if visit(pattern[start:j]) {
          return true
        }
        i = j
        break
      }
    }
  }
  return false
}

func classHasRange(content string) bool {
  for i := 1; i+1 < len(content); i++ {
    if content[i] == '\\' {
      i++
      continue
    }
    if content[i] == '-' {
      return true
    }
  }
  return false
}

// regexpCaseFoldScanLimit bounds the fold scan of a single character range.
// The scan exists to prove a range case-*invariant*, so refusing to walk an
// unbounded one costs at most a missed report, never a wrong one. Every script
// block that carries cased letters is orders of magnitude narrower than this.
const regexpCaseFoldScanLimit = 0x20000

// regexpNodeIsCaseVariant reports whether toggling the `i` flag can change what
// the node matches.
//
// Mirrors eslint-plugin-regexp's `isCaseVariant(pattern, flags, false)`, which
// judges a character class element by element rather than as a whole set: an
// element that the flag widens keeps the flag alive even when the class already
// spells both cases out, so `/[a-zA-Z]/i` is case-variant.
//
// The analysis is one-sided. Every construct it cannot settle from the AST --
// a Unicode property escape, a `v`-mode set-notation class, a range too wide to
// fold-scan, an unmodeled node -- counts as case-variant, so the rule stays
// quiet rather than order a load-bearing flag deleted. Backreferences are
// case-variant for a real reason and not merely out of caution: `i` makes the
// backreference comparison itself case-insensitive, so `/(.)\1/i` matches "aA"
// while `/(.)\1/` does not, without a single letter in the source.
func regexpNodeIsCaseVariant(node regexNode, unicodeMode bool) bool {
  switch n := node.(type) {
  case nil:
    return false
  case *regexCharNode:
    return regexpCharIsCaseVariant(n, unicodeMode)
  case *regexClassRangeNode:
    return regexpClassRangeIsCaseVariant(n, unicodeMode)
  case *regexClassNode:
    for _, expression := range n.Expressions {
      if regexpNodeIsCaseVariant(expression, unicodeMode) {
        return true
      }
    }
    return false
  case *regexAlternativeNode:
    for _, expression := range n.Expressions {
      if regexpNodeIsCaseVariant(expression, unicodeMode) {
        return true
      }
    }
    return false
  case *regexDisjunctionNode:
    return regexpNodeIsCaseVariant(n.Left, unicodeMode) ||
      regexpNodeIsCaseVariant(n.Right, unicodeMode)
  case *regexGroupNode:
    // A group name is never matched against the input, so `/(?<year>\d{4})/i`
    // stays case-invariant despite the letters in `year`.
    return regexpNodeIsCaseVariant(n.Expression, unicodeMode)
  case *regexRepetitionNode:
    return regexpNodeIsCaseVariant(n.Expression, unicodeMode)
  case *regexAssertionNode:
    switch n.Kind {
    case "^", "$":
      return false
    case "\\b", "\\B":
      // Word boundaries are defined in terms of `\w`, which grows by U+017F and
      // U+212A under `iu`/`iv`.
      return unicodeMode
    }
    return regexpNodeIsCaseVariant(n.Assertion, unicodeMode)
  case *regexBackreferenceNode:
    // `i` canonicalizes the backreference comparison itself.
    return true
  case *regexUnicodePropertyNode:
    // Whether `\p{...}` is closed under case folding needs the Unicode property
    // tables, not the source text: `\p{Lu}` moves under `i`, `\p{Nd}` does not.
    return true
  case *regexClassSetNode:
    // A `v`-mode set-notation class is kept verbatim by the parser, so its
    // members are not available to judge.
    return true
  }
  // An unmodeled node keeps the flag: silence beats a wrong deletion.
  return true
}

// regexpCharIsCaseVariant reports whether the `i` flag widens what a single
// character node matches.
func regexpCharIsCaseVariant(char *regexCharNode, unicodeMode bool) bool {
  if !char.codePointIsNaN() {
    return regexpRuneIsCaseVariant(rune(char.CodePoint))
  }
  switch char.Kind {
  case "meta":
    switch char.Value {
    case "\\w", "\\W":
      // `\w` is the one character set the flag moves: in Unicode mode
      // Canonicalize folds U+017F and U+212A onto `s` and `k`, so they join the
      // word characters (and leave `\W`).
      return unicodeMode
    case "\\d", "\\D", "\\s", "\\S", ".", "\\b":
      return false
    }
  case "decimal":
    // Annex B `\8` and `\9` match the bare digits.
    return false
  case "control":
    // `\cX` is a control code point, but a dangling `\c` matches the two
    // characters `\` and `c`, and that `c` re-cases.
    return char.Value == "\\c"
  }
  return true
}

// regexpClassRangeIsCaseVariant reports whether the `i` flag widens a character
// class range. A range is case-variant as soon as it contains one code point
// with a case-folded counterpart, because that counterpart may sit outside the
// range.
func regexpClassRangeIsCaseVariant(node *regexClassRangeNode, unicodeMode bool) bool {
  if node.From == nil || node.To == nil {
    return true
  }
  if node.From.codePointIsNaN() || node.To.codePointIsNaN() {
    // Annex B reads `[\d-z]` as three independent members while the AST still
    // models it as a range, so judge the two ends on their own.
    return regexpCharIsCaseVariant(node.From, unicodeMode) ||
      regexpCharIsCaseVariant(node.To, unicodeMode)
  }
  low, high := rune(node.From.CodePoint), rune(node.To.CodePoint)
  if low > high || high-low >= regexpCaseFoldScanLimit {
    return true
  }
  for r := low; r <= high; r++ {
    if regexpRuneIsCaseVariant(r) {
      return true
    }
  }
  return false
}

// regexpRuneIsCaseVariant reports whether a code point has any case-folded
// counterpart.
//
// unicode.SimpleFold walks the simple case-folding orbit, which is exactly the
// equivalence class ECMAScript's Canonicalize builds in `u`/`v` mode. Legacy
// mode canonicalizes more narrowly -- it keeps U+017F, U+212A and U+00DF apart
// from their ASCII or uppercase partners -- so there the orbit is a superset,
// and the rule at worst keeps quiet about a flag it could have reported.
func regexpRuneIsCaseVariant(r rune) bool {
  return unicode.SimpleFold(r) != r
}

// regexpNodeHasLineAnchor reports whether the pattern contains a `^` or `$`
// assertion, the only thing the `m` flag re-defines. Character classes are not
// descended into: `^` and `$` are literal characters in there.
func regexpNodeHasLineAnchor(node regexNode) bool {
  switch n := node.(type) {
  case *regexAssertionNode:
    if n.Kind == "^" || n.Kind == "$" {
      return true
    }
    return regexpNodeHasLineAnchor(n.Assertion)
  case *regexAlternativeNode:
    for _, expression := range n.Expressions {
      if regexpNodeHasLineAnchor(expression) {
        return true
      }
    }
    return false
  case *regexDisjunctionNode:
    return regexpNodeHasLineAnchor(n.Left) || regexpNodeHasLineAnchor(n.Right)
  case *regexGroupNode:
    return regexpNodeHasLineAnchor(n.Expression)
  case *regexRepetitionNode:
    return regexpNodeHasLineAnchor(n.Expression)
  }
  return false
}

func init() {
  Register(regexpSourceRule{name: "regexp/no-control-character", check: func(parts regexpLiteralParts) bool {
    return regexContainsControl(parts.raw)
  }})
  Register(regexpSourceRule{name: "regexp/no-dupe-characters-character-class", check: regexpHasDuplicateClassCharacter})
  Register(regexpSourceRule{name: "regexp/no-empty-alternative", check: regexpHasEmptyAlternative})
  Register(regexpSourceRule{name: "regexp/no-empty-capturing-group", check: regexpHasEmptyCapturingGroup})
  Register(regexpSourceRule{name: "regexp/no-empty-character-class", check: regexpHasEmptyCharacterClass})
  Register(regexpSourceRule{name: "regexp/no-empty-group", check: regexpHasEmptyGroup})
  Register(regexpSourceRule{name: "regexp/no-empty-lookarounds-assertion", check: regexpHasEmptyLookaround})
  Register(regexpSourceRule{name: "regexp/no-misleading-unicode-character", check: func(parts regexpLiteralParts) bool {
    return regexHasSurrogatePair(parts.raw)
  }})
  Register(regexpSourceRule{name: "regexp/no-useless-character-class", check: regexpHasUselessCharacterClass})
  Register(regexpNoUselessEscapeAlias{})
  Register(regexpSourceRule{name: "regexp/no-useless-flag", check: regexpHasUselessFlag})
  Register(regexpSourceRule{name: "regexp/no-useless-quantifier", check: regexpHasUselessQuantifier})
  Register(regexpSourceRule{name: "regexp/no-useless-two-nums-quantifier", check: regexpHasUselessTwoNumsQuantifier})
  Register(regexpSourceRule{name: "regexp/no-zero-quantifier", check: regexpHasZeroQuantifier})
  Register(regexpSourceRule{name: "regexp/prefer-d", check: regexpHasPreferD})
  Register(regexpSourceRule{name: "regexp/prefer-plus-quantifier", check: regexpHasPlusQuantifierCandidate})
  Register(regexpSourceRule{name: "regexp/prefer-question-quantifier", check: regexpHasQuestionQuantifierCandidate})
  Register(regexpSourceRule{name: "regexp/prefer-star-quantifier", check: regexpHasStarQuantifierCandidate})
  Register(regexpSourceRule{name: "regexp/prefer-w", check: regexpHasPreferW})
  Register(regexpSourceRule{name: "regexp/require-unicode-regexp", check: regexpNeedsUnicodeFlag})
  Register(regexpSourceRule{name: "regexp/require-unicode-sets-regexp", check: regexpNeedsUnicodeSetsFlag})
  Register(regexpSourceRule{name: "regexp/sort-flags", check: regexpFlagsUnsorted})
}
