# Verification record — 2026-10-27

Scope: local, dependency-free contract tests for the experimental Effect TSGO → TTSC Checker port. This record does **not** assert full native Checker compilation, npm installation, LSP integration or Effect runtime execution.

## Environment

- Node.js v22.16.0
- Go v1.23.2 (below TTSC Go 1.26 and Effect Go 1.27 requirements)
- Local container has no network access to GitHub/npm. No actual TS-Go/TTSC/Effect compiler binary was installed or invoked.
- The tested scripts and test files have identical Git blob SHA values to the files in the experimental fork branch:
  - scripts/preflight.mjs: bdebb9bc32750c72fe2220973353b0794e56b916
  - test/preflight.test.mjs: 9a009a0b53ca223cec7cadceccfd6b88f7b97a9b
  - scripts/probe-compat.mjs: 09bd71f18d0462a0e2215bf221ef8712587ae440
  - test/probe-compat.test.mjs: 29f7d2a66529ef6b84977715f393a4c0c947a940
- scripts/materialize-patched-tsgo.sh: ec6f2168864ae688898c1b5d6dd7371bc7c0bd89. Shell syntax (bash -n) and help invocation checked; patch application itself was **not** executed.

## Actually executed commands

1. node --test test/preflight.test.mjs test/probe-compat.test.mjs — **7 PASS / 0 FAIL**, against temporary test-double compiler processes and temporary project manifests.
2. bash -n scripts/materialize-patched-tsgo.sh — **PASS**, syntax only.
3. bash scripts/materialize-patched-tsgo.sh --help — expected usage output, zero exit.
4. GitHub source inspection — pinned TTSC and Effect Go/module ref mismatch (7fc57c.. vs 2bd066..), archived 14 source patch files and code/fixture files persisted to isolated branch.

## Contract checks that passed

- TTSC June TypeScript-Go ref blocks Effect Checker port preflight.
- Matching TypeScript-Go ref only gives preflight readiness, not port/compatibility PASS.
- Missing required Checker patch blocks preflight.
- Both TTSC and Effect diagnostic tool observations remain independently identifiable.
- Invalid Effect diagnostic JSON is not treated as empty diagnostic list.
- Missing compiler executable is reported as unavailable, not PASS.
- Effect config that checks zero files due to disabled/unselected plugin does not yield a false PASS.

## Not executed / unresolved

- Patch application to a local TypeScript-Go checkout.
- Go package build with Effect Checker patch, TTSC shim regeneration, checker callbacks and Effect rule runner.
- TypeScript named-plugin regression test in the actual TTSC test runner.
- Fixture dependency installation, TTSC check/emit, Effect diagnostic CLI and Node Effect runtime tests.
- LSP cancellation, UTF-16, dirty-buffer, Code Actions, data races, cache invalidation, graphs and Evidence regression.
- Independent code review, CI, PR merge or user approval.

Do not treat this record as proof of an end-to-end Effect/TTSC implementation. The next hard gate is Go/compiler dependency alignment in a real development environment.
