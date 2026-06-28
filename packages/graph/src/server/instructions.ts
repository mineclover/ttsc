/**
 * The guidance delivered in the MCP initialize response. It is the only place
 * the agent is told how to use the graph; nothing is written to its config
 * files. Keep it short; the per-tool descriptions carry the detail.
 */
export const instructions = `
For TypeScript source questions that can be answered from symbols, imports,
types, calls, declarations, references, or source ranges,
inspect_typescript_source_flow is the structured evidence channel. It is suited
to code-flow questions such as how a request, decorator, dependency injection
chain, IPC/process/channel message, value, event, API call, render, validation,
relation load, or worker message moves through TypeScript code. Use shell or
file reads instead for scripts, configs, docs, generated output, exact source
text, or non-TypeScript files.

Repository guidance files, READMEs, and prose docs are documentation evidence,
not TypeScript source evidence. Do not read them to answer a source-flow
question unless the user asks about those documents.

In a TypeScript workspace, questions about how repository concepts communicate,
dispatch, handle, render, validate, load, persist, or propagate are TypeScript
source questions unless the user asks for configuration, prose docs, exact text,
generated output, or non-TypeScript files.

For architecture or flow questions, explain the central path. Do not exhaust
variants, tests, adapters, generated clients, or terminal implementation
branches unless the user asks for those branches.

If you call inspect_typescript_source_flow for an answer, finish that answer from
graph results only. Do not call shell search or file-read tools in the same
answer to recover symbols, call targets, line numbers, tests, branch details, or
source bodies. If the graph cannot settle a source-body detail, answer with the
returned sourceSpan and the specific gap; do not open the file.

Work in this order. First, plan a 1-3 call budget. Second, call entrypoints for
the user question or lookup only when the user named a concrete symbol. Third,
use trace for flow, caller, dependency, lifecycle, render, request, validation,
and impact questions. Fourth, use details only for selected handles whose
signature, members, direct calls, direct types, or sourceSpan are needed. Stop
when file, symbol, relation, and range evidence is enough. Four graph calls is a
hard stop; a fifth returns escape so you can answer from current evidence.

Fill arguments in order: question, graphNeed, draft, review, request. These
are a short planning checklist, not answer prose. First name the evidence that
will be enough to answer. Then draft the smallest operation. Then review whether
this call should run or stop. Use review.decision="inspect" only when graph
evidence is still the next source. Use review.decision="escape" when evidence
is enough, the user needs source text/non-TypeScript facts, or another graph
call would be overfetch. review.finish is "answer", "range", or "clarify".

Checklist examples:

- First call: graphNeed="central handles for the main path"; draft={reason:
  "find ranked entrypoints before tracing", type:"entrypoints"}; review={reason:
  "no graph evidence yet; default entrypoints is smallest", decision:"inspect",
  finish:"answer"}.
- Follow-up: graphNeed="caller/dependency path from selected handle";
  draft={reason:"trace gives relationships without source reads", type:"trace"};
  review={reason:"not duplicate; capped trace is enough", decision:"inspect",
  finish:"answer"}.
- Stop: graphNeed="current graph result has names and ranges"; draft={reason:
  "source body would require file text", type:"escape"}; review={reason:"answer
  from graph or cite range; shell would duplicate source reading",
  decision:"escape", finish:"range"}.

Request selection:

- entrypoints: first call for a natural behavior or architecture question. Use
  default limits. It gives ranked handles and a small orientation slice.
- lookup: only for a concrete class, function, method, property, type, or dotted
  handle when no handle is already known.
- trace: preferred second call for flow questions. Use path mode when both ends
  are known. Use reverse for callers. Keep open traces shallow and capped.
- details: inspect one to three selected handles. Prefer default limits. Do not
  raise neighborLimit unless a previous graph result was truncated and the
  missing relation is named.
- overview: only for a broad public API or layer map. Do not use it for normal
  flow questions.
- escape: no-op route when graph should stop. Use it to finish with current
  evidence, a sourceSpan gap, or a clarification request. It is not permission
  to run shell for the same answer.

Copy exact names and ranges from returned nodes, references, aliases, trace
steps, evidence, and sourceSpan fields. Do not invent implementation text. If a
source body is needed, report the smallest returned sourceSpan and stop.
`.trim();
