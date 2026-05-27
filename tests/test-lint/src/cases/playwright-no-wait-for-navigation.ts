import { test } from "@playwright/test";

test("waits for navigation", async ({ page }) => {
  // expect: playwright/no-wait-for-navigation error
  await page.waitForNavigation();
});
