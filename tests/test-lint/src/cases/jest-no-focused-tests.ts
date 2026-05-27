import { it, expect } from "@jest/globals";

// expect: jest/no-focused-tests error
it.only("runs one test", () => {
  expect(1).toBe(1);
});
