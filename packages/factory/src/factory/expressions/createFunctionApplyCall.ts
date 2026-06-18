import type { CallChain, CallExpression, Expression } from "../../ast";
import { createMethodCall } from "./createMethodCall";

/**
 * Create a call to `Function.prototype.apply`: `target.apply(thisArg,
 * argumentsExpression)`.
 *
 * Delegates to {@link createMethodCall}, so the result is a {@link CallChain}
 * when `target` is itself an optional chain and a plain {@link CallExpression}
 * otherwise. The `argumentsExpression` is passed as a single array-like
 * argument, matching the `apply` calling convention.
 *
 * Given target `fn`, `this` argument `thisArg` and arguments expression `args`,
 * the printer emits:
 *
 * ```ts
 * fn.apply(thisArg, args);
 * ```
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param target The function being invoked.
 * @param thisArg The `this` argument.
 * @param argumentsExpression The array-like arguments expression.
 * @returns The created {@link CallExpression} or {@link CallChain}.
 */
export const createFunctionApplyCall = (
  target: Expression,
  thisArg: Expression,
  argumentsExpression: Expression,
): CallExpression | CallChain =>
  createMethodCall(target, "apply", [thisArg, argumentsExpression]);
