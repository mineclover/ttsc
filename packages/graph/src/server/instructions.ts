/**
 * The guidance delivered in the MCP initialize response. It is the only place
 * the agent is told how to use the graph; nothing is written to its config
 * files. Keep it short; the per-tool descriptions carry the detail.
 */
export const instructions = `
For TypeScript code questions, call
inspect_typescript_graph_before_shell_reading before any shell search or source
read. Do not start with ls, rg, cat, or Get-Content to locate TypeScript files,
symbols, dependencies, or line anchors. Ask the graph first, then answer from
graph evidence and sourceSpan citations.

The graph is a TypeScript index, not an answer writer. Fill thinking before
each call, then choose one request.type: find entrypoints, lookup symbols, trace
dependency paths, inspect selected symbols, or summarize the project. If more
TypeScript evidence is needed, make another graph request instead of switching
to shell search.

The graph already knows resolved symbols, dependency edges, evidence spans,
decorators, stable handles, source bodies, and sourceSpan line anchors.

Request types:

- find_question_entrypoints: first call for a natural-language code question.
  It returns ranked symbols, direct mentions, and small dependency orientation
  without source bodies.
- lookup_symbols: targeted symbol search for a class, method, function,
  property, or type when you do not already have its handle.
- trace_dependency_path: call/type/dependency flow for "how A reaches B",
  lifecycle, request-flow, rendering-flow, validation-flow, and impact
  questions.
- inspect_symbol_details: signatures, members, direct calls, direct types,
  dependency neighbors, or narrow source/sourceSpan reads for selected handles.
- summarize_project: source-free architecture map for layers, hotspots, counts,
  and public API.

For a flow question, use find_question_entrypoints once, then
trace_dependency_path before inspect_symbol_details. Keep broad dependency maps
separate from source reads. When source is true, neighbor options are ignored.

Copy exact names from returned nodes, references, aliases, evidence snippets,
sourceSpan anchors, and trace steps. Do not use shell only to recover TypeScript
line numbers already returned by graph evidence.

Package scripts, config files, generated output, and exact text searches remain
valid shell/file-read cases because they are outside the symbol graph.
`.trim();
