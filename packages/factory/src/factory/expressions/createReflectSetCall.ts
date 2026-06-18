import type { CallChain, CallExpression, Expression } from "../../ast";
import { createGlobalMethodCall } from "./createGlobalMethodCall";

/**
 * Create a `Reflect.set(target, propertyKey, value, receiver?)` call.
 *
 * Builds the global `Reflect.set` method call. The fourth `receiver` argument
 * is only included when supplied, so the call has three or four arguments
 * accordingly. The result is a chain call when `target` carries optional-chain
 * context, otherwise a plain call.
 *
 * With `target` of `t`, `propertyKey` of `"k"`, and `value` of `v`, the printer
 * emits:
 *
 * ```ts
 * Reflect.set(t, "k", v);
 * ```
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param target The target object.
 * @param propertyKey The property key expression.
 * @param value The value to set.
 * @param receiver The optional receiver expression.
 * @returns The created {@link CallExpression} or {@link CallChain}.
 */
export const createReflectSetCall = (
  target: Expression,
  propertyKey: Expression,
  value: Expression,
  receiver?: Expression,
): CallExpression | CallChain =>
  createGlobalMethodCall(
    "Reflect",
    "set",
    receiver
      ? [target, propertyKey, value, receiver]
      : [target, propertyKey, value],
  );
