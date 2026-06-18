import type { Node } from "../Node";

/**
 * A synthetic list of sibling nodes that has no syntax of its own.
 *
 * Built by {@link factory.createSyntaxList}.
 *
 * @author Jeongho Nam - https://github.com/samchon
 */
export interface SyntaxList {
  /** Discriminant tag; always `"SyntaxList"`. */
  kind: "SyntaxList";

  /** The child nodes. */
  children: readonly Node[];
}
