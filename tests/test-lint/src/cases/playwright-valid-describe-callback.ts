import { test } from "@playwright/test";

// expect: playwright/valid-describe-callback error
test.describe("suite", async () => {});
