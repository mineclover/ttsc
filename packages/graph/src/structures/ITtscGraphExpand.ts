/**
 * The resolved nodes and their source the `graph_expand` tool returns for a set
 * of handles.
 */
export interface ITtscGraphExpand {
  nodes: ITtscGraphExpand.INode[];

  /** Handles that resolved to no node. */
  unknown: string[];
}
export namespace ITtscGraphExpand {
  /** Which handles to expand, and whether to include their neighbors. */
  export interface IProps {
    /**
     * Node ids to expand, exactly as another tool returned them. Pass every
     * handle you need in one call.
     */
    handles: string[];

    /**
     * Also list each node's direct dependencies and dependents (the symbols it
     * uses and the symbols that use it).
     *
     * @default false
     */
    neighbors?: boolean;
  }

  /** One expanded node: its declaration source and optional neighbors. */
  export interface INode {
    id: string;
    name: string;
    kind: string;
    file: string;
    /** The declaration source, sliced from the node's evidence span. */
    source?: string;
    /** True when `source` was cut at the line cap. */
    truncated?: boolean;
    /** Symbols this node uses (outgoing dependency edges). */
    dependsOn?: IReference[];
    /** Symbols that use this node (incoming dependency edges). */
    dependedOnBy?: IReference[];
  }

  /** A dependency neighbor of an expanded node and the edge that links them. */
  export interface IReference {
    id: string;
    name: string;
    kind: string;
    /** The edge kind connecting the two (`calls`, `type_ref`, …). */
    relation: string;
  }
}
