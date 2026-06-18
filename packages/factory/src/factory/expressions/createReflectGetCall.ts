import type { CallChain, CallExpression, Expression } from "../../ast";
import { createGlobalMethodCall } from "./createGlobalMethodCall";

/**
 * Create a `Reflect.get(target, propertyKey, receiver?)` call.
 *
 * Builds the global `Reflect.get` method call. The third `receiver` argument is
 * only included when supplied, so the call has two or three arguments
 * accordingly. The result is a chain call when `target` carries optional-chain
 * context, otherwise a plain call.
 *
 * With `target` of `t` and `propertyKey` of `"k"`, the printer emits:
 *
 * ```ts
 * Reflect.get(t, "k");
 * ```
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param target The target object.
 * @param propertyKey The property key expression.
 * @param receiver The optional receiver expression.
 * @returns The created {@link CallExpression} or {@link CallChain}.
 */
export const createReflectGetCall = (
  target: Expression,
  propertyKey: Expression,
  receiver?: Expression,
): CallExpression | CallChain =>
  createGlobalMethodCall(
    "Reflect",
    "get",
    receiver ? [target, propertyKey, receiver] : [target, propertyKey],
  );
