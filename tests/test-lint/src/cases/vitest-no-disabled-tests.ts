// expect: vitest/no-disabled-tests error
test.skip("temporarily ignored", () => {
  expect(value).toBe(1);
});
