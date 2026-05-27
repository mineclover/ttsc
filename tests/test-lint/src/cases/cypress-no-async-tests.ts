// expect: cypress/no-async-tests error
it("saves", async () => {
  await cy.get("button");
});
