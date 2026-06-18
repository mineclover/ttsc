import type { CallChain, CallExpression, Expression } from "../../ast";
import { createMethodCall } from "./createMethodCall";

/**
 * Create a call to `Function.prototype.bind`: `target.bind(thisArg,
 * ...argumentsList)`.
 *
 * Delegates to {@link createMethodCall}, so the result is a {@link CallChain}
 * when `target` is itself an optional chain and a plain {@link CallExpression}
 * otherwise. The `thisArg` becomes the first argument and the bound arguments
 * follow in order.
 *
 * Given target `fn`, `this` argument `thisArg` and bound arguments `a`, `b`,
 * the printer emits:
 *
 * ```ts
 * fn.bind(thisArg, a, b);
 * ```
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param target The function being bound.
 * @param thisArg The `this` argument.
 * @param argumentsList The bound arguments.
 * @returns The created {@link CallExpression} or {@link CallChain}.
 */
export const createFunctionBindCall = (
  target: Expression,
  thisArg: Expression,
  argumentsList: readonly Expression[],
): CallExpression | CallChain =>
  createMethodCall(target, "bind", [thisArg, ...argumentsList]);
