import { describe, beforeEach } from "@jest/globals";

describe("suite", () => {
  beforeEach(() => {});
  // expect: jest/no-duplicate-hooks error
  beforeEach(() => {});
});
