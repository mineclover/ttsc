/**
 * Verifies testing-library/prefer-screen-queries: destructured render-result queries are rejected.
 *
 * Locks the destructuring branch: a `getByText` (etc.) bound from `render()`
 * should be replaced by the equivalent `screen.*` query.
 *
 * 1. Import `render` from Testing Library.
 * 2. Destructure `getByText` from a `render(...)` call and invoke it.
 * 3. Assert the matching diagnostic on the destructured call.
 */
import { render } from "@testing-library/react";

declare const node: unknown;

function testCase() {
  const { getByText } = render(node as never);
  // expect: testing-library/prefer-screen-queries error
  getByText("Save");
}

void testCase;
