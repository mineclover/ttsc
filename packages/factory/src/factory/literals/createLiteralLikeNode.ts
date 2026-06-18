import type {
  BigIntLiteral,
  NoSubstitutionTemplateLiteral,
  NumericLiteral,
  RegularExpressionLiteral,
  StringLiteral,
} from "../../ast";
import { createRegularExpressionLiteral } from "../expressions/createRegularExpressionLiteral";
import { createBigIntLiteral } from "./createBigIntLiteral";
import { createNoSubstitutionTemplateLiteral } from "./createNoSubstitutionTemplateLiteral";
import { createNumericLiteral } from "./createNumericLiteral";
import { createStringLiteral } from "./createStringLiteral";

/** The literal node kinds dispatched by {@link createLiteralLikeNode}. */
export type LiteralLikeNodeKind =
  | "NumericLiteral"
  | "BigIntLiteral"
  | "StringLiteral"
  | "RegularExpressionLiteral"
  | "NoSubstitutionTemplateLiteral";

/** The literal nodes produced by {@link createLiteralLikeNode}. */
export type LiteralLikeNode =
  | NumericLiteral
  | BigIntLiteral
  | StringLiteral
  | RegularExpressionLiteral
  | NoSubstitutionTemplateLiteral;

/**
 * Create a literal node by kind: a dispatcher that forwards to the matching
 * literal factory.
 *
 * The `kind` selects which literal factory runs and `text` is passed through as
 * its content. The supported kinds are `NumericLiteral`, `BigIntLiteral`,
 * `StringLiteral`, `RegularExpressionLiteral`, and
 * `NoSubstitutionTemplateLiteral`. Each delegate keeps its own behavior, so a
 * `BigIntLiteral` still gains the trailing `n` and a `StringLiteral` is still
 * double-quoted by default.
 *
 * Jsx text kinds from the legacy compiler are intentionally omitted; only the
 * literal kinds modeled by this package are handled.
 *
 * With `kind` of `StringLiteral` and `text` of `hi`, this prints:
 *
 * ```ts
 * "hi";
 * ```
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param kind The literal node kind.
 * @param text The literal text.
 * @returns The created literal node.
 */
export const createLiteralLikeNode = (
  kind: LiteralLikeNodeKind,
  text: string,
): LiteralLikeNode => {
  switch (kind) {
    case "NumericLiteral":
      return createNumericLiteral(text);
    case "BigIntLiteral":
      return createBigIntLiteral(text);
    case "StringLiteral":
      return createStringLiteral(text);
    case "RegularExpressionLiteral":
      return createRegularExpressionLiteral(text);
    case "NoSubstitutionTemplateLiteral":
      return createNoSubstitutionTemplateLiteral(text);
  }
};
