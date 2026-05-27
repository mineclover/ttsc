import { test, expect } from "@jest/globals";

test("checks conditionally", () => {
  if (Math.random() > 0.5) {
    // expect: jest/no-conditional-expect error
    expect(true).toBe(true);
  }
});
