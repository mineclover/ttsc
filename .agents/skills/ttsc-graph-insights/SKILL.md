---
name: ttsc-graph-insights
description: "Use ttsc graph MCP for TypeScript architecture insight, code-tour summaries, docs drift checks, lightweight change review, PR/comment drafting, and dependency-impact risk analysis. Trigger when working with inspect_typescript_graph, @ttsc/graph, compiler graph output, TypeScript symbol dependency analysis, changed-symbol review, architecture documentation from graph evidence, or pre-bug risk detection from call/type/public API impact."
---

# TTSC Graph Insights

## Principle

Use `inspect_typescript_graph` as compiler-resolved evidence, not as a search hint. Its returned symbols, edges, signatures, spans, tests, anchors, and `next` are the answer surface. Do not read files to re-confirm graph facts when `next.action` is `answer`.

Use `ttsc-graph dump` only when the user asks for whole-graph batch analysis, visualization, custom scoring, or persistent indexes. Use MCP responses for agent-facing insight, reviews, comments, and docs.

## Request Routing

Choose the smallest graph request that can answer the task:

- Repository tour, architecture orientation, read-next, runtime flow: `tour`.
- Public API, layers, hotspots: `overview`.
- Narrow question with unknown starting symbol: `entrypoints`.
- Concrete class/function/member/type name: `lookup`, then `details` only if needed.
- Caller/callee/path/impact analysis: `trace`.
- Selected symbol shape, signatures, members, direct deps: `details`.
- Non-TypeScript files, exact source body, configs, generated outputs, or already answered graph facts: `escape`.

Follow the returned `next`. If it says `answer`, stop graph expansion and write the result.

## Router Routing

Use `ttsc-graph-router` when the question spans configured repositories or needs
cache, checkpoint, review, lint, artifact, or visualization projections.

1. Select `repoId` with `graph_router_list_repos`; use `graph_router_health` for
   runtime and cache readiness.
2. Use `graph_router_inspect` as the routed `inspect_typescript_graph` surface.
3. Treat current-cache tools as saved-file views. The router fingerprints source
   and configuration inputs automatically; set `refresh: true` when a fresh dump
   is required in the same operation.
4. Use checkpoint tools only for fixed historical evidence. Never present a
   checkpoint as the current source snapshot.
5. Keep projection summaries separate from compiler facts:
   - change review: `graph_router_changed_impact`, `graph_router_document_symbols`
   - dependency policy: `graph_router_related_nodes`, `graph_router_lint_rules`
   - architecture views: `graph_router_scope_tree`, `graph_router_repo_map`,
     `graph_router_file_visualization`, `graph_router_file_dependencies`
   - artifact contract: `graph_router_artifact_stats`,
     `graph_router_artifact_validate`
6. For dependency projections, report inbound as "affected by" and outbound as
   "affects". Preserve returned truncation counts instead of treating omitted
   output as omitted graph relationships.

## Workflows

### Architecture Insight

1. Use `tour` for flow-oriented questions or `overview` for layer/API/hotspot questions.
2. Summarize central entrypoints first.
3. Report `primaryFlow.steps`, `nearby`, `tests`, and `answerAnchors`.
4. Keep each claim tied to graph-provided symbol ids or file/line anchors.

### Docs Drift

1. Identify the documented flow, API, or module boundary.
2. Use `tour` or `overview` to get current graph anchors.
3. Compare documented anchors/names with current `answerAnchors`, `publicApi`, or `layers`.
4. Mark drift as:
   - `stale`: documented symbol/anchor no longer appears.
   - `incomplete`: graph reports important nearby/test/public API anchors missing from docs.
   - `ok`: docs align with graph evidence.
5. Do not rewrite large docs unless requested; produce a patch-sized update plan first.

### Lightweight Change Review

1. Start from changed TypeScript files or user-named symbols.
2. Resolve ambiguous names with `lookup` or `entrypoints`.
3. For each selected symbol, use `details` for shape and direct deps.
4. Use `trace` with `direction:"impact"` or `direction:"reverse"` for dependents and tests.
5. Produce review notes ordered by risk, not by file order.

### PR / Inline Comment Drafting

1. Convert graph findings into short comments with one concrete anchor.
2. Prefer comments about impact, missing tests, public API exposure, or stale docs.
3. Avoid speculative wording beyond graph evidence. Say "graph shows" for graph facts and "risk" for inference.
4. Keep comments actionable: ask for a test, doc update, or reviewer attention to a named dependent.

### Pre-Bug Risk Detection

Use graph evidence to flag risk, not to claim a bug exists. High-value signals:

- Changed symbol is exported/public API.
- Reverse/impact trace reaches many dependents or tests.
- Direct type dependencies include broad interfaces, discriminated unions, or shared DTOs.
- Runtime trace touches persistence, auth, request handlers, rendering, or collaboration paths.
- Tests are absent from `tour.tests` or impact trace for a central flow.
- A generated/external boundary appears in the path and the user expects authored-code behavior.

Load `references/risk-rules.md` when producing a risk score or gating recommendation.
Load `references/output-templates.md` when producing repeatable docs/review/comment output.

## Output Rules

- Separate compiler facts from inference.
- Keep source anchors as coordinates; do not paste source bodies unless explicitly requested.
- If `truncated` is present, mention the cap and avoid expanding every branch.
- If graph evidence is insufficient, say what outside evidence is needed: tests, typecheck, runtime smoke, schema diff, or manual source read.
- For batch visualization, recommend `ttsc-graph dump` and map `nodes[].id` to vertices and `edges[].from -> edges[].to` to directed edges.
