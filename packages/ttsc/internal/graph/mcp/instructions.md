`ttsc-graph` is the TypeScript compiler's relationship graph: symbols, source, calls, callers, types.

- For code-flow questions, prefer `graph_explore` before broad grep/read; one query gives map plus source.
- Query flows, not one name: owner + action + nouns, e.g. `repository find manager query builder`.
- Do not fan out through files or grep to trace calls the graph shows.
- Do not re-read returned source; read only for no match, signatures, or non-TS files.
- Use `graph_diagnostics` for file errors.
