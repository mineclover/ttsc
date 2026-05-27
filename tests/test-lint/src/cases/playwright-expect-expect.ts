import { test } from "@playwright/test";

// expect: playwright/expect-expect error
test("loads page", async ({ page }) => {
  await page.goto("/");
});
