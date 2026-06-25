# ttsc-graph

A compiler-resolved index of this TypeScript program: every declaration with its file and line, and the call and type relationships between them. Query it before grep or read. It has already located the symbols and connected them, so answer from its results: when a result names a declaration you need to read, expand that handle for the body instead of opening the file. Re-reading or grepping a file the graph has already mapped repeats work it did for you.

Pick the entry tool by what you already know:

- `query_exports`: orient. The exported public surface, for onboarding or an unclear entry point.
- `query_path`: an exact path, when you know the start and end symbols, or an ordered chain. It returns the ordered nodes, the runtime edges between them, and the off-path helpers each step calls, all with handles.
- `query_nodes`: discover. The relationships around a symbol, owner, or action when you do not yet know the endpoints.
- `expand_nodes`: the source of handles a graph tool returned. This is how you read code; pass every handle you need in one call.
- `query_files`: a roster of declarations and adjacent files for known paths.
- `query_diagnostics`: current compiler and lint diagnostics.

The graph mirrors the current program; after an edit, query again rather than reuse an old result.

Fall back to grep or read only when the graph cannot answer: no matching node, non-TypeScript files, generated output, or a literal text search.
