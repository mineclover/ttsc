`ttsc-graph` is the TS compiler graph: symbols, source, calls, callers, types, diagnostics.

- Use `graph_explore` before broad grep/read for TS architecture, code-flow, and impact.
- It returns nodes, source, edges, and blast radius faster than file search.
- Query named symbols/files/domain nouns; broad onboarding may ask orientation once.
- Re-query for missing/narrower symbol/file, or after source edits.
- Avoid chasing every edge or re-reading returned source.
- Use diagnostics for file errors.
