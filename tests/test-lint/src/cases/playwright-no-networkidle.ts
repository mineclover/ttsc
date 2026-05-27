import { test } from "@playwright/test";

test("waits for network", async ({ page }) => {
  // expect: playwright/no-networkidle error
  await page.waitForLoadState("networkidle");
});
