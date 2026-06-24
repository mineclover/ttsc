Use `ttsc-graph` for TypeScript relationship questions: how an API, symbol, feature, diagnostic, request, or data path reaches its implementation. Start with `graph_explore`; ask one broad query naming owner, action, and domain nouns. Answer from returned source, edges, callers, and types. Use grep/read only for no match, omitted source, edited source, or non-TS context. Use `graph_diagnostics` for file errors.

`ttsc-graph` is a compiler-resolved relationship graph, not a text search index. It is best when the useful answer depends on calls, callers, type references, ownership, or blast radius rather than keyword occurrence.

Prefer one well-shaped query over several tiny probes. Include the public owner or subsystem when known, the action being traced, and nouns from the domain. Re-query when the returned graph reveals a better target or when a necessary symbol is missing.

Trust returned line-numbered source and relationships. Shell search is still appropriate for generated files, config, documentation, non-TypeScript assets, or source changed after the graph snapshot.
