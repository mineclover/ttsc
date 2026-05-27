// expect: vitest/no-conditional-expect error
it("checks conditionally", () => {
  if (ready) {
    expect(value).toBe(1);
  }
});
