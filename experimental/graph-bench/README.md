# @ttsc/graph benchmark

Measures `@ttsc/graph` on a real project, the way codegraph reports coverage: how long the resident `Program` takes to load, how cheap graph extraction is on top of that already-built `Program`, the node and edge counts, and the "fair coverage" (share of symbol-bearing source files with at least one resolved cross-file edge).

The counts and coverage are deterministic. The timings are indicative and only trustworthy on a quiet host (see `.codex/skills/benchmark`); a CI run shows the shape, not a publishable figure.

## Run

```bash
# Default target: packages/ttsc (this repo's launcher source)
node experimental/graph-bench/bench.mjs

# Any project, with run count
node experimental/graph-bench/bench.mjs --project=/abs/path/to/project --tsconfig=tsconfig.json --runs=5
```

It builds the `cmd/graphbench` metrics binary once, runs it `--runs` times (plus one warmup), and writes `report.json` next to this file.

## What it reports

The layout below is illustrative; run it for your own figures.

```
Result (counts deterministic; timings indicative):
  source files:  152
  nodes:         1843 (271 external boundary leaves)
  edges:         2104 (heritage 96, value-call 1187, type-ref 821)
  fair coverage: 86.8% (132/152 symbol-bearing files cross-linked)
  load time:     980 ms (median)
  graph build:   41 ms (median), 4.2% on top of the load it rides
```

The `graph build ... % on top of the load it rides` line is the point: extraction is a small fraction of the compile the graph rides, which is why a resident server answers queries without a fresh compile.
