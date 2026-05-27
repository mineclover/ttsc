import { test, expect } from "@playwright/test";

test("throws", async () => {
  // expect: playwright/require-to-throw-message error
  expect(() => {
    throw new Error("boom");
  }).toThrow();
});
