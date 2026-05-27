// expect: vitest/no-conditional-tests error
if (process.env.CI) {
  test("ci only", () => expect(true).toBe(true));
}
