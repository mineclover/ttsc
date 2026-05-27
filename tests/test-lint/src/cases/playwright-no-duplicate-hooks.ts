import { test } from "@playwright/test";

test.beforeEach(async () => {});

// expect: playwright/no-duplicate-hooks error
test.beforeEach(async () => {});
