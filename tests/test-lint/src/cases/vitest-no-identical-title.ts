// expect: vitest/no-identical-title error
describe("math", () => {
  test("adds", () => expect(add()).toBe(1));
  test("adds", () => expect(add()).toBe(2));
});
