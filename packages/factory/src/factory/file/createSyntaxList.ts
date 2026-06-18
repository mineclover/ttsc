import type { Node, SyntaxList } from "../../ast";
import { make } from "../internal/make";

/**
 * Create a {@link SyntaxList}: a synthetic node that holds a flat sequence of
 * child nodes.
 *
 * This is a simplified, emit-internal node. The printer renders it by emitting
 * each child and joining them with a single space, with no surrounding
 * punctuation. It is mostly useful as a generic container where the legacy
 * compiler would group adjacent nodes.
 *
 * Given identifiers `a` and `b` as children, this prints:
 *
 * ```ts
 * a b
 * ```
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param children The child nodes.
 * @returns The created {@link SyntaxList}.
 */
export const createSyntaxList = (children: readonly Node[]): SyntaxList =>
  make("SyntaxList", { children });
