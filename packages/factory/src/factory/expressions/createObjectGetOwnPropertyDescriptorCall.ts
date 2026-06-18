import type { CallChain, CallExpression, Expression } from "../../ast";
import { createGlobalMethodCall } from "./createGlobalMethodCall";

/**
 * Create an `Object.getOwnPropertyDescriptor(target, propertyName)` call.
 *
 * Builds the global `Object.getOwnPropertyDescriptor` method call from the two
 * given argument expressions. The result is a chain call when `target` carries
 * optional-chain context, otherwise a plain call.
 *
 * With `target` of `obj` and `propertyName` of `"x"`, the printer emits:
 *
 * ```ts
 * Object.getOwnPropertyDescriptor(obj, "x");
 * ```
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param target The target object.
 * @param propertyName The property name expression.
 * @returns The created {@link CallExpression} or {@link CallChain}.
 */
export const createObjectGetOwnPropertyDescriptorCall = (
  target: Expression,
  propertyName: Expression,
): CallExpression | CallChain =>
  createGlobalMethodCall("Object", "getOwnPropertyDescriptor", [
    target,
    propertyName,
  ]);
