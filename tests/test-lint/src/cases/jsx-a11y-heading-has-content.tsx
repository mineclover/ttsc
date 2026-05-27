/**
 * Verifies jsx-a11y/heading-has-content: heading tags need announcable
 * content.
 *
 * Pins the "empty heading" branch — assistive tech cannot announce
 * what does not exist, so the rule rejects an `<h1>` (or peer heading)
 * with no children.
 *
 * 1. Render an `<h1>` with no children.
 * 2. Lint flags the empty heading.
 */
export const X = () => <h1 />;
