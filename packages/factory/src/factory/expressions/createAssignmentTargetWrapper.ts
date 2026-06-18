import type {
  Expression,
  Identifier,
  PropertyAccessExpression,
} from "../../ast";
import { createParameterDeclaration } from "../clauses/createParameterDeclaration";
import { createSetAccessorDeclaration } from "../declarations/createSetAccessorDeclaration";
import { createBlock } from "../statements/createBlock";
import { createExpressionStatement } from "../statements/createExpressionStatement";
import { createObjectLiteralExpression } from "./createObjectLiteralExpression";
import { createParenthesizedExpression } from "./createParenthesizedExpression";
import { createPropertyAccessExpression } from "./createPropertyAccessExpression";

/**
 * Create an assignment-target wrapper: a parenthesized object literal with a
 * single `value` setter whose `.value` property is then accessed.
 *
 * This turns an arbitrary expression into a valid assignment target. Assigning
 * to the resulting `.value` property runs the setter body, which evaluates
 * `expression`. The object literal is wrapped in explicit parentheses to work
 * around a V8 regression. The result is a {@link PropertyAccessExpression} built
 * from {@link createParenthesizedExpression},
 * {@link createObjectLiteralExpression} and a
 * {@link createSetAccessorDeclaration}.
 *
 * Given setter parameter `v` and expression `expr`, the printer emits:
 *
 * ```ts
 * ({
 *   set value(v) {
 *     expr;
 *   },
 * }).value;
 * ```
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param paramName The setter parameter name.
 * @param expression The expression evaluated on assignment.
 * @returns The created {@link PropertyAccessExpression}.
 */
export const createAssignmentTargetWrapper = (
  paramName: Identifier,
  expression: Expression,
): PropertyAccessExpression =>
  createPropertyAccessExpression(
    createParenthesizedExpression(
      createObjectLiteralExpression([
        createSetAccessorDeclaration(
          undefined,
          "value",
          [
            createParameterDeclaration(
              undefined,
              undefined,
              paramName,
              undefined,
              undefined,
              undefined,
            ),
          ],
          createBlock([createExpressionStatement(expression)]),
        ),
      ]),
    ),
    "value",
  );
