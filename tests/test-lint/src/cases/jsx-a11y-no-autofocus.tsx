/**
 * Verifies jsx-a11y/no-autofocus: ban the `autoFocus` JSX attribute.
 *
 * Pins the "auto-focus on mount" branch — `autoFocus` steals focus on
 * page load and disorients keyboard and screen-reader users, so the
 * rule rejects the attribute regardless of element.
 *
 * 1. Render an `<input>` with the `autoFocus` attribute.
 * 2. Lint flags the prohibited autofocus.
 */
export const X = () => <input autoFocus />;
