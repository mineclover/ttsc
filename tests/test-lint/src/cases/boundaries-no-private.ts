/**
 * Fixture for `boundaries/no-private`.
 *
 * The rule rejects imports of configured private files from outside their
 * element. It needs project-level configuration to fire; the Go test under
 * `packages/lint/test/rules/boundaries/` pins the actual contract.
 *
 * This fixture documents the rule id in the consumer corpus tree. It declares
 * no `// expect:` annotations, so the corpus runner skips it.
 */
export {};
