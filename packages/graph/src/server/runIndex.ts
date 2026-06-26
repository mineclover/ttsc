import { TtscGraphMemory } from "../model/TtscGraphMemory";
import { ITtscGraphEdge } from "../structures/ITtscGraphEdge";
import { ITtscGraphIndex } from "../structures/ITtscGraphIndex";
import { ITtscGraphNode } from "../structures/ITtscGraphNode";
import { resolveGraphHandle } from "./resolveHandle";
import { signatureOf } from "./runExpand";
import { runQuery } from "./runQuery";

const DEFAULT_LIMIT = 8;
const MAX_LIMIT = 20;
const DEFAULT_NEIGHBORS = 4;
const MAX_NEIGHBORS = 8;
const MAX_SEEDS = 6;
const STRUCTURAL_KINDS = new Set<string>(["contains", "exports", "imports"]);

/**
 * Build the first source-free index for a code question. The result gives the
 * model stable handles, declaration signatures, and direct graph context. It is
 * deliberately not a source reader; source remains opt-in through
 * graph_expand.
 */
export function runIndex(
  graph: TtscGraphMemory,
  props: ITtscGraphIndex.IProps,
): ITtscGraphIndex {
  const query = props.query.trim();
  const limit = bound(props.limit, DEFAULT_LIMIT, 1, MAX_LIMIT);
  const neighborLimit = bound(
    props.neighbors,
    DEFAULT_NEIGHBORS,
    0,
    MAX_NEIGHBORS,
  );

  const queryResult = runQuery(graph, { query, limit });
  const hits = queryResult.hits.map((hit) => ({ ...hit }));

  const mentions = directMentions(query).map((handle) => {
    const resolved = resolveGraphHandle(graph, handle, 6);
    const mention: ITtscGraphIndex.IMention = { handle };
    if (resolved.node !== undefined)
      mention.node = nodeOf(graph, resolved.node);
    if (resolved.candidates !== undefined) {
      mention.candidates = resolved.candidates.map((node) =>
        nodeOf(graph, node),
      );
    }
    return mention;
  });

  const seeds: ITtscGraphNode[] = [];
  const seen = new Set<string>();
  const addSeed = (node: ITtscGraphNode | undefined): void => {
    if (node === undefined || seen.has(node.id)) return;
    seen.add(node.id);
    seeds.push(node);
  };
  for (const mention of mentions) {
    if (mention.node !== undefined) addSeed(graph.node(mention.node.id));
  }
  for (const hit of hits) addSeed(graph.node(hit.id));

  let truncated = seeds.length > MAX_SEEDS;
  const neighborhood: ITtscGraphIndex.INeighborhood[] = [];
  for (const seed of seeds.slice(0, MAX_SEEDS)) {
    const outgoing = refs(graph, graph.outgoing(seed.id), "to", neighborLimit);
    const incoming = refs(
      graph,
      graph.incoming(seed.id),
      "from",
      neighborLimit,
    );
    if (outgoing.truncated || incoming.truncated) truncated = true;
    neighborhood.push({
      ...nodeOf(graph, seed),
      dependsOn: outgoing.items,
      dependedOnBy: incoming.items,
    });
  }

  return {
    query,
    hits,
    mentions,
    neighborhood,
    next: {
      expand: seeds.slice(0, MAX_SEEDS).map((node) => node.id),
      traceFrom: mentions
        .map((mention) => mention.node?.id)
        .filter((id): id is string => id !== undefined),
    },
    ...(truncated ? { truncated: true } : {}),
  };
}

function nodeOf(
  graph: TtscGraphMemory,
  node: ITtscGraphNode,
): ITtscGraphIndex.INode {
  const out: ITtscGraphIndex.INode = {
    id: node.id,
    name: node.qualifiedName ?? node.name,
    kind: node.kind,
    file: node.file,
  };
  if (node.evidence?.startLine !== undefined)
    out.line = node.evidence.startLine;
  const signature = signatureOf(graph.project, node);
  if (signature !== undefined) out.signature = signature;
  return out;
}

function refOf(
  node: ITtscGraphNode,
  edge: ITtscGraphEdge,
): ITtscGraphIndex.IReference {
  const out: ITtscGraphIndex.IReference = {
    id: node.id,
    name: node.qualifiedName ?? node.name,
    kind: node.kind,
    file: node.file,
    relation: edge.kind,
  };
  if (node.evidence?.startLine !== undefined)
    out.line = node.evidence.startLine;
  return out;
}

function refs(
  graph: TtscGraphMemory,
  edges: readonly ITtscGraphEdge[],
  end: "to" | "from",
  limit: number,
): { items: ITtscGraphIndex.IReference[]; truncated: boolean } {
  const items: ITtscGraphIndex.IReference[] = [];
  const seen = new Set<string>();
  let available = 0;
  for (const edge of edges) {
    if (STRUCTURAL_KINDS.has(edge.kind)) continue;
    const other = graph.node(end === "to" ? edge.to : edge.from);
    if (other === undefined || other.kind === "file") continue;
    const key = `${edge.kind}:${other.id}`;
    if (seen.has(key)) continue;
    seen.add(key);
    available++;
    if (items.length < limit) items.push(refOf(other, edge));
  }
  return { items, truncated: available > items.length };
}

function directMentions(query: string): string[] {
  const handles = new Set<string>();
  for (const match of query.matchAll(/`([^`]+)`/g)) {
    const handle = normalizeHandle(match[1] ?? "");
    if (handle !== undefined) handles.add(handle);
  }
  for (const match of query.matchAll(
    /\b[A-Za-z_$][\w$]*(?:\.[A-Za-z_$][\w$]*)+\b/g,
  )) {
    const handle = normalizeHandle(match[0]);
    if (handle !== undefined) handles.add(handle);
  }
  return [...handles];
}

function normalizeHandle(raw: string): string | undefined {
  const value = raw.trim();
  return /^[A-Za-z_$][\w$]*(?:\.[A-Za-z_$][\w$]*)*$/.test(value)
    ? value
    : undefined;
}

function bound(
  value: number | undefined,
  fallback: number,
  min: number,
  max: number,
): number {
  const n = value === undefined || !Number.isFinite(value) ? fallback : value;
  return Math.max(min, Math.min(max, Math.floor(n)));
}
