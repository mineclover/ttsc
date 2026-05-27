/**
 * Verifies testing-library/no-promise-in-fire-event: Promise-returning args are rejected.
 *
 * Pins the argument traversal inside `fireEvent.*`: passing an awaited
 * `findBy*` (or any awaited expression) targets a Promise value that
 * `fireEvent` cannot use.
 *
 * 1. Import `fireEvent` and `screen` from Testing Library.
 * 2. Pass `await screen.findByRole(...)` into `fireEvent.click(...)`.
 * 3. Assert the matching diagnostic.
 */
import { fireEvent, screen } from "@testing-library/react";

async function testCase() {
  // expect: testing-library/no-promise-in-fire-event error
  fireEvent.click(await screen.findByRole("button"));
}

void testCase;
