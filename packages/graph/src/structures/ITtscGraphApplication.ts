import { ITtscGraphExpand } from "./ITtscGraphExpand";
import { ITtscGraphIndex } from "./ITtscGraphIndex";
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
   * The first source-free index for a code question: ranked symbols, exact code
   * handles mentioned in the query, declaration signatures, and direct
   * dependency context. Use this before reading source.
   *
   * @param props The natural code question or search phrase
   * @returns Compact graph coordinates and dependency context
   */
  graph_index(props: ITtscGraphIndex.IProps): ITtscGraphIndex;

  /**
   * The project's architecture — folder layers, dependency hotspots, and the
   * public API. Call first to orient on an unfamiliar codebase.
   *
   * @param props Which facet to project
   * @returns The requested architecture facets
   */
  graph_overview(props: ITtscGraphOverview.IProps): ITtscGraphOverview;

  /**
   * The declared shape of the given symbols: each one's signature, and for a
   * class/interface/namespace its members. Handles may be ids or dotted symbol
   * names. Set `source: true` to also read a specific body, `neighbors: true`
   * to list what it uses and what uses it.
   *
   * @param props The handles to expand
   * @returns The resolved nodes, and any handles that did not resolve
   */
  graph_expand(props: ITtscGraphExpand.IProps): ITtscGraphExpand;

  /**
   * Find any symbol — class, function, method, or field — by name or
   * description. Each hit comes with its signature, so the query alone often
   * answers the question, and `next.expand` gives handles for source
   * follow-up.
   *
   * @param props The query and result cap
   * @returns Ranked hits with handles
   */
  graph_query(props: ITtscGraphQuery.IProps): ITtscGraphQuery;

  /**
   * Follow dependency flow from a symbol: `forward` to what it uses, `reverse`
   * to what uses it, or `impact` to the public API and tests a change reaches.
   * Give `from`/`to` as ids or dotted names to get the path between two symbols
   * in one call — how A reaches B.
   *
   * @param props The start, optional target, direction, and bounds
   * @returns The ordered hops and reached nodes, or candidates for an ambiguous
   *   start
   */
  graph_trace(props: ITtscGraphTrace.IProps): ITtscGraphTrace;
}
