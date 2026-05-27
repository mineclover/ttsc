/**
 * Verifies jsx-a11y/anchor-has-content: anchors need accessible content.
 *
 * Pins the "empty anchor" branch — an `<a>` with no text, no
 * `aria-label`, and no labelled child gives assistive tech nothing to
 * announce, so the rule rejects it.
 *
 * 1. Render an `<a>` with an `href` but no children or label.
 * 2. Lint flags the empty anchor.
 */
export const X = () => <a href="/docs" />;
