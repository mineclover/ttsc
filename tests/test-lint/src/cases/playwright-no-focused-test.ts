import { test } from "@playwright/test";

// expect: playwright/no-focused-test error
test.only("focuses one case", async ({ page }) => {
  await page.goto("/");
});
