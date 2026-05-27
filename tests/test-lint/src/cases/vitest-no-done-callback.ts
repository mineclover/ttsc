// expect: vitest/no-done-callback error
it("uses callback", (done) => {
  done();
});
