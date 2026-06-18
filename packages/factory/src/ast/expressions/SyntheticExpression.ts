import type { Node } from "../Node";
import type { TypeNode } from "../types/TypeNode";

/**
 * A synthetic placeholder expression that stands in for a value of a known type
 * during transformation. It emits nothing.
 *
 * Built by {@link factory.createSyntheticExpression}.
 *
 * @author Jeongho Nam - https://github.com/samchon
 */
export interface SyntheticExpression {
  /** Discriminant tag; always `"SyntheticExpression"`. */
  kind: "SyntheticExpression";

  /** The type this placeholder stands in for. */
  type: TypeNode;

  /** Whether the placeholder represents a spread element. */
  isSpread: boolean;

  /** The node the tuple element name was derived from, if any. */
  tupleNameSource?: Node;
}
