import { TestValidator } from "@nestia/e2e";
import factory, { type Expression, SyntaxKind } from "@ttsc/factory";
import ts from "ts-legacy";

import { id, print, reparse } from "../../internal/helpers";

const call = (expression: Expression): Expression =>
  factory.createCallExpression(expression, undefined, []);
const construct = (target: Expression): Expression =>
  factory.createNewExpression(target, undefined, []);

/**
 * Verifies new-expression parenthesizer: sees a call through non-null
 * assertions, template tags, and optional chains on the target's spine.
 *
 * Boundary cases of the stop-at-call walk in `TsPrinter.newExpressionTarget`:
 * the walk must traverse every construct that keeps the call on the printed
 * left edge, not just plain accesses. `new f()!.bar()`, `new f()`x`()`, and the
 * illegal `new f()?.bar()` would each re-bind (or reject) the arguments; the
 * tag-over-identifier twin guards against wrapping every tagged template used
 * as a target.
 *
 * 1. Print `new` expressions targeting `f()!.bar` (non-null in the chain), `f()`x`
 *    ` (tagged template whose tag is a call), and `f()?.bar` (optional chain
 *    over a call — a runtime `TypeError` shape, but the printed text must still
 *    parse back to the same AST).
 * 2. Assert each target is parenthesized, and the tag-over-identifier twin `tag`x`
 *    ` stays bare.
 * 3. Re-parse each output with the legacy compiler and assert the top-level
 *    expression is still a `NewExpression`.
 */
export const test_new_expression_target_spine_operator_parentheses =
  (): void => {
    const printed = {
      "non-null in the chain": print(
        construct(
          factory.createPropertyAccessExpression(
            factory.createNonNullExpression(call(id("f"))),
            "bar",
          ),
        ),
      ),
      "template tag is a call": print(
        construct(
          factory.createTaggedTemplateExpression(
            call(id("f")),
            undefined,
            factory.createNoSubstitutionTemplateLiteral("x"),
          ),
        ),
      ),
      "optional chain over a call": print(
        construct(
          factory.createPropertyAccessChain(
            call(id("f")),
            factory.createToken(SyntaxKind.QuestionDotToken),
            id("bar"),
          ),
        ),
      ),
      "template tag is an identifier": print(
        construct(
          factory.createTaggedTemplateExpression(
            id("tag"),
            undefined,
            factory.createNoSubstitutionTemplateLiteral("x"),
          ),
        ),
      ),
    };
    TestValidator.equals(
      "non-null in the chain",
      printed["non-null in the chain"],
      "new (f()!.bar)()",
    );
    TestValidator.equals(
      "template tag is a call",
      printed["template tag is a call"],
      "new (f()`x`)()",
    );
    TestValidator.equals(
      "optional chain over a call",
      printed["optional chain over a call"],
      "new (f()?.bar)()",
    );
    TestValidator.equals(
      "template tag is an identifier",
      printed["template tag is an identifier"],
      "new tag`x`()",
    );
    for (const [title, source] of Object.entries(printed))
      TestValidator.equals(
        `${title} re-parses as new`,
        reparse(source).kind,
        ts.SyntaxKind.NewExpression,
      );
  };
