import type {
  CallChain,
  CallExpression,
  Expression,
  Identifier,
} from "../../ast";
import { createCallChain } from "./createCallChain";
import { createCallExpression } from "./createCallExpression";
import { createPropertyAccessChain } from "./createPropertyAccessChain";
import { createPropertyAccessExpression } from "./createPropertyAccessExpression";

/**
 * Create a method call: access a member on `object` and call it with
 * `argumentsList`.
 *
 * The method name is looked up on `object` via a property access, then invoked.
 * When `object` is itself an optional-chain call (`CallChain`), the property
 * access and the call are emitted as chain nodes so the optionality propagates
 * down the chain. Otherwise a plain property access and call are produced.
 *
 * With `object` of `obj`, `methodName` of `doThing`, and arguments `a` and `b`,
 * the printer emits:
 *
 * ```ts
 * obj.doThing(a, b);
 * ```
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param object The receiver object.
 * @param methodName The method name to access and call.
 * @param argumentsList The call arguments.
 * @returns The created {@link CallExpression} or {@link CallChain}.
 */
export const createMethodCall = (
  object: Expression,
  methodName: string | Identifier,
  argumentsList: readonly Expression[],
): CallExpression | CallChain =>
  object.kind === "CallChain"
    ? createCallChain(
        createPropertyAccessChain(object, undefined, methodName),
        undefined,
        undefined,
        argumentsList,
      )
    : createCallExpression(
        createPropertyAccessExpression(object, methodName),
        undefined,
        argumentsList,
      );
