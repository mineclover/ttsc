import { ITtscGraphDecorator } from "./ITtscGraphDecorator";

/** Targeted symbol lookup when a concrete name or handle is being resolved. */
export interface ITtscGraphLookup {
  /** Discriminator for targeted symbol lookup. */
  type: "lookup";

  hits: ITtscGraphLookup.IHit[];

  /** How to use this source-free result before another tool or final answer. */
  guide: string;
}
export namespace ITtscGraphLookup {
  /** Find a concrete class, method, function, property, type, or dotted handle. */
  export interface IRequest {
    /** Discriminator for targeted symbol lookup. */
    type: "lookup";

    /**
     * What to find, in natural language and code vocabulary mixed freely: a
     * symbol name, a dotted member (`OrderService.create`), or a phrase
     * (`shopping order create`, `repository find relations`). Exact names are
     * not required, but this is not a second broad entrypoints call. Use it
     * when a named handle is missing or ambiguous.
     */
    query: string;

    /**
     * Maximum hits to return.
     *
     * Prefer the default. Large hit lists usually mean the query is too broad;
     * refine the name instead of raising this.
     *
     * @default 5
     */
    limit?: number;
  }

  /** One ranked hit with a handle to follow via `details` or `trace`. */
  export interface IHit {
    id: string;
    name: string;
    kind: string;
    file: string;
    /** 1-based declaration line, when known. */
    line?: number;
    /**
     * The hit's declaration signature, so you can often answer without
     * requesting details.
     */
    signature?: string;
    /** Decorators written on this declaration, when any. */
    decorators?: ITtscGraphDecorator[];
    /** Relative relevance; higher is a better match. */
    score: number;
  }
}
