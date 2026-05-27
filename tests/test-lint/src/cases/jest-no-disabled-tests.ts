import { test, expect } from "@jest/globals";

// expect: jest/no-disabled-tests error
test.skip("does not run", () => {
  expect(1).toBe(1);
});
