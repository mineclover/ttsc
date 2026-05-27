import { test } from "@playwright/test";

test("steps", async () => {
  await test.step("outer", async () => {
    // expect: playwright/no-nested-step error
    await test.step("inner", async () => {});
  });
});
