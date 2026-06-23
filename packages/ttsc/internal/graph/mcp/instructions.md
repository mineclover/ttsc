`ttsc-graph` is a graph of this codebase resolved by the TypeScript compiler: every symbol and the calls between them.

- For any question about how the code works, call `graph_explore` first, before grep or reading files
- Name every symbol the question touches in one query; it returns their source and call graph, usually the whole answer
- Querying one symbol at a time, or grepping to trace calls, wastes turns the graph already saves
- `graph_diagnostics` reports a file's type and lint errors
