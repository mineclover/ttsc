import { test } from "@playwright/test";

test("clicks selector", async ({ page }) => {
  // expect: playwright/prefer-locator error
  await page.click("button");
});
