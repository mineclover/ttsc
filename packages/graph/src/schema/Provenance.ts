/**
 * How a node or edge was derived — the trust signal that keeps inferred
 * relationships from being read as compiler fact.
 *
 * - `checker-resolved`: the in-process TypeScript-Go checker resolved the symbol
 *   on both ends. This is the graph's core guarantee and the reason it is not a
 *   tree-sitter or text index.
 * - `framework-derived`: synthesized from a framework convention (a NestJS route
 *   decorator, a Next.js file route). Grounded in real syntax but not a checker
 *   symbol resolution.
 * - `heuristic`: an opt-in best-effort inference (a callback or event bridge).
 *   Never part of a default trace; always visibly marked.
 */
export type Provenance = "checker-resolved" | "framework-derived" | "heuristic";
