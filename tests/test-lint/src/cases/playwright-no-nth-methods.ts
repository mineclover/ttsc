import { test } from "@playwright/test";

test("uses position", async ({ page }) => {
  // expect: playwright/no-nth-methods error
  page.getByRole("button").nth(0);
});
