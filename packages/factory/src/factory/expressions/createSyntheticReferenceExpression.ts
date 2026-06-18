import type { Expression, SyntheticReferenceExpression } from "../../ast";
import { make } from "../internal/make";

/**
 * Create a {@link SyntheticReferenceExpression}: a transform helper that pairs
 * an expression with a captured `this` argument.
 *
 * `expression` is the underlying reference and `thisArg` is the `this` value
 * captured alongside it, for example when a method call is split apart during
 * lowering. Only `expression` is emitted; `thisArg` is bookkeeping the printer
 * ignores.
 *
 * With `expression` of `a` and any `thisArg`, the printer emits:
 *
 * ```ts
 * a;
 * ```
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param expression The underlying expression.
 * @param thisArg The captured `this` argument.
 * @returns The created {@link SyntheticReferenceExpression}.
 */
export const createSyntheticReferenceExpression = (
  expression: Expression,
  thisArg: Expression,
): SyntheticReferenceExpression =>
  make("SyntheticReferenceExpression", { expression, thisArg });
