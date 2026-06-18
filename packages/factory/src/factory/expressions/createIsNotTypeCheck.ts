import type { BinaryExpression, Expression } from "../../ast";
import { createStringLiteral } from "../literals/createStringLiteral";
import { createNull } from "../names/createNull";
import { createStrictInequality } from "./createStrictInequality";
import { createTypeOfExpression } from "./createTypeOfExpression";
import { createVoidZero } from "./createVoidZero";

/**
 * Create a negated runtime type check: a strict-inequality
 * {@link BinaryExpression} that is true when `value` is not of the given kind.
 *
 * For the tag `"null"` it compares `value !== null`, and for `"undefined"` it
 * compares `value !== void 0`. For any other tag it compares `typeof value !==
 * "tag"`, built from {@link createTypeOfExpression} and a string literal. All
 * three forms use {@link createStrictInequality}, so the printer surrounds the
 * `!==` operator with a single space on each side.
 *
 * Given value `value` and tag `"string"`, the printer emits:
 *
 * ```ts
 * typeof value !== "string";
 * ```
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param value The value to test.
 * @param tag The `typeof` tag (or `"null"` / `"undefined"`).
 * @returns The created {@link BinaryExpression}.
 */
export const createIsNotTypeCheck = (
  value: Expression,
  tag: string,
): BinaryExpression =>
  tag === "null"
    ? createStrictInequality(value, createNull())
    : tag === "undefined"
      ? createStrictInequality(value, createVoidZero())
      : createStrictInequality(
          createTypeOfExpression(value),
          createStringLiteral(tag),
        );
