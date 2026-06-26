/** The compact source-free index returned by `graph_index` for a code question. */
export interface ITtscGraphIndex {
  /** The original query the index was built for. */
  query: string;

  /** Ranked symbols relevant to the query. */
  hits: ITtscGraphIndex.IHit[];

  /** Code handles written directly in the query, resolved when possible. */
  mentions: ITtscGraphIndex.IMention[];

  /** Direct dependency context for the resolved mentions and highest hits. */
  neighborhood: ITtscGraphIndex.INeighborhood[];

  /** Follow-up handles for deeper graph calls. */
  next: ITtscGraphIndex.INext;

  /** True when result caps hid additional seeds or references. */
  truncated?: boolean;
}

export namespace ITtscGraphIndex {
  /**
   * Ask the graph for the first index an agent should read before opening
   * source: ranked symbols, exact mentioned handles, and nearby dependency
   * edges.
   */
  export interface IProps {
    /**
     * A natural code question or search phrase. Mix prose with code handles,
     * for example `how Repository.find loads relations` or
     * `SelectQueryBuilder.setFindOptions join aliases`.
     */
    query: string;

    /**
     * Maximum ranked hits to return.
     *
     * @default 8
     */
    limit?: number;

    /**
     * Maximum direct dependencies and dependents to return per indexed symbol.
     *
     * @default 4
     */
    neighbors?: number;
  }

  /** A compact symbol coordinate, optionally with its declaration signature. */
  export interface INode {
    id: string;
    name: string;
    kind: string;
    file: string;
    /** 1-based declaration line, when known. */
    line?: number;
    /** Declaration head, included only for indexed symbols. */
    signature?: string;
  }

  /** One ranked search hit. */
  export interface IHit extends INode {
    /** Relative relevance; higher is a better match. */
    score: number;
  }

  /** A code handle written in the query, with its resolution status. */
  export interface IMention {
    handle: string;
    node?: INode;
    candidates?: INode[];
  }

  /** Direct dependency context around one indexed symbol. */
  export interface INeighborhood extends INode {
    dependsOn: IReference[];
    dependedOnBy: IReference[];
  }

  /** One neighboring symbol and the relationship leading to it. */
  export interface IReference {
    id: string;
    name: string;
    kind: string;
    file: string;
    /** 1-based declaration line, when known. */
    line?: number;
    relation: string;
  }

  /** Tool-call handles suggested by this first index. */
  export interface INext {
    /** Pass these ids to `graph_expand`, with `source: true` only when needed. */
    expand: string[];
    /** Pass these ids to `graph_trace` when following dependency flow. */
    traceFrom: string[];
  }
}
