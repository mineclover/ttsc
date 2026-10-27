# Effect TSGO Checker → TTSC port (experimental)

Status: **source staging + compatibility probe implemented; native Checker port and real TTSC/Effect execution NOT verified**. This is a separate experiment, not a change to TTSC's production compiler, upstream submodule, published packages, or default CI.

## Authority and source pins

| Source | Pinned revision / runtime | Role |
| --- | --- | --- |
| TTSC | samchon/ttsc 7d316520811e354cad8b522a0cc895740a34c4a4 (v0.33.1) | Consumer and port target |
| TTSC's TypeScript-Go | 7fc57c0050639e01409766f41a77a10d03106297, Go 1.26 | In-process Program/Checker dependency |
| Effect-TS/tsgo | 45b0ce63d2583733d4deb9bf9a5b38bc629b33a6 (@effect/tsgo 0.59.0) | Imported source patches and intended semantics |
| Effect TypeScript-Go 7.0.2 | 2bd066d87f5bafd315be9f40889d0a60b9e58e0b, Go 1.27 | Source baseline required by Effect patches |

Microsoft TypeScript-Go's 2bd066... commit is **144 commits after** TTSC's pinned 7fc57c... commit. The same npm TypeScript version does not mean the two Go Checkers have matching source or internal APIs. The Effect-Go shim namespace (github.com/microsoft/TypeScript/tsc/shim/*) is different from TTSC's (github.com/microsoft/typescript-go/shim/*). Regenerating/aligning shims is part of the port, not a symbolic package alias.

## Imported patch candidates

The upstream-patches/ folder mirrors **14 original source patches** unmodified from the pinned Effect commit. They are **staged inputs**, not an applied or compiled TTSC patch.

| Concern | Patch numbers | Expected responsibility |
| --- | --- | --- |
| Type inference observation and checker state | 002, 003, 004, 005, 023, 033 | after-file callback, expression inference callback, instantiated generic signatures, type relation errors, per-checker EffectLinks |
| Compilation diagnostic semantics | 006, 009 | Effect diagnostic suppression isolation and exit/noEmitOnError policy |
| Config/incremental type analysis | 007, 013, 025, 028, 030 | Effect config parsing, inheritance, validation, incremental invalidation |
| Native CLI bootstrap | 001 | Effect hook initialization and separate diagnostics command |

Not yet staged / dependency review needed: Effect diagnostic codegen patch 008; AST/JSDoc stability patch 031; checker-adjacent helpers 014; LSP/fourslash/completions/code actions patches 010, 011, 012, 015, 017, 018, 021, 022, 024, 027. Full Effect rule engine (etscore, typeparser, internal/rules/rulerunner, fixables, refactors and shim generator) is **not** in this patch folder.

Original source license: MIT, Copyright (c) 2026 Effect. Retain their license and attribution if adapted. Reference: https://github.com/Effect-TS/tsgo/tree/45b0ce63d2583733d4deb9bf9a5b38bc629b33a6/_patches/typescript-go.

## Native Checker port: required implementation boundary

1. **Align the real Go compiler source baseline** in an isolated branch/build. Do not mutate a Go module cache or TTSC production toolchain. Account for TypeScript-Go's 144-commit drift, Go 1.27 requirement, TTSC's maintained shim aliases and generated declarations. Preserve upstream TTSC changes made after this tag for any future promotion.
2. **Add/check Checker callback semantics**: per-file post-check diagnostics, expression checked-type observation (including intermediate inference), type-relation error collection, instantiated signature type args and per-Checker state. Validate checker count and callback lifetime; Effect's upstream registrars are process-global/single-slot and should not silently override other TTSC contributors.
3. **Port Effect diagnostics engine** using the same rule names, type semantics, disable/enable directives and severities. Distinguish compiler TypeScript findings from Effect 377xxx findings; a TypeScript suppression directive must not accidentally suppress Effect-only diagnostics.
4. **Align configuration**: TTSC's transform plugin descriptors vs Effect's name-based plugin configuration cannot be loaded as equivalent entries. Preserve other TypeScript language-service plugin entries while selecting only TTSC native entries. Ensure inherited plugin settings and incremental cache invalidation behave deterministically.
5. **Connect to TTSC check/build/fix flow** via supported plugin/driver boundary. An LSP-only merge does not satisfy build-time verification. Only after CLI parity passes, enable Effect hover/quick fix/refactor in ttscserver proxy with position-encoding and dirty-buffer tests.

This work should ultimately be implemented in TTSC's real Go compiler integration; copying patches alone is insufficient.

## Debugging and verification compatibility matrix

| Plane | Owner | Evidence and failure policy |
| --- | --- | --- |
| Authoring and LSP | TTSC editor + upstream Effect LSP | Document version, UTF-16 range, diagnostic origin, stale withdrawal, QuickFix/refactor IDs |
| Static/check | TTSC native Program/Checker, @ttsc/lint, Effect rules | Same source+tsconfig revision; per-rule outcomes, compiler exit codes; missing rule execution != zero findings |
| Evidence/graph | Existing TTSC Evidence/graph providers | Same Program/checkpoint identity; Effect finding may link to graph node only with independently verified anchor |
| Build/emit | TTSC compiler | Effect diagnostic severity mapping must not bypass build gates or alter emitted JS unexpectedly |
| Runtime | Node running actual Effect 4 programs | Success, typed failure, laziness, Context/Layer, generator, resource finalizer, interruption behavior |
| Runtime debugger | Effect Runtime DevTools (separate tool) | Fiber/Context/Span/pause-on-defect as independently observed capabilities; static Checker port alone does not enable them |
| AI fix loop | Agent/tool consumer, work-infra receipt | Before/after source fingerprint, exact finding, applied code-action provenance, runtime test, non-stale proof |

Effect v4 real-use fixtures have separate TTSC and Effect tsconfig files. The combined config reproduces the current plugin-selection collision. Do not treat the combined config as supported until fixed and observed.

## Commands and what they prove

- Readiness (expected BLOCKED on TTSC v0.33.1): node experimental/effect-checker-port/scripts/preflight.mjs
- Contract tests without external dependencies: node --test experimental/effect-checker-port/test/*.test.mjs
- In fixture directory after isolated dependency installation: npm run check:ttsc; npm run check:effect
- Separate diagnostic probe: npm run probe
- TTSC emitted JS runtime assertion: npm run test:runtime

The probe inspects the same source/config files before and after calling both engines; it distinguishes Effect diagnostics, disabled plugin, malformed JSON, absent binary, subprocess failures, and optional Node runtime tests. It emits **observed-needs-compatibility-review**, not a blanket compatibility PASS.

No installer, binary build, real Effect diagnostic or runtime integration test was run in the authoring environment. Node-level test-double contract tests are not evidence of a functioning port.

## Acceptance for progressing past source staging

- [ ] Native TypeScript-Go Go/shim toolchain aligns to a declared pinned source; patch application checked and built.
- [ ] Effect Checker-specific rules run **inside the TTSC check/build pipeline**, not only as an external diagnostic comparison.
- [ ] Normal Effect program compiles through TTSC and runs; floatingEffect and generator/context invalid cases are flagged with correct origin/severity.
- [ ] Error/noEmitOnError and TS-ignore/Effect-directive behavior matches Effect TSGO.
- [ ] Watch/edit and LSP daemon restarts invalidate stale diagnostics; callback state does not bleed across projects/checker generations.
- [ ] Existing TTSC lint/Evidence, graph, generation, formatting and non-Effect files remain compatible and regressions are checked.
- [ ] High-level results are tied to immutable source/check/config/toolchain identities; no inherited PASS from previous runs.

Do not update or merge the fork master automatically. This branch starts from the exact v0.33.1 source consumed by ttsc-ex; fork master is older and will require a deliberate source-history reconciliation before a PR or release.
