import { test } from "@playwright/test";

test("forces click", async ({ page }) => {
  // expect: playwright/no-force-option error
  await page.getByRole("button").click({ force: true });
});
