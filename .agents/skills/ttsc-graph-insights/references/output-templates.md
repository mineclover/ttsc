# Output Templates

## Insight Summary

```md
### Graph Insight

Central entrypoints:
- `<symbol>` (`<kind>`) at `<file>:<line>`

Primary flow:
- `<step summary from graph>`

Nearby paths:
- `<anchor reason>`: `<symbol>` at `<file>:<line>`

Tests / usage anchors:
- `<test or usage anchor>`

Inference:
- `<short interpretation separated from compiler facts>`
```

## Change Review

```md
### Graph-Based Review

Changed symbol:
- `<id>` — `<signature>` at `<file>:<line>`

Impact:
- Public API: `<yes/no/unknown>`
- Dependents: `<names from trace/details>`
- Tests reached: `<anchors or none>`

Risk: `<high|medium|low|info>`
Reason:
- `<compiler fact>`
- `<risk inference>`

Suggested verification:
- `<test/typecheck/smoke/doc update>`
```

## PR Comment

```md
Graph evidence shows `<symbol>` is used by `<dependent>` at `<file>:<line>`.
Risk inference: this change may affect `<flow/API/path>`.
Suggested verification: `<specific test or doc update>`.
```

## Docs Drift Note

```md
### Docs Drift Check

Documented topic: `<flow/API/module>`

Current graph anchors:
- `<symbol>` at `<file>:<line>`

Drift:
- `<stale|incomplete|ok>` — `<reason>`

Patch recommendation:
- `<minimal doc update>`
```
