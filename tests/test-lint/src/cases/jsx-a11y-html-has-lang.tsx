/**
 * Verifies jsx-a11y/html-has-lang: `<html>` needs a non-empty `lang`.
 *
 * Pins the "missing lang" branch — screen readers fall back to the
 * user's default voice without a `lang` attribute, which mispronounces
 * content, so the rule requires the attribute to exist and be non-empty.
 *
 * 1. Render an `<html>` element with no `lang` attribute.
 * 2. Lint flags the missing language declaration.
 */
export const X = () => <html />;
