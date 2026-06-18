import type { CallChain, CallExpression, Expression } from "../../ast";
import { createMethodCall } from "./createMethodCall";

/**
 * Create a call to `Function.prototype.call`: `target.call(thisArg,
 * ...argumentsList)`.
 *
 * Delegates to {@link createMethodCall}, so the result is a {@link CallChain}
 * when `target` is itself an optional chain and a plain {@link CallExpression}
 * otherwise. The `thisArg` becomes the first argument and the call arguments
 * follow in order.
 *
 * Given target `fn`, `this` argument `thisArg` and arguments `a`, `b`, the
 * printer emits:
 *
 * ```ts
 * fn.call(thisArg, a, b);
 * ```
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param target The function being invoked.
 * @param thisArg The `this` argument.
 * @param argumentsList The call arguments.
 * @returns The created {@link CallExpression} or {@link CallChain}.
 */
export const createFunctionCallCall = (
  target: Expression,
  thisArg: Expression,
  argumentsList: readonly Expression[],
): CallExpression | CallChain =>
  createMethodCall(target, "call", [thisArg, ...argumentsList]);
