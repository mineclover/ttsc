import assert from "node:assert/strict";
import { test } from "node:test";
import * as Effect from "effect/Effect";
import { observation, increment, incrementTwice, runIncrement } from "../dist/program.js";

test("Effect.sync is lazy until runPromise", async () => {
  assert.equal(observation.calls, 0);
  assert.equal(await runIncrement(), 1);
  assert.equal(observation.calls, 1);
});

test("Effect.gen executes yield* in expected order", async () => {
  const actual = await Effect.runPromise(incrementTwice);
  assert.deepEqual(actual, { first: 2, second: 3 });
  assert.equal(observation.calls, 3);
});
