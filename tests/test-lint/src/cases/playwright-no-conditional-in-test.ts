import { test } from "@playwright/test";

test("branches", async ({ page }) => {
  const ready = await page.isVisible("main");
  // expect: playwright/no-conditional-in-test error
  if (ready) {
    await page.click("button");
  }
});
