/**
 * Fixture for `boundaries/external`.
 *
 * The rule restricts external package imports by package/specifier pattern.
 * It needs project-level `disallow` config to fire; the Go test under
 * `packages/lint/test/rules/boundaries/` pins the actual contract.
 *
 * This fixture documents the rule id in the consumer corpus tree. It declares
 * no `// expect:` annotations, so the corpus runner skips it.
 */
export {};
