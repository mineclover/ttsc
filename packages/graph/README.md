# `@ttsc/graph`

![banner of @ttsc/graph](https://ttsc.dev/og.jpg)

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/samchon/ttsc/blob/master/LICENSE) [![NPM Version](https://img.shields.io/npm/v/@ttsc/graph.svg)](https://www.npmjs.com/package/@ttsc/graph) [![NPM Downloads](https://img.shields.io/npm/dm/@ttsc/graph.svg)](https://www.npmjs.com/package/@ttsc/graph) [![Build Status](https://github.com/samchon/ttsc/workflows/test/badge.svg)](https://github.com/samchon/ttsc/actions?query=workflow%3Atest) [![Guide Documents](https://img.shields.io/badge/Guide-Documents-forestgreen)](https://ttsc.dev/docs/graph) [![Discord Badge](https://img.shields.io/badge/discord-samchon-d91965?style=flat&labelColor=5866f2&logo=discord&logoColor=white&link=https://discord.gg/E94XhzrUCZ)](https://discord.gg/E94XhzrUCZ)

A code-graph and diagnostics server for coding agents. Part of the [`ttsc`](https://ttsc.dev) toolchain, it hands an AI agent a map of your codebase that the TypeScript compiler itself resolved — so the agent can answer "how does this work?" without opening file after file.

It speaks the [Model Context Protocol (MCP)](https://modelcontextprotocol.io), the standard coding agents (Claude Code, Codex, …) use to call tools. You point your agent at it once; the agent gains two tools, `graph_explore` and `graph_diagnostics`.

## Demonstration

Ask your agent an architecture question — say _"how does the editor render a shape?"_. Without a graph it greps, opens a file, follows an import, opens another, and fans out across a dozen files. With `@ttsc/graph` it calls `graph_explore` once and gets the answer back:

```
graph_explore { "query": "render shape" }

class ShapeRenderer  src/render/shape.ts:18
  -> function rasterize (value-call)
  -> interface Shape (type-ref)
  -> class Canvas (value-call)        # new Canvas()
  <- class Editor (value-call)        # who calls into it
  blast radius: 9 transitive dependent(s)
  18  export class ShapeRenderer {
  19    constructor(private readonly canvas: Canvas) {}
  20    draw(shape: Shape): void {
  ...
```

Every arrow is resolved by the real type checker, not guessed from text. A call through a barrel re-export lands on the file that actually declares the symbol; a method-to-method call and a `<Component />` use are real edges; a `node_modules` or `.d.ts` type is marked as an external boundary the graph does not walk into. The agent reads the relationships and the source in one call instead of fanning out.

On codegraph's own agent-cost benchmark this cuts an agent's tokens by ~70% and tool calls by ~83% (averaged across two repositories), with the agent reading zero files. See the [benchmark](https://ttsc.dev/docs/graph/benchmark) for the full numbers and method.

## Install

```bash
npm install -D ttsc @ttsc/graph typescript@rc
```

`ttsc` is the host: `@ttsc/graph` runs the native graph server that ships inside ttsc's platform package, so install `ttsc` alongside it (the same pair as `@ttsc/lint`). There is no Go toolchain to install and nothing else to set up.

## Configure your agent

Add the server to your MCP client once. For Claude Code (in `.mcp.json`, or via `claude mcp add`):

```json
{
  "mcpServers": {
    "ttsc-graph": {
      "command": "npx",
      "args": ["-y", "@ttsc/graph"]
    }
  }
}
```

Start your agent from the project root so the server finds your `tsconfig.json`, or point it explicitly:

```bash
npx @ttsc/graph --cwd ./packages/app --tsconfig tsconfig.json
```

Usage guidance is delivered to the agent in the MCP `initialize` response. The server never writes into your `CLAUDE.md`, `AGENTS.md`, or any agent config file.

## Tools

| Tool | What it answers |
| --- | --- |
| `graph_explore` | _"What relates to this symbol? How does this work?"_ Give it a symbol name or a file fragment; it returns the matching declarations, their checker-resolved relationships (calls, types, heritage — in both directions), a blast-radius count, and the source. |
| `graph_diagnostics` | _"What's wrong with this file?"_ The TypeScript compiler's type errors for one file, with the exact `tsc` code and location. |

## How it works

`@ttsc/graph` rides ttsc's in-process TypeScript-Go compiler. The compiler type-checks your project once and keeps the result warm; the graph is read straight off that, so every node and edge is the compiler's own resolution. That is why the relationships are exact where a text- or tree-sitter-based tool can only guess. The server is a thin layer over that warm compiler, so a query is a method call on an already-built checker, not a fresh compile.

## Environment

- `TTSC_GRAPH_BINARY`: absolute path to a native graph binary, overriding the default platform resolution.
