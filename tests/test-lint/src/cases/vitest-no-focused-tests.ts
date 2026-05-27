// expect: vitest/no-focused-tests error
test.only("focused", () => {
  expect(value).toBe(1);
});
