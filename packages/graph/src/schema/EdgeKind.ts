/**
 * The relationship a directed edge encodes between two {@link IGraphNode}s.
 *
 * Structural edges (`contains`, `exports`, `imports`) come from the declaration
 * pass. Value and type edges (`calls`, `accesses`, `instantiates`, `type_ref`,
 * `extends`, `implements`, `overrides`) are resolved by the checker.
 * `decorates` carries the decorator fact a framework pass reads. `renders`,
 * `handles_route`, and `tests` are higher-level relationships a framework or
 * convention pass adds.
 *
 * Every edge is tagged with a {@link Provenance} and {@link Confidence}, so a
 * consumer can separate checker-resolved fact from framework or heuristic
 * inference regardless of kind.
 */
export type EdgeKind =
  | "contains"
  | "exports"
  | "imports"
  | "calls"
  | "accesses"
  | "instantiates"
  | "type_ref"
  | "extends"
  | "implements"
  | "overrides"
  | "decorates"
  | "renders"
  | "handles_route"
  | "tests";
