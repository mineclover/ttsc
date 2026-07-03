import { ITtscGraphEvidence } from "./ITtscGraphEvidence";

/**
 * Declaration metadata attached to a graph node without interpreting a specific
 * domain convention. Producers emit the written source/name/value facts;
 * consumers decide whether an annotation means a feature boundary, layer,
 * semantic traversal cut, or another tool-specific marker.
 */
export interface ITtscGraphAnnotation {
  /**
   * Where the annotation came from. `jsdoc` is emitted by the native graph
   * producer today; the wider union leaves room for compatible plugin and
   * manifest producers without changing the node shape.
   */
  source: "jsdoc" | "decorator" | "manifest" | "plugin";

  /** The annotation name as written or normalized by its producer. */
  name: string;

  /** Optional grouping namespace, when the producer can identify one. */
  namespace?: string;

  /** Parsed values in source order. Empty when the annotation is marker-only. */
  values: string[];

  /** The source span of the annotation itself, when available. */
  evidence?: ITtscGraphEvidence;
}
