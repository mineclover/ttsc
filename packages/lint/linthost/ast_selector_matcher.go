package linthost

import (
  "math"
  "strconv"
  "strings"

  shimast "github.com/microsoft/typescript-go/shim/ast"
  shimscanner "github.com/microsoft/typescript-go/shim/scanner"
)

// matchASTSelector returns every node selected below root. A configured
// selector reports each node once even when multiple alternatives or subject
// paths reach it.
func matchASTSelector(root *shimast.Node, selector *astSelector) []*shimast.Node {
  if root == nil || selector == nil {
    return nil
  }
  subjects := astSelectorSubjects(selector, selector)
  selected := make(map[*shimast.Node]struct{})
  ordered := make([]*shimast.Node, 0)
  walkDescendants(root, func(node *shimast.Node) {
    if !astSelectorMatchesNode(node, selector, nil) {
      return
    }
    if len(subjects) == 0 {
      if _, duplicate := selected[node]; !duplicate {
        selected[node] = struct{}{}
        ordered = append(ordered, node)
      }
      return
    }
    for candidate := node; candidate != nil; candidate = selectorParent(candidate) {
      for _, subject := range subjects {
        if !astSelectorMatchesNode(candidate, subject, nil) {
          continue
        }
        if _, duplicate := selected[candidate]; !duplicate {
          selected[candidate] = struct{}{}
          ordered = append(ordered, candidate)
        }
        break
      }
    }
  })
  return ordered
}

func astSelectorSubjects(selector, ancestor *astSelector) []*astSelector {
  if selector == nil {
    return nil
  }
  subjects := make([]*astSelector, 0)
  if selector.subject {
    subjects = append(subjects, ancestor)
  }
  if selector.left != nil {
    subjects = append(subjects, astSelectorSubjects(selector.left, selector.left)...)
  }
  if selector.right != nil {
    subjects = append(subjects, astSelectorSubjects(selector.right, ancestor)...)
  }
  for _, child := range selector.selectors {
    subjects = append(subjects, astSelectorSubjects(child, ancestor)...)
  }
  return subjects
}

func astSelectorMatchesNode(node *shimast.Node, selector *astSelector, scopeRoot *shimast.Node) bool {
  if node == nil || selector == nil {
    return false
  }
  switch selector.kind {
  case astSelectorWildcard:
    return true
  case astSelectorNodeType:
    return astSelectorMatchesNodeType(node, selector.name)
  case astSelectorExactNode:
    return scopeRoot != nil && node == scopeRoot
  case astSelectorAttribute:
    return astSelectorMatchesAttribute(node, selector)
  case astSelectorField:
    return astSelectorMatchesField(node, selector.name)
  case astSelectorMatches:
    for _, alternative := range selector.selectors {
      if astSelectorMatchesNode(node, alternative, scopeRoot) {
        return true
      }
    }
    return false
  case astSelectorCompound:
    for _, part := range selector.selectors {
      if !astSelectorMatchesNode(node, part, scopeRoot) {
        return false
      }
    }
    return true
  case astSelectorNot:
    for _, excluded := range selector.selectors {
      if astSelectorMatchesNode(node, excluded, scopeRoot) {
        return false
      }
    }
    return true
  case astSelectorHas:
    found := false
    walkDescendants(node, func(candidate *shimast.Node) {
      if found {
        return
      }
      for _, nested := range selector.selectors {
        if astSelectorMatchesNode(candidate, nested, node) {
          found = true
          return
        }
      }
    })
    return found
  case astSelectorChild:
    parent := selectorParent(node)
    return parent != nil &&
      astSelectorMatchesNode(node, selector.right, scopeRoot) &&
      astSelectorMatchesNode(parent, selector.left, scopeRoot)
  case astSelectorDescendant:
    if !astSelectorMatchesNode(node, selector.right, scopeRoot) {
      return false
    }
    for parent := selectorParent(node); parent != nil; parent = selectorParent(parent) {
      if astSelectorMatchesNode(parent, selector.left, scopeRoot) {
        return true
      }
      if scopeRoot != nil && parent == scopeRoot {
        break
      }
    }
    return false
  case astSelectorSibling:
    if astSelectorMatchesNode(node, selector.right, scopeRoot) {
      for _, sibling := range selectorSiblings(node, false) {
        if astSelectorMatchesNode(sibling, selector.left, scopeRoot) {
          return true
        }
      }
    }
    if selector.left != nil && selector.left.subject && astSelectorMatchesNode(node, selector.left, scopeRoot) {
      for _, sibling := range selectorSiblings(node, true) {
        if astSelectorMatchesNode(sibling, selector.right, scopeRoot) {
          return true
        }
      }
    }
    return false
  case astSelectorAdjacent:
    if astSelectorMatchesNode(node, selector.right, scopeRoot) {
      if sibling := selectorAdjacent(node, false); sibling != nil {
        return astSelectorMatchesNode(sibling, selector.left, scopeRoot)
      }
    }
    if selector.left != nil && selector.left.subject && astSelectorMatchesNode(node, selector.left, scopeRoot) {
      if sibling := selectorAdjacent(node, true); sibling != nil {
        return astSelectorMatchesNode(sibling, selector.right, scopeRoot)
      }
    }
    return false
  case astSelectorNthChild:
    return selectorChildIndex(node, false) == selector.index
  case astSelectorNthLastChild:
    return selectorChildIndex(node, true) == selector.index
  case astSelectorClass:
    return astSelectorMatchesClass(node, selector.name)
  }
  return false
}

func selectorParent(node *shimast.Node) *shimast.Node {
  if node == nil {
    return nil
  }
  return node.Parent
}

func selectorChildren(node *shimast.Node) []*shimast.Node {
  if node == nil {
    return nil
  }
  children := make([]*shimast.Node, 0)
  node.ForEachChild(func(child *shimast.Node) bool {
    if child != nil {
      children = append(children, child)
    }
    return false
  })
  return children
}

func selectorSiblings(node *shimast.Node, following bool) []*shimast.Node {
  parent := selectorParent(node)
  if parent == nil {
    return nil
  }
  children := selectorChildren(parent)
  index := -1
  for i, child := range children {
    if child == node {
      index = i
      break
    }
  }
  if index < 0 {
    return nil
  }
  if following {
    return children[index+1:]
  }
  return children[:index]
}

func selectorAdjacent(node *shimast.Node, following bool) *shimast.Node {
  siblings := selectorSiblings(node, following)
  if len(siblings) == 0 {
    return nil
  }
  if following {
    return siblings[0]
  }
  return siblings[len(siblings)-1]
}

func selectorChildIndex(node *shimast.Node, fromEnd bool) int {
  parent := selectorParent(node)
  if parent == nil {
    return 0
  }
  children := selectorChildren(parent)
  for index, child := range children {
    if child != node {
      continue
    }
    if fromEnd {
      return len(children) - index
    }
    return index + 1
  }
  return 0
}

func astSelectorMatchesNodeType(node *shimast.Node, requested string) bool {
  if node == nil {
    return false
  }
  native := strings.TrimPrefix(node.Kind.String(), "Kind")
  if strings.EqualFold(native, requested) {
    return true
  }
  switch strings.ToLower(requested) {
  case "program":
    return node.Kind == shimast.KindSourceFile
  case "arrowfunctionexpression":
    return node.Kind == shimast.KindArrowFunction
  case "objectexpression":
    return node.Kind == shimast.KindObjectLiteralExpression
  case "arrayexpression":
    return node.Kind == shimast.KindArrayLiteralExpression
  case "objectpattern":
    return node.Kind == shimast.KindObjectBindingPattern
  case "arraypattern":
    return node.Kind == shimast.KindArrayBindingPattern
  case "variabledeclarator":
    return node.Kind == shimast.KindVariableDeclaration
  case "memberexpression":
    return node.Kind == shimast.KindPropertyAccessExpression || node.Kind == shimast.KindElementAccessExpression
  case "assignmentexpression":
    if node.Kind != shimast.KindBinaryExpression {
      return false
    }
    expression := node.AsBinaryExpression()
    return expression != nil && expression.OperatorToken != nil && isAssignmentOperator(expression.OperatorToken.Kind)
  case "logicalexpression":
    if node.Kind != shimast.KindBinaryExpression {
      return false
    }
    expression := node.AsBinaryExpression()
    if expression == nil || expression.OperatorToken == nil {
      return false
    }
    switch expression.OperatorToken.Kind {
    case shimast.KindAmpersandAmpersandToken, shimast.KindBarBarToken, shimast.KindQuestionQuestionToken:
      return true
    }
    return false
  case "updateexpression":
    if node.Kind == shimast.KindPrefixUnaryExpression {
      expression := node.AsPrefixUnaryExpression()
      return expression != nil && (expression.Operator == shimast.KindPlusPlusToken || expression.Operator == shimast.KindMinusMinusToken)
    }
    if node.Kind == shimast.KindPostfixUnaryExpression {
      expression := node.AsPostfixUnaryExpression()
      return expression != nil && (expression.Operator == shimast.KindPlusPlusToken || expression.Operator == shimast.KindMinusMinusToken)
    }
    return false
  case "unaryexpression":
    if node.Kind == shimast.KindPrefixUnaryExpression {
      expression := node.AsPrefixUnaryExpression()
      return expression != nil && expression.Operator != shimast.KindPlusPlusToken && expression.Operator != shimast.KindMinusMinusToken
    }
    return node.Kind == shimast.KindDeleteExpression ||
      node.Kind == shimast.KindTypeOfExpression ||
      node.Kind == shimast.KindVoidExpression
  case "literal":
    return astSelectorIsLiteral(node)
  case "thisexpression":
    return node.Kind == shimast.KindThisKeyword
  case "super":
    return node.Kind == shimast.KindSuperKeyword
  case "templateliteral":
    return node.Kind == shimast.KindTemplateExpression || node.Kind == shimast.KindNoSubstitutionTemplateLiteral
  case "property":
    switch node.Kind {
    case shimast.KindPropertyAssignment,
      shimast.KindShorthandPropertyAssignment,
      shimast.KindSpreadAssignment:
      return true
    case shimast.KindMethodDeclaration,
      shimast.KindGetAccessor,
      shimast.KindSetAccessor:
      return node.Parent != nil && node.Parent.Kind == shimast.KindObjectLiteralExpression
    }
  case "restelement":
    return node.Kind == shimast.KindSpreadElement || node.Kind == shimast.KindSpreadAssignment ||
      node.Kind == shimast.KindBindingElement && node.AsBindingElement() != nil && node.AsBindingElement().DotDotDotToken != nil
  case "tsasexpression":
    return node.Kind == shimast.KindAsExpression
  case "tstypeassertion":
    return node.Kind == shimast.KindTypeAssertionExpression
  case "tssatisfiesexpression":
    return node.Kind == shimast.KindSatisfiesExpression
  case "tsnonnullexpression":
    return node.Kind == shimast.KindNonNullExpression
  }
  return false
}

func astSelectorIsLiteral(node *shimast.Node) bool {
  if node == nil {
    return false
  }
  switch node.Kind {
  case shimast.KindNumericLiteral,
    shimast.KindBigIntLiteral,
    shimast.KindStringLiteral,
    shimast.KindRegularExpressionLiteral,
    shimast.KindNoSubstitutionTemplateLiteral,
    shimast.KindTrueKeyword,
    shimast.KindFalseKeyword,
    shimast.KindNullKeyword:
    return true
  }
  return false
}

func astSelectorMatchesClass(node *shimast.Node, class string) bool {
  switch strings.ToLower(class) {
  case "statement":
    native := strings.TrimPrefix(node.Kind.String(), "Kind")
    return strings.HasSuffix(native, "Statement") || strings.HasSuffix(native, "Declaration")
  case "expression":
    return astSelectorIsExpression(node)
  case "declaration":
    return astSelectorIsDeclaration(node)
  case "function":
    return isFunctionLikeKind(node)
  case "pattern":
    return astSelectorIsPattern(node) || astSelectorIsExpression(node)
  default:
    return false
  }
}

func astSelectorIsExpression(node *shimast.Node) bool {
  if node == nil {
    return false
  }
  native := strings.TrimPrefix(node.Kind.String(), "Kind")
  if strings.HasSuffix(native, "Expression") || strings.HasSuffix(native, "Literal") {
    return true
  }
  switch node.Kind {
  case shimast.KindIdentifier:
    return node.Parent == nil || node.Parent.Kind != shimast.KindMetaProperty
  case shimast.KindPrivateIdentifier,
    shimast.KindMetaProperty,
    shimast.KindThisKeyword,
    shimast.KindSuperKeyword,
    shimast.KindArrowFunction:
    return true
  }
  return false
}

func astSelectorIsDeclaration(node *shimast.Node) bool {
  if node == nil {
    return false
  }
  return strings.HasSuffix(strings.TrimPrefix(node.Kind.String(), "Kind"), "Declaration")
}

func astSelectorIsPattern(node *shimast.Node) bool {
  if node == nil {
    return false
  }
  return strings.HasSuffix(strings.TrimPrefix(node.Kind.String(), "Kind"), "Pattern")
}

type astSelectorRuntimeValueKind uint8

const (
  astSelectorRuntimeMissing astSelectorRuntimeValueKind = iota
  astSelectorRuntimeNull
  astSelectorRuntimeString
  astSelectorRuntimeNumber
  astSelectorRuntimeBigInt
  astSelectorRuntimeBoolean
  astSelectorRuntimeObject
  astSelectorRuntimeNode
  astSelectorRuntimeNodes
)

type astSelectorRuntimeValue struct {
  kind    astSelectorRuntimeValueKind
  text    string
  number  float64
  boolean bool
  node    *shimast.Node
  nodes   []*shimast.Node
}

func astSelectorMatchesAttribute(node *shimast.Node, selector *astSelector) bool {
  value := astSelectorPathValue(astSelectorRuntimeValue{kind: astSelectorRuntimeNode, node: node}, strings.Split(selector.name, "."))
  if selector.operator == "" {
    return value.kind != astSelectorRuntimeMissing && value.kind != astSelectorRuntimeNull
  }
  matched := false
  switch selector.value.kind {
  case astSelectorValueRegexp:
    if selector.operator == "!=" {
      return selector.value.regexp == nil || !selector.value.regexp.MatchString(astSelectorRuntimeStringValue(value))
    }
    matched = value.kind == astSelectorRuntimeString && selector.value.regexp != nil && selector.value.regexp.MatchString(value.text)
  case astSelectorValueType:
    matched = selector.value.literal == astSelectorRuntimeType(value)
  case astSelectorValueLiteral:
    if selector.operator == "=" || selector.operator == "!=" {
      matched = astSelectorRuntimeStringValue(value) == selector.value.literal
    } else {
      return astSelectorCompare(value, selector.value, selector.operator)
    }
  default:
    return false
  }
  if selector.operator == "!=" {
    return !matched
  }
  return matched
}

func astSelectorCompare(value astSelectorRuntimeValue, expected astSelectorValue, operator string) bool {
  if value.kind == astSelectorRuntimeString && expected.number == nil {
    switch operator {
    case ">":
      return value.text > expected.literal
    case ">=":
      return value.text >= expected.literal
    case "<":
      return value.text < expected.literal
    case "<=":
      return value.text <= expected.literal
    }
    return false
  }
  left, leftNumber := astSelectorRuntimeNumberValue(value)
  right := 0.0
  rightNumber := false
  if expected.number != nil {
    right = *expected.number
    rightNumber = true
  } else if parsed, err := strconv.ParseFloat(expected.literal, 64); err == nil {
    right = parsed
    rightNumber = true
  }
  if !leftNumber || !rightNumber || math.IsNaN(left) || math.IsNaN(right) {
    return false
  }
  switch operator {
  case ">":
    return left > right
  case ">=":
    return left >= right
  case "<":
    return left < right
  case "<=":
    return left <= right
  }
  return false
}

func astSelectorRuntimeType(value astSelectorRuntimeValue) string {
  switch value.kind {
  case astSelectorRuntimeString:
    return "string"
  case astSelectorRuntimeNumber:
    return "number"
  case astSelectorRuntimeBigInt:
    return "bigint"
  case astSelectorRuntimeBoolean:
    return "boolean"
  case astSelectorRuntimeObject, astSelectorRuntimeNode, astSelectorRuntimeNodes, astSelectorRuntimeNull:
    return "object"
  default:
    return "undefined"
  }
}

func astSelectorRuntimeStringValue(value astSelectorRuntimeValue) string {
  switch value.kind {
  case astSelectorRuntimeString:
    return value.text
  case astSelectorRuntimeNumber:
    return strconv.FormatFloat(value.number, 'g', -1, 64)
  case astSelectorRuntimeBigInt:
    return value.text
  case astSelectorRuntimeBoolean:
    return strconv.FormatBool(value.boolean)
  case astSelectorRuntimeNull:
    return "null"
  case astSelectorRuntimeMissing:
    return "undefined"
  case astSelectorRuntimeNodes:
    if len(value.nodes) == 0 {
      return ""
    }
    return strings.TrimSuffix(strings.Repeat("[object Object],", len(value.nodes)), ",")
  }
  return "[object Object]"
}

func astSelectorRuntimeNumberValue(value astSelectorRuntimeValue) (float64, bool) {
  switch value.kind {
  case astSelectorRuntimeNumber:
    return value.number, true
  case astSelectorRuntimeString:
    parsed, err := strconv.ParseFloat(value.text, 64)
    return parsed, err == nil
  case astSelectorRuntimeBoolean:
    if value.boolean {
      return 1, true
    }
    return 0, true
  case astSelectorRuntimeNull:
    return 0, true
  }
  return 0, false
}

func astSelectorPathValue(value astSelectorRuntimeValue, path []string) astSelectorRuntimeValue {
  current := value
  for _, field := range path {
    switch current.kind {
    case astSelectorRuntimeNode:
      current = astSelectorNodeProperty(current.node, field)
    case astSelectorRuntimeNodes:
      if field == "length" {
        current = astSelectorRuntimeValue{kind: astSelectorRuntimeNumber, number: float64(len(current.nodes))}
        continue
      }
      index, err := strconv.Atoi(field)
      if err != nil || index < 0 || index >= len(current.nodes) {
        return astSelectorRuntimeValue{}
      }
      current = astSelectorNode(current.nodes[index])
    case astSelectorRuntimeString:
      if field != "length" {
        return astSelectorRuntimeValue{}
      }
      current = astSelectorRuntimeValue{kind: astSelectorRuntimeNumber, number: float64(len([]rune(current.text)))}
    default:
      return astSelectorRuntimeValue{}
    }
  }
  return current
}

func astSelectorMatchesField(node *shimast.Node, name string) bool {
  path := strings.Split(name, ".")
  ancestor := node
  for range path {
    ancestor = selectorParent(ancestor)
    if ancestor == nil {
      return false
    }
  }
  return astSelectorNodeInPath(node, astSelectorRuntimeValue{kind: astSelectorRuntimeNode, node: ancestor}, path)
}

func astSelectorNodeInPath(node *shimast.Node, value astSelectorRuntimeValue, path []string) bool {
  if len(path) == 0 {
    return value.kind == astSelectorRuntimeNode && value.node == node
  }
  if value.kind != astSelectorRuntimeNode {
    return false
  }
  next := astSelectorNodeProperty(value.node, path[0])
  if next.kind == astSelectorRuntimeNodes {
    for _, child := range next.nodes {
      if astSelectorNodeInPath(node, astSelectorRuntimeValue{kind: astSelectorRuntimeNode, node: child}, path[1:]) {
        return true
      }
    }
    return false
  }
  return astSelectorNodeInPath(node, next, path[1:])
}

func astSelectorNodeProperty(node *shimast.Node, field string) astSelectorRuntimeValue {
  if node == nil {
    return astSelectorRuntimeValue{kind: astSelectorRuntimeNull}
  }
  switch field {
  case "type":
    return astSelectorString(strings.TrimPrefix(node.Kind.String(), "Kind"))
  case "name":
    if name := identifierText(node); name != "" {
      return astSelectorString(name)
    }
    if node.Kind == shimast.KindPrivateIdentifier {
      return astSelectorString(node.AsPrivateIdentifier().Text)
    }
    nameNode := node.Name()
    if name := identifierText(nameNode); name != "" {
      return astSelectorString(name)
    }
    if nameNode != nil && nameNode.Kind == shimast.KindPrivateIdentifier {
      return astSelectorString(nameNode.AsPrivateIdentifier().Text)
    }
    if value := astSelectorLiteralValue(nameNode); value.kind != astSelectorRuntimeMissing {
      return value
    }
    return astSelectorRuntimeValue{kind: astSelectorRuntimeNull}
  case "value":
    if node.Kind == shimast.KindPropertyAssignment {
      return astSelectorNode(node.AsPropertyAssignment().Initializer)
    }
    if node.Kind == shimast.KindShorthandPropertyAssignment {
      return astSelectorNode(node.AsShorthandPropertyAssignment().Name())
    }
    if node.Kind == shimast.KindBindingElement {
      return astSelectorNode(node.Name())
    }
    return astSelectorLiteralValue(node)
  case "raw":
    if astSelectorIsLiteral(node) {
      if file := shimast.GetSourceFileOfNode(node); file != nil {
        return astSelectorString(nodeText(file, node))
      }
    }
    return astSelectorRuntimeValue{}
  case "operator":
    if operator := astSelectorNodeOperator(node); operator != "" {
      return astSelectorString(operator)
    }
    return astSelectorRuntimeValue{}
  case "kind":
    if kind := astSelectorVariableKind(node); kind != "" {
      return astSelectorString(kind)
    }
    return astSelectorRuntimeValue{}
  case "async":
    if node.FunctionLikeData() == nil {
      return astSelectorRuntimeValue{}
    }
    return astSelectorBoolean(hasModifier(node, shimast.KindAsyncKeyword))
  case "generator":
    if node.FunctionLikeData() == nil {
      return astSelectorRuntimeValue{}
    }
    return astSelectorBoolean(astSelectorIsGenerator(node))
  case "static":
    if !astSelectorSupportsStatic(node) {
      return astSelectorRuntimeValue{}
    }
    if node.Kind == shimast.KindClassStaticBlockDeclaration {
      return astSelectorBoolean(true)
    }
    return astSelectorBoolean(hasModifier(node, shimast.KindStaticKeyword))
  case "readonly":
    if !astSelectorSupportsReadonly(node) {
      return astSelectorRuntimeValue{}
    }
    return astSelectorBoolean(hasModifier(node, shimast.KindReadonlyKeyword))
  case "declare":
    if !astSelectorIsDeclaration(node) {
      return astSelectorRuntimeValue{}
    }
    return astSelectorBoolean(hasModifier(node, shimast.KindDeclareKeyword))
  case "optional":
    optional, applicable := astSelectorIsOptional(node)
    if !applicable {
      return astSelectorRuntimeValue{}
    }
    return astSelectorBoolean(optional)
  case "computed":
    computed, applicable := astSelectorIsComputed(node)
    if !applicable {
      return astSelectorRuntimeValue{}
    }
    return astSelectorBoolean(computed)
  case "prefix":
    switch node.Kind {
    case shimast.KindPrefixUnaryExpression,
      shimast.KindDeleteExpression,
      shimast.KindTypeOfExpression,
      shimast.KindVoidExpression:
      return astSelectorBoolean(true)
    case shimast.KindPostfixUnaryExpression:
      return astSelectorBoolean(false)
    }
    return astSelectorRuntimeValue{}
  case "id":
    return astSelectorNode(node.Name())
  case "params":
    if function := node.FunctionLikeData(); function != nil && function.Parameters != nil {
      return astSelectorNodes(function.Parameters.Nodes)
    }
  case "body":
    if body := node.Body(); body != nil {
      return astSelectorNode(body)
    }
    return astSelectorStatementBody(node)
  case "callee":
    if node.Kind == shimast.KindCallExpression {
      if call := node.AsCallExpression(); call != nil {
        return astSelectorNode(call.Expression)
      }
    }
    if node.Kind == shimast.KindNewExpression {
      if expression := node.AsNewExpression(); expression != nil {
        return astSelectorNode(expression.Expression)
      }
    }
  case "arguments":
    if node.Kind == shimast.KindCallExpression {
      if call := node.AsCallExpression(); call != nil && call.Arguments != nil {
        return astSelectorNodes(call.Arguments.Nodes)
      }
    }
    if node.Kind == shimast.KindNewExpression {
      if expression := node.AsNewExpression(); expression != nil && expression.Arguments != nil {
        return astSelectorNodes(expression.Arguments.Nodes)
      }
    }
  case "object":
    if node.Kind == shimast.KindPropertyAccessExpression {
      return astSelectorNode(node.AsPropertyAccessExpression().Expression)
    }
    if node.Kind == shimast.KindElementAccessExpression {
      return astSelectorNode(node.AsElementAccessExpression().Expression)
    }
  case "property":
    if node.Kind == shimast.KindPropertyAccessExpression {
      return astSelectorNode(node.AsPropertyAccessExpression().Name())
    }
    if node.Kind == shimast.KindElementAccessExpression {
      return astSelectorNode(node.AsElementAccessExpression().ArgumentExpression)
    }
  case "left":
    if node.Kind == shimast.KindBinaryExpression {
      return astSelectorNode(node.AsBinaryExpression().Left)
    }
    if node.Kind == shimast.KindQualifiedName {
      return astSelectorNode(node.AsQualifiedName().Left)
    }
  case "right":
    if node.Kind == shimast.KindBinaryExpression {
      return astSelectorNode(node.AsBinaryExpression().Right)
    }
    if node.Kind == shimast.KindQualifiedName {
      return astSelectorNode(node.AsQualifiedName().Right)
    }
  case "argument":
    return astSelectorUnaryArgument(node)
  case "expression":
    return astSelectorExpressionChild(node)
  case "test":
    return astSelectorTestChild(node)
  case "consequent":
    if node.Kind == shimast.KindIfStatement {
      return astSelectorNode(node.AsIfStatement().ThenStatement)
    }
    if node.Kind == shimast.KindConditionalExpression {
      return astSelectorNode(node.AsConditionalExpression().WhenTrue)
    }
    if node.Kind == shimast.KindCaseClause || node.Kind == shimast.KindDefaultClause {
      clause := node.AsCaseOrDefaultClause()
      if clause != nil && clause.Statements != nil {
        return astSelectorNodes(clause.Statements.Nodes)
      }
    }
  case "alternate":
    if node.Kind == shimast.KindIfStatement {
      return astSelectorNode(node.AsIfStatement().ElseStatement)
    }
    if node.Kind == shimast.KindConditionalExpression {
      return astSelectorNode(node.AsConditionalExpression().WhenFalse)
    }
  case "init":
    if node.Kind == shimast.KindVariableDeclaration {
      return astSelectorNode(node.AsVariableDeclaration().Initializer)
    }
    if node.Kind == shimast.KindForStatement {
      return astSelectorNode(node.AsForStatement().Initializer)
    }
  case "update":
    if node.Kind == shimast.KindForStatement {
      return astSelectorNode(node.AsForStatement().Incrementor)
    }
  case "declarations":
    if node.Kind == shimast.KindVariableDeclarationList {
      list := node.AsVariableDeclarationList()
      if list != nil && list.Declarations != nil {
        return astSelectorNodes(list.Declarations.Nodes)
      }
    }
    if node.Kind == shimast.KindVariableStatement {
      statement := node.AsVariableStatement()
      if statement != nil && statement.DeclarationList != nil {
        list := statement.DeclarationList.AsVariableDeclarationList()
        if list != nil && list.Declarations != nil {
          return astSelectorNodes(list.Declarations.Nodes)
        }
      }
    }
  case "elements":
    if node.Kind == shimast.KindArrayLiteralExpression {
      expression := node.AsArrayLiteralExpression()
      if expression != nil && expression.Elements != nil {
        return astSelectorNodes(expression.Elements.Nodes)
      }
    }
    if node.Kind == shimast.KindObjectBindingPattern || node.Kind == shimast.KindArrayBindingPattern {
      pattern := node.AsBindingPattern()
      if pattern != nil && pattern.Elements != nil {
        return astSelectorNodes(pattern.Elements.Nodes)
      }
    }
  case "properties":
    if node.Kind == shimast.KindObjectLiteralExpression {
      expression := node.AsObjectLiteralExpression()
      if expression != nil && expression.Properties != nil {
        return astSelectorNodes(expression.Properties.Nodes)
      }
    }
    if node.Kind == shimast.KindObjectBindingPattern {
      pattern := node.AsBindingPattern()
      if pattern != nil && pattern.Elements != nil {
        return astSelectorNodes(pattern.Elements.Nodes)
      }
    }
  case "key":
    return astSelectorPropertyKey(node)
  case "source":
    if node.Kind == shimast.KindImportDeclaration {
      return astSelectorNode(node.AsImportDeclaration().ModuleSpecifier)
    }
    if node.Kind == shimast.KindExportDeclaration {
      return astSelectorNode(node.AsExportDeclaration().ModuleSpecifier)
    }
  case "statements":
    return astSelectorStatements(node)
  }
  return astSelectorRuntimeValue{}
}

func astSelectorString(value string) astSelectorRuntimeValue {
  return astSelectorRuntimeValue{kind: astSelectorRuntimeString, text: value}
}

func astSelectorBoolean(value bool) astSelectorRuntimeValue {
  return astSelectorRuntimeValue{kind: astSelectorRuntimeBoolean, boolean: value}
}

func astSelectorLiteralValue(node *shimast.Node) astSelectorRuntimeValue {
  if node == nil {
    return astSelectorRuntimeValue{kind: astSelectorRuntimeNull}
  }
  switch node.Kind {
  case shimast.KindStringLiteral, shimast.KindNoSubstitutionTemplateLiteral:
    if data := node.LiteralLikeData(); data != nil {
      return astSelectorString(data.Text)
    }
  case shimast.KindRegularExpressionLiteral:
    return astSelectorRuntimeValue{kind: astSelectorRuntimeObject}
  case shimast.KindNumericLiteral:
    if data := node.LiteralLikeData(); data != nil {
      text := strings.ReplaceAll(data.Text, "_", "")
      if number, err := strconv.ParseFloat(text, 64); err == nil {
        return astSelectorRuntimeValue{kind: astSelectorRuntimeNumber, number: number}
      }
      if number, err := strconv.ParseUint(text, 0, 64); err == nil {
        return astSelectorRuntimeValue{kind: astSelectorRuntimeNumber, number: float64(number)}
      }
    }
  case shimast.KindBigIntLiteral:
    if data := node.LiteralLikeData(); data != nil {
      return astSelectorRuntimeValue{
        kind: astSelectorRuntimeBigInt,
        text: strings.TrimSuffix(strings.ReplaceAll(data.Text, "_", ""), "n"),
      }
    }
  case shimast.KindTrueKeyword:
    return astSelectorBoolean(true)
  case shimast.KindFalseKeyword:
    return astSelectorBoolean(false)
  case shimast.KindNullKeyword:
    return astSelectorRuntimeValue{kind: astSelectorRuntimeNull}
  }
  return astSelectorRuntimeValue{}
}

func astSelectorNode(node *shimast.Node) astSelectorRuntimeValue {
  if node == nil {
    return astSelectorRuntimeValue{kind: astSelectorRuntimeNull}
  }
  return astSelectorRuntimeValue{kind: astSelectorRuntimeNode, node: node}
}

func astSelectorNodes(nodes []*shimast.Node) astSelectorRuntimeValue {
  if nodes == nil {
    return astSelectorRuntimeValue{kind: astSelectorRuntimeNull}
  }
  return astSelectorRuntimeValue{kind: astSelectorRuntimeNodes, nodes: nodes}
}

func astSelectorNodeOperator(node *shimast.Node) string {
  if node == nil {
    return ""
  }
  kind := shimast.KindUnknown
  switch node.Kind {
  case shimast.KindBinaryExpression:
    if expression := node.AsBinaryExpression(); expression != nil && expression.OperatorToken != nil {
      kind = expression.OperatorToken.Kind
    }
  case shimast.KindPrefixUnaryExpression:
    kind = node.AsPrefixUnaryExpression().Operator
  case shimast.KindPostfixUnaryExpression:
    kind = node.AsPostfixUnaryExpression().Operator
  case shimast.KindDeleteExpression:
    return "delete"
  case shimast.KindTypeOfExpression:
    return "typeof"
  case shimast.KindVoidExpression:
    return "void"
  }
  return astSelectorOperatorText(kind)
}

func astSelectorOperatorText(kind shimast.Kind) string {
  return shimscanner.TokenToString(kind)
}

func astSelectorVariableKind(node *shimast.Node) string {
  if node == nil {
    return ""
  }
  list := node
  if node.Kind == shimast.KindVariableStatement {
    statement := node.AsVariableStatement()
    if statement == nil {
      return ""
    }
    list = statement.DeclarationList
  } else if node.Kind == shimast.KindVariableDeclaration {
    list = node.Parent
  }
  if list == nil || list.Kind != shimast.KindVariableDeclarationList {
    return ""
  }
  flags := list.Flags & shimast.NodeFlagsBlockScoped
  switch flags {
  case shimast.NodeFlagsLet:
    return "let"
  case shimast.NodeFlagsConst:
    return "const"
  case shimast.NodeFlagsUsing:
    return "using"
  case shimast.NodeFlagsAwaitUsing:
    return "await using"
  default:
    return "var"
  }
}

func astSelectorIsGenerator(node *shimast.Node) bool {
  if node == nil || !isFunctionLikeKind(node) {
    return false
  }
  body := node.BodyData()
  return body != nil && body.AsteriskToken != nil
}

func astSelectorIsOptional(node *shimast.Node) (bool, bool) {
  if node == nil {
    return false, false
  }
  switch node.Kind {
  case shimast.KindPropertyAccessExpression:
    return node.AsPropertyAccessExpression().QuestionDotToken != nil, true
  case shimast.KindElementAccessExpression:
    return node.AsElementAccessExpression().QuestionDotToken != nil, true
  case shimast.KindCallExpression:
    return node.AsCallExpression().QuestionDotToken != nil, true
  case shimast.KindTaggedTemplateExpression:
    return node.AsTaggedTemplateExpression().QuestionDotToken != nil, true
  case shimast.KindParameter,
    shimast.KindNamedTupleMember,
    shimast.KindMethodDeclaration,
    shimast.KindShorthandPropertyAssignment,
    shimast.KindMethodSignature,
    shimast.KindPropertySignature,
    shimast.KindPropertyAssignment,
    shimast.KindPropertyDeclaration,
    shimast.KindEnumMember,
    shimast.KindGetAccessor,
    shimast.KindSetAccessor:
    return node.QuestionToken() != nil, true
  }
  return false, false
}

func astSelectorSupportsStatic(node *shimast.Node) bool {
  if node == nil {
    return false
  }
  switch node.Kind {
  case shimast.KindClassStaticBlockDeclaration:
    return true
  case shimast.KindPropertyDeclaration,
    shimast.KindMethodDeclaration,
    shimast.KindGetAccessor,
    shimast.KindSetAccessor:
    return node.Parent != nil &&
      (node.Parent.Kind == shimast.KindClassDeclaration || node.Parent.Kind == shimast.KindClassExpression)
  }
  return false
}

func astSelectorSupportsReadonly(node *shimast.Node) bool {
  if node == nil {
    return false
  }
  switch node.Kind {
  case shimast.KindParameter,
    shimast.KindPropertyDeclaration,
    shimast.KindPropertySignature,
    shimast.KindIndexSignature:
    return true
  }
  return false
}

func astSelectorIsComputed(node *shimast.Node) (bool, bool) {
  if node == nil {
    return false, false
  }
  switch node.Kind {
  case shimast.KindElementAccessExpression, shimast.KindComputedPropertyName:
    return true, true
  case shimast.KindPropertyAccessExpression:
    return false, true
  case shimast.KindBindingElement:
    element := node.AsBindingElement()
    return element != nil && element.PropertyName != nil && element.PropertyName.Kind == shimast.KindComputedPropertyName, true
  case shimast.KindPropertyAssignment,
    shimast.KindShorthandPropertyAssignment,
    shimast.KindMethodDeclaration,
    shimast.KindMethodSignature,
    shimast.KindGetAccessor,
    shimast.KindSetAccessor,
    shimast.KindPropertyDeclaration,
    shimast.KindPropertySignature,
    shimast.KindEnumMember:
    name := node.Name()
    return name != nil && name.Kind == shimast.KindComputedPropertyName, true
  }
  return false, false
}

func astSelectorStatementBody(node *shimast.Node) astSelectorRuntimeValue {
  if node == nil {
    return astSelectorRuntimeValue{}
  }
  switch node.Kind {
  case shimast.KindSourceFile:
    file := node.AsSourceFile()
    if file != nil && file.Statements != nil {
      return astSelectorNodes(file.Statements.Nodes)
    }
  case shimast.KindBlock:
    block := node.AsBlock()
    if block != nil && block.Statements != nil {
      return astSelectorNodes(block.Statements.Nodes)
    }
  case shimast.KindModuleBlock:
    block := node.AsModuleBlock()
    if block != nil && block.Statements != nil {
      return astSelectorNodes(block.Statements.Nodes)
    }
  case shimast.KindIfStatement:
    return astSelectorNode(node.AsIfStatement().ThenStatement)
  case shimast.KindWhileStatement:
    return astSelectorNode(node.AsWhileStatement().Statement)
  case shimast.KindDoStatement:
    return astSelectorNode(node.AsDoStatement().Statement)
  case shimast.KindForStatement:
    return astSelectorNode(node.AsForStatement().Statement)
  case shimast.KindForInStatement, shimast.KindForOfStatement:
    return astSelectorNode(node.AsForInOrOfStatement().Statement)
  case shimast.KindWithStatement:
    return astSelectorNode(node.AsWithStatement().Statement)
  }
  return astSelectorRuntimeValue{}
}

func astSelectorUnaryArgument(node *shimast.Node) astSelectorRuntimeValue {
  if node == nil {
    return astSelectorRuntimeValue{}
  }
  switch node.Kind {
  case shimast.KindPrefixUnaryExpression:
    return astSelectorNode(node.AsPrefixUnaryExpression().Operand)
  case shimast.KindPostfixUnaryExpression:
    return astSelectorNode(node.AsPostfixUnaryExpression().Operand)
  case shimast.KindDeleteExpression:
    return astSelectorNode(node.AsDeleteExpression().Expression)
  case shimast.KindTypeOfExpression:
    return astSelectorNode(node.AsTypeOfExpression().Expression)
  case shimast.KindVoidExpression:
    return astSelectorNode(node.AsVoidExpression().Expression)
  case shimast.KindAwaitExpression:
    return astSelectorNode(node.AsAwaitExpression().Expression)
  case shimast.KindYieldExpression:
    return astSelectorNode(node.AsYieldExpression().Expression)
  case shimast.KindReturnStatement:
    return astSelectorNode(node.AsReturnStatement().Expression)
  case shimast.KindThrowStatement:
    return astSelectorNode(node.AsThrowStatement().Expression)
  case shimast.KindSpreadElement:
    return astSelectorNode(node.AsSpreadElement().Expression)
  case shimast.KindSpreadAssignment:
    return astSelectorNode(node.AsSpreadAssignment().Expression)
  case shimast.KindBindingElement:
    element := node.AsBindingElement()
    if element != nil && element.DotDotDotToken != nil {
      return astSelectorNode(element.Name())
    }
  }
  return astSelectorRuntimeValue{}
}

func astSelectorExpressionChild(node *shimast.Node) astSelectorRuntimeValue {
  if node == nil {
    return astSelectorRuntimeValue{}
  }
  switch node.Kind {
  case shimast.KindExpressionStatement:
    return astSelectorNode(node.AsExpressionStatement().Expression)
  case shimast.KindParenthesizedExpression:
    return astSelectorNode(node.AsParenthesizedExpression().Expression)
  case shimast.KindAsExpression:
    return astSelectorNode(node.AsAsExpression().Expression)
  case shimast.KindTypeAssertionExpression:
    return astSelectorNode(node.AsTypeAssertion().Expression)
  case shimast.KindSatisfiesExpression:
    return astSelectorNode(node.AsSatisfiesExpression().Expression)
  case shimast.KindNonNullExpression:
    return astSelectorNode(node.AsNonNullExpression().Expression)
  case shimast.KindAwaitExpression:
    return astSelectorNode(node.AsAwaitExpression().Expression)
  case shimast.KindYieldExpression:
    return astSelectorNode(node.AsYieldExpression().Expression)
  }
  return astSelectorRuntimeValue{}
}

func astSelectorTestChild(node *shimast.Node) astSelectorRuntimeValue {
  if node == nil {
    return astSelectorRuntimeValue{}
  }
  switch node.Kind {
  case shimast.KindIfStatement:
    return astSelectorNode(node.AsIfStatement().Expression)
  case shimast.KindWhileStatement:
    return astSelectorNode(node.AsWhileStatement().Expression)
  case shimast.KindDoStatement:
    return astSelectorNode(node.AsDoStatement().Expression)
  case shimast.KindForStatement:
    return astSelectorNode(node.AsForStatement().Condition)
  case shimast.KindConditionalExpression:
    return astSelectorNode(node.AsConditionalExpression().Condition)
  }
  return astSelectorRuntimeValue{}
}

func astSelectorPropertyKey(node *shimast.Node) astSelectorRuntimeValue {
  if node == nil {
    return astSelectorRuntimeValue{}
  }
  switch node.Kind {
  case shimast.KindPropertyAssignment:
    return astSelectorNode(node.AsPropertyAssignment().Name())
  case shimast.KindShorthandPropertyAssignment:
    return astSelectorNode(node.AsShorthandPropertyAssignment().Name())
  case shimast.KindBindingElement:
    element := node.AsBindingElement()
    if element.PropertyName != nil {
      return astSelectorNode(element.PropertyName)
    }
    return astSelectorNode(element.Name())
  }
  return astSelectorNode(node.Name())
}

func astSelectorStatements(node *shimast.Node) astSelectorRuntimeValue {
  if node == nil {
    return astSelectorRuntimeValue{}
  }
  switch node.Kind {
  case shimast.KindSourceFile:
    file := node.AsSourceFile()
    if file != nil && file.Statements != nil {
      return astSelectorNodes(file.Statements.Nodes)
    }
  case shimast.KindBlock:
    block := node.AsBlock()
    if block != nil && block.Statements != nil {
      return astSelectorNodes(block.Statements.Nodes)
    }
  case shimast.KindModuleBlock:
    block := node.AsModuleBlock()
    if block != nil && block.Statements != nil {
      return astSelectorNodes(block.Statements.Nodes)
    }
  case shimast.KindCaseClause, shimast.KindDefaultClause:
    clause := node.AsCaseOrDefaultClause()
    if clause != nil && clause.Statements != nil {
      return astSelectorNodes(clause.Statements.Nodes)
    }
  }
  return astSelectorRuntimeValue{}
}
