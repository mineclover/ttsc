import type { Node, SyntheticExpression, TypeNode } from "../../ast";
import { make } from "../internal/make";

/**
 * Create a {@link SyntheticExpression}: a type-only placeholder used by the
 * checker, with no source form.
 *
 * The node stands in for a value of `type` during analysis, for example a tuple
 * element. `isSpread` marks it as a spread placeholder, and `tupleNameSource`
 * records the node a tuple element name was derived from. This node carries no
 * printable syntax.
 *
 * Because it is a checker-internal placeholder, the printer emits nothing for
 * it; the output is the empty string:
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param type The type this placeholder stands in for.
 * @param isSpread Whether the placeholder represents a spread element.
 * @param tupleNameSource The node the tuple element name was derived from, if
 *   any.
 * @returns The created {@link SyntheticExpression}.
 */
export const createSyntheticExpression = (
  type: TypeNode,
  isSpread: boolean = false,
  tupleNameSource?: Node,
): SyntheticExpression =>
  make("SyntheticExpression", { type, isSpread, tupleNameSource });
