/**
 * Verifies regexp/no-useless-flag: flag that the literal does not exercise.
 *
 * Pins the branch that flags flags whose effect is unobservable on the
 * pattern. The classic case is `/\d+/i` where the `i` (ignore-case) flag does
 * nothing because the pattern contains no letters whose case could vary.
 *
 * 1. Declare a regex literal with a flag that cannot change its match.
 * 2. Assert it is flagged.
 */
// expect: regexp/no-useless-flag error
const value = /\d+/i;

JSON.stringify(value);
