# ttsc-graph

A compiler-resolved graph of TypeScript relationships: calls, callers, types, ownership, blast radius. It mirrors the code, so an edit changes it; re-query after editing, not from an old result.

- **How does code connect?** Reach for `graph_explore` before grep: one broad query (owner + action + nouns); answer from its result.
- **A file's errors, or all of them?** Call `graph_diagnostics` (omit the file for all).
- **Graph can't answer** (no match, omitted source, non-TypeScript)? grep/read.

## Call `graph_explore` for relationship questions

**Call it first to trace how code connects:** callers, callees, type flow, ownership, blast radius.

- It returns a compiler-resolved snapshot: line-numbered source, calls, callers, types, blast radius.
- grep/read cannot assemble that, because the answer depends on resolved relationships, not on where a keyword appears.

## Shape one broad query, not one-symbol probes

**Put the whole flow into a single `query`: owner + action + domain nouns.** Example: `repository find manager query builder`.

- Avoid lone-symbol lookups and generic navigation words; a flow-shaped query lets the graph rank the central nodes for you.
- One well-aimed query usually beats several narrow ones, but there is no call limit.

## Fall back to grep/read when the graph cannot answer

**Switch to ordinary search and file reads when the graph does not fit:**

- No node matches the symbol.
- The result omits source you still need.
- Non-TypeScript context: config, generated output, docs, JSON, other languages.
- You need every literal occurrence of a string.

`ttsc-graph` is a relationship graph, not a text index. Keyword counting and full-file dumps belong to grep and read.

## Use `graph_diagnostics` to triage one file

**Pass a file path (absolute, or a trailing fragment like `src/main.ts`) to get that file's errors:**

- TypeScript type errors,
- the project's `@ttsc/lint` rule violations,
- transform-plugin (typia, nestia) findings,

each with its code and location exactly as ttsc reports them. It inspects one file's problems; tracing relationships stays `graph_explore`'s job.

## Re-query freely

**Re-query freely; never ration calls to a fixed number.** Do it when:

- a better entry node appears in the result,
- a needed symbol was missing,
- you have edited a file since the last query.

The one trap is reusing an earlier result: it predates any edit you made after it, so query again instead of trusting what is already in your context.

## Final checklist

- Relationship or flow question? Use `graph_explore` with one broad owner + action + noun query.
- A file's errors, or the whole project's? Use `graph_diagnostics` with a path, or none for everything.
- No match, omitted source, non-TypeScript, or literal text search? Use grep/read.
- Edited a file since exploring? Query again; the result re-checks your edit.
