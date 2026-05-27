import { describe, expect } from "@jest/globals";

describe("suite", () => {
  // expect: jest/no-standalone-expect error
  expect(1).toBe(1);
});
