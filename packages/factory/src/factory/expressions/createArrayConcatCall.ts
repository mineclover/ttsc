import type { CallChain, CallExpression, Expression } from "../../ast";
import { createMethodCall } from "./createMethodCall";

/**
 * Create a call to the `concat` method on an array: `array.concat(...)`.
 *
 * Delegates to {@link createMethodCall}, so the result is a {@link CallChain}
 * when `array` is itself an optional chain and a plain {@link CallExpression}
 * otherwise. The arguments are printed comma separated inside the parentheses.
 *
 * Given array `a` and arguments `b`, `c`, the printer emits:
 *
 * ```ts
 * a.concat(b, c);
 * ```
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param array The array expression.
 * @param argumentsList The arguments to concatenate.
 * @returns The created {@link CallExpression} or {@link CallChain}.
 */
export const createArrayConcatCall = (
  array: Expression,
  argumentsList: readonly Expression[],
): CallExpression | CallChain =>
  createMethodCall(array, "concat", argumentsList);
