/**
 * Verifies jsx-a11y/autocomplete-valid: `autocomplete` token must match
 * the HTML vocabulary and the surrounding `type`.
 *
 * Pins the "type/token mismatch" branch — `autocomplete="url"` and
 * `type="email"` belong to disjoint vocabularies, so the rule rejects
 * the pair.
 *
 * 1. Render an `<input type="email">` with `autocomplete="url"`.
 * 2. Lint flags the token as invalid for the input type.
 */
export const X = () => <input type="email" autoComplete="url" />;
