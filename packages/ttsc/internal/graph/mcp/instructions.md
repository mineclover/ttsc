`ttsc-graph` is the live TS compiler graph: symbols, source, calls, callers, types, diagnostics.

- Use `graph_explore` for TS architecture, code-flow, and impact questions.
- It returns relevant nodes, source, edges, and blast radius faster than broad file search.
- Query named symbols/files/domain nouns; broad onboarding may ask orientation once.
- Re-query for missing/narrower symbol/file, or after source edits.
- Avoid chasing every edge or re-reading returned source.
- Use diagnostics for file errors.
