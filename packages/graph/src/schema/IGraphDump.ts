import { IGraphDiagnostic } from "./IGraphDiagnostic";
import { IGraphEdge } from "./IGraphEdge";
import { IGraphNode } from "./IGraphNode";

/**
 * The whole-graph export `ttscgraph dump` writes and the MCP server loads — the
 * wire contract between the Go fact-builder and the TypeScript graph engine.
 *
 * It is the complete graph with none of the per-response caps the MCP tools
 * apply: every node and edge the build resolved. The server parses it once at
 * startup (typia-validated) into an in-memory resident graph and answers every
 * tool call from that warm model; the bundled 3D viewer reduces the same dump.
 *
 * Paths in `project` and `tsconfig` are absolute; `file` fields on nodes,
 * edges, and diagnostics are project-relative. `schemaVersion` is bumped on any
 * breaking shape change so a consumer can refuse a dump it cannot read.
 */
export interface IGraphDump {
  /** The dump shape version; bumped on any breaking change. */
  schemaVersion: number;

  /** Absolute path of the project root the graph was built for. */
  project: string;

  /** The tsconfig the program was loaded from, relative to `project`. */
  tsconfig: string;

  /** Every node the build recorded. */
  nodes: IGraphNode[];

  /** Every edge the build resolved. */
  edges: IGraphEdge[];

  /**
   * Fused compiler and plugin diagnostics, when diagnostics were collected.
   * Absent when the dump was built without a diagnostics pass.
   */
  diagnostics?: IGraphDiagnostic[];
}

/**
 * The current {@link IGraphDump.schemaVersion}. Bump this in lockstep with the
 * Go `dump.go` writer whenever the dump shape changes in a breaking way.
 */
export const GRAPH_SCHEMA_VERSION = 2;
