// expect: vitest/no-test-return-statement error
test("returns", () => {
  return buildValue();
});
