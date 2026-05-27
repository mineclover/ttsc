/**
 * Fixture for `boundaries/entry-point`.
 *
 * The rule requires imports into an element to target its configured public
 * entry files. It needs project-level configuration to fire; the Go test under
 * `packages/lint/test/rules/boundaries/` pins the actual contract.
 *
 * This fixture documents the rule id in the consumer corpus tree. It declares
 * no `// expect:` annotations, so the corpus runner skips it.
 */
export {};
