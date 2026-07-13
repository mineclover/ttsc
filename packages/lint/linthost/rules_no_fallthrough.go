// noFallthrough: `switch` cases whose end is reachable and that spill into
// the next `case` / `default` label without an intentional-fallthrough
// comment.
//
// The rule mirrors ESLint's `no-fallthrough` semantics:
//
//   - Reachability is decided by a structured statement completion
//     analysis ("can this statement list complete normally?") that
//     composes blocks, `if/else`, loops, `switch`, labeled statements,
//     and `try/catch/finally`. Nested function and class bodies are
//     never entered, so a callback's `return` does not terminate the
//     enclosing case. The analysis is deliberately implemented here
//     instead of reusing the checker's flow nodes: binder flow graphs
//     require a bound Program plus checker-side reachability walks,
//     while this rule must stay AST-only so the engine can keep running
//     it in the parallel no-checker lane.
//   - An intentional fallthrough is marked by the last comment before
//     the next case label (or, when the case body is exactly one block
//     statement, the last comment before that block's closing brace)
//     matching the fallthrough comment pattern. The default pattern is
//     ESLint's `/falls?\s?through/i`. Directive comments
//     (`eslint-disable-next-line …`, `lint-disable …`, `globals …`) are
//     never markers even when their text happens to match.
//   - `commentPattern`, `allowEmptyCase`, and
//     `reportUnusedFallthroughComment` options match the upstream
//     schema. A custom `commentPattern` replaces the default marker
//     pattern entirely, exactly like ESLint.
//
// https://eslint.org/docs/latest/rules/no-fallthrough
package linthost

import (
  "regexp"
  "strconv"
  "strings"

  shimast "github.com/microsoft/typescript-go/shim/ast"
  shimscanner "github.com/microsoft/typescript-go/shim/scanner"
)

// noFallthroughDefaultCommentPattern is ESLint's DEFAULT_FALLTHROUGH_COMMENT
// (`/falls?\s?through/iu`) translated to RE2. It accepts `fall through`,
// `falls through`, `fallthrough`, `fallsthrough`, any letter case.
var noFallthroughDefaultCommentPattern = regexp.MustCompile(`(?i)falls?\s?through`)

// noFallthroughDirectivePattern ports ESLint's shared `directivesPattern`
// and extends it with the `lint-*` directive family this host also
// recognizes (see directives.go). A directive comment is configuration,
// not documentation, so it never counts as an intentional-fallthrough
// marker even when its text matches the comment pattern (e.g.
// `// eslint-enable no-fallthrough`).
var noFallthroughDirectivePattern = regexp.MustCompile(
  `^(?:eslint(?:-env|-enable|-disable(?:(?:-next)?-line)?)?|lint-(?:enable|disable(?:(?:-next)?-line)?)|exported|globals?)(?:\s|$)`,
)

// noFallthroughOptions mirrors the upstream rule's options schema.
type noFallthroughOptions struct {
  // CommentPattern replaces the default fallthrough marker pattern when
  // non-empty. Compiled as an RE2 expression against the comment's inner
  // text; an invalid pattern falls back to the default marker pattern so
  // a config typo cannot silently disable marker recognition.
  CommentPattern string `json:"commentPattern"`
  // AllowEmptyCase suppresses the blank-line heuristic for empty cases:
  // by default an empty case separated from the next label by at least
  // one blank line is treated as an accidental fallthrough.
  AllowEmptyCase bool `json:"allowEmptyCase"`
  // ReportUnusedFallthroughComment reports fallthrough marker comments on
  // cases that cannot actually fall through (e.g. after a `break`).
  ReportUnusedFallthroughComment bool `json:"reportUnusedFallthroughComment"`
}

type noFallthrough struct{}

func (noFallthrough) Name() string           { return "no-fallthrough" }
func (noFallthrough) Visits() []shimast.Kind { return []shimast.Kind{shimast.KindSwitchStatement} }
func (noFallthrough) Check(ctx *Context, node *shimast.Node) {
  sw := node.AsSwitchStatement()
  if sw == nil || sw.CaseBlock == nil || ctx.File == nil {
    return
  }
  block := sw.CaseBlock.AsCaseBlock()
  if block == nil || block.Clauses == nil {
    return
  }
  var opts noFallthroughOptions
  ctx.DecodeOptions(&opts)
  pattern := noFallthroughDefaultCommentPattern
  if opts.CommentPattern != "" {
    if custom, err := regexp.Compile(opts.CommentPattern); err == nil {
      pattern = custom
    }
  }
  clauses := block.Clauses.Nodes
  for i := 0; i+1 < len(clauses); i++ {
    clause := clauses[i].AsCaseOrDefaultClause()
    next := clauses[i+1]
    if clause == nil || next == nil {
      continue
    }
    stmts := clauseStatements(clause)
    // ESLint calls this `isSwitchExitReachable`: the code path segment at
    // the end of the case's consequent is reachable, i.e. the statement
    // list can complete normally.
    endReachable := statementListCompletion(stmts).normal
    // The last case never falls through (there is no next label), which
    // the `i+1 < len(clauses)` loop bound already guarantees.
    fallsThrough := endReachable &&
      (len(stmts) > 0 ||
        (!opts.AllowEmptyCase && noFallthroughHasBlankLinesBetween(ctx.File, clauses[i], next)))
    commentPos, commentEnd, hasMarker := noFallthroughMarkerComment(ctx.File, clause, next, pattern)
    if fallsThrough && !hasMarker {
      label := "case"
      if next.Kind == shimast.KindDefaultClause {
        label = "default"
      }
      ctx.Report(next, "Expected a 'break' statement before '"+label+"'.")
    } else if opts.ReportUnusedFallthroughComment && !endReachable && hasMarker {
      ctx.ReportRange(commentPos, commentEnd, "Found a comment that would permit fallthrough, but case cannot fall through.")
    }
  }
}

// clauseStatements returns a clause's statement list, tolerating nil
// intermediates on malformed parses.
func clauseStatements(clause *shimast.CaseOrDefaultClause) []*shimast.Node {
  if clause == nil || clause.Statements == nil {
    return nil
  }
  return clause.Statements.Nodes
}

// noFallthroughHasBlankLinesBetween reports whether at least one blank
// line separates the end of `clause` from the first token of `next`
// (comments between them are skipped, exactly like ESLint's
// `getTokenAfter`). ESLint treats a blank-line gap after an empty case as
// an accidental fallthrough unless `allowEmptyCase` is set.
func noFallthroughHasBlankLinesBetween(file *shimast.SourceFile, clause, next *shimast.Node) bool {
  text := file.Text()
  nextToken := shimscanner.SkipTrivia(text, next.Pos())
  endLine := shimscanner.GetECMALineOfPosition(file, clause.End())
  nextLine := shimscanner.GetECMALineOfPosition(file, nextToken)
  return nextLine > endLine+1
}

// noFallthroughMarkerComment locates ESLint's intentional-fallthrough
// marker for the transition `clause` → `next` and returns its source
// range. Two positions are eligible, checked in upstream order:
//
//  1. When the case body is exactly one block statement, the last
//     comment before the block's closing brace.
//  2. The last comment between the end of `clause` and the `case` /
//     `default` keyword of `next`.
//
// Only the LAST comment of each region counts (`getCommentsBefore(...).pop()`
// upstream), so an unrelated comment after the marker invalidates it.
func noFallthroughMarkerComment(
  file *shimast.SourceFile,
  clause *shimast.CaseOrDefaultClause,
  next *shimast.Node,
  pattern *regexp.Regexp,
) (pos, end int, ok bool) {
  text := file.Text()
  stmts := clauseStatements(clause)
  if len(stmts) == 1 && stmts[0] != nil && stmts[0].Kind == shimast.KindBlock {
    if block := stmts[0].AsBlock(); block != nil {
      from := blockInteriorStart(text, block)
      closeBrace := block.End() - 1 // exclusive End covers the `}` token
      if from >= 0 && closeBrace > from {
        if p, e, found := lastCommentInTrivia(text, from, closeBrace); found &&
          isNoFallthroughMarker(text[p:e], pattern) {
          return p, e, true
        }
      }
    }
  }
  nextToken := shimscanner.SkipTrivia(text, next.Pos())
  if p, e, found := lastCommentInTrivia(text, clause.End(), nextToken); found &&
    isNoFallthroughMarker(text[p:e], pattern) {
    return p, e, true
  }
  return 0, 0, false
}

// blockInteriorStart returns the offset just past a block's last interior
// token: after the final statement when the block has one, otherwise
// after the opening `{`. The region from there to the closing brace is
// trivia-only. Returns -1 when the block shape is malformed.
func blockInteriorStart(text string, block *shimast.Block) int {
  if block.Statements != nil {
    if nodes := block.Statements.Nodes; len(nodes) > 0 {
      last := nodes[len(nodes)-1]
      if last == nil {
        return -1
      }
      return last.End()
    }
  }
  openBrace := shimscanner.SkipTrivia(text, block.Pos())
  if openBrace < 0 || openBrace >= len(text) || text[openBrace] != '{' {
    return -1
  }
  return openBrace + 1
}

// isNoFallthroughMarker reports whether one raw comment token (delimiters
// included) is an intentional-fallthrough marker: its inner text matches
// the configured pattern and it is not a directive comment. Mirrors
// upstream `isFallThroughComment`, which tests the un-trimmed comment
// value against the pattern but the trimmed value against the directive
// pattern.
func isNoFallthroughMarker(raw string, pattern *regexp.Regexp) bool {
  value := commentInnerText(raw)
  return pattern.MatchString(value) &&
    !noFallthroughDirectivePattern.MatchString(strings.TrimSpace(value))
}

// commentInnerText strips the comment delimiters from a raw comment token
// and returns the interior verbatim — no trimming, no JSDoc `*` handling —
// matching ESLint's `comment.value`. (stripCommentDelimiters is not reused
// here because it trims and strips JSDoc stars, which would let anchored
// custom patterns match text ESLint would reject.)
func commentInnerText(raw string) string {
  switch {
  case strings.HasPrefix(raw, "//"):
    return raw[2:]
  case strings.HasPrefix(raw, "/*"):
    inner := raw[2:]
    inner = strings.TrimSuffix(inner, "*/")
    return inner
  default:
    return raw
  }
}

// lastCommentInTrivia scans the trivia-only region text[from:to) and
// returns the source range of the last comment token in it. Both regions
// this rule inspects — between a clause end and the next case keyword,
// and between a block's last token and its closing brace — contain only
// whitespace and comments by construction, so a local scan is exact. The
// scan bails out on any unexpected non-trivia byte instead of guessing.
func lastCommentInTrivia(text string, from, to int) (pos, end int, ok bool) {
  if from < 0 || to > len(text) || from >= to {
    return 0, 0, false
  }
  i := from
  for i < to {
    c := text[i]
    switch {
    case c == ' ' || c == '\t' || c == '\r' || c == '\n' || c == '\v' || c == '\f':
      i++
    case c == '/' && i+1 < to && text[i+1] == '/':
      start := i
      i += 2
      for i < to && text[i] != '\n' && text[i] != '\r' {
        i++
      }
      pos, end, ok = start, i, true
    case c == '/' && i+1 < to && text[i+1] == '*':
      start := i
      i += 2
      for i < to {
        if text[i] == '*' && i+1 < to && text[i+1] == '/' {
          i += 2
          break
        }
        i++
      }
      pos, end, ok = start, i, true
    case c >= 0x80:
      // Exotic Unicode whitespace (NBSP, U+2028, …) is legal trivia;
      // skip the whole rune. Any other non-ASCII byte here would mean
      // the region was not trivia-only, and the rune skip is still the
      // safest way to keep scanning without splitting a code point.
      i++
      for i < to && text[i]&0xC0 == 0x80 {
        i++
      }
    default:
      // Not a trivia byte: the region assumption was violated
      // (malformed parse). Stop rather than misattribute comments.
      return pos, end, ok
    }
  }
  return pos, end, ok
}

// caseCompletion records every way a statement can finish: normally, by
// return or throw, or through a labeled/unlabeled break or continue. The
// empty-string key stands for the unlabeled form. Enclosing constructs
// consume the completions they own and propagate the rest.
type caseCompletion struct {
  normal    bool
  returns   bool
  throws    bool
  breaks    map[string]struct{}
  continues map[string]struct{}
}

func (c *caseCompletion) addBreak(label string) {
  if c.breaks == nil {
    c.breaks = map[string]struct{}{}
  }
  c.breaks[label] = struct{}{}
}

func (c *caseCompletion) addContinue(label string) {
  if c.continues == nil {
    c.continues = map[string]struct{}{}
  }
  c.continues[label] = struct{}{}
}

func (c *caseCompletion) hasBreak(label string) bool {
  _, ok := c.breaks[label]
  return ok
}

func (c *caseCompletion) hasContinue(label string) bool {
  _, ok := c.continues[label]
  return ok
}

func (c *caseCompletion) removeBreak(label string)    { delete(c.breaks, label) }
func (c *caseCompletion) removeContinue(label string) { delete(c.continues, label) }

// mergeAbrupt unions the other completion's abrupt paths into this one.
// normal is deliberately untouched; each composition rule computes it.
func (c *caseCompletion) mergeAbrupt(other caseCompletion) {
  c.returns = c.returns || other.returns
  c.throws = c.throws || other.throws
  for label := range other.breaks {
    c.addBreak(label)
  }
  for label := range other.continues {
    c.addContinue(label)
  }
}

// hasCompletion reports whether at least one control-flow path leaves the
// statement. An infinite loop with no reachable escape has no completion,
// which also means a following finally block is unreachable.
func (c *caseCompletion) hasCompletion() bool {
  return c.normal || c.returns || c.throws || len(c.breaks) > 0 || len(c.continues) > 0
}

// statementListCompletion runs the completion analysis over a statement
// list executed in order. Once a statement cannot complete normally the
// remainder of the list is unreachable: it neither restores normal
// completion nor contributes escapes (an unreachable `break` can never
// execute, matching ESLint's unreachable-segment tracking).
func statementListCompletion(stmts []*shimast.Node) caseCompletion {
  out := caseCompletion{normal: true}
  for _, stmt := range stmts {
    if stmt == nil {
      continue
    }
    if !out.normal {
      break
    }
    r := statementCompletion(stmt, nil)
    out.mergeAbrupt(r)
    out.normal = r.normal
  }
  return out
}

// statementCompletion computes how a single statement can complete.
// `labels` carries the label names bound to this statement by directly
// wrapping labeled statements, so loops can absorb `continue L` / `break L`
// aimed at themselves. Evaluated expressions are inspected only for the
// throwable nodes that can reach a catch; nested function, class-field,
// and static-block code paths remain isolated from the enclosing case.
func statementCompletion(stmt *shimast.Node, labels []string) caseCompletion {
  if stmt == nil {
    return caseCompletion{normal: true}
  }
  switch stmt.Kind {
  case shimast.KindReturnStatement:
    ret := stmt.AsReturnStatement()
    return caseCompletion{
      returns: true,
      throws:  ret != nil && executableNodePotentiallyThrows(ret.Expression),
    }
  case shimast.KindThrowStatement:
    // Evaluating the operand may throw first, and the statement itself
    // always produces a throw completion. Both reach the same catch edge.
    return caseCompletion{throws: true}
  case shimast.KindBreakStatement:
    out := caseCompletion{}
    out.addBreak(identifierText(stmt.AsBreakStatement().Label))
    return out
  case shimast.KindContinueStatement:
    out := caseCompletion{}
    out.addContinue(identifierText(stmt.AsContinueStatement().Label))
    return out
  case shimast.KindBlock:
    block := stmt.AsBlock()
    if block == nil || block.Statements == nil {
      return caseCompletion{normal: true}
    }
    return statementListCompletion(block.Statements.Nodes)
  case shimast.KindIfStatement:
    return ifCompletion(stmt.AsIfStatement())
  case shimast.KindLabeledStatement:
    return labeledCompletion(stmt.AsLabeledStatement(), labels)
  case shimast.KindWhileStatement:
    s := stmt.AsWhileStatement()
    constTrue, constFalse := literalTruthiness(s.Expression)
    return loopCompletion(
      s.Statement,
      labels,
      constTrue,
      constFalse,
      false,
      executableNodePotentiallyThrows(s.Expression),
      false,
    )
  case shimast.KindDoStatement:
    s := stmt.AsDoStatement()
    constTrue, constFalse := literalTruthiness(s.Expression)
    return loopCompletion(
      s.Statement,
      labels,
      constTrue,
      constFalse,
      true,
      executableNodePotentiallyThrows(s.Expression),
      false,
    )
  case shimast.KindForStatement:
    s := stmt.AsForStatement()
    constTrue, constFalse := literalTruthiness(s.Condition)
    if s.Condition == nil {
      // `for (;;)` loops forever unless something breaks out.
      constTrue, constFalse = true, false
    }
    out := loopCompletion(
      s.Statement,
      labels,
      constTrue,
      constFalse,
      false,
      executableNodePotentiallyThrows(s.Condition),
      executableNodePotentiallyThrows(s.Incrementor),
    )
    out.throws = out.throws || executableNodePotentiallyThrows(s.Initializer)
    return out
  case shimast.KindForInStatement, shimast.KindForOfStatement:
    // The iterated collection may be empty, so the loop always offers
    // normal completion — identical to a non-constant loop test.
    s := stmt.AsForInOrOfStatement()
    out := loopCompletion(s.Statement, labels, false, false, false, false, false)
    out.throws = out.throws ||
      executableNodePotentiallyThrows(s.Initializer) ||
      executableNodePotentiallyThrows(s.Expression)
    return out
  case shimast.KindSwitchStatement:
    return switchCompletion(stmt.AsSwitchStatement(), labels)
  case shimast.KindTryStatement:
    return tryCompletion(stmt.AsTryStatement())
  case shimast.KindWithStatement:
    s := stmt.AsWithStatement()
    out := statementCompletion(s.Statement, nil)
    out.throws = out.throws || executableNodePotentiallyThrows(s.Expression)
    return out
  default:
    // Leaf statements complete normally, but their evaluated expressions
    // can still enter an enclosing catch. The walker excludes type syntax
    // and nested function/class execution contexts.
    return caseCompletion{normal: true, throws: executableNodePotentiallyThrows(stmt)}
  }
}

// ifCompletion: an `if` without `else` can always complete normally (the
// condition may be false); with an `else`, normal completion requires at
// least one branch to complete normally. Conditions are not constant-
// folded, matching ESLint's code path analysis which folds only loop
// tests.
func ifCompletion(s *shimast.IfStatement) caseCompletion {
  if s == nil {
    return caseCompletion{normal: true}
  }
  then := statementCompletion(s.ThenStatement, nil)
  if s.ElseStatement == nil {
    then.normal = true
    then.throws = then.throws || executableNodePotentiallyThrows(s.Expression)
    return then
  }
  els := statementCompletion(s.ElseStatement, nil)
  out := caseCompletion{
    normal: then.normal || els.normal,
    throws: executableNodePotentiallyThrows(s.Expression),
  }
  out.mergeAbrupt(then)
  out.mergeAbrupt(els)
  return out
}

// labeledCompletion: `L: stmt` completes normally when the body does, or
// when a `break L` escapes the body (control resumes right after the
// labeled statement). The label is added to the body's label set so a
// directly-labeled loop can absorb `continue L` itself.
func labeledCompletion(s *shimast.LabeledStatement, labels []string) caseCompletion {
  if s == nil {
    return caseCompletion{normal: true}
  }
  name := identifierText(s.Label)
  if name == "" {
    return statementCompletion(s.Statement, labels)
  }
  inner := statementCompletion(s.Statement, append(labels, name))
  if inner.hasBreak(name) {
    inner.normal = true
  }
  inner.removeBreak(name)
  inner.removeContinue(name)
  return inner
}

// loopCompletion is the shared engine for `while`, `do/while`, `for`,
// `for-in`, and `for-of` bodies. Loop tests are constant-folded only for
// simple literals (`true`, `1`, `"x"`, `null`, …), matching ESLint's
// `getBooleanValueIfSimpleConstant`. The loop absorbs unlabeled breaks
// and continues plus the labeled forms naming the loop itself (via
// `labels`); everything else escapes to the enclosing construct.
func loopCompletion(
  body *shimast.Node,
  labels []string,
  constTrue, constFalse, isDoWhile bool,
  testThrows, incrementThrows bool,
) caseCompletion {
  if constFalse && !isDoWhile {
    // The body never runs, so nothing inside it (including breaks) can
    // execute. The test is still evaluated once before normal completion.
    return caseCompletion{normal: true, throws: testThrows}
  }
  r := statementCompletion(body, nil)
  exitByBreak := r.hasBreak("")
  iterationEnds := r.normal || r.hasContinue("")
  for _, l := range labels {
    exitByBreak = exitByBreak || r.hasBreak(l)
    iterationEnds = iterationEnds || r.hasContinue(l)
  }
  if !isDoWhile || iterationEnds {
    // while/for tests run before the first iteration. A do/while test is
    // reachable only when the body reaches the iteration boundary.
    r.throws = r.throws || testThrows
  }
  if iterationEnds {
    r.throws = r.throws || incrementThrows
  }
  var normal bool
  switch {
  case constTrue:
    // Infinite loop: only a break reaches the code after it.
    normal = exitByBreak
  case isDoWhile:
    // The body runs at least once; afterwards the test can fail.
    normal = exitByBreak || iterationEnds
  default:
    // The test may fail before the first iteration.
    normal = true
  }
  r.removeBreak("")
  r.removeContinue("")
  for _, l := range labels {
    r.removeBreak(l)
    r.removeContinue(l)
  }
  r.normal = normal
  return r
}

// switchCompletion: a nested `switch` completes normally when it has no
// `default` (the discriminant may match nothing), when its final clause
// completes normally (falling off the end), or when any reachable `break`
// targets it. Labeled breaks naming the switch (via `labels`) are
// absorbed the same way; `continue` never targets a switch and passes
// through.
func switchCompletion(s *shimast.SwitchStatement, labels []string) caseCompletion {
  if s == nil || s.CaseBlock == nil {
    return caseCompletion{normal: true}
  }
  block := s.CaseBlock.AsCaseBlock()
  if block == nil || block.Clauses == nil || len(block.Clauses.Nodes) == 0 {
    return caseCompletion{normal: true}
  }
  clauses := block.Clauses.Nodes
  hasDefault := false
  exitByBreak := false
  lastNormal := false
  out := caseCompletion{throws: executableNodePotentiallyThrows(s.Expression)}
  for i, clauseNode := range clauses {
    if clauseNode == nil {
      continue
    }
    if clauseNode.Kind == shimast.KindDefaultClause {
      hasDefault = true
    }
    clause := clauseNode.AsCaseOrDefaultClause()
    if clause != nil {
      out.throws = out.throws || executableNodePotentiallyThrows(clause.Expression)
    }
    r := statementListCompletion(clauseStatements(clause))
    if r.hasBreak("") {
      exitByBreak = true
    }
    for _, l := range labels {
      if r.hasBreak(l) {
        exitByBreak = true
      }
    }
    if i == len(clauses)-1 {
      lastNormal = r.normal
    }
    out.mergeAbrupt(r)
  }
  out.removeBreak("")
  for _, l := range labels {
    out.removeBreak(l)
  }
  out.normal = !hasDefault || lastNormal || exitByBreak
  return out
}

// tryCompletion composes normal, return, throw, break, and continue paths.
// A catch contributes only when the protected block has a reachable throw
// edge. A reachable finally runs for every completion from try/catch; its
// abrupt completions override the incoming completion, while an ordinary
// finally path preserves it.
func tryCompletion(s *shimast.TryStatement) caseCompletion {
  if s == nil {
    return caseCompletion{normal: true}
  }
  tryC := blockNodeCompletion(s.TryBlock)
  var catchC caseCompletion
  hasCatch := false
  if s.CatchClause != nil {
    if clause := s.CatchClause.AsCatchClause(); clause != nil {
      hasCatch = true
      catchC = blockNodeCompletion(clause.Block)
    }
  }
  main := tryC
  if hasCatch && tryC.throws {
    // The catch consumes every throw from the protected block and replaces
    // those paths with the catch block's possible completions.
    main.throws = false
    main.normal = tryC.normal || catchC.normal
    main.mergeAbrupt(catchC)
  }
  if s.FinallyBlock == nil {
    return main
  }
  if !main.hasCompletion() {
    // No path leaves the protected region (for example, a closed infinite
    // loop), so execution never enters the finally block.
    return caseCompletion{}
  }
  finallyC := blockNodeCompletion(s.FinallyBlock)
  out := caseCompletion{}
  if finallyC.normal {
    out.mergeAbrupt(main)
  }
  out.mergeAbrupt(finallyC)
  out.normal = main.normal && finallyC.normal
  return out
}

// blockNodeCompletion analyzes a *BlockNode (try/catch/finally bodies).
func blockNodeCompletion(node *shimast.Node) caseCompletion {
  if node == nil {
    return caseCompletion{normal: true}
  }
  block := node.AsBlock()
  if block == nil || block.Statements == nil {
    return caseCompletion{normal: true}
  }
  return statementListCompletion(block.Statements.Nodes)
}

// executableNodePotentiallyThrows mirrors the throwable-node boundary in
// ESLint's code-path analyzer. A reachable value-reference identifier,
// member access, call, import call, or construction can enter the nearest
// catch. Type syntax and separately analyzed function/class execution
// contexts never leak throw edges into the enclosing statement.
func executableNodePotentiallyThrows(node *shimast.Node) bool {
  if node == nil || shimast.IsTypeNode(node) {
    return false
  }
  switch node.Kind {
  case shimast.KindFunctionDeclaration,
    shimast.KindFunctionExpression,
    shimast.KindArrowFunction,
    shimast.KindClassStaticBlockDeclaration:
    return false
  case shimast.KindMethodDeclaration,
    shimast.KindConstructor,
    shimast.KindGetAccessor,
    shimast.KindSetAccessor,
    shimast.KindPropertyDeclaration:
    return classElementHeaderPotentiallyThrows(node)
  case shimast.KindIdentifier:
    return noFallthroughIdentifierIsReference(node)
  case shimast.KindPropertyAccessExpression,
    shimast.KindElementAccessExpression,
    shimast.KindCallExpression,
    shimast.KindNewExpression:
    return true
  }
  found := false
  node.ForEachChild(func(child *shimast.Node) bool {
    if executableNodePotentiallyThrows(child) {
      found = true
      return true
    }
    return false
  })
  return found
}

// classElementHeaderPotentiallyThrows keeps class member bodies and field
// initializers in their own code paths while retaining immediately evaluated
// decorators and computed property names in the enclosing class evaluation.
func classElementHeaderPotentiallyThrows(node *shimast.Node) bool {
  if node == nil {
    return false
  }
  if modifiers := node.Modifiers(); modifiers != nil {
    for _, modifier := range modifiers.Nodes {
      if executableNodePotentiallyThrows(modifier) {
        return true
      }
    }
  }
  return executableNodePotentiallyThrows(node.Name())
}

// noFallthroughIdentifierIsReference excludes declaration names, property
// names, and statement labels. Every remaining identifier is a value read,
// matching the first-throwable-node rule used by ESLint inside try blocks.
func noFallthroughIdentifierIsReference(node *shimast.Node) bool {
  if node == nil || node.Parent == nil {
    return true
  }
  parent := node.Parent
  switch parent.Kind {
  case shimast.KindPropertyAccessExpression:
    access := parent.AsPropertyAccessExpression()
    return access == nil || access.Name() != node
  case shimast.KindPropertyAssignment:
    assignment := parent.AsPropertyAssignment()
    return assignment == nil || assignment.Name() != node
  case shimast.KindBindingElement:
    element := parent.AsBindingElement()
    return element == nil ||
      (element.Name() != node && element.PropertyName != node)
  case shimast.KindVariableDeclaration:
    declaration := parent.AsVariableDeclaration()
    return declaration == nil || declaration.Name() != node
  case shimast.KindParameter:
    parameter := parent.AsParameterDeclaration()
    return parameter == nil || parameter.Name() != node
  case shimast.KindFunctionDeclaration,
    shimast.KindFunctionExpression,
    shimast.KindClassDeclaration,
    shimast.KindClassExpression,
    shimast.KindMethodDeclaration,
    shimast.KindPropertyDeclaration,
    shimast.KindGetAccessor,
    shimast.KindSetAccessor,
    shimast.KindEnumDeclaration,
    shimast.KindEnumMember,
    shimast.KindModuleDeclaration:
    return parent.Name() != node
  case shimast.KindLabeledStatement:
    statement := parent.AsLabeledStatement()
    return statement == nil || statement.Label != node
  case shimast.KindBreakStatement:
    statement := parent.AsBreakStatement()
    return statement == nil || statement.Label != node
  case shimast.KindContinueStatement:
    statement := parent.AsContinueStatement()
    return statement == nil || statement.Label != node
  }
  return true
}

// literalTruthiness folds a loop test into a constant when it is a simple
// literal, mirroring ESLint's `getBooleanValueIfSimpleConstant` (ESTree
// `Literal` nodes only — template literals and negations are NOT folded).
// Returns (false, false) for non-literal or unparseable expressions.
func literalTruthiness(expr *shimast.Node) (constTrue, constFalse bool) {
  expr = stripParens(expr)
  if expr == nil {
    return false, false
  }
  switch expr.Kind {
  case shimast.KindTrueKeyword:
    return true, false
  case shimast.KindFalseKeyword, shimast.KindNullKeyword:
    return false, true
  case shimast.KindRegularExpressionLiteral:
    return true, false
  case shimast.KindStringLiteral:
    if lit := expr.AsStringLiteral(); lit != nil {
      return lit.Text != "", lit.Text == ""
    }
  case shimast.KindNumericLiteral:
    if lit := expr.AsNumericLiteral(); lit != nil {
      if value, err := strconv.ParseFloat(lit.Text, 64); err == nil {
        return value != 0, value == 0
      }
    }
  case shimast.KindBigIntLiteral:
    if lit := expr.AsBigIntLiteral(); lit != nil {
      return bigIntLiteralTruthiness(lit.Text)
    }
  }
  return false, false
}

// bigIntLiteralTruthiness determines zero/non-zero without converting the
// value to a fixed-width integer. The scanner can preserve any accepted
// radix prefix, digit separators, and arbitrarily wide values.
func bigIntLiteralTruthiness(text string) (constTrue, constFalse bool) {
  if !strings.HasSuffix(text, "n") {
    return false, false
  }
  digits := strings.ReplaceAll(strings.TrimSuffix(text, "n"), "_", "")
  base := byte(10)
  if len(digits) >= 2 && digits[0] == '0' {
    switch digits[1] {
    case 'b', 'B':
      base, digits = 2, digits[2:]
    case 'o', 'O':
      base, digits = 8, digits[2:]
    case 'x', 'X':
      base, digits = 16, digits[2:]
    }
  }
  if digits == "" {
    return false, false
  }
  nonZero := false
  for i := 0; i < len(digits); i++ {
    digit := digits[i]
    value := byte(255)
    switch {
    case digit >= '0' && digit <= '9':
      value = digit - '0'
    case digit >= 'a' && digit <= 'f':
      value = digit - 'a' + 10
    case digit >= 'A' && digit <= 'F':
      value = digit - 'A' + 10
    }
    if value >= base {
      return false, false
    }
    nonZero = nonZero || value != 0
  }
  return nonZero, !nonZero
}

func init() {
  Register(noFallthrough{})
}
