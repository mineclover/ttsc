import type {
  NoSubstitutionTemplateLiteral,
  TemplateHead,
  TemplateMiddle,
  TemplateTail,
} from "../../ast";
import { createNoSubstitutionTemplateLiteral } from "./createNoSubstitutionTemplateLiteral";
import { createTemplateHead } from "./createTemplateHead";
import { createTemplateMiddle } from "./createTemplateMiddle";
import { createTemplateTail } from "./createTemplateTail";

/**
 * The template literal node kinds dispatched by
 * {@link createTemplateLiteralLikeNode}.
 */
export type TemplateLiteralLikeNodeKind =
  | "TemplateHead"
  | "TemplateMiddle"
  | "TemplateTail"
  | "NoSubstitutionTemplateLiteral";

/** The template literal nodes produced by {@link createTemplateLiteralLikeNode}. */
export type TemplateLiteralLikeNode =
  | TemplateHead
  | TemplateMiddle
  | TemplateTail
  | NoSubstitutionTemplateLiteral;

/**
 * Create a template literal span node by kind: a dispatcher that forwards to
 * the matching template literal factory.
 *
 * The `kind` selects which factory runs and `text` is its cooked content. The
 * supported kinds are `TemplateHead`, `TemplateMiddle`, `TemplateTail`, and
 * `NoSubstitutionTemplateLiteral`. The optional `rawText` reaches only the
 * `TemplateHead` delegate; the middle, tail, and no-substitution factories take
 * cooked text alone.
 *
 * The `templateFlags` parameter from the legacy compiler is accepted for
 * signature parity but ignored, as this package does not model raw template
 * scan flags.
 *
 * With `kind` of `TemplateHead` and `text` of `head`, this prints the opening
 * span up to the first substitution:
 *
 * ```ts
 * `head${
 * ```
 *
 * @author Jeongho Nam - https://github.com/samchon
 * @param kind The template literal node kind.
 * @param text The cooked text.
 * @param rawText The raw text, if any.
 * @param _templateFlags Ignored; kept for signature parity.
 * @returns The created template literal node.
 */
export const createTemplateLiteralLikeNode = (
  kind: TemplateLiteralLikeNodeKind,
  text: string,
  rawText?: string,
  _templateFlags?: number,
): TemplateLiteralLikeNode => {
  switch (kind) {
    case "TemplateHead":
      return createTemplateHead(text, rawText);
    case "TemplateMiddle":
      return createTemplateMiddle(text);
    case "TemplateTail":
      return createTemplateTail(text);
    case "NoSubstitutionTemplateLiteral":
      return createNoSubstitutionTemplateLiteral(text);
  }
};
