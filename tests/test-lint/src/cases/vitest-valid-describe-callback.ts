// expect: vitest/valid-describe-callback error
describe("suite", async () => {
  test("case", () => expect(value).toBe(1));
});
