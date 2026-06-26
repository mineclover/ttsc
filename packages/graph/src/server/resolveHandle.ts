import { TtscGraphMemory } from "../model/TtscGraphMemory";
import { ITtscGraphNode } from "../structures/ITtscGraphNode";

export interface IResolvedGraphHandle {
  node?: ITtscGraphNode;
  candidates?: ITtscGraphNode[];
}

/** Resolve a tool handle as an id, exact symbol name, or dotted suffix. */
export function resolveGraphHandle(
  graph: TtscGraphMemory,
  handle: string,
  candidateLimit = 12,
): IResolvedGraphHandle {
  const byId = graph.node(handle);
  if (byId !== undefined) return { node: byId };

  const exact = graph.symbols(handle);
  if (exact.length === 1) return { node: exact[0] };
  if (exact.length > 1) return { candidates: exact.slice(0, candidateLimit) };

  if (handle.includes(".")) {
    const suffix = `.${handle}`;
    const suffixMatches = graph.nodes.filter(
      (node) =>
        node.kind !== "file" && node.qualifiedName?.endsWith(suffix) === true,
    );
    if (suffixMatches.length === 1) return { node: suffixMatches[0] };
    if (suffixMatches.length > 1) {
      return { candidates: suffixMatches.slice(0, candidateLimit) };
    }
  }
  return {};
}
