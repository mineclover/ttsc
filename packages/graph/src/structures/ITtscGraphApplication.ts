import { ITtscGraphExpand } from "./ITtscGraphExpand";
import { ITtscGraphOverview } from "./ITtscGraphOverview";
import { ITtscGraphQuery } from "./ITtscGraphQuery";
import { ITtscGraphTrace } from "./ITtscGraphTrace";

/**
 * The MCP tool surface of `@ttsc/graph`, as a typed application.
 *
 * Each method is one MCP tool; its name is the tool name and its parameter
 * object becomes the tool's JSON schema once `typia.llm.controller` reflects
 * this interface. `TtscGraphApplication` implements it over the resident
 * graph.
 */
export interface ITtscGraphApplication {
  /**
   * A compact architecture map — folder layers, the dependency hotspots, and the
   * public API (exported types, ranked) — with no source read. Call this first
   * on an unfamiliar project instead of listing files or opening the entry
   * module: the layers are the structure and the public API is the entry surface.
   *
   * @param props Which facet to project
   * @returns The requested architecture facets
   */
  graph_overview(props: ITtscGraphOverview.IProps): ITtscGraphOverview;

  /**
   * Get the declared shape of symbols the graph located, by their handles: each
   * one's signature, and for a class/interface/namespace/file its member outline
   * (every member with its own signature) — the resolved structure, not inlined
   * source. This is the graph's edge over reading files; it is compiler-resolved
   * and authoritative, so don't grep or reopen the file to confirm it. Pass every
   * handle in ONE call. Set `source: true` only for the few leaf methods or
   * functions whose actual control-flow logic you must read; `neighbors: true`
   * adds what each symbol uses and what uses it.
   *
   * @param props The handles to expand
   * @returns The resolved nodes with source, and any handles that did not
   *   resolve
   */
  graph_expand(props: ITtscGraphExpand.IProps): ITtscGraphExpand;

  /**
   * Find any symbol in this project — a class, function, method, even a single
   * field — by name or a plain-language description, instead of grepping or
   * listing files. The project is fully indexed, so this resolves what `rg` would
   * but ranked by name, subword, path, and centrality, and returns handles to
   * expand or trace. Reach for this first when you need to locate code.
   *
   * @param props The query and result cap
   * @returns Ranked hits with handles
   */
  graph_query(props: ITtscGraphQuery.IProps): ITtscGraphQuery;

  /**
   * Follow dependency flow from a symbol — `forward` to what it uses, `reverse`
   * to what uses it, `impact` to the public API and tests a change reaches —
   * instead of chasing calls by grep. Real call/type edges only (structural
   * edges excluded); returns ordered hops with handles.
   *
   * @param props The start, direction, and bounds
   * @returns The ordered hops and reached nodes, or candidates for an ambiguous
   *   start
   */
  graph_trace(props: ITtscGraphTrace.IProps): ITtscGraphTrace;
}
