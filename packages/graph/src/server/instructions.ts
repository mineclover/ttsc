/**
 * The guidance delivered in the MCP initialize response. It is the only place
 * the agent is told how to use the graph; nothing is written to its config
 * files. Keep it short; the per-tool descriptions carry the detail.
 */
export const instructions = `
This TypeScript project is indexed by the compiler. Use the graph as the first
source of codebase coordinates: it has already resolved symbols and the
relationships between them.

- graph_index: first call for source questions. It returns ranked symbols,
  declaration signatures, directly mentioned code handles, and nearby dependency
  context without source bodies.
- graph_overview: the architecture, layers, hotspots, and public API.
- graph_query: find any symbol by name or description; each hit carries its
  signature, usually enough to answer without expanding. Use next.expand when
  you need bodies.
- graph_trace: follow a flow forward, reverse, or to its impact; or give dotted
  from/to names for the path between two symbols, how A reaches B, in one call.
- graph_expand: a symbol's declared shape, its signature, and a container's
  members; source:true to read a specific TypeScript body after graph_index,
  graph_query, or graph_trace has located it.

Answer in as few calls as you can. For TypeScript declarations, read source with
graph_expand(source:true) on the resolved handles. Use shell or file reads only
for a non-TypeScript file, generated output, or an exact literal text search that
is not represented as symbols or edges.
`.trim();
