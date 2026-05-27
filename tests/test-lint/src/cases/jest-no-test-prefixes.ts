import { fit, expect } from "@jest/globals";

// expect: jest/no-test-prefixes error
fit("focuses", () => {
  expect(1).toBe(1);
});
