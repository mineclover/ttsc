/**
 * Verifies jsx-a11y/lang: `<html lang>` must be a valid BCP-47 tag.
 *
 * Pins the "present but invalid tag" branch — this rule is the
 * superset of `html-has-lang` that also catches a non-empty `lang`
 * whose value is not a recognized IETF tag, so `lang="foo"` is rejected.
 *
 * 1. Render an `<html>` element with `lang="foo"`.
 * 2. Lint flags the invalid language tag.
 */
export const X = () => <html lang="foo" />;
