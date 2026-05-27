/**
 * Fixture for `boundaries/element-types`.
 *
 * The rule enforces allowed dependency directions between configured source-
 * path element types. It needs project-level `elements` and `rules` config to
 * fire, so a single virtual fixture cannot produce a diagnostic; the Go test
 * under `packages/lint/test/rules/boundaries/` pins the actual contract.
 *
 * This fixture documents the rule id in the consumer corpus tree. It declares
 * no `// expect:` annotations, so the corpus runner skips it.
 */
export {};
