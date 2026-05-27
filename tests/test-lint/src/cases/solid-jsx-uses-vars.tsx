/**
 * Verifies solid/jsx-uses-vars: compatibility-only rule emits no diagnostics.
 *
 * Pins the deliberate no-op for `solid/jsx-uses-vars`. ESLint uses this rule
 * to mark JSX identifiers as variable reads for its unused-variable pass;
 * @ttsc/lint accepts the rule for config compatibility but does not run an
 * unused-variable marker pass.
 *
 * 1. Import Solid so the rule family activates.
 * 2. Define a component used only as a JSX tag and assert no diagnostic fires.
 */
import { createSignal } from "solid-js";

createSignal(0);
const Button = () => <button />;

// Negative: no diagnostic; the rule is accepted for config compatibility only.
const tree = <Button />;

JSON.stringify({ tree });
