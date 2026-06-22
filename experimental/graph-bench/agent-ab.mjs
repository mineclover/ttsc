// Agent-cost A/B for @ttsc/graph, the codegraph "X% cheaper / fewer tokens"
// benchmark. For each structural task it runs the Claude Code CLI headless twice:
// once with the @ttsc/graph MCP server configured, once without, and compares
// cost, output tokens, turns, and wall time.
//
// This spends real Claude credits and is non-deterministic; run it on a quiet
// host with enough runs to take a median (codegraph uses 4). It is deliberately
// NOT wired into CI. Requires `claude` and `go` on PATH.
//
// Usage:
//   node experimental/graph-bench/agent-ab.mjs                 # packages/ttsc, 1 run
//   node experimental/graph-bench/agent-ab.mjs --project=/abs --runs=4

import cp from "node:child_process";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";

const here = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(here, "..", "..");
const ttscDir = path.join(repoRoot, "packages", "ttsc");

const args = parseArgs(process.argv.slice(2));
const project = path.resolve(args.project ?? ttscDir);
const tsconfig = args.tsconfig ?? "tsconfig.json";
const runs = Number(args.runs ?? 1);

// Structural questions a coding agent would otherwise answer by grepping and
// reading across files. graph_explore answers each in one call.
const TASKS = [
  "List every function that calls `resolveBinary`, with the file path of each caller. Be concise.",
  "What is the blast radius of the `TtscService` class (how many declarations transitively depend on it), and name a few? Be concise.",
  "Which file declares `resolveTtscserverBinary`, and which functions call it? Be concise.",
];

const goRoot = path.join(os.homedir(), "go-sdk", "go", "bin");
const goEnv = {
  ...process.env,
  PATH: fs.existsSync(goRoot)
    ? `${goRoot}${path.delimiter}${process.env.PATH ?? ""}`
    : process.env.PATH,
};

const binary = path.join(
  os.tmpdir(),
  `ttscgraph-ab-${process.pid}${process.platform === "win32" ? ".exe" : ""}`,
);
console.log("Building ttscgraph...");
runOrThrow("go", ["build", "-o", binary, "./cmd/ttscgraph"], ttscDir, goEnv);

const mcpConfig = path.join(os.tmpdir(), `graph-mcp-${process.pid}.json`);
fs.writeFileSync(
  mcpConfig,
  JSON.stringify({
    mcpServers: {
      "ttsc-graph": {
        command: binary,
        args: ["--stdio", "--cwd", project, "--tsconfig", tsconfig],
      },
    },
  }),
);

const arms = [
  { name: "baseline", extra: [] },
  { name: "graph", extra: ["--mcp-config", mcpConfig] },
];

console.log(
  `\nAgent-cost A/B on ${path.relative(repoRoot, project) || project} — ${TASKS.length} tasks x ${runs} run(s) x 2 arms\n`,
);

const results = { baseline: [], graph: [] };
let spent = 0;
for (let t = 0; t < TASKS.length; t++) {
  for (const arm of arms) {
    for (let r = 0; r < runs; r++) {
      const metric = runClaude(TASKS[t], arm.extra);
      results[arm.name].push(metric);
      spent += metric.costUsd;
      console.log(
        `  task ${t + 1} ${arm.name.padEnd(8)}: $${metric.costUsd.toFixed(3)}, ` +
          `${metric.outputTokens} out tok, ${metric.numTurns} turns, ` +
          `${(metric.durationMs / 1000).toFixed(1)}s  [running $${spent.toFixed(2)}]`,
      );
    }
  }
}

const summary = {};
for (const arm of ["baseline", "graph"]) {
  summary[arm] = {
    costUsd: sum(results[arm].map((m) => m.costUsd)),
    outputTokens: sum(results[arm].map((m) => m.outputTokens)),
    inputTokens: sum(results[arm].map((m) => m.inputTokens)),
    cacheReadTokens: sum(results[arm].map((m) => m.cacheReadTokens)),
    numTurns: sum(results[arm].map((m) => m.numTurns)),
    durationMs: sum(results[arm].map((m) => m.durationMs)),
  };
}

const delta = (g, b) => (b === 0 ? 0 : ((g - b) / b) * 100);
console.log(`\nTotals (sum over ${TASKS.length} tasks x ${runs} run, graph vs baseline):`);
report("cost", `$${summary.graph.costUsd.toFixed(3)}`, `$${summary.baseline.costUsd.toFixed(3)}`, delta(summary.graph.costUsd, summary.baseline.costUsd));
report("output tokens", summary.graph.outputTokens, summary.baseline.outputTokens, delta(summary.graph.outputTokens, summary.baseline.outputTokens));
report("turns", summary.graph.numTurns, summary.baseline.numTurns, delta(summary.graph.numTurns, summary.baseline.numTurns));
report("wall time", `${(summary.graph.durationMs / 1000).toFixed(1)}s`, `${(summary.baseline.durationMs / 1000).toFixed(1)}s`, delta(summary.graph.durationMs, summary.baseline.durationMs));
console.log(`\nTotal spend this run: $${spent.toFixed(2)}`);

const reportPath = path.join(here, "agent-ab-report.json");
fs.writeFileSync(
  reportPath,
  `${JSON.stringify({ project: path.relative(repoRoot, project), runs, tasks: TASKS, results, summary }, null, 2)}\n`,
);
console.log(`Report: ${path.relative(repoRoot, reportPath)}`);

try {
  fs.rmSync(binary, { force: true });
  fs.rmSync(mcpConfig, { force: true });
} catch {
  /* best effort */
}

function report(label, graph, baseline, pct) {
  const sign = pct <= 0 ? "" : "+";
  console.log(`  ${label.padEnd(14)} baseline ${baseline}  ->  graph ${graph}  (${sign}${pct.toFixed(1)}%)`);
}

function runClaude(task, extra) {
  const result = cp.spawnSync(
    "claude",
    ["-p", "--output-format", "json", "--permission-mode", "bypassPermissions", ...extra],
    {
      cwd: project,
      // Prompt over stdin (not argv) so spaces/backticks need no shell quoting;
      // shell:true lets Windows resolve the `claude` .cmd shim.
      input: task,
      encoding: "utf8",
      windowsHide: true,
      shell: true,
      maxBuffer: 64 * 1024 * 1024,
      timeout: 240_000,
    },
  );
  if (result.error) throw result.error;
  if (result.status !== 0) {
    throw new Error(`claude failed (${result.status}): ${(result.stderr || result.stdout || "").slice(0, 600)}`);
  }
  const json = JSON.parse(result.stdout);
  const usage = json.usage ?? {};
  return {
    costUsd: json.total_cost_usd ?? 0,
    outputTokens: usage.output_tokens ?? 0,
    inputTokens: usage.input_tokens ?? 0,
    cacheReadTokens: usage.cache_read_input_tokens ?? 0,
    numTurns: json.num_turns ?? 0,
    durationMs: json.duration_ms ?? 0,
  };
}

function runOrThrow(command, commandArgs, cwd, env) {
  const result = cp.spawnSync(command, commandArgs, { cwd, env, encoding: "utf8", windowsHide: true });
  if (result.error) throw result.error;
  if (result.status !== 0) {
    throw new Error(`${command} ${commandArgs.join(" ")} failed (${result.status})\n${result.stderr ?? ""}`);
  }
  return result.stdout ?? "";
}

function sum(values) {
  return values.reduce((a, b) => a + b, 0);
}

function parseArgs(argv) {
  const out = {};
  for (const arg of argv) {
    const match = /^--([^=]+)=(.*)$/.exec(arg);
    if (match) out[match[1]] = match[2];
  }
  return out;
}
