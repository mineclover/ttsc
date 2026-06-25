import { GraphModel } from "../model/GraphModel";
import { IGraphNode } from "../schema";

/** Trace how execution and dependency flow to or from a symbol or route. */
export interface ITraceProps {
  /**
   * Where to start: a node id from another tool, or a symbol/route name. An
   * ambiguous name returns its candidates instead of a trace.
   */
  from: string;

  /**
   * `forward` follows what the start uses (callees, instantiations, renders);
   * `reverse` follows what uses the start (callers); `impact` is a reverse
   * trace that flags the public API, routes, and tests a change would reach.
   *
   * @default "forward"
   */
  direction?: "forward" | "reverse" | "impact";

  /**
   * How many hops deep to follow.
   *
   * @default 6
   */
  maxDepth?: number;

  /**
   * Cap on reached nodes; the trace stops and marks itself truncated past it.
   *
   * @default 60
   */
  maxNodes?: number;
}

export interface ITraceResult {
  /** The resolved start node, or undefined when `from` matched nothing. */
  start?: ITraceNode;
  direction: string;
  /** Edges traversed, in breadth-first order. */
  hops: ITraceHop[];
  /** Unique nodes reached (excluding the start), each with its depth and roles. */
  reached: ITraceNode[];
  /** True when the trace hit maxNodes or maxDepth and more flow exists. */
  truncated: boolean;
  /** When `from` was an ambiguous name, the matches to disambiguate with. */
  candidates?: ITraceNode[];
}

export interface ITraceHop {
  from: string;
  to: string;
  kind: string;
  /** Hops from the start (1 = direct). */
  depth: number;
}

export interface ITraceNode {
  id: string;
  name: string;
  kind: string;
  file: string;
  /** Hops from the start, on a reached node. */
  depth?: number;
  /** Why this node matters to an impact trace: `exported`, `route`, `test`. */
  roles?: string[];
}

const DEFAULT_DEPTH = 6;
const DEFAULT_MAX_NODES = 60;

/**
 * Breadth-first trace along the dependency graph. Structural
 * (contains/exports/imports) and heuristic edges are excluded so the path is
 * real call/type flow; forward walks callees, reverse and impact walk callers.
 * Impact additionally tags each reached node's role so the blast radius on the
 * public surface is legible.
 */
export function runTrace(graph: GraphModel, props: ITraceProps): ITraceResult {
  const direction = props.direction ?? "forward";
  const maxDepth = props.maxDepth ?? DEFAULT_DEPTH;
  const maxNodes = props.maxNodes ?? DEFAULT_MAX_NODES;
  const reverse = direction === "reverse" || direction === "impact";

  const start = resolveStart(graph, props.from);
  if (start.candidates) {
    return {
      direction,
      hops: [],
      reached: [],
      truncated: false,
      candidates: start.candidates.map((n) => summary(n)),
    };
  }
  if (start.node === undefined) {
    return { direction, hops: [], reached: [], truncated: false };
  }

  const hops: ITraceHop[] = [];
  const reached = new Map<string, ITraceNode>();
  const visited = new Set<string>([start.node.id]);
  let queue: Array<{ id: string; depth: number }> = [
    { id: start.node.id, depth: 0 },
  ];
  let truncated = false;

  while (queue.length > 0) {
    const next: Array<{ id: string; depth: number }> = [];
    for (const { id, depth } of queue) {
      if (depth >= maxDepth) {
        truncated = true;
        continue;
      }
      const edges = reverse ? graph.incoming(id) : graph.outgoing(id);
      for (const edge of edges) {
        if (!traversable(edge.kind, edge.provenance)) continue;
        const otherId = reverse ? edge.from : edge.to;
        const other = graph.node(otherId);
        if (other === undefined || other.kind === "file") continue;
        hops.push({
          from: edge.from,
          to: edge.to,
          kind: edge.kind,
          depth: depth + 1,
        });
        if (visited.has(otherId)) continue;
        if (reached.size >= maxNodes) {
          truncated = true;
          continue;
        }
        visited.add(otherId);
        reached.set(otherId, summary(other, depth + 1));
        next.push({ id: otherId, depth: depth + 1 });
      }
    }
    queue = next;
  }

  return {
    start: summary(start.node),
    direction,
    hops,
    reached: [...reached.values()],
    truncated,
  };
}

/** Resolve `from` to a single node, or report ambiguous-name candidates. */
function resolveStart(
  graph: GraphModel,
  from: string,
): { node?: IGraphNode; candidates?: IGraphNode[] } {
  const byId = graph.node(from);
  if (byId !== undefined) return { node: byId };
  const named = graph.named(from).filter((n) => n.kind !== "file");
  if (named.length === 1) return { node: named[0] };
  if (named.length > 1) return { candidates: named.slice(0, 12) };
  return {};
}

/** A node summary; roles are attached when present so impact reads at a glance. */
function summary(node: IGraphNode, depth?: number): ITraceNode {
  const roles: string[] = [];
  if (node.exported) roles.push("exported");
  if (node.kind === "route") roles.push("route");
  if (isTestFile(node.file)) roles.push("test");
  const out: ITraceNode = {
    id: node.id,
    name: node.qualifiedName ?? node.name,
    kind: node.kind,
    file: node.file,
  };
  if (depth !== undefined) out.depth = depth;
  if (roles.length > 0) out.roles = roles;
  return out;
}

/** An edge the trace should follow: real dependency, not structure or guess. */
function traversable(kind: string, provenance: string): boolean {
  if (provenance === "heuristic") return false;
  return kind !== "contains" && kind !== "exports" && kind !== "imports";
}

function isTestFile(file: string): boolean {
  return (
    /(^|\/)(test|tests|__tests__|spec)\//.test(file) ||
    /\.(test|spec)\.[cm]?tsx?$/.test(file)
  );
}
