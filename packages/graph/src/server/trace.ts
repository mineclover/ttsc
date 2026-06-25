import { TtscGraphMemory } from "../model/TtscGraphMemory";
import { ITtscGraphNode } from "../structures/ITtscGraphNode";
import { ITtscGraphTrace } from "../structures/ITtscGraphTrace";

const DEFAULT_DEPTH = 6;
const DEFAULT_MAX_NODES = 60;

/**
 * Breadth-first trace along the dependency graph. Structural
 * (contains/exports/imports) and heuristic edges are excluded so the path is
 * real call/type flow; forward walks callees, reverse and impact walk callers.
 * Impact additionally tags each reached node's role so the blast radius on the
 * public surface is legible.
 */
export function runTrace(
  graph: TtscGraphMemory,
  props: ITtscGraphTrace.IProps,
): ITtscGraphTrace {
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

  const hops: ITtscGraphTrace.IHop[] = [];
  const reached = new Map<string, ITtscGraphTrace.INode>();
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
  graph: TtscGraphMemory,
  from: string,
): { node?: ITtscGraphNode; candidates?: ITtscGraphNode[] } {
  const byId = graph.node(from);
  if (byId !== undefined) return { node: byId };
  const named = graph.named(from).filter((n) => n.kind !== "file");
  if (named.length === 1) return { node: named[0] };
  if (named.length > 1) return { candidates: named.slice(0, 12) };
  return {};
}

/** A node summary; roles are attached when present so impact reads at a glance. */
function summary(node: ITtscGraphNode, depth?: number): ITtscGraphTrace.INode {
  const roles: string[] = [];
  if (node.exported) roles.push("exported");
  if (node.kind === "route") roles.push("route");
  if (isTestFile(node.file)) roles.push("test");
  const out: ITtscGraphTrace.INode = {
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
