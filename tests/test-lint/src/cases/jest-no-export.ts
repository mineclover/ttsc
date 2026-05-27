import { test, expect } from "@jest/globals";

// expect: jest/no-export error
export const helper = 1;

test("uses helper", () => {
  expect(helper).toBe(1);
});
