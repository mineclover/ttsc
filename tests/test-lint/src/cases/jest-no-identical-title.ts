import { describe, it, expect } from "@jest/globals";

describe("math", () => {
  it("adds", () => {
    expect(1 + 1).toBe(2);
  });

  // expect: jest/no-identical-title error
  it("adds", () => {
    expect(2 + 2).toBe(4);
  });
});
