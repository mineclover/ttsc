import type { CallChain, CallExpression, Expression } from "../../ast";
import { createIdentifier } from "../names/createIdentifier";
import { createMethodCall } from "./createMethodCall";

/**
 * Create a call to a method of a global object: `GlobalObject.method(...)`.
 *
 * The global object name is turned into an {@link Identifier} and the call is
 * built with {@link createMethodCall}, so the result is a {@link CallChain} or a
 * plain {@link CallExpression} depending on the receiver. The arguments are
 * printed comma separated inside the parentheses.
 *
 * Given global object `"Object"`, method `"keys"` and argument `obj`, the
 * printer emits:
 *
 * ```ts
 * Object.keys(obj);
 * ```
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param globalObjectName The global object name (e.g. `"Object"`).
 * @param methodName The method name.
 * @param argumentsList The call arguments.
 * @returns The created {@link CallExpression} or {@link CallChain}.
 */
export const createGlobalMethodCall = (
  globalObjectName: string,
  methodName: string,
  argumentsList: readonly Expression[],
): CallExpression | CallChain =>
  createMethodCall(
    createIdentifier(globalObjectName),
    methodName,
    argumentsList,
  );
