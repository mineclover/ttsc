/**
 * Verifies testing-library/prefer-user-event: `fireEvent` calls are rejected when user-event covers them.
 *
 * Pins the `fireEvent.<name>` check: events such as `click` have an equivalent
 * user-event method that better simulates real user interactions.
 *
 * 1. Import `fireEvent` and `screen` from Testing Library.
 * 2. Call `fireEvent.click(...)`.
 * 3. Assert the matching diagnostic.
 */
import { fireEvent, screen } from "@testing-library/react";

function testCase() {
  // expect: testing-library/prefer-user-event error
  fireEvent.click(screen.getByText("Save"));
}

void testCase;
