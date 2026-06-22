# @ttsc/graph benchmark

Two benchmarks, mirroring the two codegraph publishes: a structural one (`bench.mjs`, coverage and counts) and an agent-cost A/B (`agent-ab.mjs`, "X% cheaper / fewer tokens").

## Structural benchmark (`bench.mjs`)

Measures `@ttsc/graph` on a real project, the way codegraph reports coverage: how long the resident `Program` takes to load, how cheap graph extraction is on top of that already-built `Program`, the node and edge counts, and the "fair coverage" (share of symbol-bearing source files with at least one resolved cross-file edge).

The counts and coverage are deterministic. The timings are indicative and only trustworthy on a quiet host (see `.codex/skills/benchmark`); a CI run shows the shape, not a publishable figure.

### Run

```bash
# Default target: packages/ttsc (this repo's launcher source)
node experimental/graph-bench/bench.mjs

# Any project, with run count
node experimental/graph-bench/bench.mjs --project=/abs/path/to/project --tsconfig=tsconfig.json --runs=5
```

It builds the `cmd/graphbench` metrics binary once, runs it `--runs` times (plus one warmup), and writes `report.json` next to this file.

### What it reports

A CI run against this repo's `packages/ttsc` (53 source files) reported:

```
Result (counts deterministic; timings indicative):
  source files:  53
  nodes:         575 (66 external boundary leaves)
  edges:         1402 (heritage 2, value-call 1016, type-ref 384)
  fair coverage: 92.2% (47/51 symbol-bearing files cross-linked)
  load time:     81 ms (median)
  graph build:   37 ms (median), 45.7% on top of the load it rides
```

Read the coverage as the codegraph-style flex: 92.2% of symbol-bearing files have at least one checker-resolved cross-file edge. The `graph build ... % on top of the load it rides` line is honest about cost: on a small project the walk is a real fraction of the (already fast) load, and the ratio shrinks as type-checking dominates on larger trees. The point is not that extraction is free, but that it rides the `Program` the compiler already built, so the server answers queries without a second compile or an external language-server round-trip.

## Agent-cost A/B (`agent-ab.mjs`)

The benchmark codegraph leads with: does an agent spend less when it has the graph? For each structural task it runs the Claude Code CLI headless twice, once with the `@ttsc/graph` MCP server configured and once without, and compares cost, output tokens, turns, and wall time. It spends real Claude credits, is non-deterministic, and is not wired into CI; run it on a quiet host with enough runs to take a median. Requires `claude` and `go` on `PATH`.

```bash
node experimental/graph-bench/agent-ab.mjs                  # packages/ttsc, 1 run
node experimental/graph-bench/agent-ab.mjs --project=/abs --runs=4
```

An indicative single run over three structural questions about `packages/ttsc` (callers of a function, blast radius of a class, where a symbol is declared and what calls it):

```
Totals (graph vs baseline):
  cost           baseline $0.576  ->  graph $0.490  (-15.0%)
  output tokens  baseline 5073    ->  graph 3727    (-26.5%)
  turns          baseline 18      ->  graph 18      (0.0%)
  wall time      baseline 107.9s  ->  graph 70.5s   (-34.6%)
```

This is one run (`--runs=1`), so treat it as a shape, not a published figure: roughly 15% cheaper and a quarter fewer output tokens with the graph, in line with codegraph's reported ~16% cheaper. Take a median of several runs for a number worth quoting.
