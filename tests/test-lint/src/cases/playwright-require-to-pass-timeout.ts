import { test, expect } from "@playwright/test";

test("passes eventually", async () => {
  // expect: playwright/require-to-pass-timeout error
  await expect(async () => {}).toPass();
});
