package checker

// Compile-time guard: the declared-type / base-chain walk that typia's
// plain.classify uses to detect an inherited `#private` field must compose
// end-to-end through the shim.
//
// #246 exposed a CONSUMER of the base chain (Checker_getBaseTypes) without the
// one PRODUCER that keeps the walk alive across a generic boundary
// (Checker_getDeclaredTypeOfSymbol): getBaseTypes nil-derefs on a generic
// `Reference` base such as the `Mid<T>` in `class Sub extends Mid<string>`
// (internally `t.AsInterfaceType()` is nil), so the walk dead-ends and a
// `#private` field inherited through that boundary stays invisible. classify
// then field-copies the class and crashes at runtime — undetected. Like the
// 0.15.5 Signature gap, the closure auditor cannot see this: getDeclaredType-
// OfSymbol is an UNEXPORTED *Checker method, and the auditor only fails on
// partial enums / reachable funcs / escaping types — never on consumer demand
// for an internal helper by name (it lists them in a non-failing triage pool).
//
// This lives in a NORMAL (non-_test) file for the same reason as
// signature_api_completeness.go: the shim/checker package is a nested Go module
// no CI job runs `go test` against, so a _test.go guard would never compile.
// As ordinary package code it is type-checked by every build that links the
// shim — typia's native plugin and ttsc's lint engine (the latter via
// `pnpm test:go`) — so dropping, renaming, or signature-breaking any leg below
// (e.g. an upstream typescript-go bump) turns those builds red. It never runs
// (assigned to blank, never called).
var _ = func(c *Checker, instanceType *Type) {
  // Walk the base chain of a class instance type. For a plain
  // ClassOrInterface base this is a direct getBaseTypes step.
  for _, base := range Checker_getBaseTypes(c, instanceType) {
    // The base may be a generic `Reference` (e.g. `Mid<string>`), on which
    // getBaseTypes would nil-deref. Resolve it to its declared (instance) type
    // first: base -> its name symbol -> declared type, which IS a
    // ClassOrInterface and is therefore safe to feed back into getBaseTypes.
    // This is the leg #246 adds; without it the walk dead-ends here.
    baseSymbol := Type_getTypeNameSymbol(base)
    declared := Checker_getDeclaredTypeOfSymbol(c, baseSymbol)

    // Continue the walk from the resolved declared type to reach an ancestor's
    // `#private` field through the generic boundary.
    for _, grand := range Checker_getBaseTypes(c, declared) {
      _ = grand
    }
  }
}
