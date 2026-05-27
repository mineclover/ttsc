import { test } from "@playwright/test";

// expect: playwright/no-hooks error
test.beforeEach(async () => {});
