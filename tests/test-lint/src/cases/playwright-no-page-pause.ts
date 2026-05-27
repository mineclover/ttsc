import { test } from "@playwright/test";

test("debugs page", async ({ page }) => {
  // expect: playwright/no-page-pause error
  await page.pause();
});
