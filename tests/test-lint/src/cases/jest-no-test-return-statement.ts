import { test } from "@jest/globals";

test("returns", () => {
  // expect: jest/no-test-return-statement error
  return 1;
});
