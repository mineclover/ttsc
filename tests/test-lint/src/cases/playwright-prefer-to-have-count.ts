import { test, expect } from "@playwright/test";

test("counts", async ({ page }) => {
  // expect: playwright/prefer-to-have-count error
  expect(await page.locator("li").count()).toBe(2);
});
