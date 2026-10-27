import * as Effect from "effect/Effect";

// Success oracle: creating an Effect must be lazy; executing it must run once.
export const observation = { calls: 0 };
export const increment = Effect.sync(() => ++observation.calls);
export const incrementTwice = Effect.gen(function* () {
  const first = yield* increment;
  const second = yield* increment;
  return { first, second };
});
export const runIncrement = () => Effect.runPromise(increment);
