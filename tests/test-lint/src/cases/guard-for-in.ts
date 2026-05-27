// Positive: an unguarded `for...in` body walks the prototype chain and
// processes inherited keys exactly the same as own keys.
function dumpAll(obj: Record<string, unknown>): void {
  // expect: guard-for-in error
  for (const key in obj) {
    console.log(key, obj[key]);
  }
}

// Positive: a guard that lives below another statement is not the very
// first statement of the body, so the inherited-key check still leaks
// the work above it.
function dumpAfterEffect(obj: Record<string, unknown>): void {
  // expect: guard-for-in error
  for (const key in obj) {
    console.log("scanning", key);
    if (Object.hasOwn(obj, key)) {
      console.log(obj[key]);
    }
  }
}

// Negative: `Object.hasOwn(obj, key)` immediately guards the body.
function dumpWithHasOwn(obj: Record<string, unknown>): void {
  for (const key in obj) {
    if (Object.hasOwn(obj, key)) {
      console.log(key, obj[key]);
    }
  }
}

// Negative: `Object.prototype.hasOwnProperty.call(obj, key)` is the
// older guard form and is accepted on the same terms.
function dumpWithHasOwnPropertyCall(obj: Record<string, unknown>): void {
  for (const key in obj) {
    if (Object.prototype.hasOwnProperty.call(obj, key)) {
      console.log(key, obj[key]);
    }
  }
}

// Negative: a `continue` guarded by the negated check is the canonical
// early-skip pattern and is also accepted.
function dumpWithEarlyContinue(obj: Record<string, unknown>): void {
  for (const key in obj) {
    if (!Object.hasOwn(obj, key)) {
      continue;
    }
    console.log(key, obj[key]);
  }
}

JSON.stringify({
  dumpAll: dumpAll({ a: 1 }),
  dumpAfterEffect: dumpAfterEffect({ a: 1 }),
  dumpWithHasOwn: dumpWithHasOwn({ a: 1 }),
  dumpWithHasOwnPropertyCall: dumpWithHasOwnPropertyCall({ a: 1 }),
  dumpWithEarlyContinue: dumpWithEarlyContinue({ a: 1 }),
});
