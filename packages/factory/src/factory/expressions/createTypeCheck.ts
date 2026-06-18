import type { BinaryExpression, Expression } from "../../ast";
import { createStringLiteral } from "../literals/createStringLiteral";
import { createNull } from "../names/createNull";
import { createStrictEquality } from "./createStrictEquality";
import { createTypeOfExpression } from "./createTypeOfExpression";
import { createVoidZero } from "./createVoidZero";

/**
 * Create a runtime type check: a strict equality comparison that tests `value`
 * against a tag.
 *
 * The tag selects the comparison. `"null"` compares `value` directly against
 * `null`. `"undefined"` compares it against `void 0`. Any other tag compares
 * `typeof value` against the tag as a string literal. The result is always a
 * strict equality (`===`) binary expression.
 *
 * With `value` of `x` and `tag` of `"string"`, the printer emits:
 *
 * ```ts
 * typeof x === "string";
 * ```
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param value The value to test.
 * @param tag The `typeof` tag (or `"null"` / `"undefined"`).
 * @returns The created {@link BinaryExpression}.
 */
export const createTypeCheck = (
  value: Expression,
  tag: string,
): BinaryExpression =>
  tag === "null"
    ? createStrictEquality(value, createNull())
    : tag === "undefined"
      ? createStrictEquality(value, createVoidZero())
      : createStrictEquality(
          createTypeOfExpression(value),
          createStringLiteral(tag),
        );
