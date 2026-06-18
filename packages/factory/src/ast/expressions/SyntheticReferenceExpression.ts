import type { Expression } from "./Expression";

/**
 * A synthetic expression that pairs a value expression with the `this` argument
 * it should be invoked against. It emits as its underlying expression.
 *
 * Built by {@link factory.createSyntheticReferenceExpression}.
 *
 * @author Jeongho Nam - https://github.com/samchon
 */
export interface SyntheticReferenceExpression {
  /** Discriminant tag; always `"SyntheticReferenceExpression"`. */
  kind: "SyntheticReferenceExpression";

  /** The underlying expression. */
  expression: Expression;

  /** The captured `this` argument. */
  thisArg: Expression;
}
