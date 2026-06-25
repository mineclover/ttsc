import { ITtscGraphEvidence } from "./ITtscGraphEvidence";
import { TtscGraphConfidence } from "./TtscGraphConfidence";
import { TtscGraphEdgeKind } from "./TtscGraphEdgeKind";
import { TtscGraphProvenance } from "./TtscGraphProvenance";

/**
 * A directed relationship from one {@link ITtscGraphNode} to another, both named
 * by `id`. The triple `(from, to, kind)` is unique; a repeated relationship
 * keeps the first source-order evidence.
 *
 * `provenance` and `confidence` are mandatory so every edge declares how it was
 * derived and how far to trust it: a `calls` edge the checker resolved is
 * `checker-resolved`/`high`. The `framework-derived` and `heuristic`
 * provenances are reserved for inference layers a consumer may add on top.
 */
export interface ITtscGraphEdge {
  /** Node id the relationship originates from. */
  from: string;

  /** Node id the relationship points to. */
  to: string;

  /** The relationship kind. */
  kind: TtscGraphEdgeKind;

  /** How the edge was derived. */
  provenance: TtscGraphProvenance;

  /** How much to trust the edge. */
  confidence: TtscGraphConfidence;

  /** The source expression that produced the edge, for display and expansion. */
  evidence?: ITtscGraphEvidence;
}
