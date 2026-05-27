import { test, expect } from "@jest/globals";

test("throws", () => {
  // expect: jest/require-to-throw-message error
  expect(() => { throw new Error("x"); }).toThrow();
});
