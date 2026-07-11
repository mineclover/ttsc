# Risk Rules

Use these rules to turn `inspect_typescript_graph` evidence into review severity. Treat scores as prioritization, not proof of defects.

## Inputs

- `details.nodes[]`: signature, members, calls, types, dependsOn, dependedOnBy.
- `trace`: direction, hops, reached nodes, path, steps, truncated.
- `overview`: hotspots and publicApi.
- `tour`: entrypoints, primaryFlow, nearby, tests, answerAnchors.

## Severity

`high`
: Public API or central runtime path changed, and impact/reverse trace reaches multiple dependents or user-facing tests are absent.

`medium`
: Shared internal symbol changed with callers, type dependents, or persistence/auth/rendering/collaboration involvement.

`low`
: Local symbol with limited dependents and nearby tests or clear anchors.

`info`
: Documentation/read-next insight with no direct change risk.

## Risk Signals

- `exported` or appears in `overview.publicApi`: raise one level.
- `trace.truncated`: mention incomplete impact surface; do not assume all dependents are known.
- `direction:"impact"` reaches test anchors: cite them as verification candidates.
- No tests reached for a central flow: ask for targeted test evidence.
- `calls` includes persistence/auth/network/render/collab vocabulary: flag integration risk.
- `types` includes shared DTO/config/schema names: flag contract risk.
- External nodes only matter when the question is about library/API boundaries.

## Review Language

Use:
- "Graph evidence shows..."
- "Risk inference: ..."
- "Suggested verification: ..."

Avoid:
- "This is broken" unless a failing test or runtime error proves it.
- "No impact" when traces are capped or the graph returned `truncated`.
