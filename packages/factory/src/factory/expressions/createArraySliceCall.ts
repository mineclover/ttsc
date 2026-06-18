import type { CallChain, CallExpression, Expression } from "../../ast";
import { createNumericLiteral } from "../literals/createNumericLiteral";
import { createMethodCall } from "./createMethodCall";

/**
 * Create a call to the `slice` method on an array: `array.slice(start)`.
 *
 * Delegates to {@link createMethodCall}, so the result is a {@link CallChain}
 * when `array` is itself an optional chain and a plain {@link CallExpression}
 * otherwise. A numeric `start` is wrapped with {@link createNumericLiteral};
 * when `start` is omitted the call is emitted with no arguments.
 *
 * Given array `a` and `start` of `1`, the printer emits:
 *
 * ```ts
 * a.slice(1);
 * ```
 *
 * With `start` omitted it emits `a.slice()`.
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param array The array expression.
 * @param start The optional start index.
 * @returns The created {@link CallExpression} or {@link CallChain}.
 */
export const createArraySliceCall = (
  array: Expression,
  start?: number | Expression,
): CallExpression | CallChain =>
  createMethodCall(
    array,
    "slice",
    start === undefined
      ? []
      : [typeof start === "number" ? createNumericLiteral(start) : start],
  );
