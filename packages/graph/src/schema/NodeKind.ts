/**
 * What a graph node represents.
 *
 * The symbol kinds (`file` through `parameter`) are declarations the TypeScript
 * program owns and the checker resolves. `route` and `component` are _virtual_
 * nodes the framework pass synthesizes (provenance `framework-derived`); they
 * have no single declaring symbol. `external_symbol` is a dependency-boundary
 * leaf the workspace references but does not declare — the graph keeps it as a
 * named endpoint without expanding the dependency's internals.
 *
 * Used as the `kind` discriminant on {@link IGraphNode}.
 */
export type NodeKind =
  | "file"
  | "package"
  | "namespace"
  | "module"
  | "function"
  | "class"
  | "interface"
  | "type"
  | "enum"
  | "variable"
  | "method"
  | "property"
  | "parameter"
  | "route"
  | "component"
  | "external_symbol";
