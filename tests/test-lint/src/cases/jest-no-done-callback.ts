import { test, expect } from "@jest/globals";

// expect: jest/no-done-callback error
test("finishes later", done => {
  expect(1).toBe(1);
  done();
});
