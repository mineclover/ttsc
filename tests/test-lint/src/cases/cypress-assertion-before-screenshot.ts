// expect: cypress/assertion-before-screenshot error
cy.get("[data-cy=dialog]").should("be.visible");
cy.screenshot();
