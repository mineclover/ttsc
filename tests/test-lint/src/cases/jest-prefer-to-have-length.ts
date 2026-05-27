import { test, expect } from "@jest/globals";

test("checks length", () => {
  const values = [1, 2, 3];
  // expect: jest/prefer-to-have-length error
  expect(values.length).toBe(3);
});
