import type { Expression } from "../../ast";
import { createVoidZero } from "./createVoidZero";

/**
 * The split of a call expression's callee into the value to invoke
 * ({@link target}) and the receiver to bind as `this` ({@link thisArg}).
 */
export interface CallBinding {
  /** The expression to invoke. */
  target: Expression;
  /** The `this` argument to apply. */
  thisArg: Expression;
}

/**
 * Split a callee {@link Expression} into a `{ target, thisArg }` binding,
 * mirroring the legacy `ts.factory.createCallBinding`.
 *
 * This is a source-construction simplification of the compiler's emit helper:
 * the `super` / helper-name / temp-variable-caching paths (which require emit
 * flags and a transformation context) are omitted, so `recordTempVariable`,
 * `languageVersion` and `cacheIdentifiers` are accepted for signature parity
 * but not used. For a property / element access the object becomes the
 * `thisArg`; otherwise `thisArg` is `void 0`.
 *
 * The returned binding is not a node and is not printed directly. For the
 * callee `obj.method`, `target` prints as `obj.method` and `thisArg` prints as
 * the receiver:
 *
 * ```ts
 * obj;
 * ```
 *
 * For a non-access callee `fn`, `target` is `fn` and `thisArg` prints as `void
 * 0`.
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param expression The callee expression.
 * @param recordTempVariable Accepted for parity; unused.
 * @param languageVersion Accepted for parity; unused.
 * @param cacheIdentifiers Accepted for parity; unused.
 * @returns The `{ target, thisArg }` binding.
 */
export const createCallBinding = (
  expression: Expression,
  recordTempVariable?: (node: Expression) => void,
  languageVersion?: number,
  cacheIdentifiers?: boolean,
): CallBinding => {
  void recordTempVariable;
  void languageVersion;
  void cacheIdentifiers;
  if (
    expression.kind === "PropertyAccessExpression" ||
    expression.kind === "ElementAccessExpression"
  )
    return { target: expression, thisArg: expression.expression };
  return { target: expression, thisArg: createVoidZero() };
};
