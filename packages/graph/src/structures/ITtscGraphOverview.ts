/**
 * A compact, source-read-free architecture map of the project — the result of
 * the `graph_overview` tool.
 */
export interface ITtscGraphOverview {
  /** Absolute project root. */
  project: string;

  /** Size of the graph. */
  counts: ITtscGraphOverview.ICounts;

  /** Folder layering, largest first. */
  layers?: ITtscGraphOverview.ILayer[];

  /** Highest-dependency symbols, busiest first. */
  hotspots?: ITtscGraphOverview.IHotspot[];

  /** Export surface by file, widest first. */
  publicApi?: ITtscGraphOverview.IPublicApi[];
}
export namespace ITtscGraphOverview {
  /** Which architecture facets `graph_overview` should return. */
  export interface IProps {
    /**
     * The facet to project, or `all` for every facet. `layers` is the folder
     * layering, `hotspots` the highest-dependency symbols, `publicApi` the
     * export surface.
     *
     * @default "all"
     */
    aspect?: "all" | "layers" | "hotspots" | "publicApi";
  }

  /** Size of the graph by node/edge totals and per-kind node counts. */
  export interface ICounts {
    files: number;
    nodes: number;
    edges: number;
    /** Node count per kind. */
    byKind: Record<string, number>;
  }

  /** One folder layer: its source files and export surface. */
  export interface ILayer {
    /** Directory, project-relative. */
    dir: string;
    /** Distinct source files under it. */
    files: number;
    /** Exported symbols declared under it. */
    exported: number;
  }

  /** A high-dependency symbol with its non-structural fan-in and fan-out. */
  export interface IHotspot {
    id: string;
    name: string;
    kind: string;
    file: string;
    /** Non-structural edges pointing at this symbol. */
    fanIn: number;
    /** Non-structural edges leaving this symbol. */
    fanOut: number;
  }

  /** The exported symbols a single file contributes to the public surface. */
  export interface IPublicApi {
    file: string;
    /** Exported symbol names declared in the file (capped). */
    symbols: string[];
  }
}
