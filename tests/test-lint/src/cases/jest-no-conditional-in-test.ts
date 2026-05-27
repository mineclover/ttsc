import { test, expect } from "@jest/globals";

test("branches", () => {
  // expect: jest/no-conditional-in-test error
  if (Math.random()) {
    expect(1).toBe(1);
  }
});
