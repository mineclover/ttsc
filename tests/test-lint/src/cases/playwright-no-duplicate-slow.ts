import { test } from "@playwright/test";

test("marks slow twice", async () => {
  test.slow();
  // expect: playwright/no-duplicate-slow error
  test.slow();
});
