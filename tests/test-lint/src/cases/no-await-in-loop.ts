declare function getPromise(): Promise<number>;
declare function getAsyncIterator(): AsyncIterableIterator<number>;

// Positive: await inside `for` loop body.
async function inForLoop(): Promise<number> {
  let total = 0;
  for (let i = 0; i < 3; i++) {
    // expect: no-await-in-loop error
    total += await getPromise();
  }
  return total;
}

// Positive: await inside `while` loop body.
async function inWhileLoop(): Promise<number> {
  let total = 0;
  let i = 0;
  while (i < 3) {
    // expect: no-await-in-loop error
    total += await getPromise();
    i++;
  }
  return total;
}

// Positive: await inside `for ... of` (plain, not for-await).
async function inForOfLoop(items: number[]): Promise<number> {
  let total = 0;
  for (const _ of items) {
    // expect: no-await-in-loop error
    total += await getPromise();
  }
  return total;
}

// Negative: await inside `for await ... of` — exempt by design.
async function inForAwaitOfLoop(): Promise<number> {
  let total = 0;
  for await (const value of getAsyncIterator()) {
    total += value;
  }
  return total;
}

// Negative: await inside a nested non-loop async function.
async function awaitNotInLoop(): Promise<number> {
  return await getPromise();
}

// Negative: await inside an inner async function — the inner function
// is the boundary, the outer loop should not be charged. The for body
// itself contains no await.
async function nestedClosureInsideLoop(): Promise<void> {
  for (let i = 0; i < 3; i++) {
    const inner = async (): Promise<number> => {
      return await getPromise();
    };
    JSON.stringify(inner);
  }
}

JSON.stringify({
  inForLoop,
  inWhileLoop,
  inForOfLoop,
  inForAwaitOfLoop,
  awaitNotInLoop,
  nestedClosureInsideLoop,
});
