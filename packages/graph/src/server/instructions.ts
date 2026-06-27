/**
 * The guidance delivered in the MCP initialize response. It is the only place
 * the agent is told how to use the graph; nothing is written to its config
 * files. Keep it short; the per-tool descriptions carry the detail.
 */
export const instructions = `
Use the graph tools as a staged TypeScript index, not as an answer writer. Get
handles with question_entrypoints or symbol_lookup, trace relationships with
dependency_path, inspect selected symbols with symbol_details, and read source
only for the few bodies whose logic decides the answer. Keep source reads
separate from dependency maps. Use shell reads only when the graph lacks the
needed non-TypeScript file, generated output, or literal text.

The graph already knows resolved symbols, edges, evidence spans, decorators,
and stable handles. Prefer it before shell-reading TypeScript source.

- question_entrypoints: first call for a natural-language code question. Use it
  once to get ranked starting symbols, direct mentions, and a small dependency
  orientation slice without source bodies.
- symbol_lookup: targeted symbol search. Use when you need a specific class,
  method, function, property, or type and do not already have its handle.
- dependency_path: dependency and call/type flow. Use it for "how A reaches B",
  lifecycle, request-flow, rendering-flow, validation-flow, and impact questions.
- symbol_details: selected symbol details. Use it for signatures, members,
  direct calls, direct types, dependency neighbors, and the few source bodies
  whose implementation decides the answer. It does not write answer text.
- project_overview: project-wide architecture map. Use for layers, hotspots,
  counts, and public API; not for a specific code question.

For a flow question, call question_entrypoints once, then dependency_path before
symbol_details. Read source only for the one or two leaf bodies whose logic is
needed.

Do not batch source:true across a path. Use symbol_details(neighbors:true)
without source to map dependencies; use symbol_details(source:true) without
neighbors to read bodies. When source is true, neighbor options are ignored.

Copy exact names from returned nodes, references, aliases, evidence snippets,
and trace steps. Prefer graph evidence and sourceSpan line anchors over shell
reads for citations.

Package scripts, config files, generated output, and exact text searches remain
valid shell/file-read cases because they are outside the symbol graph.
`.trim();
