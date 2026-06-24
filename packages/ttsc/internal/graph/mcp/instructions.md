`ttsc-graph` is the TypeScript compiler's relationship graph: symbols, calls, callers, source, types.

- For code-flow questions, call `graph_explore` before grep/read/shell.
- Query concrete domain nouns or symbols; drop generic words like code, method, request, main.
- If results are broad/generic, re-query once with the best returned symbol or file.
- Answer from graph; read only for no match, signature-only, or non-TS files.
- Use `graph_diagnostics` for file errors.
