import type { CallChain, CallExpression, Expression } from "../../ast";
import { createGlobalMethodCall } from "./createGlobalMethodCall";

/**
 * Create an `Object.defineProperty(target, propertyName, attributes)` call.
 *
 * Builds the global `Object.defineProperty` method call from the three given
 * argument expressions, in order. The result is a chain call when `target`
 * carries optional-chain context, otherwise a plain call.
 *
 * With `target` of `obj`, `propertyName` of `"x"`, and `attributes` of `desc`,
 * the printer emits:
 *
 * ```ts
 * Object.defineProperty(obj, "x", desc);
 * ```
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param target The target object.
 * @param propertyName The property name expression.
 * @param attributes The property descriptor expression.
 * @returns The created {@link CallExpression} or {@link CallChain}.
 */
export const createObjectDefinePropertyCall = (
  target: Expression,
  propertyName: Expression,
  attributes: Expression,
): CallExpression | CallChain =>
  createGlobalMethodCall("Object", "defineProperty", [
    target,
    propertyName,
    attributes,
  ]);
