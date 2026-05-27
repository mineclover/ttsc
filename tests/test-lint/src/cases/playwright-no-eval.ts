import { test } from "@playwright/test";

test("evaluates against element", async ({ page }) => {
  // expect: playwright/no-eval error
  await page.$eval("#root", (el) => el.textContent);
});
